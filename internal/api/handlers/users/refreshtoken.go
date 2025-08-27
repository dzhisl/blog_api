package users

import (
	"net/http"

	"example.com/m/internal/api/auth"
	"example.com/m/internal/api/utils"
	"example.com/m/internal/storage"
	"example.com/m/internal/types"
	"example.com/m/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type refreshTokenRequest struct {
	AccessToken string `json:"access_token"`
}

func RefreshTokenHandler(c *gin.Context) {
	var req refreshTokenRequest
	ctx := c.Request.Context()
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Debug(ctx, "invalid request body", zap.Error(err))
		c.JSON(utils.FormInvalidRequestResponse())
		return
	}
	claims, err := auth.ValidateToken(req.AccessToken, types.TokenRefresh, true)
	if err != nil {
		c.JSON(utils.FormErrResponse(400, "invalid refresh token"))
		return
	}
	t, err := storage.GetJwtToken(req.AccessToken)
	if err != nil {
		logger.Error(ctx, "failed to get token from db", zap.Error(err))
		c.JSON(utils.FormInternalErrResponse())
		return
	}
	logger.Debug(ctx, "received token from db", zap.Any("", t))
	if !t.Active {
		logger.Warn(ctx, "inactive jwt token", zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}
	user, err := storage.GetUser(claims.UserID)
	if err != nil {
		logger.Error(ctx, "failed to get user from DB", zap.Error(err))
		c.JSON(utils.FormInternalErrResponse())
		return
	}
	newAccessToken, err := auth.GenerateToken(*user, types.TokenAccess)
	if err != nil {
		logger.Error(ctx, "failed to create jwt token", zap.Error(err))
		c.JSON(utils.FormInternalErrResponse())
		return
	}
	if err = storage.StoreToken(newAccessToken); err != nil {
		logger.Error(ctx, "failed to store jwt token in DB", zap.Error(err))
		c.JSON(utils.FormInternalErrResponse())
		return
	}

	c.JSON(200, signUpResponse{
		AccessToken: newAccessToken,
	})
}
