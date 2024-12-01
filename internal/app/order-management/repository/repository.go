package repository

import (
	"context"

	"github.com/allnightmarel0Ng/albums/internal/domain/repository"
)

type OrderManagementRepository interface {
	AddToOrder(ctx context.Context, userID, albumID int) error
	RemoveFromOrder(ctx context.Context, userID, albumID int) error
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
