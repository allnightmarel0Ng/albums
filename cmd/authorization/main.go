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
	"github.com/allnightmarel0Ng/albums/internal/infrastructure/redis"
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

	client := redis.NewClient(fmt.Sprintf("redis:%s", conf.RedisPort), "", 0)
	defer func() {
		err = client.Close()
		if err != nil {
			log.Fatalf("unable to close redis connection: %s", err.Error())
		}
	}()
	if client.Ping(context.Background()) != nil {
		log.Fatalf("unable to connect to redis: %s", err.Error())
	}

	repo := repository.NewAuthorizationRepository(domainRepository.NewUserRepository(db), client)
	useCase := usecase.NewAuthorizationUseCase(repo, []byte(conf.JwtSecretKey))
	handler := handler.NewAuthorizationHandler(useCase)

	router := gin.Default()
	router.GET("/authorize", handler.HandleAuthorization)
	router.GET("/authenticate", handler.HandleAuthentication)
	router.POST("/logout", handler.HandleLogout)
	router.POST("/registration", handler.HandleRegistration)

	log.Fatal(http.ListenAndServe(":"+conf.AuthorizationPort, router))
}
