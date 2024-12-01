package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/allnightmarel0Ng/albums/internal/app/order-management/handler"
	"github.com/allnightmarel0Ng/albums/internal/app/order-management/repository"
	"github.com/allnightmarel0Ng/albums/internal/app/order-management/usecase"
	"github.com/allnightmarel0Ng/albums/internal/config"
	domainRepository "github.com/allnightmarel0Ng/albums/internal/domain/repository"
	"github.com/allnightmarel0Ng/albums/internal/infrastructure/postgres"
	"github.com/gin-gonic/gin"
)

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err.Error())
	}

	db, err := postgres.NewDatabase(context.Background(), fmt.Sprintf("postgresql://%s:%s@postgres:%s/%s?sslmode=disable", conf.PostgresUser, conf.PostgresPassword, conf.PostgresPort, conf.PostgresDb))
	if err != nil {
		log.Fatalf("unable to establish db connection: %s", err.Error())
	}
	defer db.Close()

	repo := repository.NewOrderManagementRepository(domainRepository.NewOrderRepository(db))
	useCase := usecase.NewOrderManagementUseCase(repo)
	handler := handler.NewOrderManagementHandler(useCase)

	router := gin.Default()
	router.POST("/add", handler.HandleAdd)
	router.POST("/remove", handler.HandleRemove)
	router.GET("/orders/:id", handler.HandleOrders)

	log.Fatal(http.ListenAndServe(":"+conf.OrderManagementPort, router))
}
