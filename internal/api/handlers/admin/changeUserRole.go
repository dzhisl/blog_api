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

type changeUserRoleRequest struct {
	NewRole types.Role `json:"new_role"`
}

func ChangeUserRoleHandler(c *gin.Context) {
	ctx := c.Request.Context()
	var req changeUserRoleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Debug(ctx, "invalid request body", zap.Error(err))
		c.JSON(utils.FormInvalidRequestResponse())
		return
	}

	userIDStr := c.Param("user_id")

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(utils.FormErrResponse(400, "invalid user_id"))
		return
	}

	if !types.IsValidRole(req.NewRole) {
		c.JSON(utils.FormErrResponse(400, "invalid role"))
		return
	}

	if err = storage.ChangeUserRole(userID, req.NewRole); err != nil {
		logger.Error(ctx, "failed to change user role", zap.Error(err))
		c.JSON(utils.FormInternalErrResponse())
		return
	}

	c.JSON(utils.FormResponse("changed user role"))

}
