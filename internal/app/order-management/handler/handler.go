package handler

import (
	"errors"
	"net/http"
	"strconv"

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
		utils.Send(c, &api.ErrorResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	response := o.useCase.AddAlbumToUserOrder(request)
	if response != nil {
		utils.Send(c, response)
		return
	}

	c.String(http.StatusOK, "")
}

func (o *orderManagementHandler) HandleRemove(c *gin.Context) {
	request, err := parseRequestBody(c)
	if err != nil {
		utils.Send(c, &api.ErrorResponse{
			Code:  http.StatusBadRequest,
			Error: err.Error(),
		})
		return
	}

	response := o.useCase.RemoveAlbumFromUserOrder(request)
	if response != nil {
		utils.Send(c, response)
		return
	}

	c.String(http.StatusOK, "")
}

func (o *orderManagementHandler) HandleOrders(c *gin.Context) {
	id, err := utils.GetParam(c, "id")
	if err != nil {
		utils.Send(c, &api.ErrorResponse{
			Code:  http.StatusBadRequest,
			Error: "invalid id parameter",
		})
		return
	}

	unpaidOnlyStr := c.DefaultQuery("unpaidOnly", "false")
	unpaidOnly, err := strconv.ParseBool(unpaidOnlyStr)
	if err != nil {
		utils.Send(c, &api.ErrorResponse{
			Code:  http.StatusBadRequest,
			Error: "invalid 'unPaidOnly' parameter",
		})
		return
	}

	response := o.useCase.UserOrder(id, unpaidOnly)
	if response != nil {
		utils.Send(c, response)
		return
	}

	c.String(http.StatusOK, "")
}

func parseRequestBody(c *gin.Context) (api.OrderActionRequest, error) {
	var request api.OrderActionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		return api.OrderActionRequest{}, errInvalidRequestBody
	}

	return request, nil
}
