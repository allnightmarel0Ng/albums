package repository

import (
	"github.com/allnightmarel0Ng/albums/internal/domain/model"
	"github.com/allnightmarel0Ng/albums/internal/domain/repository"
)

type GatewayRepository interface {
	GetAllAlbums() ([]model.Album, error)
}

type gatewayRepository struct {
	albums repository.AlbumRepository
}

func NewGatewayRepository(albums repository.AlbumRepository) GatewayRepository {
	return &gatewayRepository{
		albums: albums,
	}
}

func (g *gatewayRepository) GetAllAlbums() ([]model.Album, error) {
	return g.albums.GetAllAlbumsLike("")
}
