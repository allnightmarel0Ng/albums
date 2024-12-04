package repository

import (
	"context"

	"github.com/allnightmarel0Ng/albums/internal/domain/model"
	"github.com/allnightmarel0Ng/albums/internal/domain/repository"
)

type ProfileRepository interface {
	GetUserProfile(ctx context.Context, id int) (model.User, []model.Album, error)
	GetArtistProfile(ctx context.Context, id int) (model.Artist, []model.Album, error)
}

type profileRepository struct {
	users   repository.UserRepository
	albums  repository.AlbumRepository
	artists repository.ArtistRepository
}

func NewProfileRepository(users repository.UserRepository, albums repository.AlbumRepository, artists repository.ArtistRepository) ProfileRepository {
	return &profileRepository{
		users:   users,
		albums:  albums,
		artists: artists,
	}
}

func (p *profileRepository) GetUserProfile(ctx context.Context, id int) (model.User, []model.Album, error) {
	select {
	case <-ctx.Done():
		return model.User{}, nil, ctx.Err()
	default:
		user, err := p.users.GetUser(ctx, id)
		if err != nil {
			return model.User{}, nil, err
		}

		purchased, err := p.albums.GetUsersPurchasedAlbums(ctx, id)
		if err != nil {
			return model.User{}, nil, err
		}

		return user, purchased, nil
	}
}

func (p *profileRepository) GetArtistProfile(ctx context.Context, id int) (model.Artist, []model.Album, error) {
	select {
	case <-ctx.Done():
		return model.Artist{}, nil, ctx.Err()
	default:
		albums, err := p.albums.GetArtistsAlbums(ctx, id)
		if err != nil {
			return model.Artist{}, nil, err
		}

		artist, err := p.artists.GetArtistByID(ctx, id)
		if err != nil {
			return model.Artist{}, nil, err
		}

		for i := 0; i < len(albums); i++ {
			albums[i].Author = model.Artist{}
		}

		return artist, albums, nil
	}
}
