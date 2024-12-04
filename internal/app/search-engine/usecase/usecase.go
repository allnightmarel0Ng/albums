package usecase

import (
	"net/http"

	"github.com/allnightmarel0Ng/albums/internal/app/search-engine/repository"
	"github.com/allnightmarel0Ng/albums/internal/domain/api"
	"github.com/allnightmarel0Ng/albums/internal/utils"
)

type SearchEngineUseCase interface {
	SearchEntities(query string) api.Response
	RandomEntities(artistsCount, albumsCount uint) api.Response
}

type searchEngineUseCase struct {
	repo repository.SearchEngineRepository
}

func NewSearchEngineUseCase(repo repository.SearchEngineRepository) SearchEngineUseCase {
	return &searchEngineUseCase{
		repo: repo,
	}
}

func (s *searchEngineUseCase) SearchEntities(query string) api.Response {
	query = utils.SearchLikeString(query)

	ctx, cancel := utils.DeadlineContext(5)
	defer cancel()

	artists, albums, err := s.repo.GetEntitiesLikeName(ctx, query)
	if err != nil {
		return &api.ErrorResponse{
			Code:  http.StatusInternalServerError,
			Error: "db unexpected error",
		}
	}

	return &api.SearchEngineResponse{
		Code:    http.StatusOK,
		Artists: artists,
		Albums:  albums,
	}
}

func (s *searchEngineUseCase) RandomEntities(artistsCount, albumsCount uint) api.Response {
	ctx, cancel := utils.DeadlineContext(5)
	defer cancel()

	artists, albums, err := s.repo.GetRandomNEntities(ctx, artistsCount, albumsCount)
	if err != nil {
		return &api.ErrorResponse{
			Code:  http.StatusInternalServerError,
			Error: "db unexpected error",
		}
	}

	return &api.SearchEngineResponse{
		Code:    http.StatusOK,
		Artists: artists,
		Albums:  albums,
	}
}
