package repository

import (
	"context"
	"fmt"

	"github.com/allnightmarel0Ng/albums/internal/domain/model"
	"github.com/allnightmarel0Ng/albums/internal/infrastructure/postgres"
)

const (
	selectUsersPurchasedAlbumsSQL =
	/* sql */ `SELECT
					a.id,
					a.name,
					ar.id,
					ar.name,
					ar.genre,
					ar.image_url,
					a.image_url,
					a.price,
					t.id,
					t.name,
					t.number
				FROM public.purchased_albums AS pu
				JOIN public.albums AS a ON pu.album_id = a.id
				JOIN public.artists AS ar ON a.artist_id = ar.id
				RIGHT JOIN public.tracks AS t ON t.album_id = a.id
				WHERE pu.user_id = $1
				ORDER BY a.name, t.number;`

	selectAlbumsLikeNameSQL =
	/* sql */ `SELECT 
					a.id,
					a.name,
					ar.id,
					ar.name,
					ar.genre,
					ar.image_url,
					a.image_url,
					a.price,
					t.id,
					t.name,
					t.number
				FROM public.tracks AS t
				LEFT JOIN public.albums AS a ON t.album_id = a.id
				JOIN public.artists AS ar ON a.artist_id = ar.id
				WHERE a.name ILIKE $1
				ORDER BY t.number;`

	selectArtistsAlbumsSQL =
	/* sql */ `SELECT
					a.id,
					a.name,
					ar.id,
					ar.name,
					ar.genre,
					ar.image_url,
					a.image_url,
					a.price,
					t.id,
					t.name,
					t.number
				FROM public.tracks AS t
				LEFT JOIN public.albums AS a ON t.album_id = a.id
				JOIN public.artists AS ar ON a.artist_id = ar.id
				WHERE ar.id = $1
				ORDER BY t.number;`

	selectRandomNAlbumsSQL =
	/* sql */ `SELECT 
					a.id,
					a.name,
					ar.id,
					ar.name,
					ar.genre,
					ar.image_url,
					a.image_url,
					a.price,
					t.id,
					t.name,
					t.number
				FROM public.tracks AS t
				LEFT JOIN public.albums AS a ON t.album_id = a.id
				JOIN public.artists AS ar ON a.artist_id = ar.id
				ORDER BY RANDOM()
				LIMIT $1;`

	deleteAlbumSQL =
	/* sql */ `DELETE FROM public.albums
				WHERE id = $1;`

	selectAlbumByIdSQL =
	/* sql */ `SELECT
					a.id,
					a.name,
					ar.id,
					ar.name,
					ar.genre,
					ar.image_url,
					a.image_url,
					a.price,
					t.id,
					t.name,
					t.number
				FROM public.tracks AS t
				LEFT JOIN public.albums AS a ON t.album_id = a.id
				JOIN public.artists AS ar ON a.artist_id = ar.id
				WHERE a.id = $1
				ORDER BY t.number;`

	selectAlbumNameSQL =
	/* sql */ `SELECT name
				FROM public.albums
				WHERE id = $1;`
)

type AlbumRepository interface {
	GetUsersPurchasedAlbums(ctx context.Context, userID int) ([]model.Album, error)
	GetAlbumsLikeName(ctx context.Context, name string) ([]model.Album, error)
	GetArtistsAlbums(ctx context.Context, artistID int) ([]model.Album, error)
	GetRandomNAlbums(ctx context.Context, count uint) ([]model.Album, error)
	DeleteAlbum(ctx context.Context, albumID int) error
	GetAlbumByID(ctx context.Context, albumID int) (model.Album, error)
	GetAlbumName(ctx context.Context, albumID int) (string, error)
}

type albumRepository struct {
	db postgres.Database
}

func NewAlbumRepository(db postgres.Database) AlbumRepository {
	return &albumRepository{
		db: db,
	}
}

func albumsFromRows(rows postgres.Rows) ([]model.Album, error) {
	albumsMap := make(map[int]*model.Album)
	var ids []int

	for rows.Next() {
		var (
			album  model.Album
			track  model.Track
			author model.Artist
		)

		err := rows.Scan(&album.ID, &album.Name, &author.ID,
			&author.Name, &author.Genre, &author.ImageURL, &album.ImageURL, &album.Price,
			&track.ID, &track.Name, &track.Number)
		if err != nil {
			return nil, err
		}

		album.Author = &author

		_, ok := albumsMap[album.ID]
		if !ok {
			album.Tracks = []model.Track{track}
			albumsMap[album.ID] = &album
			ids = append(ids, album.ID)
		} else {
			albumsMap[album.ID].Tracks = append(albumsMap[album.ID].Tracks, track)
		}
	}

	result := make([]model.Album, len(ids))
	for i, id := range ids {
		result[i] = *albumsMap[id]
	}

	return result, nil
}

func (a *albumRepository) GetUsersPurchasedAlbums(ctx context.Context, userID int) ([]model.Album, error) {
	rows, err := a.db.Query(ctx, selectUsersPurchasedAlbumsSQL, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return albumsFromRows(rows)
}

func (a *albumRepository) GetAlbumsLikeName(ctx context.Context, name string) ([]model.Album, error) {
	rows, err := a.db.Query(ctx, selectAlbumsLikeNameSQL, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return albumsFromRows(rows)
}

func (a *albumRepository) GetArtistsAlbums(ctx context.Context, artistID int) ([]model.Album, error) {
	rows, err := a.db.Query(ctx, selectArtistsAlbumsSQL, artistID)
	if err != nil {
		return nil, err
	}

	return albumsFromRows(rows)
}

func (a *albumRepository) GetRandomNAlbums(ctx context.Context, count uint) ([]model.Album, error) {
	rows, err := a.db.Query(ctx, selectRandomNAlbumsSQL, count)
	if err != nil {
		return nil, err
	}

	return albumsFromRows(rows)
}

func (a *albumRepository) DeleteAlbum(ctx context.Context, albumID int) error {
	return a.db.Exec(ctx, deleteAlbumSQL, albumID)
}

func (a *albumRepository) GetAlbumByID(ctx context.Context, albumID int) (model.Album, error) {
	rows, err := a.db.Query(ctx, selectAlbumByIdSQL, albumID)
	if err != nil {
		return model.Album{}, err
	}
	defer rows.Close()

	result, err := albumsFromRows(rows)
	if err != nil {
		return model.Album{}, err
	}

	if len(result) != 1 {
		return model.Album{}, fmt.Errorf("album not found")
	}

	return result[0], nil
}

func (a *albumRepository) GetAlbumName(ctx context.Context, id int) (string, error) {
	var result string
	err := a.db.QueryRow(ctx, selectAlbumNameSQL, id).Scan(&result)
	return result, err
}
