package usecase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"

	"github.com/allnightmarel0Ng/albums/internal/domain/api"
	"github.com/allnightmarel0Ng/albums/internal/infrastructure/kafka"
	"github.com/allnightmarel0Ng/albums/internal/utils"
)

type GatewayUseCase interface {
	Authentication(authHeader string) (int, []byte)
	Logout(authHeader string) (int, []byte)
	Register(body io.Reader) (int, []byte)

	MainPage(body io.Reader) (int, []byte)
	Search(body io.Reader) (int, []byte)

	UserProfile(jsonWebToken string) (int, []byte)
	ArtistProfile(params string) (int, []byte)
	AlbumProfile(params string) (int, []byte)

	AddToOrder(authHeader string, albumID int) (int, []byte)
	RemoveFromOrder(authHeader string, albumID int) (int, []byte)
	UserOrders(jsonWebToken string) (int, []byte)

	Deposit(authHeader string, diff uint) api.Response
	Buy(authHeader string) api.Response

	Logs(authHeader string, params string) (int, []byte)
	DeleteAlbum(authHeader string, params string) (int, []byte)
	SaveDump(authHeader string) (int, []byte)
	LoadDump(authHeader, filePath string) (int, []byte)
	AuthorizeAdmin(authHeader string) (int, []byte)
}

type gatewayUseCase struct {
	producer *kafka.Producer

	authorizationPort   string
	profilePort         string
	orderManagementPort string
	searchEnginePort    string
	adminPanelPort      string

	jwtSecretKey string

	postgresUser     string
	postgresPassword string
	postgresPort     string
	postgresDB       string
}

func NewGatewayUseCase(
	producer *kafka.Producer,
	authorizationPort,
	profilePort,
	orderManagementPort,
	searchEnginePort,
	adminPanelPort,
	jwtSecretKey,
	postgresUser,
	postgresPassword,
	postgresPort,
	postgresDB string) GatewayUseCase {
	return &gatewayUseCase{
		producer:            producer,
		authorizationPort:   authorizationPort,
		profilePort:         profilePort,
		orderManagementPort: orderManagementPort,
		searchEnginePort:    searchEnginePort,
		adminPanelPort:      adminPanelPort,
		jwtSecretKey:        jwtSecretKey,

		postgresUser:     postgresUser,
		postgresPassword: postgresPassword,
		postgresPort:     postgresPort,
		postgresDB:       postgresDB,
	}
}

func (g *gatewayUseCase) Authentication(authHeader string) (int, []byte) {
	return utils.RequestAndParseResponse("GET", fmt.Sprintf("http://authorization:%s/authenticate", g.authorizationPort), authHeader, nil)
}

func (g *gatewayUseCase) Logout(authHeader string) (int, []byte) {
	return utils.RequestAndParseResponse("POST", fmt.Sprintf("http://authorization:%s/logout", g.authorizationPort), authHeader, nil)
}

func (g *gatewayUseCase) Register(body io.Reader) (int, []byte) {
	return utils.RequestAndParseResponse("POST", fmt.Sprintf("http://authorization:%s/registration", g.authorizationPort), "", body)
}

func (g *gatewayUseCase) UserProfile(authHeader string) (int, []byte) {
	authorizationResponse := utils.Authorize(authHeader, g.authorizationPort)
	if authorizationResponse.GetCode() != http.StatusOK {
		raw, _ := json.Marshal(authorizationResponse)
		return authorizationResponse.GetCode(), raw
	}

	claims := authorizationResponse.(*api.AuthorizationResponse)

	return utils.RequestAndParseResponse("GET", fmt.Sprintf("http://profile:%s/users/%d", g.profilePort, claims.ID), "", nil)
}

func (g *gatewayUseCase) ArtistProfile(params string) (int, []byte) {
	return utils.RequestAndParseResponse("GET", fmt.Sprintf("http://profile:%s/artists/%s", g.profilePort, params), "", nil)
}

func (g *gatewayUseCase) AlbumProfile(params string) (int, []byte) {
	return utils.RequestAndParseResponse("GET", fmt.Sprintf("http://profile:%s/albums/%s", g.profilePort, params), "", nil)
}

func (g *gatewayUseCase) AddToOrder(authHeader string, albumID int) (int, []byte) {
	return g.orderAction(albumID, authHeader, "add")
}

func (g *gatewayUseCase) RemoveFromOrder(authHeader string, albumID int) (int, []byte) {
	return g.orderAction(albumID, authHeader, "remove")
}

