package repository

import (
	"context"

	"github.com/allnightmarel0Ng/albums/internal/domain/model"
	"github.com/allnightmarel0Ng/albums/internal/infrastructure/postgres"
)

const (
	selectArtistByID =
	/* sql */ `SELECT
					id,
					name,
					genre,
					image_url
				FROM public.artists
				WHERE id = $1;`
)

type ArtistRepository interface {
	GetArtistByID(ctx context.Context, id int) (model.Artist, error)
}

type artistRepository struct {
	db postgres.Database
}

func NewArtistRepository(db postgres.Database) ArtistRepository {
	return &artistRepository{
		db: db,
	}
}

func (a *artistRepository) GetArtistByID(ctx context.Context, id int) (model.Artist, error) {
	var result model.Artist
	err := a.db.QueryRow(ctx, selectArtistByID, id).Scan(&result.ID, &result.Name, &result.Genre, &result.ImageURL)
	return result, err
}
