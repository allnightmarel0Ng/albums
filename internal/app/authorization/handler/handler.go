package handler

import (
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

	utils.Send(c, a.useCase.Logout(authData[len("Bearer "):]))
}
