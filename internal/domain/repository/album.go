package repository

import (
	"strings"

	"github.com/allnightmarel0Ng/albums/internal/domain/model"
	"github.com/allnightmarel0Ng/albums/internal/infrastructure/postgres"
)

const (
	geAlbumsLikeNameSQL = /* sql */ `SELECT 
							a.name,
							a.profile_picture_url,
							a.bio,
							al.id,
							al.name,
							al.release_date,
							al.cover_art_url,
							al.price,
							al.genre,
							t.name,
							t.duration,
							t.audio_file_url
						FROM public.tracks AS t
						LEFT JOIN public.albums AS al ON t.album_id = al.id
						JOIN public.artists AS a ON al.artist_id = a.id
						WHERE al.name LIKE $1
						ORDER BY al.id;`
)

type AlbumRepository interface {
	GetAllAlbumsLike(name string) ([]model.Album, error)
}

type albumRepository struct {
	db postgres.Database
}

func NewAlbumRepository(db postgres.Database) AlbumRepository {
	return &albumRepository{
		db: db,
	}
}

func (a *albumRepository) GetAllAlbumsLike(name string) ([]model.Album, error) {
	tokens := strings.Split(name, " ")

	var sb strings.Builder
	sb.WriteByte('%')

	for _, token := range tokens {
		_, err := sb.Write([]byte(token))
		if err != nil {
			return nil, err
		}
	}

	if len(tokens) != 0 {
		sb.WriteByte('%')
	}

	rows, err := a.db.Query(geAlbumsLikeNameSQL, sb.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.Album
	prevAlbumID := -1
	var currentAlbumTracks []model.Track

	for rows.Next() {
		var (
			track   model.Track
			album   model.Album
			albumID int
		)

		err = rows.Scan(&album.Artist.Name, &album.Artist.ProfilePictureUrl, &album.Artist.Bio, &albumID, &album.Name, &album.ReleaseDate, &album.CoverArtUrl, &album.Price, &album.Genre, &track.Name, &track.Duration, &track.AudioFileUrl)
		if err != nil {
			return nil, err
		}

		currentAlbumTracks = append(currentAlbumTracks, track)

		if albumID != prevAlbumID && prevAlbumID != -1 {
			album.Tracks = make([]model.Track, len(currentAlbumTracks))
			copy(album.Tracks, currentAlbumTracks)
			result = append(result, album)
			currentAlbumTracks = []model.Track{}
		}

		prevAlbumID = albumID
	}

	return result, nil
}
