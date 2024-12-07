package main

import (
	"context"
	"fmt"
	"log"

	"github.com/allnightmarel0Ng/albums/internal/app/money-operations/handler"
	"github.com/allnightmarel0Ng/albums/internal/app/money-operations/repository"
	"github.com/allnightmarel0Ng/albums/internal/app/money-operations/usecase"
	"github.com/allnightmarel0Ng/albums/internal/config"
	domainRepository "github.com/allnightmarel0Ng/albums/internal/domain/repository"
	"github.com/allnightmarel0Ng/albums/internal/infrastructure/kafka"
	"github.com/allnightmarel0Ng/albums/internal/infrastructure/postgres"
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

	if err = c.SubscribeTopics([]string{"money-operations"}); err != nil {
		log.Fatalf("unable to subscribe to topic %s", err.Error())
	}

	p, err := kafka.NewProducer(fmt.Sprintf("kafka:%s", conf.KafkaPort), 1)
	if err != nil {
		log.Fatalf("unable to create a producer: %s", err.Error())
	}
	defer p.Close()

	db, err := postgres.NewDatabase(context.Background(), fmt.Sprintf("postgresql://%s:%s@postgres:%s/%s?sslmode=disable", conf.PostgresUser, conf.PostgresPassword, conf.PostgresPort, conf.PostgresDb))
	if err != nil {
		log.Fatalf("unable to establish db connection: %s", err.Error())
	}
	defer db.Close()

	repo := repository.NewMoneyOperationsRepository(domainRepository.NewUserRepository(db))
	useCase := usecase.NewMoneyOperationsUseCase(repo, p)
	handler := handler.NewMoneyOperationsHandler(useCase, c)
	handler.Handle()
}
