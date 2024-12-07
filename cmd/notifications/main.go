package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/allnightmarel0Ng/albums/internal/app/notifications/handler"
	"github.com/allnightmarel0Ng/albums/internal/app/notifications/usecase"
	"github.com/allnightmarel0Ng/albums/internal/config"
	"github.com/allnightmarel0Ng/albums/internal/infrastructure/kafka"
)

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("unable to load config: %s", err.Error())
	}

	c, err := kafka.NewConsumer(fmt.Sprintf("kafka:%s", conf.KafkaPort), "consumers")
	if err != nil {
		log.Fatalf("unable to create consumer: %s", err.Error())
	}
	defer c.Close()

	if err = c.SubscribeTopics([]string{"notifications"}); err != nil {
		log.Fatalf("unable to subscribe to topic %s", err.Error())
	}

	useCase := usecase.NewNotificationsUseCase(c)
	handler := handler.NewNotificationsHandler(useCase, conf.AuthorizationPort)

	http.HandleFunc("/ws", handler.HandleNotifications)

	go useCase.Consume()

	log.Fatal(http.ListenAndServe(":"+conf.NotificationsPort, nil))
}
