package handler

import (
	"net/http"
	"strings"

	"github.com/allnightmarel0Ng/albums/internal/app/gateway/usecase"
	"github.com/allnightmarel0Ng/albums/internal/domain/api"
	"github.com/allnightmarel0Ng/albums/internal/utils"
	"github.com/gin-gonic/gin"
)

type GatewayHandler interface {
	HandleLogin(c *gin.Context)
	HandleMainPage(c *gin.Context)
	HandleSearch(c *gin.Context)
	HandleUserProfile(c *gin.Context)
	HandleArtistProfile(c *gin.Context)
	HandleOrderAdd(c *gin.Context)
	HandleOrderRemove(c *gin.Context)
	HandleOrders(c *gin.Context)
	HandleDeposit(c *gin.Context)
	HandleBuy(c *gin.Context)
}

type gatewayHandler struct {
	useCase usecase.GatewayUseCase
}

func NewGatewayHandler(useCase usecase.GatewayUseCase) GatewayHandler {
	return &gatewayHandler{
		useCase: useCase,
	}
}

func (g *gatewayHandler) HandleLogin(c *gin.Context) {
	code, raw := g.useCase.Authentication(c.GetHeader("Authorization"))
	utils.SendRaw(c, code, raw)
}

func (g *gatewayHandler) HandleMainPage(c *gin.Context) {
	code, raw := g.useCase.MainPage(c.Request.Body)
	utils.SendRaw(c, code, raw)
}

func (g *gatewayHandler) HandleSearch(c *gin.Context) {
	code, raw := g.useCase.Search(c.Request.Body)
	utils.SendRaw(c, code, raw)
}

func (g *gatewayHandler) HandleUserProfile(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		utils.Send(c, &api.AuthorizationResponse{
			Code:  http.StatusBadRequest,
			Error: "wrong auth token",
		})
		return
	}

	code, raw := g.useCase.UserProfile(authHeader)
	utils.SendRaw(c, code, raw)
}

func (g *gatewayHandler) HandleArtistProfile(c *gin.Context) {
	id, err := utils.GetIDParam(c)
	if err != nil {
		utils.Send(c, &api.ErrorResponse{
			Code:  http.StatusBadRequest,
			Error: "invalid id parameter",
		})
		return
	}

	code, raw := g.useCase.ArtistProfile(id)
	utils.SendRaw(c, code, raw)
}

func (g *gatewayHandler) HandleOrderAdd(c *gin.Context) {
	handleOrderAction(c, g.useCase.AddToOrder)
}

func (g *gatewayHandler) HandleOrderRemove(c *gin.Context) {
	handleOrderAction(c, g.useCase.RemoveFromOrder)
}

func handleOrderAction(c *gin.Context, callback func(int, string) (int, []byte)) {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		utils.Send(c, &api.AuthorizationResponse{
			Code:  http.StatusBadRequest,
			Error: "wrong auth token",
		})
		return
	}

	id, err := utils.GetIDParam(c)
	if err != nil {
		utils.Send(c, &api.ErrorResponse{
			Code:  http.StatusBadRequest,
			Error: "invalid id parameter",
		})
		return
	}

	code, raw := callback(id, authHeader)
	utils.SendRaw(c, code, raw)
}

func (g *gatewayHandler) HandleOrders(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		utils.Send(c, &api.AuthorizationResponse{
			Code:  http.StatusBadRequest,
			Error: "wrong auth token",
		})
		return
	}

	code, raw := g.useCase.UserOrders(authHeader)
	utils.SendRaw(c, code, raw)
}

func (g *gatewayHandler) HandleDeposit(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		utils.Send(c, &api.AuthorizationResponse{
			Code:  http.StatusBadRequest,
			Error: "wrong auth token",
		})
		return
	}

	var request api.DepositRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.Send(c, &api.ErrorResponse{
			Code:  http.StatusBadRequest,
			Error: "invalid body in request",
		})
		return
	}

	response := g.useCase.Deposit(authHeader, request.Money)
	if response != nil {
		utils.Send(c, response)
		return
	}

	c.String(http.StatusOK, "")
}

func (g *gatewayHandler) HandleBuy(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		utils.Send(c, &api.AuthorizationResponse{
			Code:  http.StatusBadRequest,
			Error: "wrong auth token",
		})
		return
	}

	response := g.useCase.Buy(authHeader)
	if response != nil {
		utils.Send(c, response)
		return
	}

	c.String(http.StatusOK, "")
}
