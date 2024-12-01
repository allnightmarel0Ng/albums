package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/allnightmarel0Ng/albums/internal/app/order-management/usecase"
	"github.com/allnightmarel0Ng/albums/internal/domain/api"
	"github.com/allnightmarel0Ng/albums/internal/utils"
	"github.com/gin-gonic/gin"
)

var errInvalidRequestBody = errors.New("invalid request body")

type OrderManagementHandler interface {
	HandleAdd(c *gin.Context)
	HandleRemove(c *gin.Context)
	HandleOrders(c *gin.Context)
}

type orderManagementHandler struct {
	useCase usecase.OrderManagementUseCase
}

func NewOrderManagementHandler(useCase usecase.OrderManagementUseCase) OrderManagementHandler {
	return &orderManagementHandler{
		useCase: useCase,
	}
}

func (o *orderManagementHandler) HandleAdd(c *gin.Context) {
	request, err := parseRequestBody(c)
	if err != nil {
		utils.Send(c, &api.OrderActionResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	o.useCase.AddAlbumToUserOrder(request)
}

func (o *orderManagementHandler) HandleRemove(c *gin.Context) {
	request, err := parseRequestBody(c)
	if err != nil {
		utils.Send(c, &api.OrderActionResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	o.useCase.RemoveAlbumFromUserOrder(request)
}

func (o *orderManagementHandler) HandleOrders(c *gin.Context) {
	id, err := utils.GetIDParam(c)
	if err != nil {
		utils.Send(c, &api.UserOrdersResponse{
			Code:  http.StatusBadRequest,
			Error: "invalid id parameter",
		})
		return
	}

	utils.Send(c, o.useCase.UserOrder(id))
}

func parseRequestBody(c *gin.Context) (api.OrderActionRequest, error) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return api.OrderActionRequest{}, errInvalidRequestBody
	}
	defer c.Request.Body.Close()

	var request api.OrderActionRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		return api.OrderActionRequest{}, errInvalidRequestBody
	}

	return request, nil
}
