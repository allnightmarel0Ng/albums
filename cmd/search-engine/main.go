package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/allnightmarel0Ng/albums/internal/app/search-engine/handler"
	"github.com/allnightmarel0Ng/albums/internal/app/search-engine/repository"
	"github.com/allnightmarel0Ng/albums/internal/app/search-engine/usecase"
	"github.com/allnightmarel0Ng/albums/internal/config"
	domainRepository "github.com/allnightmarel0Ng/albums/internal/domain/repository"
	"github.com/allnightmarel0Ng/albums/internal/infrastructure/postgres"
	"github.com/gin-gonic/gin"
)

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("unable to load config %s:", err.Error())
	}

	db, err := postgres.NewDatabase(context.Background(), fmt.Sprintf("postgresql://%s:%s@postgres:%s/%s?sslmode=disable", conf.PostgresUser, conf.PostgresPassword, conf.PostgresPort, conf.PostgresDb))
	if err != nil {
		log.Fatalf("unable to establish db connection: %s", err.Error())
	}
	defer db.Close()

	repo := repository.NewSearchEngineRepository(domainRepository.NewArtistRepository(db), domainRepository.NewAlbumRepository(db))
	usecase := usecase.NewSearchEngineUseCase(repo)
	handler := handler.NewSearchEngineHandler(usecase)

	router := gin.Default()
	router.POST("/random", handler.HandleRandom)
	router.POST("/search", handler.HandleSearch)

	log.Fatal(http.ListenAndServe(":"+conf.SearchEnginePort, router))
}
