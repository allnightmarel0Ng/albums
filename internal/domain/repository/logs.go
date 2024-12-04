package repository

import (
	"context"

	"github.com/allnightmarel0Ng/albums/internal/domain/model"
	"github.com/allnightmarel0Ng/albums/internal/infrastructure/postgres"
)

const (
	selectLogsSQL =
	/* sql */ `SELECT b.id,
					u.id,
					u.email,
					u.is_admin,
					u.nickname,
					u.balance,
					u.image_url,
					a.id,
					a.name,
					ar.name,
					ar.genre,
					ar.image_url,
					a.image_url,
					a.price,
					b.logging_time
				FROM public.buy_logs AS b
				JOIN public.users AS u ON u.id = b.user_id
				JOIN public.albums AS a ON a.id = b.album_id
				LIMIT $2
				OFFSET $1;`

	selectLogsCountSQL =
	/* sql */ `SELECT COUNT(*) FROM public.buy_logs;`
)

type LogsRepository interface {
	GetLogs(ctx context.Context, offset, limit uint) ([]model.BuyLog, error)
	GetLogsCount(ctx context.Context) (uint, error)
}

type logsRepository struct {
	db postgres.Database
}

func NewLogsRepository(db postgres.Database) LogsRepository {
	return &logsRepository{
		db: db,
	}
}

func (l *logsRepository) GetLogs(ctx context.Context, offset, limit uint) ([]model.BuyLog, error) {
	rows, err := l.db.Query(ctx, selectLogsSQL, offset, limit)
	if err != nil {
		return nil, err
	}

	var logs []model.BuyLog

	for rows.Next() {
		var log model.BuyLog
		var user model.User
		var album model.Album
		var artistName, artistGenre, artistImageURL string

		err := rows.Scan(
			&log.ID,
			&user.ID,
			&user.Email,
			&user.IsAdmin,
			&user.Nickname,
			&user.Balance,
			&user.ImageURL,
			&album.ID,
			&album.Name,
			&artistName,
			&artistGenre,
			&artistImageURL,
			&album.ImageURL,
			&album.Price,
			&log.LoggingTime,
		)
		if err != nil {
			return nil, err
		}

		album.Author = model.Artist{
			Name:     artistName,
			Genre:    artistGenre,
			ImageURL: artistImageURL,
		}

		log.Buyer = user
		log.Album = album

		logs = append(logs, log)
	}

	return logs, nil
}

func (l *logsRepository) GetLogsCount(ctx context.Context) (uint, error) {
	var result uint
	err := l.db.QueryRow(ctx, selectLogsCountSQL).Scan(&result)
	return result, err
}
