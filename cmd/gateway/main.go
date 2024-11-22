package main

import (
	"log"
	"net/http"

	"github.com/allnightmarel0Ng/albums/internal/app/gateway/handler"
	"github.com/allnightmarel0Ng/albums/internal/app/gateway/usecase"
	"github.com/allnightmarel0Ng/albums/internal/config"
	"github.com/gin-gonic/gin"
)

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("unable to load config: %s", err.Error())
	}

	useCase := usecase.NewGatewayUseCase(conf.AuthorizationPort)
	handler := handler.NewGatewayHandler(useCase)

	router := gin.Default()
	router.GET("/login", handler.HandleLogin)

	log.Fatal(http.ListenAndServe(":"+conf.GatewayPort, router))
}
