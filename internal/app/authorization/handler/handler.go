package handler

import (
	"log"
	"net/http"
	"strings"

	"github.com/allnightmarel0Ng/albums/internal/app/authorization/usecase"
	"github.com/allnightmarel0Ng/albums/internal/domain/api"
	"github.com/allnightmarel0Ng/albums/internal/utils"
	"github.com/gin-gonic/gin"
)

type AuthorizationHandler interface {
	HandleAuthentication(c *gin.Context)
	HandleAuthorization(c *gin.Context)
	HandleLogout(c *gin.Context)
	HandleRegistration(c *gin.Context)
}

type authorizationHandler struct {
	useCase usecase.AuthorizationUseCase
}

func NewAuthorizationHandler(useCase usecase.AuthorizationUseCase) AuthorizationHandler {
	return &authorizationHandler{
		useCase: useCase,
	}
}

func (a *authorizationHandler) HandleAuthentication(c *gin.Context) {
	authData := c.GetHeader("Authorization")
	log.Print(authData)
	if authData == "" || !strings.HasPrefix(authData, "Basic ") {
		utils.Send(c, &api.AuthorizationResponse{
			Code:  http.StatusBadRequest,
			Error: "bad authorization base64 token",
		})
		return
	}

	utils.Send(c, a.useCase.Authenticate(authData[len("Basic "):]))
}

func (a *authorizationHandler) HandleAuthorization(c *gin.Context) {
	authData := c.GetHeader("Authorization")
	if authData == "" || !strings.HasPrefix(authData, "Bearer ") {
		utils.Send(c, &api.ErrorResponse{
			Code:  http.StatusBadRequest,
			Error: "bad authorization base64 token",
		})
		return
	}

	utils.Send(c, a.useCase.Authorize(authData[len("Bearer "):]))
}

func (a *authorizationHandler) HandleLogout(c *gin.Context) {
	authData := c.GetHeader("Authorization")
	if authData == "" || !strings.HasPrefix(authData, "Bearer ") {
		utils.Send(c, &api.ErrorResponse{
			Code:  http.StatusBadRequest,
			Error: "bad authorization base64 token",
		})
		return
	}

	response := a.useCase.Logout(authData[len("Bearer "):])
	if response != nil {
		utils.Send(c, response)
		return
	}

	c.String(http.StatusOK, "")
}

func (a *authorizationHandler) HandleRegistration(c *gin.Context) {
	var request api.RegistrationRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.Send(c, &api.ErrorResponse{
			Code:  http.StatusBadRequest,
			Error: "invalid request fields",
		})
		log.Print(err.Error())
		return
	}

	response := a.useCase.Register(request)
	if response != nil {
		utils.Send(c, response)
		return
	}

	c.String(http.StatusOK, "")
}
