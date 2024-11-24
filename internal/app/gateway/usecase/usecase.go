package usecase

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/allnightmarel0Ng/albums/internal/app/gateway/repository"
	"github.com/allnightmarel0Ng/albums/internal/domain/model"
	"github.com/golang-jwt/jwt/v4"
)

type GatewayUseCase interface {
	Authorization(authHeader string) model.Response
	MainPage(authHeader string) model.Response
}

type gatewayUseCase struct {
	repo              repository.GatewayRepository
	authorizationPort string
	jwtSecretKey      string
}

func NewGatewayUseCase(repo repository.GatewayRepository, authorizationPort, jwtSecretKey string) GatewayUseCase {
	return &gatewayUseCase{
		repo:              repo,
		authorizationPort: authorizationPort,
		jwtSecretKey:      jwtSecretKey,
	}
}

func (g *gatewayUseCase) interprocessCommunicationError() model.Response {
	return &model.AuthorizationResponse{
		Code:  http.StatusInternalServerError,
		Error: "interprocess communication error",
	}
}

func (g *gatewayUseCase) Authorization(authHeader string) model.Response {
	request, err := http.NewRequest("GET", fmt.Sprintf("http://authorization:%s/", g.authorizationPort), nil)
	if err != nil {
		return g.interprocessCommunicationError()
	}

	request.Header.Set("Authorization", authHeader)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return g.interprocessCommunicationError()
	}
	defer response.Body.Close()

	var result model.AuthorizationResponse
	json.NewDecoder(response.Body).Decode(&result)
	log.Printf("got from auth microservice: %v", result)
	return &result
}

func (g *gatewayUseCase) MainPage(jwtToken string) model.Response {
	data := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(jwtToken, data, func(token *jwt.Token) (interface{}, error) {
		return []byte(g.jwtSecretKey), nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return &model.AlbumsResponse{
				Code:  http.StatusUnauthorized,
				Error: "invalid token signature",
			}
		}
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return &model.AlbumsResponse{
					Code:  http.StatusUnauthorized,
					Error: "token has expired",
				}
			}
		}
		return &model.AlbumsResponse{
			Code:  http.StatusBadRequest,
			Error: "unparsable jwt",
		}
	}

	if !token.Valid {
		return &model.AlbumsResponse{
			Code:  http.StatusUnauthorized,
			Error: "invalid token",
		}
	}

	exp, ok := data["exp"]
	if !ok {
		return &model.AlbumsResponse{
			Code:  http.StatusUnprocessableEntity,
			Error: "invalid jwt token: missing 'exp' claim",
		}
	}

	expFloat, ok := exp.(float64)
	if !ok {
		return &model.AlbumsResponse{
			Code:  http.StatusUnprocessableEntity,
			Error: "invalid jwt token: 'exp' claim is not a number",
		}
	}

	if time.Now().After(time.Unix(int64(expFloat), 0)) {
		return &model.AlbumsResponse{
			Code:  http.StatusUnauthorized,
			Error: "authorization time has expired",
		}
	}

	role, ok := data["role"]
	if !ok {
		return &model.AlbumsResponse{
			Code:  http.StatusUnprocessableEntity,
			Error: "invalid jwt token",
		}
	}

	switch role.(string) {
	case "customer":
		return g.customerMainPage()
	default:
		return &model.AlbumsResponse{
			Code:  http.StatusNotImplemented,
			Error: "not yet implemented",
		}
	}
}

func (g *gatewayUseCase) customerMainPage() model.Response {
	albums, err := g.repo.GetAllAlbums()
	if err != nil {
		return &model.AlbumsResponse{
			Code:  http.StatusInternalServerError,
			Error: "retrieving from db error",
		}
	}

	if albums == nil {
		return &model.AlbumsResponse{
			Code:  http.StatusNotFound,
			Error: "no albums found",
		}
	}

	return &model.AlbumsResponse{
		Code: http.StatusOK,
		Data: albums,
	}
}
