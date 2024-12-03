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
	// HandleMainPage(c *gin.Context)
	HandleUserProfile(c *gin.Context)
	HandleArtistProfile(c *gin.Context)
	HandleOrderAdd(c *gin.Context)
	HandleOrderRemove(c *gin.Context)
	HandleOrders(c *gin.Context)
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
	utils.Send(c, g.useCase.Authentication(c.GetHeader("Authorization")))
}

// func (g *gatewayHandler) HandleMainPage(c *gin.Context) {
// 	authHeader := c.GetHeader("Authorization")

// 	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
// 		utils.Send(c, &api.AuthorizationResponse{
// 			Code:  http.StatusBadRequest,
// 			Error: "wrong auth token",
// 		})
// 	}

// 	utils.Send(c, g.useCase.MainPage(authHeader[len("Bearer "):]))
// }

func (g *gatewayHandler) HandleUserProfile(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		utils.Send(c, &api.AuthorizationResponse{
			Code:  http.StatusBadRequest,
			Error: "wrong auth token",
		})
		return
	}

	utils.Send(c, g.useCase.UserProfile(authHeader[len("Bearer "):]))
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

	utils.Send(c, g.useCase.ArtistProfile(id))
}

func (g *gatewayHandler) HandleOrderAdd(c *gin.Context) {
	handleOrderAction(c, g.useCase.AddToOrder)
}

func (g *gatewayHandler) HandleOrderRemove(c *gin.Context) {
	handleOrderAction(c, g.useCase.RemoveFromOrder)
}

func handleOrderAction(c *gin.Context, callback func(int, string) api.Response) {
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

	utils.Send(c, callback(id, authHeader[len("Bearer "):]))
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

	utils.Send(c, g.useCase.UserOrders(authHeader[len("Bearer "):]))
}
