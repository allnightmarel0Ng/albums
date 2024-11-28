package utils

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/allnightmarel0Ng/albums/internal/domain/api"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func Send(c *gin.Context, response api.Response) {
	k, v := response.GetKeyValue()
	c.JSON(response.GetCode(), gin.H{
		"code": response.GetCode(),
		k:      v,
	})
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

func GetJWTClaims(jwtToken string, secretKey string) (api.JWTClaims, error) {
	data := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(jwtToken, data, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return api.JWTClaims{}, err
	}

	if !token.Valid {
		return api.JWTClaims{}, errors.New("invalid token")
	}

	var result api.JWTClaims
	idFloat, err := SafelyCastJWTClaim[float64](data, "id")
	if err != nil {
		return result, err
	}
	result.ID = int(idFloat)

	expFloat, err := SafelyCastJWTClaim[float64](data, "exp")
	if err != nil {
		return result, err
	}
	result.Exp = int64(expFloat)

	if time.Now().After(time.Unix(result.Exp, 0)) {
		return result, errors.New("authorization token has expired")
	}

	result.IsAdmin, err = SafelyCastJWTClaim[bool](data, "isAdmin")

	return result, err
}

func InterserviceCommunicationError() api.Response {
	return &api.ErrorResponse{
		Code:  http.StatusInternalServerError,
		Error: "interservice communication error",
	}
}

func DeadlineContext(seconds int) (context.Context, context.CancelFunc) {
	return context.WithDeadline(context.Background(), time.Now().Add(2 * time.Second))
}