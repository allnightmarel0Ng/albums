package repository

import (
	"context"

	"github.com/allnightmarel0Ng/albums/internal/domain/model"
	"github.com/allnightmarel0Ng/albums/internal/domain/repository"
)

type OrderManagementRepository interface {
	AddToOrder(ctx context.Context, userID, albumID int) error
	RemoveFromOrder(ctx context.Context, userID, albumID int) error
	UserOrder(ctx context.Context, userID int, unpaidOnly bool) ([]model.Order, error)
}

type orderManagementRepository struct {
	orders repository.OrderRepository
}

func NewOrderManagementRepository(orders repository.OrderRepository) OrderManagementRepository {
	return &orderManagementRepository{
		orders: orders,
	}
}

func (o *orderManagementRepository) AddToOrder(ctx context.Context, userID, albumID int) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return o.orders.AddAlbumToUserOrder(ctx, userID, albumID)
	}
}

func (o *orderManagementRepository) RemoveFromOrder(ctx context.Context, userID, albumID int) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return o.orders.DeleteAlbumFromUserOrder(ctx, userID, albumID)
	}
}

func (o *orderManagementRepository) UserOrder(ctx context.Context, userID int, unpaidOnly bool) ([]model.Order, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return o.orders.GetUserOrders(ctx, userID, unpaidOnly)
	}
}
