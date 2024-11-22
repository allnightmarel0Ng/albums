package handler

import (
	"net/http"
	"strings"

	"github.com/allnightmarel0Ng/albums/internal/app/authorization/usecase"
	"github.com/allnightmarel0Ng/albums/internal/utils"
	"github.com/gin-gonic/gin"
)

type AuthorizationHandler interface {
	HandleAuthorization(c *gin.Context)
}

type authorizationHandler struct {
	useCase usecase.AuthorizationUseCase
}

func NewAuthorizationHandler(useCase usecase.AuthorizationUseCase) AuthorizationHandler {
	return &authorizationHandler{
		useCase: useCase,
	}
}

func (a *authorizationHandler) HandleAuthorization(c *gin.Context) {
	authData := c.GetHeader("Authorization")
	if authData == "" || !strings.HasPrefix(authData, "Basic ") {
		utils.Send(c, http.StatusBadRequest, "error", "bad authorization base64 token")
		return
	}

	jwt, code, err := a.useCase.Authorize(authData[len("Basic"):])
	if err != nil {
		utils.Send(c, code, "error", err.Error())
		return
	}

	utils.Send(c, code, "jwt", jwt)
}
