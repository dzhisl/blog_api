package handlers

import (
	"net/http"

	"example.com/m/internal/api/utils"
	"example.com/m/internal/storage"
	"example.com/m/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type invalidateTokenRequest struct {
	Token string `json:"token"`
}

type invalidateAllTokensRequest struct {
	UserID int `json:"user_id"`
}

// InvalidateTokenHandler invalidates a specific token
func InvalidateTokenHandler(c *gin.Context) {
	var req invalidateTokenRequest
	ctx := c.Request.Context()

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Debug(ctx, "invalid request body", zap.Error(err))
		c.JSON(utils.FormInvalidRequestResponse())
		return
	}

	if err := storage.InvalidateToken(req.Token); err != nil {
		logger.Error(ctx, "failed to invalidate token", zap.Error(err))
		c.JSON(utils.FormErrResponse(http.StatusBadRequest, err.Error()))
		return
	}
	c.JSON(utils.FormResponse("token invalidated successfully"))
}

// InvalidateAllTokensHandler invalidates all tokens for a user
func InvalidateAllTokensHandler(c *gin.Context) {
	var req invalidateAllTokensRequest
	ctx := c.Request.Context()

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Debug(ctx, "invalid request body", zap.Error(err))
		c.JSON(utils.FormInvalidRequestResponse())
		return
	}

	if err := storage.InvalidateAllTokensByUser(req.UserID); err != nil {
		logger.Error(ctx, "failed to invalidate all tokens", zap.Error(err))
		c.JSON(utils.FormErrResponse(http.StatusBadRequest, err.Error()))
		return
	}

	c.JSON(utils.FormResponse("all tokens invalidated successfully"))
}
