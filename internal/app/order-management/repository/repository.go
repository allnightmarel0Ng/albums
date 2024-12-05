package repository

import (
	"context"
	"errors"

	"github.com/allnightmarel0Ng/albums/internal/domain/model"
	"github.com/allnightmarel0Ng/albums/internal/domain/repository"
)

var (
	ErrDatabaseCommunication = errors.New("db communication error")
	ErrAlreadyInOrder        = errors.New("album is already in order")
	ErrNotInOrder            = errors.New("album not in order")
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
		orders, err := o.UserOrder(ctx, userID, true)
		if err != nil {
			return ErrDatabaseCommunication
		}

		if len(orders) == 1 {
			found := false
			for i := 0; i < len(orders[0].Albums); i++ {
				if albumID == orders[0].Albums[0].ID {
					found = true
				}
			}

			if !found {
				return ErrNotInOrder
			}
		}

		return o.orders.AddAlbumToUserOrder(ctx, userID, albumID)
	}
}

func (o *orderManagementRepository) RemoveFromOrder(ctx context.Context, userID, albumID int) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		orders, err := o.UserOrder(ctx, userID, true)
		if err != nil {
			return ErrDatabaseCommunication
		}

		if len(orders) == 1 {
			for i := 0; i < len(orders[0].Albums); i++ {
				if albumID == orders[0].Albums[0].ID {
					return ErrAlreadyInOrder
				}
			}
		}

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
