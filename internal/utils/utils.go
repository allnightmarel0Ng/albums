package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/allnightmarel0Ng/albums/internal/domain/api"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func SendRaw(c *gin.Context, code int, response []byte) {
	c.Data(code, "application/json", response)
}

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

func RequestAndParseResponse(method, url, auth string, body io.Reader) (int, []byte) {
	ctx, cancel := DeadlineContext(10)
	defer cancel()

	response, err := Request(ctx, method, url, auth, body)
	if err != nil {
		return InterserviceCommunicationErrorRaw()
	}
	defer response.Body.Close()

	result, err := io.ReadAll(response.Body)
	if err != nil {
		return InterserviceCommunicationErrorRaw()
	}
	return response.StatusCode, result
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

func InterserviceCommunicationErrorRaw() (int, []byte) {
	result, _ := json.Marshal(api.ErrorResponse{
		Code:  http.StatusInternalServerError,
		Error: "interservice communication error",
	})

	return http.StatusInternalServerError, result
}

func DeadlineContext(seconds int) (context.Context, context.CancelFunc) {
	return context.WithDeadline(context.Background(), time.Now().Add(time.Duration(seconds)*time.Second))
}

func GetParam[T any](c *gin.Context, name string) (T, error) {
	var result T

	paramStr, ok := c.Params.Get(name)
	if !ok {
		return result, fmt.Errorf("param not found")
	}

	switch any(result).(type) {
	case int:
		val, err := strconv.Atoi(paramStr)
		if err != nil {
			return result, err
		}
		result = any(val).(T)
	default:
		return result, fmt.Errorf("unsupported type: %T", result)
	}

	return result, nil
}

func SearchLikeString(str string) string {
	result := strings.Replace(str, " ", "%", -1)
	result = "%" + result
	result += "%"
	result = strings.ToLower(result)
	return result
}
