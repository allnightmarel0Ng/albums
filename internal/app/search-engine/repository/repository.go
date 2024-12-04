package repository

import (
	"context"

	"github.com/allnightmarel0Ng/albums/internal/domain/model"
	"github.com/allnightmarel0Ng/albums/internal/domain/repository"
)

type SearchEngineRepository interface {
	GetRandomNEntities(ctx context.Context, artistsCount uint, albumsCount uint) ([]model.Artist, []model.Album, error)
	GetEntitiesLikeName(ctx context.Context, name string) ([]model.Artist, []model.Album, error)
}

type searchEngineRepository struct {
	artists repository.ArtistRepository
	albums  repository.AlbumRepository
}

func NewSearchEngineRepository(artists repository.ArtistRepository, albums repository.AlbumRepository) SearchEngineRepository {
	return &searchEngineRepository{
		artists: artists,
		albums:  albums,
	}
}

func (s *searchEngineRepository) GetRandomNEntities(ctx context.Context, artistsCount uint, albumsCount uint) ([]model.Artist, []model.Album, error) {
	select {
	case <-ctx.Done():
		return nil, nil, ctx.Err()
	default:
		artists, err := s.artists.GetRandomNArtists(ctx, artistsCount)
		if err != nil {
			return nil, nil, err
		}

		albums, err := s.albums.GetRandomNAlbums(ctx, albumsCount)
		return artists, albums, err
	}
}

func (s *searchEngineRepository) GetEntitiesLikeName(ctx context.Context, name string) ([]model.Artist, []model.Album, error) {
	select {
	case <-ctx.Done():
		return nil, nil, ctx.Err()
	default:
		artists, err := s.artists.GetArtistsLikeName(ctx, name)
		if err != nil {
			return nil, nil, err
		}

		albums, err := s.albums.GetAlbumsLikeName(ctx, name)
		return artists, albums, err
	}
}
