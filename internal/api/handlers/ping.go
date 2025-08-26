package handlers

import (
	"example.com/m/pkg/logger"
	"github.com/gin-gonic/gin"
)

func PingHandler(c *gin.Context) {
	ctx := c.Request.Context()
	logger.Debug(ctx, "pong")

	c.JSON(200, gin.H{
		"message": "pong",
	})
}
