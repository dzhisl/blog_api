package users

import (
	"net/http"
	"strconv"

	"example.com/m/internal/api/utils"
	"example.com/m/internal/storage"
	"example.com/m/internal/storage/models"
	"example.com/m/internal/types"
	"example.com/m/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// get user by user id
func GetUserHandler(c *gin.Context) {
	ctx := c.Request.Context()
	var userObj *models.User
	claims, err := utils.GetClaims(c)
	if err != nil {
		logger.Error(ctx, "failed to get claims from ctx", zap.Error(err))
		c.JSON(utils.FormInternalErrResponse())
		return
	}

	if c.Query("username") != "" {
		username := c.Query("username")
		user, err := storage.GetUserByUsername(username)
		if err != nil {
			logger.Debug(ctx, "invalid username param", zap.String("username", username), zap.Error(err))
			c.JSON(utils.FormErrResponse(http.StatusBadRequest, "invalid username"))
			return
		}
		userObj = user
	} else if c.Query("user_id") != "" {
		param := c.Query("user_id")
		userID, err := strconv.Atoi(param)
		if err != nil || userID <= 0 {
			logger.Debug(ctx, "invalid user_id param", zap.String("user_id", param), zap.Error(err))
			c.JSON(utils.FormErrResponse(http.StatusBadRequest, "invalid user_id"))
			return
		}
		user, err := storage.GetUser(userID)
		if err != nil {
			logger.Error(ctx, "failed to fetch user", zap.Int("user_id", userID), zap.Error(err))
			c.JSON(utils.FormErrResponse(http.StatusNotFound, "user not found"))
			return
		}
		userObj = user
	} else {
		c.JSON(utils.FormErrResponse(http.StatusNotFound, "provide either username or user_id"))
		return
	}

	if types.CompareRoles(types.Role(claims.Role), types.RoleModerator) >= 0 {
		c.JSON(http.StatusOK, userObj)
		return
	}

	c.JSON(http.StatusOK, models.User{
		ID:        userObj.ID,
		Username:  userObj.Username,
		FirstName: userObj.FirstName,
		CreatedAt: userObj.CreatedAt,
	})
	return
}
