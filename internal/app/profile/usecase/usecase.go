package usecase

import (
	"context"
	"log"
	"net/http"

	"github.com/allnightmarel0Ng/albums/internal/app/profile/repository"
	"github.com/allnightmarel0Ng/albums/internal/domain/api"
	"github.com/allnightmarel0Ng/albums/internal/utils"
)

type ProfileUseCase interface {
	GetUserProfile(id int) api.Response
	GetArtistProfile(id int) api.Response
	GetAlbumProfile(id int) api.Response
	GetAlbumOwnersIds(albumID int) api.Response
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
	ctx, cancel := utils.DeadlineContext(5)
	defer cancel()

	user, purchased, err := p.repo.GetUserProfile(ctx, id)
	if err != nil {
		switch err {
		case context.DeadlineExceeded:
			return &api.ErrorResponse{
				Code:  http.StatusInternalServerError,
				Error: "database communication error",
			}
		default:
			return &api.ErrorResponse{
				Code:  http.StatusNotFound,
				Error: "unable to find such profile",
			}
		}
	}

	return &api.UserProfileResponse{
		Code:      http.StatusOK,
		User:      user,
		Purchased: purchased,
	}
}

func (p *profileUseCase) GetArtistProfile(id int) api.Response {
	ctx, cancel := utils.DeadlineContext(5)
	defer cancel()

	artist, albums, err := p.repo.GetArtistProfile(ctx, id)
	if err != nil {
		switch err {
		case context.DeadlineExceeded:
			return &api.ErrorResponse{
				Code:  http.StatusInternalServerError,
				Error: "database communication error",
			}
		default:
			return &api.ErrorResponse{
				Code:  http.StatusNotFound,
				Error: "unable to find such profile",
			}
		}
	}

	return &api.ArtistProfileResponse{
		Code:   http.StatusOK,
		Artist: artist,
		Albums: albums,
	}
}

func (p *profileUseCase) GetAlbumProfile(id int) api.Response {
	ctx, cancel := utils.DeadlineContext(5)
	defer cancel()

	album, err := p.repo.GetAlbumProfile(ctx, id)
	if err != nil {
		switch err {
		case context.DeadlineExceeded:
			return &api.ErrorResponse{
				Code:  http.StatusInternalServerError,
				Error: "database communication error",
			}
		default:
			return &api.ErrorResponse{
				Code:  http.StatusNotFound,
				Error: "unable to find such profile",
			}
		}
	}

	return &api.AlbumProfileResponse{
		Code:  http.StatusOK,
		Album: album,
	}
}

func (p *profileUseCase) GetAlbumOwnersIds(albumID int) api.Response {
	ctx, cancel := utils.DeadlineContext(5)
	defer cancel()

	name, ids, err := p.repo.GetAlbumOwnersIds(ctx, albumID)
	if err != nil {
		switch err {
		case context.DeadlineExceeded:
			return &api.ErrorResponse{
				Code:  http.StatusInternalServerError,
				Error: "database communication error",
			}
		default:
			return &api.ErrorResponse{
				Code:  http.StatusNotFound,
				Error: "unable to find such album",
			}
		}
	}

	log.Print(ids)
	return &api.AlbumOwnersResponse{
		Code:      http.StatusOK,
		Ids:       ids,
		AlbumName: name,
	}
}
