package repository

import (
	"context"

	"github.com/allnightmarel0Ng/albums/internal/domain/model"
	"github.com/allnightmarel0Ng/albums/internal/domain/repository"
)

type AdminPanelRepository interface {
	GetBuyLogsAndCount(ctx context.Context, offset, limit uint) (uint, []model.BuyLog, error)
	DeleteAlbum(ctx context.Context, albumID int) error
}

type adminPanelRepository struct {
	albums repository.AlbumRepository
	logs   repository.LogsRepository
}

func NewAdminPanelRepository(albums repository.AlbumRepository, logs repository.LogsRepository) AdminPanelRepository {
	return &adminPanelRepository{
		albums: albums,
		logs:   logs,
	}
}

func (a *adminPanelRepository) GetBuyLogsAndCount(ctx context.Context, offset, limit uint) (uint, []model.BuyLog, error) {
	select {
	case <-ctx.Done():
		return 0, nil, ctx.Err()
	default:
		logs, err := a.logs.GetLogs(ctx, offset, limit)
		if err != nil {
			return 0, nil, err
		}

		count, err := a.logs.GetLogsCount(ctx)
		return count, logs, err
	}
}

func (a *adminPanelRepository) DeleteAlbum(ctx context.Context, albumID int) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return a.albums.DeleteAlbum(ctx, albumID)
	}
}
