package handler

import (
	"net/http"

	"github.com/allnightmarel0Ng/albums/internal/app/profile/usecase"
	"github.com/allnightmarel0Ng/albums/internal/domain/api"
	"github.com/allnightmarel0Ng/albums/internal/utils"
	"github.com/gin-gonic/gin"
)

type ProfileHandler interface {
	HandleArtistProfile(c *gin.Context)
	HandleUserProfile(c *gin.Context)
}

type profileHandler struct {
	useCase usecase.ProfileUseCase
}

func NewProfileHandler(useCase usecase.ProfileUseCase) ProfileHandler {
	return &profileHandler{
		useCase: useCase,
	}
}

func (p *profileHandler) HandleArtistProfile(c *gin.Context) {
	id, err := utils.GetIDParam(c)
	if err != nil {
		utils.Send(c, &api.ErrorResponse{
			Code:  http.StatusBadRequest,
			Error: "unable to parse id param: " + err.Error(),
		})
		return
	}

	utils.Send(c, p.useCase.GetArtistProfile(id))
}

func (p *profileHandler) HandleUserProfile(c *gin.Context) {
	id, err := utils.GetIDParam(c)
	if err != nil {
		utils.Send(c, &api.ErrorResponse{
			Code:  http.StatusBadRequest,
			Error: "unable to parse id param: " + err.Error(),
		})
		return
	}

	utils.Send(c, p.useCase.GetUserProfile(id))
}
