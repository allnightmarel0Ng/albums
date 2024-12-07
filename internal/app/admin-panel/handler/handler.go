package handler

import (
	"net/http"
	"strconv"

	"github.com/allnightmarel0Ng/albums/internal/app/admin-panel/usecase"
	"github.com/allnightmarel0Ng/albums/internal/domain/api"
	"github.com/allnightmarel0Ng/albums/internal/utils"
	"github.com/gin-gonic/gin"
)

type AdminPanelHandler interface {
	HandleBuyLogs(c *gin.Context)
	HandleDeleteAlbum(c *gin.Context)
}

type adminPanelHandler struct {
	useCase usecase.AdminPanelUseCase
}

func NewAdminPanelHandler(useCase usecase.AdminPanelUseCase) AdminPanelHandler {
	return &adminPanelHandler{
		useCase: useCase,
	}
}

func (a *adminPanelHandler) HandleBuyLogs(c *gin.Context) {
	paramStr, ok := c.Params.Get("pageNumber")
	if !ok {
		utils.Send(c, &api.ErrorResponse{
			Code:  http.StatusBadRequest,
			Error: "invalid 'pageNumber' parameter",
		})
		return
	}

	pageNumber, err := strconv.ParseUint(paramStr, 10, 64)
	if err != nil {
		utils.Send(c, &api.ErrorResponse{
			Code:  http.StatusBadRequest,
			Error: "invalid 'pageNumber' parameter",
		})
		return
	}

	pageSizeStr := c.DefaultQuery("pageSize", "10")
	pageSize, err := strconv.ParseUint(pageSizeStr, 10, 64)
	if err != nil {
		utils.Send(c, &api.ErrorResponse{
			Code:  http.StatusBadRequest,
			Error: "invalid 'pageSize' parameter",
		})
		return
	}

	utils.Send(c, a.useCase.Logs(uint(pageNumber), uint(pageSize)))
}

func (a *adminPanelHandler) HandleDeleteAlbum(c *gin.Context) {
	id, err := utils.GetParam(c, "id")
	if err != nil {
		utils.Send(c, &api.ErrorResponse{
			Code:  http.StatusBadRequest,
			Error: "invalid 'id' parameter",
		})
		return
	}

	response := a.useCase.DeleteAlbum(id)
	if response != nil {
		utils.Send(c, response)
		return
	}

	c.String(http.StatusOK, "")
}
