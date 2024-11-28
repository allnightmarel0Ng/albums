package repository

import (
	"context"

	"github.com/allnightmarel0Ng/albums/internal/domain/model"
	"github.com/allnightmarel0Ng/albums/internal/domain/repository"
)

type ProfileRepository interface {
	GetUserProfile(ctx context.Context, id int) (model.User, error)
	GetArtistProfile(ctx context.Context, id int) ([]model.Album, error)
}

type profileRepository struct {
	users  repository.UserRepository
	albums repository.AlbumRepository
}

func NewProfileRepository(users repository.UserRepository, albums repository.AlbumRepository, artists repository.ArtistRepository) ProfileRepository {
	return &profileRepository{
		users:  users,
		albums: albums,
	}
}

func (p *profileRepository) GetUserProfile(ctx context.Context, id int) (model.User, error) {
	user, err := p.users.GetUser(ctx, id)
	if err != nil {
		return model.User{}, err
	}

	user.Purchased, err = p.albums.GetUsersPurchasedAlbums(ctx, id)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (p *profileRepository) GetArtistProfile(ctx context.Context, id int) ([]model.Album, error) {
	return p.albums.GetArtistsAlbums(ctx, id)
}
