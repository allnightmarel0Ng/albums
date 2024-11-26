package handler

import (
	"errors"
	"net/http"
	"strconv"

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
	id, err := getIDParam(c)
	if err != nil {
		utils.Send(c, &api.ArtistProfileResponse{
			Code:  http.StatusBadRequest,
			Error: "unable to parse id param: " + err.Error(),
		})
	}

	utils.Send(c, p.useCase.GetArtistProfile(id))
}

func (p *profileHandler) HandleUserProfile(c *gin.Context) {
	id, err := getIDParam(c)
	if err != nil {
		utils.Send(c, &api.ArtistProfileResponse{
			Code:  http.StatusBadRequest,
			Error: "unable to parse id param: " + err.Error(),
		})
	}

	utils.Send(c, p.useCase.GetUserProfile(id))
}

func getIDParam(c *gin.Context) (int, error) {
	idStr, ok := c.Params.Get("id")
	if !ok {
		return 0, errors.New("id param not found")
	}

	id, err := strconv.Atoi(idStr)
	return id, err
}
