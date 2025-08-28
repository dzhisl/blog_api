package user

import (
	"fmt"
	"strings"

	"example.com/m/internal/api/auth"
	"example.com/m/internal/api/utils"
	"example.com/m/internal/storage"
	"example.com/m/internal/storage/models"
	"example.com/m/internal/types"
	"example.com/m/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	minUsernameLength  = 3
	maxUsernameLength  = 15
	maxFirstNameLength = 20
	minPasswordLength  = 8
	maxPasswordLength  = 72
)

type signUpRequest struct {
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	Password  string `json:"password"`
}

type signUpResponse struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

func SignUpHandler(c *gin.Context) {
	var req signUpRequest
	ctx := c.Request.Context()
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Debug(ctx, "invalid request body", zap.Error(err))
		c.JSON(utils.FormInvalidRequestResponse())
		return
	}

	if len(req.Username) > maxUsernameLength || len(req.Username) < minUsernameLength {
		logger.Debug(ctx, "invalid username length")
		c.JSON(utils.FormErrResponse(400, fmt.Sprintf("username length should be between %d and %d symbols", minUsernameLength, maxUsernameLength)))
		return
	}

	if len(req.FirstName) > maxFirstNameLength {
		logger.Debug(ctx, "invalid first name length")
		c.JSON(utils.FormErrResponse(400, fmt.Sprintf("first name length should be less than %d symbols", maxFirstNameLength)))
		return
	}

	if len(req.Password) < minPasswordLength || len(req.Password) > maxPasswordLength {
		logger.Debug(ctx, "invalid first name length")
		c.JSON(utils.FormErrResponse(400, fmt.Sprintf("password should be between %d and %d symbols", minPasswordLength, maxPasswordLength)))
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		logger.Error(ctx, "error hashing password", zap.Error(err))
		c.JSON(utils.FormInternalErrResponse())
		return
	}

	user := models.NewUserObject(hashedPassword, req.Username, req.FirstName, types.RoleUser, types.StatusOk)

	err = storage.CreateUser(user)
	if err != nil {
		logger.Error(ctx, "failed to create user", zap.Error(err), zap.Any("user object", user))
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			c.JSON(utils.FormErrResponse(400, "user with this username already exists"))
			return
		}

		c.JSON(utils.FormInternalErrResponse())
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
