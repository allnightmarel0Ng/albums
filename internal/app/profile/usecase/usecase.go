package usecase

import (
	"log"
	"net/http"

	"github.com/allnightmarel0Ng/albums/internal/app/profile/repository"
	"github.com/allnightmarel0Ng/albums/internal/domain/api"
	"github.com/allnightmarel0Ng/albums/internal/utils"
)

type ProfileUseCase interface {
	GetUserProfile(id int) api.Response
	GetArtistProfile(id int) api.Response
}

type profileUseCase struct {
	repo repository.ProfileRepository
}

func NewProfileUseCase(repo repository.ProfileRepository) ProfileUseCase {
	return &profileUseCase{
		repo: repo,
	}
}

func (p *profileUseCase) GetUserProfile(id int) api.Response {
	ctx, cancel := utils.DeadlineContext(2)
	defer cancel()

	user, purchased, err := p.repo.GetUserProfile(ctx, id)
	if err != nil {
		log.Print(err.Error())
		return &api.ErrorResponse{
			Code:  http.StatusNotFound,
			Error: "unable to find such profile",
		}
	}

	return &api.UserProfileResponse{
		Code:      http.StatusOK,
		User:      user,
		Purchased: purchased,
	}
}

func (p *profileUseCase) GetArtistProfile(id int) api.Response {
	ctx, cancel := utils.DeadlineContext(2)
	defer cancel()

	artist, albums, err := p.repo.GetArtistProfile(ctx, id)
	if err != nil {
		log.Print(err.Error())
		return &api.ErrorResponse{
			Code:  http.StatusNotFound,
			Error: "unable to find such artist",
		}
	}

	return &api.ArtistProfileResponse{
		Code:   http.StatusOK,
		Artist: artist,
		Albums: albums,
	}
}
