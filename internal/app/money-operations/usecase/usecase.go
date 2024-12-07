package usecase

import (
	"context"
	"log"

	"github.com/allnightmarel0Ng/albums/internal/app/money-operations/repository"
	"github.com/allnightmarel0Ng/albums/internal/domain/api"
	"github.com/allnightmarel0Ng/albums/internal/infrastructure/kafka"
	"github.com/allnightmarel0Ng/albums/internal/utils"
)

type MoneyOperationsUseCase interface {
	Deposit(id int, diff uint)
	BuyOrder(userID, albumID int)
}

type moneyOperationsUseCase struct {
	repo     repository.MoneyOperationsRepository
	producer *kafka.Producer
}

func NewMoneyOperationsUseCase(repo repository.MoneyOperationsRepository, producer *kafka.Producer) MoneyOperationsUseCase {
	return &moneyOperationsUseCase{
		repo:     repo,
		producer: producer,
	}
}

func (m *moneyOperationsUseCase) Deposit(id int, diff uint) {
	err := m.repo.Deposit(context.Background(), id, diff)
	if err != nil {
		log.Printf("unable to deposit money: %s", err.Error())
	}

	success := (err == nil)
	err = utils.ProduceNotificationMessage(api.NotificationKafkaMessage{
		Type:    api.Deposit,
		UserID:  id,
		Success: &success,
	}, m.producer)
	if err != nil {
		log.Printf("unable to produce notification message: %s", err.Error())
	}
}

func (m *moneyOperationsUseCase) BuyOrder(userID, orderID int) {
	err := m.repo.BuyOrder(context.Background(), userID, orderID)
	if err != nil {
		log.Printf("unable to deposit money: %s", err.Error())
	}

	success := (err == nil)
	err = utils.ProduceNotificationMessage(api.NotificationKafkaMessage{
		Type:    api.Buy,
		UserID:  userID,
		OrderID: orderID,
		Success: &success,
	}, m.producer)
	if err != nil {
		log.Printf("unable to produce notification message: %s", err.Error())
	}

}
