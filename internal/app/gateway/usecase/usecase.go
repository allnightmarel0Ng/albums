package usecase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/allnightmarel0Ng/albums/internal/domain/api"
	"github.com/allnightmarel0Ng/albums/internal/infrastructure/kafka"
	"github.com/allnightmarel0Ng/albums/internal/utils"
)

type GatewayUseCase interface {
	Authentication(authHeader string) (int, []byte)
	MainPage(request io.Reader) (int, []byte)
	Search(body io.Reader) (int, []byte)
	UserProfile(jsonWebToken string) (int, []byte)
	ArtistProfile(id int) (int, []byte)
	AddToOrder(albumID int, jsonWebToken string) (int, []byte)
	RemoveFromOrder(albumID int, jsonWebToken string) (int, []byte)
	UserOrders(jsonWebToken string) (int, []byte)
	Deposit(authHeader string, diff uint) api.Response
	Buy(authHeader string) api.Response
}

type gatewayUseCase struct {
	producer            *kafka.Producer
	authorizationPort   string
	profilePort         string
	orderManagementPort string
	searchEnginePort    string
	jwtSecretKey        string
}

func NewGatewayUseCase(producer *kafka.Producer, authorizationPort, profilePort, orderManagementPort, searchEnginePort, jwtSecretKey string) GatewayUseCase {
	return &gatewayUseCase{
		producer:            producer,
		authorizationPort:   authorizationPort,
		profilePort:         profilePort,
		orderManagementPort: orderManagementPort,
		searchEnginePort:    searchEnginePort,
		jwtSecretKey:        jwtSecretKey,
	}
}

func (g *gatewayUseCase) Authentication(authHeader string) (int, []byte) {
	return utils.RequestAndParseResponse("GET", fmt.Sprintf("http://authorization:%s/authenticate", g.authorizationPort), authHeader, nil)
}

func (g *gatewayUseCase) UserProfile(jsonWebToken string) (int, []byte) {
	code, authorizationResponse := g.authorize(jsonWebToken)
	if code != http.StatusOK {
		raw, _ := json.Marshal(authorizationResponse)
		return code, raw
	}

	claims := authorizationResponse.(*api.AuthorizationResponse)

	return utils.RequestAndParseResponse("GET", fmt.Sprintf("http://profile:%s/users/%d", g.profilePort, claims.ID), "", nil)
}

func (g *gatewayUseCase) ArtistProfile(id int) (int, []byte) {
	return utils.RequestAndParseResponse("GET", fmt.Sprintf("http://profile:%s/artists/%d", g.profilePort, id), "", nil)
}

func (g *gatewayUseCase) AddToOrder(albumID int, jsonWebToken string) (int, []byte) {
	return g.orderAction(albumID, jsonWebToken, "add")
}

func (g *gatewayUseCase) RemoveFromOrder(albumID int, jsonWebToken string) (int, []byte) {
	return g.orderAction(albumID, jsonWebToken, "remove")
}

func (g *gatewayUseCase) Deposit(authHeader string, diff uint) api.Response {
	if diff <= 0 {
		return &api.ErrorResponse{
			Code:  http.StatusBadRequest,
			Error: "trying to deposit negative amount of money",
		}
	}

	code, authResponse := g.authorize(authHeader)
	if code != http.StatusOK {
		return authResponse
	}

	claims := authResponse.(*api.AuthorizationResponse)

	operation := api.MoneyOperationKafkaMessage{
		Type:   "deposit",
		UserID: claims.ID,
		Diff:   diff,
	}

	raw, err := json.Marshal(operation)
	if err != nil {
		return &api.ErrorResponse{
			Code:  http.StatusInternalServerError,
			Error: "unable to deposit money",
		}
	}

	err = g.producer.Produce("money-operations", raw)
	if err != nil {
		return utils.InterserviceCommunicationError()
	}

	return nil
}

func (g *gatewayUseCase) Buy(authHeader string) api.Response {
	code, authResponse := g.authorize(authHeader)
	if code != http.StatusOK {
		return authResponse
	}

	claims := authResponse.(*api.AuthorizationResponse)

	ctx, cancel := utils.DeadlineContext(10)
	defer cancel()

	response, err := utils.Request(ctx, "GET", fmt.Sprintf("http://order-management:%s/orders/%d?unpaidOnly=true", g.orderManagementPort, claims.ID), "", nil)
	if err != nil {
		return utils.InterserviceCommunicationError()
	}

	if response.StatusCode != http.StatusOK {
		return &api.ErrorResponse{
			Code:  http.StatusExpectationFailed,
			Error: "multiple unpaid orders found or no orders found",
		}
	}

	var orderResponse api.UnpaidUserOrderResponse
	json.NewDecoder(response.Body).Decode(&orderResponse)

	operation := api.MoneyOperationKafkaMessage{
		Type:    "buy",
		UserID:  claims.ID,
		OrderID: orderResponse.Order.ID,
	}

	raw, err := json.Marshal(operation)
	if err != nil {
		return &api.ErrorResponse{
			Code:  http.StatusInternalServerError,
			Error: "unable to buy order",
		}
	}

	err = g.producer.Produce("money-operations", raw)
	if err != nil {
		return utils.InterserviceCommunicationError()
	}

	return nil
}

func (g *gatewayUseCase) orderAction(albumID int, authHeader string, action string) (int, []byte) {
	code, authorizationResponse := g.authorize(authHeader)
	if code != http.StatusOK {
		raw, _ := json.Marshal(authorizationResponse)
		return code, raw
	}

	claims := authorizationResponse.(*api.AuthorizationResponse)
	body, err := json.Marshal(api.OrderActionRequest{
		UserID:  claims.ID,
		AlbumID: albumID,
	})
	if err != nil {
		return utils.InterserviceCommunicationErrorRaw()
	}

	return utils.RequestAndParseResponse("POST", fmt.Sprintf("http://order-management:%s/%s", g.orderManagementPort, action), "", bytes.NewReader(body))
}

func (g *gatewayUseCase) UserOrders(authHeader string) (int, []byte) {
	code, authorizationResponse := g.authorize(authHeader)
	if code != http.StatusOK {
		raw, _ := json.Marshal(authorizationResponse)
		return code, raw
	}

	claims := authorizationResponse.(*api.AuthorizationResponse)

	return utils.RequestAndParseResponse("GET", fmt.Sprintf("http://order-management:%s/orders/%d", g.orderManagementPort, claims.ID), "", nil)
}

func (g *gatewayUseCase) MainPage(body io.Reader) (int, []byte) {
	return utils.RequestAndParseResponse("POST", fmt.Sprintf("http://search-engine:%s/random", g.searchEnginePort), "", body)
}

func (g *gatewayUseCase) Search(body io.Reader) (int, []byte) {
	return utils.RequestAndParseResponse("POST", fmt.Sprintf("http://search-engine:%s/search", g.searchEnginePort), "", body)
}

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

func (g *gatewayUseCase) authorize(authHeader string) (int, api.Response) {
	ctx, cancel := utils.DeadlineContext(10)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("http://authorization:%s/authorize", g.authorizationPort), nil)
	if err != nil {
		return http.StatusInternalServerError, utils.InterserviceCommunicationError()
	}
	request.Header.Set("Authorization", authHeader)

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
