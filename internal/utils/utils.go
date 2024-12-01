package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/allnightmarel0Ng/albums/internal/domain/api"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func Send(c *gin.Context, response api.Response) {
	encodedResponse, err := json.Marshal(response)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Internal server error",
		})
		return
	}

	c.Data(response.GetCode(), "application/json", encodedResponse)
}

func Request(ctx context.Context, method, url, auth string, body io.Reader) (*http.Response, error) {
	request, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", auth)

	client := &http.Client{}
	response, err := client.Do(request)
	return response, err
}

func SafelyCastJWTClaim[T any](data jwt.MapClaims, fieldName string) (T, error) {
	var result T

	raw, ok := data[fieldName]
	if !ok {
		return result, fmt.Errorf("invalid jwt token: missing '%s' claim", fieldName)
	}

	result, ok = raw.(T)
	if !ok {
		return result, fmt.Errorf("invalid jwt token: '%s' claim is not a %v", fieldName, reflect.TypeOf(result).String())
	}

	return result, nil
}

func GetJWTClaims(jsonWebToken string, secretKey string) api.Response {
	data := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(jsonWebToken, data, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return &api.AuthorizationResponse{
			Code:  http.StatusInternalServerError,
			Error: err.Error(),
		}
	}

	if !token.Valid {
		return &api.AuthorizationResponse{
			Code:  http.StatusBadRequest,
			Error: "invalid token",
		}
	}

	var result api.AuthorizationResponse
	idFloat, err := SafelyCastJWTClaim[float64](data, "id")
	if err != nil {
		return &api.AuthorizationResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		}
	}
	result.ID = int(idFloat)

	result.IsAdmin, err = SafelyCastJWTClaim[bool](data, "isAdmin")
	if err != nil {
		return &api.AuthorizationResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		}
	}

	result.Code = http.StatusOK
	return &result
}

func InterserviceCommunicationError() api.Response {
	return &api.ErrorResponse{
		Code:  http.StatusInternalServerError,
		Error: "interservice communication error",
	}
}

func DeadlineContext(seconds int) (context.Context, context.CancelFunc) {
	return context.WithDeadline(context.Background(), time.Now().Add(2*time.Second))
}

func GetIDParam(c *gin.Context) (int, error) {
	idStr, ok := c.Params.Get("id")
	if !ok {
		return 0, errors.New("id param not found")
	}

	id, err := strconv.Atoi(idStr)
	return id, err
}
