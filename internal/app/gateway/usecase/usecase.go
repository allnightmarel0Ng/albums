package usecase

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/allnightmarel0Ng/albums/internal/app/gateway/repository"
	"github.com/allnightmarel0Ng/albums/internal/domain/api"
	"github.com/allnightmarel0Ng/albums/internal/utils"
)

type GatewayUseCase interface {
	Authorization(authHeader string) api.Response
	MainPage(authHeader string) api.Response
	UserProfile(jwtToken string) api.Response
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

func (g *gatewayUseCase) Authorization(authHeader string) api.Response {
	request, err := http.NewRequest("GET", fmt.Sprintf("http://authorization:%s/", g.authorizationPort), nil)
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

	var result api.AuthorizationResponse
	json.NewDecoder(response.Body).Decode(&result)
	return &result
}

func (g *gatewayUseCase) UserProfile(jwtToken string) api.Response {
	claims, err := utils.GetJWTClaims(jwtToken, g.jwtSecretKey)
	if err != nil {
		return &api.AlbumsResponse{
			Code:  http.StatusUnauthorized,
			Error: err.Error(),
		}
	}

	request, err := http.NewRequest("GET", fmt.Sprintf("http://profile:%s/users/%d", g.profilePort, claims.ID), nil)
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
	request, err := http.NewRequest("GET", fmt.Sprintf("http://profile:%s/artists/%d", g.profilePort, id), nil)
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

func (g *gatewayUseCase) MainPage(jwtToken string) api.Response {
	claims, err := utils.GetJWTClaims(jwtToken, g.jwtSecretKey)
	if err != nil {
		return &api.AlbumsResponse{
			Code:  http.StatusUnauthorized,
			Error: err.Error(),
		}
	}

	// userID, err := utils.SafelyCastJWTClaim[float64](data, "id")
	// if err != nil {
	// 	return &api.AlbumsResponse{
	// 		Code:  http.StatusUnprocessableEntity,
	// 		Error: err.Error(),
	// 	}
	// }

	switch claims.IsAdmin {
	case false:
		return g.nonAdminMainPage()
	default:
		return &api.AlbumsResponse{
			Code:  http.StatusNotImplemented,
			Error: "not yet implemented",
		}
	}
}

func (g *gatewayUseCase) nonAdminMainPage() api.Response {
	albums, err := g.repo.GetAllAlbums()
	if err != nil {
		return &api.AlbumsResponse{
			Code:  http.StatusInternalServerError,
			Error: "retrieving from db error",
		}
	}

	if albums == nil {
		return &api.AlbumsResponse{
			Code:  http.StatusNotFound,
			Error: "no albums found",
		}
	}

	return &api.AlbumsResponse{
		Code: http.StatusOK,
		Data: albums,
	}
}
