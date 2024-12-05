package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/allnightmarel0Ng/albums/internal/app/gateway/handler"
	"github.com/allnightmarel0Ng/albums/internal/app/gateway/usecase"
	"github.com/allnightmarel0Ng/albums/internal/config"
	"github.com/allnightmarel0Ng/albums/internal/infrastructure/kafka"
	"github.com/gin-gonic/gin"
)

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("unable to load config: %s", err.Error())
	}

	p, err := kafka.NewProducer(fmt.Sprintf("kafka:%s", conf.KafkaPort))
	if err != nil {
		log.Fatalf("unable to create a producer: %s", err.Error())
	}
	defer p.Close()

	useCase := usecase.NewGatewayUseCase(p, conf.AuthorizationPort, conf.ProfilePort, conf.OrderManagementPort, conf.SearchEnginePort, conf.AdminPanelPort, conf.JwtSecretKey, conf.PostgresUser, conf.PostgresPassword, conf.PostgresPort, conf.PostgresDb)
	handler := handler.NewGatewayHandler(useCase)

	router := gin.Default()

	router.GET("/login", handler.HandleLogin)
	router.POST("/logout", handler.HandleLogout)
	router.POST("/registration", handler.HandleRegistration)

	router.POST("/", handler.HandleMainPage)
	router.POST("/search", handler.HandleSearch)

	router.GET("/profile", handler.HandleUserProfile)
	router.GET("/artists/:id", handler.HandleArtistProfile)
	router.GET("/albums/:id", handler.HandleAlbumProfile)

	router.POST("/add/:id", handler.HandleOrderAdd)
	router.POST("/remove/:id", handler.HandleOrderRemove)
	router.GET("/orders", handler.HandleOrders)

	router.POST("/deposit", handler.HandleDeposit)
	router.POST("/buy", handler.HandleBuy)

	router.GET("/admin-panel/logs/:pageNumber", handler.HandleLogs)
	router.DELETE("/admin-panel/delete/:id", handler.HandleDelete)
	router.GET("/admin-panel/save-dump", handler.HandleSaveDump)
	router.POST("/admin-panel/load-dump", handler.HandleLoadDump)

	log.Fatal(http.ListenAndServe(":"+conf.GatewayPort, router))
}
