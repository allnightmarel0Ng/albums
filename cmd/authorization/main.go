package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/allnightmarel0Ng/albums/internal/app/authorization/handler"
	"github.com/allnightmarel0Ng/albums/internal/app/authorization/repository"
	"github.com/allnightmarel0Ng/albums/internal/app/authorization/usecase"
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
	defer db.Close()

	repo := repository.NewAuthorizationRepository(domainRepository.NewUserRepository(db))
	useCase := usecase.NewAuthorizationUseCase(repo, conf.JwtSecretKey)
	handler := handler.NewAuthorizationHandler(useCase)

	router := gin.Default()
	router.GET("/", handler.HandleAuthorization)
	log.Fatal(http.ListenAndServe(":"+conf.AuthorizationPort, router))
}
