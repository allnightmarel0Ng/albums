package repository

import (
	"context"

	"github.com/allnightmarel0Ng/albums/internal/domain/repository"
)

type MoneyOperationsRepository interface {
	Deposit(ctx context.Context, id int, diff uint) error
	BuyOrder(ctx context.Context, userID, albumID int) error
}

type moneyOperationsRepository struct {
	users repository.UserRepository
}

func NewMoneyOperationsRepository(users repository.UserRepository) MoneyOperationsRepository {
	return &moneyOperationsRepository{
		users: users,
	}
}

func (m *moneyOperationsRepository) Deposit(ctx context.Context, id int, diff uint) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return m.users.ChangeBalance(ctx, id, diff)
	}
}

func (m *moneyOperationsRepository) BuyOrder(ctx context.Context, userID, orderID int) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return m.users.PayForOrder(ctx, userID, orderID)
	}
}
