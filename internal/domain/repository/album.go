package repository

import (
	"context"
	"strings"

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
				WHERE pu.user_id = $1;`

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
				WHERE a.name LIKE $1;`

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
				WHERE ar.id = $1;`
)

type AlbumRepository interface {
	GetUsersPurchasedAlbums(ctx context.Context, userID int) ([]model.Album, error)
	GetAllAlbumsLike(ctx context.Context, name string) ([]model.Album, error)
	GetArtistsAlbums(ctx context.Context, artistID int) ([]model.Album, error)
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

	for rows.Next() {
		var (
			album model.Album
			track model.Track
		)

		err := rows.Scan(&album.ID, &album.Name, &album.Author.ID,
			&album.Author.Name, &album.Author.Genre, &album.Author.ImageURL, &album.ImageURL, &album.Price,
			&track.ID, &track.Name, &track.Number)
		if err != nil {
			return nil, err
		}

		_, ok := albumsMap[album.ID]
		if !ok {
			album.Tracks = []model.Track{track}
			albumsMap[album.ID] = &album
		} else {
			albumsMap[album.ID].Tracks = append(albumsMap[album.ID].Tracks, track)
		}
	}

	result := make([]model.Album, len(albumsMap))
	i := 0
	for _, v := range albumsMap {
		result[i] = *v
		i++
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

func (a *albumRepository) GetAllAlbumsLike(ctx context.Context, name string) ([]model.Album, error) {
	tokens := strings.Split(name, " ")
	var sb strings.Builder

	sb.WriteRune('%')
	for _, token := range tokens {
		sb.WriteString(token)
	}

	if len(tokens) > 0 {
		sb.WriteRune('%')
	}

	rows, err := a.db.Query(ctx, selectAlbumsLikeNameSQL, sb.String())
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
