package handler

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/allnightmarel0Ng/albums/internal/app/money-operations/usecase"
	"github.com/allnightmarel0Ng/albums/internal/domain/api"
	"github.com/allnightmarel0Ng/albums/internal/infrastructure/kafka"
)

type MoneyOperationsHandler interface {
	Handle()
}

type moneyOperationsHandler struct {
	useCase  usecase.MoneyOperationsUseCase
	consumer *kafka.Consumer
}

func NewMoneyOperationsHandler(useCase usecase.MoneyOperationsUseCase, consumer *kafka.Consumer) MoneyOperationsHandler {
	return &moneyOperationsHandler{
		useCase:  useCase,
		consumer: consumer,
	}
}

func (m *moneyOperationsHandler) Handle() {
	m.consumer.ConsumeMessagesEternally(m.forkMessages, log.Printf, log.Printf)
}

func (m *moneyOperationsHandler) handleDeposit(userID, diff int) {
	m.useCase.Deposit(userID, diff)
}

func (m *moneyOperationsHandler) handleBuy(userID, albumID int) {
	m.useCase.BuyOrder(userID, albumID)
}

func (m *moneyOperationsHandler) forkMessages(msg []byte) error {
	var operation api.MoneyOperationKafkaMessage
	if err := json.Unmarshal(msg, &operation); err != nil {
		return err
	}

	switch operation.Type {
	case "deposit":
		go m.handleDeposit(operation.UserID, operation.Diff)
	case "buy":
		go m.handleBuy(operation.UserID, operation.Diff)
	default:
		return errors.New("unknown message type")
	}

	return nil
}
