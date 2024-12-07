package usecase

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/allnightmarel0Ng/albums/internal/domain/api"
	"github.com/allnightmarel0Ng/albums/internal/infrastructure/kafka"
)

type NotificationsUseCase interface {
	AddUser(userID int, channel chan<- *api.NotificationKafkaMessage)
	DeleteUser(userID int)
	Consume()
}

type notificationsUseCase struct {
	consumer *kafka.Consumer
	channels map[int]chan<- *api.NotificationKafkaMessage
	mu       sync.Mutex
}

func NewNotificationsUseCase(consumer *kafka.Consumer) NotificationsUseCase {
	return &notificationsUseCase{
		consumer: consumer,
		channels: make(map[int]chan<- *api.NotificationKafkaMessage),
	}
}

func (n *notificationsUseCase) AddUser(userID int, channel chan<- *api.NotificationKafkaMessage) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.channels[userID] = channel
}

func (n *notificationsUseCase) DeleteUser(userID int) {
	n.mu.Lock()
	defer n.mu.Unlock()
	delete(n.channels, userID)
}

func (n *notificationsUseCase) Consume() {
	n.consumer.ConsumeMessagesEternally(n.onConsume, log.Printf, log.Printf)
}

func (n *notificationsUseCase) onConsume(msg []byte) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	var notification api.NotificationKafkaMessage
	err := json.Unmarshal(msg, &notification)
	if err != nil {
		return err
	}

	_, ok := n.channels[notification.UserID]
	if !ok {
		return fmt.Errorf("no such user with ID: %d", notification.UserID)
	}

	n.channels[notification.UserID] <- &notification

	return nil
}
