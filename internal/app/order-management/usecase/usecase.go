package usecase

import (
	"context"
	"net/http"

	"github.com/allnightmarel0Ng/albums/internal/app/order-management/repository"
	"github.com/allnightmarel0Ng/albums/internal/domain/api"
	"github.com/allnightmarel0Ng/albums/internal/utils"
)

type OrderManagementUseCase interface {
	AddAlbumToUserOrder(request api.OrderActionRequest) api.Response
	DeleteAlbumFromUserOrder(request api.OrderActionRequest) api.Response
}

type orderManagementUseCase struct {
	repo repository.OrderManagementRepository
}

func (o *orderManagementUseCase) AddAlbumToUserOrder(request api.OrderActionRequest) api.Response {
	ctx, cancel := utils.DeadlineContext(2)
	defer cancel()

	err := o.repo.AddToOrder(ctx, request.UserID, request.AlbumID)
	if err != nil {
		switch err {
		case context.DeadlineExceeded:
			return &api.OrderActionResponse{
				Code:  http.StatusInternalServerError,
				Error: "database fail",
			}
		default:
			return &api.OrderActionResponse{
				Code:  http.StatusBadRequest,
				Error: "order creation error: album not found",
			}
		}
	}

	return &api.OrderActionResponse{
		Code: http.StatusOK,
	}
}
