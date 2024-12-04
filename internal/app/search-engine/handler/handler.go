package handler

import (
	"net/http"

	"github.com/allnightmarel0Ng/albums/internal/app/search-engine/usecase"
	"github.com/allnightmarel0Ng/albums/internal/domain/api"
	"github.com/allnightmarel0Ng/albums/internal/utils"
	"github.com/gin-gonic/gin"
)

type SearchEngineHandler interface {
	HandleSearch(c *gin.Context)
	HandleRandom(c *gin.Context)
}

type searchEngineHandler struct {
	useCase usecase.SearchEngineUseCase
}

func NewSearchEngineHandler(useCase usecase.SearchEngineUseCase) SearchEngineHandler {
	return &searchEngineHandler{
		useCase: useCase,
	}
}

func (s *searchEngineHandler) HandleSearch(c *gin.Context) {
	var request api.SearchRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.Send(c, &api.ErrorResponse{
			Code:  http.StatusBadRequest,
			Error: "invalid search request",
		})
		return
	}

	utils.Send(c, s.useCase.SearchEntities(request.Query))
}

func (s *searchEngineHandler) HandleRandom(c *gin.Context) {
	var request api.RandomEntitiesRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.Send(c, &api.ErrorResponse{
			Code:  http.StatusBadRequest,
			Error: "invalid random entities request",
		})
		return
	}

	utils.Send(c, s.useCase.RandomEntities(request.ArtistsCount, request.AlbumsCount))
}
