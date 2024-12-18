package repository

import (
	"context"

	"github.com/allnightmarel0Ng/albums/internal/domain/model"
	"github.com/allnightmarel0Ng/albums/internal/domain/repository"
)

type ProfileRepository interface {
	GetUserProfile(ctx context.Context, id int) (model.User, []model.Album, error)
	GetArtistProfile(ctx context.Context, id int) (model.Artist, []model.Album, error)
	GetAlbumProfile(ctx context.Context, id int) (model.Album, error)
	GetAlbumOwnersIds(ctx context.Context, albumId int) (string, []int, error)
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
			albums[i].Author = nil
		}

		return artist, albums, nil
	}
}

func (p *profileRepository) GetAlbumProfile(ctx context.Context, id int) (model.Album, error) {
	select {
	case <-ctx.Done():
		return model.Album{}, ctx.Err()
	default:
		return p.albums.GetAlbumByID(ctx, id)
	}
}

func (p *profileRepository) GetAlbumOwnersIds(ctx context.Context, albumID int) (string, []int, error) {
	select {
	case <-ctx.Done():
		return "", nil, ctx.Err()
	default:
		name, err := p.albums.GetAlbumName(ctx, albumID)
		if err != nil {
			return "", nil, err
		}

		ids, err := p.users.GetAlbumOwnersIds(ctx, albumID)
		return name, ids, err
	}
}
