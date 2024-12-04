package usecase

import (
	"net/http"

	"github.com/allnightmarel0Ng/albums/internal/app/admin-panel/repository"
	"github.com/allnightmarel0Ng/albums/internal/domain/api"
	"github.com/allnightmarel0Ng/albums/internal/utils"
)

type AdminPanelUseCase interface {
	Logs(pageNumber uint, pageSize uint) api.Response
	DeleteAlbum(albumID int) api.Response
}

type adminPanelUseCase struct {
	repo repository.AdminPanelRepository
}

func NewAdminPanelUseCase(repo repository.AdminPanelRepository) AdminPanelUseCase {
	return &adminPanelUseCase{
		repo: repo,
	}
}

func (a *adminPanelUseCase) Logs(pageNumber uint, pageSize uint) api.Response {
	offset := (pageNumber - 1) * pageSize
	limit := pageSize

	ctx, cancel := utils.DeadlineContext(2)
	defer cancel()

	count, logs, err := a.repo.GetBuyLogsAndCount(ctx, offset, limit)
	if err != nil {
		return &api.ErrorResponse{
			Code:  http.StatusInternalServerError,
			Error: "db error",
		}
	}

	return &api.BuyLogsResponse{
		Code:      http.StatusOK,
		Logs:      logs,
		LogsCount: count,
	}
}

func (a *adminPanelUseCase) DeleteAlbum(albumID int) api.Response {
	ctx, cancel := utils.DeadlineContext(2)
	defer cancel()

	err := a.repo.DeleteAlbum(ctx, albumID)
	if err != nil {
		return &api.ErrorResponse{
			Code:  http.StatusInternalServerError,
			Error: "db error",
		}
	}

	return nil
}
