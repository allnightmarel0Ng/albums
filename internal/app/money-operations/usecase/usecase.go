package usecase

import (
	"context"

	"github.com/allnightmarel0Ng/albums/internal/app/money-operations/repository"
)

type MoneyOperationsUseCase interface {
	Deposit(id, diff int)
	BuyOrder(userID, albumID int)
}

type moneyOperationsUseCase struct {
	repo repository.MoneyOperationsRepository
}

func NewMoneyOperationsUseCase(repo repository.MoneyOperationsRepository) MoneyOperationsUseCase {
	return &moneyOperationsUseCase{
		repo: repo,
	}
}

func (m *moneyOperationsUseCase) Deposit(id, diff int) {
	m.repo.Deposit(context.Background(), id, diff)
}

func (m *moneyOperationsUseCase) BuyOrder(userID, albumID int) {
	m.repo.BuyOrder(context.Background(), userID, albumID)
}
