package user

import (
	"fmt"

	"example.com/m/internal/api/utils"
	"example.com/m/internal/storage"
	"example.com/m/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type changeProfileReq struct {
	Username  *string `json:"username"`
	FirstName *string `json:"first_name"`
}

func ChangeProfileHandler(c *gin.Context) {
	var req changeProfileReq
	ctx := c.Request.Context()

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Debug(ctx, "invalid request body", zap.Error(err))
		c.JSON(utils.FormInvalidRequestResponse())
		return
	}

	claims, err := utils.GetClaims(c)
	if err != nil {
		logger.Error(ctx, "failed to get claims from ctx", zap.Error(err))
		c.JSON(utils.FormInternalErrResponse())
		return
	}

	// Username validation (cannot be empty)
	if req.Username != nil {
		if len(*req.Username) > maxUsernameLength || len(*req.Username) < minUsernameLength {
			logger.Debug(ctx, "invalid username length")
			c.JSON(utils.FormErrResponse(400,
				fmt.Sprintf("username length should be between %d and %d symbols",
					minUsernameLength, maxUsernameLength)))
			return
		}
	}

	// First name validation (can be empty)
	if req.FirstName != nil {
		if *req.FirstName != "" && len(*req.FirstName) > maxFirstNameLength {
			logger.Debug(ctx, "invalid first name length")
			c.JSON(utils.FormErrResponse(400,
				fmt.Sprintf("first name length should be less than %d symbols",
					maxFirstNameLength)))
			return
		}
	}

	if err = storage.ChangeUserProfile(claims.UserID, req.Username, req.FirstName); err != nil {
		logger.Error(ctx, "failed to update user profile", zap.Error(err))
		c.JSON(utils.FormInternalErrResponse())
		return
	}

	c.JSON(utils.FormResponse("updated profile"))
}
