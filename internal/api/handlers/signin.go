package handlers

import (
	"example.com/m/internal/api/auth"
	"example.com/m/internal/api/utils"
	"example.com/m/internal/storage"
	"example.com/m/internal/types"
	"example.com/m/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type signInRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func SignInHandler(c *gin.Context) {
	var req signInRequest
	ctx := c.Request.Context()
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Debug(ctx, "invalid request body", zap.Error(err))
		c.JSON(utils.FormInvalidRequestResponse())
		return
	}
	user, err := storage.GetUserByUsername(req.Username)
	if err != nil {
		logger.Debug(ctx, "attempt to login with invalid username")
		c.JSON(utils.FormErrResponse(400, "invalid credentials"))
		return
	}
	ok := auth.CheckPasswordHash(req.Password, user.PasswordHash)
	if !ok {
		logger.Debug(ctx, "attempt to login with invalid password")
		c.JSON(utils.FormErrResponse(400, "invalid credentials"))
		return
	}

	accessToken, err := auth.GenerateToken(*user, types.TokenAccess)
	if err != nil {
		logger.Error(ctx, "failed to create jwt token", zap.Error(err))
		c.JSON(utils.FormInternalErrResponse())
		return
	}
	if err = storage.StoreToken(accessToken); err != nil {
		logger.Error(ctx, "failed to store jwt token in DB", zap.Error(err))
		c.JSON(utils.FormInternalErrResponse())
		return
	}

	refreshToken, err := auth.GenerateToken(*user, types.TokenRefresh)
	if err != nil {
		logger.Error(ctx, "failed to create jwt token", zap.Error(err))
		c.JSON(utils.FormInternalErrResponse())
		return
	}

	if err = storage.StoreToken(refreshToken); err != nil {
		logger.Error(ctx, "failed to store jwt token in DB", zap.Error(err))
		c.JSON(utils.FormInternalErrResponse())
		return
	}

	resp := signUpResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	c.JSON(200, resp)
}
