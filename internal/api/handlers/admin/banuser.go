package admin

import (
	"strconv"

	"example.com/m/internal/api/utils"
	"example.com/m/internal/storage"
	"example.com/m/internal/types"
	"example.com/m/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ChangeUserStatushandler(status types.Status) func(c *gin.Context) {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		requestAuthor, err := utils.GetClaims(c)
		if err != nil {
			logger.Error(ctx, "failed to get request initiator claims", zap.Error(err))
			c.JSON(utils.FormInternalErrResponse())
			return
		}

		userIDStr := c.Param("user_id")

		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			c.JSON(utils.FormErrResponse(400, "invalid user_id"))
			return
		}
		user, err := storage.GetUser(userID)
		if err != nil {
			logger.Warn(ctx, "failed to get user", zap.Error(err))
			c.JSON(utils.FormErrResponse(400, "failed to fetch user"))
			return
		}
		if types.CompareRoles(types.Role(requestAuthor.Role), types.Role(user.UserRole)) == -1 {
			c.JSON(utils.FormErrResponse(400, "insufficient rights"))
			return
		}
		err = storage.ChangeUserStatus(userID, status)
		if err != nil {
			logger.Error(ctx, "failed to change user status", zap.Error(err))
			c.JSON(utils.FormInternalErrResponse())
			return
		}
		c.JSON(utils.FormResponse("updated user status"))
	}
}
