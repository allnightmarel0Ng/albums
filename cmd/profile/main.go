package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/allnightmarel0Ng/albums/internal/app/profile/handler"
	"github.com/allnightmarel0Ng/albums/internal/app/profile/repository"
	"github.com/allnightmarel0Ng/albums/internal/app/profile/usecase"
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
		log.Fatalf("unable to connect to database: %s", err.Error())
	}
	defer db.Close()

	repo := repository.NewProfileRepository(
		domainRepository.NewUserRepository(db),
		domainRepository.NewAlbumRepository(db),
		domainRepository.NewArtistRepository(db),
	)
	usecase := usecase.NewProfileUseCase(repo)
	handler := handler.NewProfileHandler(usecase)

	router := gin.Default()
	router.GET("/users/:id", handler.HandleUserProfile)
	router.GET("/artists/:id", handler.HandleArtistProfile)

	log.Fatal(http.ListenAndServe(":"+conf.ProfilePort, router))
}
