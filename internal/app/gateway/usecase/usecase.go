package usecase

import (
	"bytes"
	"encoding/json"
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
	AddToOrder(albumID int, jsonWebToken string) api.Response
	RemoveFromOrder(albumID int, jsonWebToken string) api.Response
	UserOrders(jsonWebToken string) api.Response
}

type gatewayUseCase struct {
	repo                repository.GatewayRepository
	authorizationPort   string
	profilePort         string
	orderManagementPort string
	jwtSecretKey        string
}

func NewGatewayUseCase(repo repository.GatewayRepository, authorizationPort, profilePort, orderManagementPort, jwtSecretKey string) GatewayUseCase {
	return &gatewayUseCase{
		repo:                repo,
		authorizationPort:   authorizationPort,
		profilePort:         profilePort,
		orderManagementPort: orderManagementPort,
		jwtSecretKey:        jwtSecretKey,
	}
}

func (g *gatewayUseCase) Authentication(authHeader string) api.Response {
	ctx, cancel := utils.DeadlineContext(10)
	defer cancel()

	response, err := utils.Request(ctx, "GET", fmt.Sprintf("http://authorization:%s/authenticate", g.authorizationPort), authHeader, nil)
	if err != nil {
		return utils.InterserviceCommunicationError()
	}
	defer response.Body.Close()

	var result api.AuthenticationResponse
	json.NewDecoder(response.Body).Decode(&result)
	return &result
}

func (g *gatewayUseCase) UserProfile(jsonWebToken string) api.Response {
	code, authorizationResponse := g.authorize(jsonWebToken)
	if code != http.StatusOK {
		return authorizationResponse
	}

	claims := authorizationResponse.(*api.AuthorizationResponse)

	ctx, cancel := utils.DeadlineContext(10)
	defer cancel()

	response, err := utils.Request(ctx, "GET", fmt.Sprintf("http://profile:%s/users/%d", g.profilePort, claims.ID), "", nil)
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

	response, err := utils.Request(ctx, "GET", fmt.Sprintf("http://profile:%s/artists/%d", g.profilePort, id), "", nil)
	if err != nil {
		return utils.InterserviceCommunicationError()
	}
	defer response.Body.Close()

	var result api.ArtistProfileResponse
	json.NewDecoder(response.Body).Decode(&result)
	return &result
}

func (g *gatewayUseCase) AddToOrder(albumID int, jsonWebToken string) api.Response {
	return g.orderAction(albumID, jsonWebToken, "add")
}

func (g *gatewayUseCase) RemoveFromOrder(albumID int, jsonWebToken string) api.Response {
	return g.orderAction(albumID, jsonWebToken, "remove")
}

func (g *gatewayUseCase) orderAction(albumID int, jsonWebToken string, action string) api.Response {
	code, authorizationResponse := g.authorize(jsonWebToken)
	if code != http.StatusOK {
		return authorizationResponse
	}

	claims := authorizationResponse.(*api.AuthorizationResponse)
	body, err := json.Marshal(api.OrderActionRequest{
		UserID:  claims.ID,
		AlbumID: albumID,
	})
	if err != nil {
		return utils.InterserviceCommunicationError()
	}

	ctx, cancel := utils.DeadlineContext(10)
	defer cancel()

	response, err := utils.Request(ctx, "POST", fmt.Sprintf("http://order-management:%s/%s", g.orderManagementPort, action), "", bytes.NewReader(body))
	if err != nil {
		return utils.InterserviceCommunicationError()
	}
	defer response.Body.Close()

	var result api.OrderActionResponse
	json.NewDecoder(response.Body).Decode(&result)
	return &result
}

func (g *gatewayUseCase) UserOrders(jsonWebToken string) api.Response {
	code, authorizationResponse := g.authorize(jsonWebToken)
	if code != http.StatusOK {
		return authorizationResponse
	}

	claims := authorizationResponse.(*api.AuthorizationResponse)

	ctx, cancel := utils.DeadlineContext(10)
	defer cancel()

	response, err := utils.Request(ctx, "GET", fmt.Sprintf("http://order-management:%s/orders/%d", g.orderManagementPort, claims.ID), "", nil)
	if err != nil {
		return utils.InterserviceCommunicationError()
	}
	defer response.Body.Close()

	var result api.UserOrdersResponse
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

func (g *gatewayUseCase) authorize(jsonWebToken string) (int, api.Response) {
	ctx, cancel := utils.DeadlineContext(10)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("http://authorization:%s/authorize", g.authorizationPort), nil)
	if err != nil {
		return http.StatusInternalServerError, utils.InterserviceCommunicationError()
	}
	request.Header.Set("Authorization", "Bearer "+jsonWebToken)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return http.StatusInternalServerError, utils.InterserviceCommunicationError()
	}
	defer response.Body.Close()

	var result api.AuthorizationResponse
	json.NewDecoder(response.Body).Decode(&result)

	if response.StatusCode != http.StatusOK {
		return http.StatusUnauthorized, &api.ErrorResponse{
			Code:  http.StatusUnauthorized,
			Error: result.Error,
		}
	}
	return http.StatusOK, &result
}
