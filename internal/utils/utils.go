package utils

import "github.com/gin-gonic/gin"

func Send(c *gin.Context, code int, propertyName, message string) {
	c.Status(code)
	c.JSON(code, gin.H{
		"code":       code,
		propertyName: message,
	})
}