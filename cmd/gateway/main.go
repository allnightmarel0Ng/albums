package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/allnightmarel0Ng/albums/internal/app/gateway/handler"
	"github.com/allnightmarel0Ng/albums/internal/app/gateway/repository"
	"github.com/allnightmarel0Ng/albums/internal/app/gateway/usecase"
	"github.com/allnightmarel0Ng/albums/internal/config"
	domainRepository "github.com/allnightmarel0Ng/albums/internal/domain/repository"
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

	repo := repository.NewGatewayRepository(domainRepository.NewAlbumRepository(db))
	useCase := usecase.NewGatewayUseCase(repo, conf.AuthorizationPort, conf.ProfilePort, conf.JwtSecretKey)
	handler := handler.NewGatewayHandler(useCase)

	router := gin.Default()
	router.GET("/login", handler.HandleLogin)
	router.GET("/", handler.HandleMainPage)
	router.GET("/artists/:id", handler.HandleArtistProfile)
	router.GET("/profile", handler.HandleUserProfile)

	log.Fatal(http.ListenAndServe(":"+conf.GatewayPort, router))
}
