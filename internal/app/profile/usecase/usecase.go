package usecase

import (
	"log"
	"net/http"

	"github.com/allnightmarel0Ng/albums/internal/app/profile/repository"
	"github.com/allnightmarel0Ng/albums/internal/domain/api"
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
	user, err := p.repo.GetUserProfile(id)
	if err != nil {
		log.Print(err.Error())
		return &api.UserProfileResponse{
			Code:  http.StatusNotFound,
			Error: "unable to find such profile",
		}
	}

	return &api.UserProfileResponse{
		Code: http.StatusOK,
		User: user,
	}
}

func (p *profileUseCase) GetArtistProfile(id int) api.Response {
	albums, err := p.repo.GetArtistProfile(id)
	if err != nil {
		log.Print(err.Error())
		return &api.ArtistProfileResponse{
			Code:  http.StatusNotFound,
			Error: "unable to find such artist",
		}
	}

	return &api.ArtistProfileResponse{
		Code:     http.StatusOK,
		Albums:   albums,
	}
}