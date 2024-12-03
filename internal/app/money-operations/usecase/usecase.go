package usecase

import (
	"context"
	"log"

	"github.com/allnightmarel0Ng/albums/internal/app/money-operations/repository"
)

type MoneyOperationsUseCase interface {
	Deposit(id int, diff uint)
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

func (m *moneyOperationsUseCase) Deposit(id int, diff uint) {
	log.Print(m.repo.Deposit(context.Background(), id, diff))
}

func (m *moneyOperationsUseCase) BuyOrder(userID, orderID int) {
	m.repo.BuyOrder(context.Background(), userID, orderID)
}
