package handler

import (
	"net/http"
	"strings"

	"github.com/allnightmarel0Ng/albums/internal/app/gateway/usecase"
	"github.com/allnightmarel0Ng/albums/internal/domain/model"
	"github.com/allnightmarel0Ng/albums/internal/utils"
	"github.com/gin-gonic/gin"
)

type GatewayHandler interface {
	HandleLogin(c *gin.Context)
	HandleMainPage(c *gin.Context)
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
	utils.Send(c, g.useCase.Authorization(c.GetHeader("Authorization")))
}

func (g *gatewayHandler) HandleMainPage(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		utils.Send(c, &model.AuthorizationResponse{
			Code:  http.StatusBadRequest,
			Error: "wrong auth token",
		})
	}

	utils.Send(c, g.useCase.MainPage(authHeader[len("Bearer "):]))
}
