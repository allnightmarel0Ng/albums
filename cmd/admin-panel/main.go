package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/allnightmarel0Ng/albums/internal/app/admin-panel/handler"
	"github.com/allnightmarel0Ng/albums/internal/app/admin-panel/repository"
	"github.com/allnightmarel0Ng/albums/internal/app/admin-panel/usecase"
	"github.com/allnightmarel0Ng/albums/internal/config"
	domainRepository "github.com/allnightmarel0Ng/albums/internal/domain/repository"
	"github.com/allnightmarel0Ng/albums/internal/infrastructure/kafka"
	"github.com/allnightmarel0Ng/albums/internal/infrastructure/postgres"
	"github.com/gin-gonic/gin"
)

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("unable to load config: %s", err.Error())
	}

	db, err := postgres.NewDatabase(context.Background(), fmt.Sprintf("postgresql://%s:%s@postgres:%s/%s?sslmode=disable", conf.PostgresUser, conf.PostgresPassword, conf.PostgresPort, conf.PostgresDb))
	if err != nil {
		log.Fatalf("unable to establish db connection: %s", err.Error())
	}
	defer db.Close()

	p, err := kafka.NewProducer(fmt.Sprintf("kafka:%s", conf.KafkaPort), 3)
	if err != nil {
		log.Fatalf("unable to create a producer: %s", err.Error())
	}
	defer p.Close()

	repo := repository.NewAdminPanelRepository(domainRepository.NewAlbumRepository(db), domainRepository.NewLogsRepository(db))
	useCase := usecase.NewAdminPanelUseCase(repo, conf.ProfilePort, p)
	handler := handler.NewAdminPanelHandler(useCase)

	router := gin.Default()
	router.GET("/logs/:pageNumber", handler.HandleBuyLogs)
	router.DELETE("/delete/:id", handler.HandleDeleteAlbum)

	log.Fatal(http.ListenAndServe(":"+conf.AdminPanelPort, router))
}
