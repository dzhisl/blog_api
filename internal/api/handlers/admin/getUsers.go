package admin

import (
	"strconv"

	"example.com/m/internal/api/utils"
	"example.com/m/internal/storage"
	"example.com/m/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetAllUsersHandler(c *gin.Context) {
	ctx := c.Request.Context()
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))

	if offset < 0 {
		offset = 0
	}
	if limit <= 0 {
		limit = 100
	}

	users, err := storage.GetUsers(offset, limit)
	if err != nil {
		logger.Error(ctx, "failed to get users from db", zap.Error(err))
		c.JSON(utils.FormInternalErrResponse())
		return
	}
	c.JSON(200, users)
}
