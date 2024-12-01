package repository

import (
	"context"

	"github.com/allnightmarel0Ng/albums/internal/domain/model"
	"github.com/allnightmarel0Ng/albums/internal/domain/repository"
)

type GatewayRepository interface {
	GetAllAlbums(ctx context.Context) ([]model.Album, error)
}

type gatewayRepository struct {
	albums repository.AlbumRepository
}

func NewGatewayRepository(albums repository.AlbumRepository) GatewayRepository {
	return &gatewayRepository{
		albums: albums,
	}
}

func (g *gatewayRepository) GetAllAlbums(ctx context.Context) ([]model.Album, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return g.albums.GetAllAlbumsLike(ctx, "")
	}
}
