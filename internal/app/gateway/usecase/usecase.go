package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/allnightmarel0Ng/albums/internal/app/gateway/repository"
	"github.com/allnightmarel0Ng/albums/internal/domain/api"
	"github.com/allnightmarel0Ng/albums/internal/utils"
)

type GatewayUseCase interface {
	Authentication(authHeader string) api.Response
	// MainPage(authHeader string) api.Response
	UserProfile(jsonWebToken string) api.Response
	ArtistProfile(id int) api.Response
}

type gatewayUseCase struct {
	repo              repository.GatewayRepository
	authorizationPort string
	profilePort       string
	jwtSecretKey      string
}

func NewGatewayUseCase(repo repository.GatewayRepository, authorizationPort, profilePort, jwtSecretKey string) GatewayUseCase {
	return &gatewayUseCase{
		repo:              repo,
		authorizationPort: authorizationPort,
		profilePort:       profilePort,
		jwtSecretKey:      jwtSecretKey,
	}
}

func (g *gatewayUseCase) Authentication(authHeader string) api.Response {
	ctx, cancel := utils.DeadlineContext(10)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("http://authorization:%s/authenticate", g.authorizationPort), nil)
	if err != nil {
		return utils.InterserviceCommunicationError()
	}
	request.Header.Set("Authorization", authHeader)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return utils.InterserviceCommunicationError()
	}
	defer response.Body.Close()

	var result api.AuthenticationResponse
	json.NewDecoder(response.Body).Decode(&result)
	return &result
}

func (g *gatewayUseCase) UserProfile(jsonWebToken string) api.Response {
	code, claims, err := g.authorize(jsonWebToken)
	if err != nil {
		return utils.InterserviceCommunicationError()
	}

	if code != http.StatusOK {
		return &api.ErrorResponse{
			Code:  code,
			Error: claims.Error,
		}
	}

	ctx, cancel := utils.DeadlineContext(10)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("http://profile:%s/users/%d", g.profilePort, claims.ID), nil)
	if err != nil {
		return utils.InterserviceCommunicationError()
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return utils.InterserviceCommunicationError()
	}
	defer response.Body.Close()

	var result api.UserProfileResponse
	json.NewDecoder(response.Body).Decode(&result)
	return &result
}

func (g *gatewayUseCase) ArtistProfile(id int) api.Response {
	ctx, cancel := utils.DeadlineContext(10)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("http://profile:%s/artists/%d", g.profilePort, id), nil)
	if err != nil {
		return utils.InterserviceCommunicationError()
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return utils.InterserviceCommunicationError()
	}
	defer response.Body.Close()

	var result api.ArtistProfileResponse
	json.NewDecoder(response.Body).Decode(&result)
	return &result
}

// func (g *gatewayUseCase) MainPage(jwtToken string) api.Response {
// 	claims, err := utils.GetJWTClaims(jwtToken, g.jwtSecretKey)
// 	if err != nil {
// 		return &api.AlbumsResponse{
// 			Code:  http.StatusUnauthorized,
// 			Error: err.Error(),
// 		}
// 	}

// 	// userID, err := utils.SafelyCastJWTClaim[float64](data, "id")
// 	// if err != nil {
// 	// 	return &api.AlbumsResponse{
// 	// 		Code:  http.StatusUnprocessableEntity,
// 	// 		Error: err.Error(),
// 	// 	}
// 	// }

// 	switch claims.IsAdmin {
// 	case false:
// 		return g.nonAdminMainPage()
// 	default:
// 		return &api.AlbumsResponse{
// 			Code:  http.StatusNotImplemented,
// 			Error: "not yet implemented",
// 		}
// 	}
// }

// func (g *gatewayUseCase) nonAdminMainPage() api.Response {
// 	ctx, cancel := utils.DeadlineContext(2)
// 	defer cancel()

// 	albums, err := g.repo.GetAllAlbums(ctx)
// 	if err != nil {
// 		return &api.AlbumsResponse{
// 			Code:  http.StatusInternalServerError,
// 			Error: "retrieving from db error",
// 		}
// 	}

// 	if albums == nil {
// 		return &api.AlbumsResponse{
// 			Code:  http.StatusNotFound,
// 			Error: "no albums found",
// 		}
// 	}

// 	return &api.AlbumsResponse{
// 		Code: http.StatusOK,
// 		Data: albums,
// 	}
// }

func (g *gatewayUseCase) authorize(jsonWebToken string) (int, api.AuthorizationResponse, error) {
	ctx, cancel := utils.DeadlineContext(10)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("http://authorization:%s/authorize", g.authorizationPort), nil)
	if err != nil {
		return http.StatusInternalServerError, api.AuthorizationResponse{}, errors.New("interservice communication error")
	}
	request.Header.Set("Authorization", "Bearer "+jsonWebToken)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return http.StatusInternalServerError, api.AuthorizationResponse{}, errors.New("interservice communication error")
	}
	defer response.Body.Close()

	var result api.AuthorizationResponse
	json.NewDecoder(response.Body).Decode(&result)
	return response.StatusCode, result, nil
}
