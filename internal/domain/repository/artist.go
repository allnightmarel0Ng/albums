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

	selectArtistsLikeName =
	/* sql */ `SELECT
					id,
					name,
					genre,
					image_url
				FROM public.artists
				WHERE name ILIKE $1;`

	selectRandomNArtistsSQL =
	/* sql */ `SELECT
					id,
					name,
					genre,
					image_url
				FROM public.artists
				ORDER BY RANDOM()
				LIMIT $1;`
)

type ArtistRepository interface {
	GetArtistByID(ctx context.Context, id int) (model.Artist, error)
	GetArtistsLikeName(ctx context.Context, name string) ([]model.Artist, error)
	GetRandomNArtists(ctx context.Context, count uint) ([]model.Artist, error)
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

func (a *artistRepository) GetArtistsLikeName(ctx context.Context, name string) ([]model.Artist, error) {
	var result []model.Artist
	rows, err := a.db.Query(ctx, selectArtistsLikeName, name)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var current model.Artist
		err := rows.Scan(&current.ID, &current.Name, &current.Genre, &current.ImageURL)
		if err != nil {
			return nil, err
		}
		result = append(result, current)
	}

	return result, nil
}

func (a *artistRepository) GetRandomNArtists(ctx context.Context, count uint) ([]model.Artist, error) {
	rows, err := a.db.Query(ctx, selectRandomNArtistsSQL, count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]model.Artist, count)
	index := 0
	for rows.Next() {
		err := rows.Scan(&result[index].ID, &result[index].Name, &result[index].Genre, &result[index].ImageURL)
		if err != nil {
			return nil, err
		}

		index++
	}

	return result, nil
}
