package utils

import (
	"github.com/allnightmarel0Ng/albums/internal/domain/model"
	"github.com/gin-gonic/gin"
)

func Send(c *gin.Context, response model.Response) {
	propertyName, propertyData := response.GetKeyValue()
	c.JSON(response.GetCode(), gin.H{
		"code":       response.GetCode(),
		propertyName: propertyData,
	})
}