func (g *gatewayUseCase) Deposit(authHeader string, diff uint) api.Response {
	authResponse := utils.Authorize(authHeader, g.authorizationPort)
	if authResponse.GetCode() != http.StatusOK {
		return authResponse
	}

	claims := authResponse.(*api.AuthorizationResponse)

	operation := api.MoneyOperationKafkaMessage{
		Type:   api.Deposit,
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
	authResponse := utils.Authorize(authHeader, g.authorizationPort)
	if authResponse.GetCode() != http.StatusOK {
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
	err = json.NewDecoder(response.Body).Decode(&orderResponse)
	if err != nil {
		return utils.InterserviceCommunicationError()
	}

	if orderResponse.Order.Orderer.Balance < orderResponse.Order.TotalPrice {
		return &api.ErrorResponse{
			Code:  http.StatusBadRequest,
			Error: "not enough money on balance",
		}
	}

	operation := api.MoneyOperationKafkaMessage{
		Type:    api.Buy,
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

func (g *gatewayUseCase) UserOrders(authHeader string) (int, []byte) {
	authorizationResponse := utils.Authorize(authHeader, g.authorizationPort)
	if authorizationResponse.GetCode() != http.StatusOK {
		raw, _ := json.Marshal(authorizationResponse)
		return authorizationResponse.GetCode(), raw
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

func (g *gatewayUseCase) Logs(authHeader string, params string) (int, []byte) {
	adminAuthorizationCode, raw := g.AuthorizeAdmin(authHeader)
	if adminAuthorizationCode != http.StatusOK {
		return adminAuthorizationCode, raw
	}

	return utils.RequestAndParseResponse("GET", fmt.Sprintf("http://admin-panel:%s/logs/%s", g.adminPanelPort, params), "", nil)
}

func (g *gatewayUseCase) DeleteAlbum(authHeader string, params string) (int, []byte) {
	adminAuthorizationCode, raw := g.AuthorizeAdmin(authHeader)
	if adminAuthorizationCode != http.StatusOK {
		return adminAuthorizationCode, raw
	}

	return utils.RequestAndParseResponse("DELETE", fmt.Sprintf("http://admin-panel:%s/delete/%s", g.adminPanelPort, params), "", nil)
}

func (g *gatewayUseCase) SaveDump(authHeader string) (int, []byte) {
	adminAuthorizationCode, raw := g.AuthorizeAdmin(authHeader)
	if adminAuthorizationCode != http.StatusOK {
		return adminAuthorizationCode, raw
	}

	cmd := exec.Command("pg_dump", "--clean", "-U", g.postgresUser, "-h", "postgres", "-p", g.postgresPort, g.postgresDB)
	cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", g.postgresPassword))

	output, err := cmd.Output()
	if err != nil {
		raw, _ := json.Marshal(api.ErrorResponse{
			Code:  http.StatusInternalServerError,
			Error: "unable to create dump",
		})
		return http.StatusInternalServerError, raw
	}

	return http.StatusOK, output
}

func (g *gatewayUseCase) LoadDump(authHeader, filePath string) (int, []byte) {
	cmd := exec.Command("psql", "-U", g.postgresUser, "-h", "postgres", "-p", g.postgresPort, g.postgresDB, "-f", filePath)

	cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", g.postgresPassword))

	output, err := cmd.CombinedOutput()
	if err != nil {
		raw, _ := json.Marshal(api.ErrorResponse{
			Code:  http.StatusInternalServerError,
			Error: "unable to load sql dump",
		})
		return http.StatusInternalServerError, raw
	}

	// log.Print(string(output))
	return http.StatusOK, output
}

func (g *gatewayUseCase) AuthorizeAdmin(authHeader string) (int, []byte) {
	authorizationResponse := utils.Authorize(authHeader, g.authorizationPort)
	if authorizationResponse.GetCode() != http.StatusOK {
		raw, _ := json.Marshal(authorizationResponse)
		return authorizationResponse.GetCode(), raw
	}

	claims := authorizationResponse.(*api.AuthorizationResponse)
	if !claims.IsAdmin {
		raw, _ := json.Marshal(api.ErrorResponse{
			Code:  http.StatusUnauthorized,
			Error: "non-admin user cannot do that",
		})
		return http.StatusUnauthorized, raw
	}

	return http.StatusOK, nil
}

func (g *gatewayUseCase) orderAction(albumID int, authHeader string, action string) (int, []byte) {
	authorizationResponse := utils.Authorize(authHeader, g.authorizationPort)
	if authorizationResponse.GetCode() != http.StatusOK {
		raw, _ := json.Marshal(authorizationResponse)
		return authorizationResponse.GetCode(), raw
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
