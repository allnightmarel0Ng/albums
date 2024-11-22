package handler

import (
	"log"
	"net/http"

	"github.com/allnightmarel0Ng/albums/internal/app/gateway/usecase"
	"github.com/allnightmarel0Ng/albums/internal/utils"
	"github.com/gin-gonic/gin"
)

type GatewayHandler interface {
	HandleLogin(c *gin.Context)
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
	response := g.useCase.Authorization(c.GetHeader("Authorization"))

	log.Printf("response: %v", response)
	
	switch {
	case response.Code == http.StatusOK:
		utils.Send(c, response.Code, "jwt", response.Jwt)
	default:
		utils.Send(c, response.Code, "error", response.Error)
	}

}
