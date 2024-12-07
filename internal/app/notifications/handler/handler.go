package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/allnightmarel0Ng/albums/internal/app/notifications/usecase"
	"github.com/allnightmarel0Ng/albums/internal/domain/api"
	"github.com/allnightmarel0Ng/albums/internal/utils"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type NotificationsHandler interface {
	HandleNotifications(w http.ResponseWriter, r *http.Request)
}

type notificationsHandler struct {
	useCase           usecase.NotificationsUseCase
	authorizationPort string
}

func NewNotificationsHandler(useCase usecase.NotificationsUseCase, authorizationPort string) NotificationsHandler {
	return &notificationsHandler{
		useCase:           useCase,
		authorizationPort: authorizationPort,
	}
}

func (n *notificationsHandler) HandleNotifications(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgrade error: %s", err.Error())
		return
	}
	defer conn.Close()

	msgType, msg, err := conn.ReadMessage()
	if err != nil {
		log.Printf("unable to read the message from client: %s", err.Error())
		return
	}

	var subscription api.NotificationSubscribeRequest
	err = json.Unmarshal(msg, &subscription)
	if err != nil {
		raw, _ := json.Marshal(api.ErrorResponse{
			Code:  http.StatusBadRequest,
			Error: "invalid subscription data",
		})
		conn.WriteMessage(msgType, raw)
		log.Printf("invalid subscription data")
		return
	}

	response := utils.Authorize("Bearer "+subscription.Jwt, n.authorizationPort)
	if response.GetCode() != http.StatusOK {
		raw, _ := json.Marshal(response)
		conn.WriteMessage(msgType, raw)
		log.Printf("unable to authorize the user")
		return
	}

	claims := response.(*api.AuthorizationResponse)
	end := make(chan bool)
	go func() {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Print("client disconnected")
			end <- true
		}
	}()

	notificationChannel := make(chan *api.NotificationKafkaMessage)
	n.useCase.AddUser(claims.ID, notificationChannel)

	run := true
	for run {
		select {
		case <-end:
			run = false
			n.useCase.DeleteUser(claims.ID)
		case notification := <-notificationChannel:
			var response api.NotificationResponse
			response.Success = *notification.Success
			response.Message = getMessage(notification)

			raw, _ := json.Marshal(response)
			err = conn.WriteMessage(msgType, raw)
			if err != nil {
				log.Println("unable to write message")
			}
			log.Print(raw)
		}
	}
}

func getMessage(notification *api.NotificationKafkaMessage) string {
	switch notification.Type {
	case api.Deposit:
		if *notification.Success {
			return "Money has been added to your account successfully"
		}
		return "Money has not been added to your account"
	case api.Buy:
		if *notification.Success {
			return fmt.Sprintf("Order %d has been paid successfully", notification.OrderID)
		}
		return fmt.Sprintf("Order %d has not been paid", notification.OrderID)
	default:
		return fmt.Sprintf("Album %s, that you owned, has been deleted", notification.AlbumName)
	}
}
