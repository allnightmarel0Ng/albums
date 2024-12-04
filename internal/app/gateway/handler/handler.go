package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/allnightmarel0Ng/albums/internal/app/gateway/usecase"
	"github.com/allnightmarel0Ng/albums/internal/domain/api"
	"github.com/allnightmarel0Ng/albums/internal/utils"
	"github.com/gin-gonic/gin"
)

type GatewayHandler interface {
	HandleLogin(c *gin.Context)
	HandleMainPage(c *gin.Context)
	HandleSearch(c *gin.Context)
	HandleUserProfile(c *gin.Context)
	HandleArtistProfile(c *gin.Context)
	HandleOrderAdd(c *gin.Context)
	HandleOrderRemove(c *gin.Context)
	HandleOrders(c *gin.Context)
	HandleDeposit(c *gin.Context)
	HandleBuy(c *gin.Context)
	HandleLogs(c *gin.Context)
	HandleDelete(c *gin.Context)
	HandleSaveDump(c *gin.Context)
	HandleLoadDump(c *gin.Context)
}

type gatewayHandler struct {
	useCase usecase.GatewayUseCase
}

func NewGatewayHandler(useCase usecase.GatewayUseCase) GatewayHandler {
	return &gatewayHandler{
		useCase: useCase,
	}
}

func (g *gatewayHandler) HandleLogin(c *gin.Context) {
	code, raw := g.useCase.Authentication(c.GetHeader("Authorization"))
	utils.SendRaw(c, code, raw)
}

func (g *gatewayHandler) HandleMainPage(c *gin.Context) {
	code, raw := g.useCase.MainPage(c.Request.Body)
	utils.SendRaw(c, code, raw)
}

func (g *gatewayHandler) HandleSearch(c *gin.Context) {
	code, raw := g.useCase.Search(c.Request.Body)
	utils.SendRaw(c, code, raw)
}

func (g *gatewayHandler) HandleUserProfile(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		utils.Send(c, &api.AuthorizationResponse{
			Code:  http.StatusBadRequest,
			Error: "wrong auth token",
		})
		return
	}

	code, raw := g.useCase.UserProfile(authHeader)
	utils.SendRaw(c, code, raw)
}

func (g *gatewayHandler) HandleArtistProfile(c *gin.Context) {
	id, err := utils.GetParam[int](c, "id")
	if err != nil {
		utils.Send(c, &api.ErrorResponse{
			Code:  http.StatusBadRequest,
			Error: "invalid id parameter",
		})
		return
	}

	code, raw := g.useCase.ArtistProfile(id)
	utils.SendRaw(c, code, raw)
}

func (g *gatewayHandler) HandleOrderAdd(c *gin.Context) {
	handleOrderAction(c, g.useCase.AddToOrder)
}

func (g *gatewayHandler) HandleOrderRemove(c *gin.Context) {
	handleOrderAction(c, g.useCase.RemoveFromOrder)
}

func (g *gatewayHandler) HandleOrders(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		utils.Send(c, &api.AuthorizationResponse{
			Code:  http.StatusBadRequest,
			Error: "wrong auth token",
		})
		return
	}

	code, raw := g.useCase.UserOrders(authHeader)
	utils.SendRaw(c, code, raw)
}

func (g *gatewayHandler) HandleDeposit(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		utils.Send(c, &api.AuthorizationResponse{
			Code:  http.StatusBadRequest,
			Error: "wrong auth token",
		})
		return
	}

	var request api.DepositRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.Send(c, &api.ErrorResponse{
			Code:  http.StatusBadRequest,
			Error: "invalid body in request",
		})
		return
	}

	response := g.useCase.Deposit(authHeader, request.Money)
	if response != nil {
		utils.Send(c, response)
		return
	}

	c.String(http.StatusOK, "")
}

func (g *gatewayHandler) HandleBuy(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		utils.Send(c, &api.AuthorizationResponse{
			Code:  http.StatusBadRequest,
			Error: "wrong auth token",
		})
		return
	}

	response := g.useCase.Buy(authHeader)
	if response != nil {
		utils.Send(c, response)
		return
	}

	c.String(http.StatusOK, "")
}

func (g *gatewayHandler) HandleLogs(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		utils.Send(c, &api.AuthorizationResponse{
			Code:  http.StatusBadRequest,
			Error: "wrong auth token",
		})
		return
	}

	params := c.Param("id")
	query := c.Request.URL.RawQuery
	if query != "" {
		params += "?" + query
	}

	code, raw := g.useCase.Logs(authHeader, params)
	utils.SendRaw(c, code, raw)
}

func (g *gatewayHandler) HandleDelete(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		utils.Send(c, &api.AuthorizationResponse{
			Code:  http.StatusBadRequest,
			Error: "wrong auth token",
		})
		return
	}

	code, raw := g.useCase.DeleteAlbum(authHeader, c.Param("id"))
	utils.SendRaw(c, code, raw)
}

func (g *gatewayHandler) HandleSaveDump(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		utils.Send(c, &api.AuthorizationResponse{
			Code:  http.StatusBadRequest,
			Error: "wrong auth token",
		})
		return
	}

	code, dump := g.useCase.SaveDump(authHeader)
	if code != http.StatusOK {
		utils.SendRaw(c, code, dump)
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.sql", strings.Replace(time.Now().String(), " ", "_", -1)))
	c.Header("Content-Type", "application/sql")
	c.Header("Content-Length", fmt.Sprintf("%d", len(dump)))

	c.Writer.Write(dump)
}

func (g *gatewayHandler) HandleLoadDump(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		utils.Send(c, &api.AuthorizationResponse{
			Code:  http.StatusBadRequest,
			Error: "wrong auth token",
		})
		return
	}

	file, err := c.FormFile("dump")
	if err != nil {
		utils.Send(c, &api.ErrorResponse{
			Code:  http.StatusBadRequest,
			Error: "get form error",
		})
		log.Print(err.Error())
		return
	}

	tempFile, err := os.CreateTemp("", "dump-*.sql")
	if err != nil {
		utils.Send(c, &api.ErrorResponse{
			Code:  http.StatusInternalServerError,
			Error: "create temp file error",
		})
		log.Print(err.Error())
		return
	}
	defer os.Remove(tempFile.Name())

	if err := c.SaveUploadedFile(file, tempFile.Name()); err != nil {
		utils.Send(c, &api.ErrorResponse{
			Code:  http.StatusInternalServerError,
			Error: "save file error",
		})
		log.Print(err.Error())
		return
	}

	code, raw := g.useCase.LoadDump(authHeader, tempFile.Name())
	if code != http.StatusOK {
		utils.SendRaw(c, code, raw)
		return
	}

	c.String(http.StatusOK, "")
}

func handleOrderAction(c *gin.Context, callback func(string, int) (int, []byte)) {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		utils.Send(c, &api.AuthorizationResponse{
			Code:  http.StatusBadRequest,
			Error: "wrong auth token",
		})
		return
	}

	id, err := utils.GetParam[int](c, "id")
	if err != nil {
		utils.Send(c, &api.ErrorResponse{
			Code:  http.StatusBadRequest,
			Error: "invalid id parameter",
		})
		return
	}

	code, raw := callback(authHeader, id)
	utils.SendRaw(c, code, raw)
}
