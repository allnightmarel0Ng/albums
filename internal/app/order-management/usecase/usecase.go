package usecase

import (
	"context"
	"log"
	"net/http"

	"github.com/allnightmarel0Ng/albums/internal/app/order-management/repository"
	"github.com/allnightmarel0Ng/albums/internal/domain/api"
	"github.com/allnightmarel0Ng/albums/internal/utils"
)

type OrderManagementUseCase interface {
	AddAlbumToUserOrder(request api.OrderActionRequest) api.Response
	RemoveAlbumFromUserOrder(request api.OrderActionRequest) api.Response
	UserOrder(userID int) api.Response
}

type orderManagementUseCase struct {
	repo repository.OrderManagementRepository
}

func NewOrderManagementUseCase(repo repository.OrderManagementRepository) OrderManagementUseCase {
	return &orderManagementUseCase{
		repo: repo,
	}
}

func (o *orderManagementUseCase) AddAlbumToUserOrder(request api.OrderActionRequest) api.Response {
	ctx, cancel := utils.DeadlineContext(5)
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

func (o *orderManagementUseCase) RemoveAlbumFromUserOrder(request api.OrderActionRequest) api.Response {
	ctx, cancel := utils.DeadlineContext(5)
	defer cancel()

	err := o.repo.RemoveFromOrder(ctx, request.UserID, request.AlbumID)
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
				Error: "order deletion error: album not found",
			}
		}
	}

	return &api.OrderActionResponse{
		Code: http.StatusOK,
	}
}

func (o *orderManagementUseCase) UserOrder(userID int) api.Response {
	ctx, cancel := utils.DeadlineContext(5)
	defer cancel()

	result, err := o.repo.UserOrder(ctx, userID)
	if err != nil {
		log.Print(err.Error())
		return &api.UserOrdersResponse{
			Code:  http.StatusInternalServerError,
			Error: "error retrieving orders from database",
		}
	}

	return &api.UserOrdersResponse{
		Code:   http.StatusOK,
		Orders: result,
	}
}
