package middleware

import (
	"net/http"
	"strings"
	"time"

	"example.com/m/internal/api/auth"
	"example.com/m/internal/api/utils"
	"example.com/m/internal/storage"
	"example.com/m/internal/types"
	"example.com/m/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func UserAuthMiddleware(c *gin.Context) {
	ctx := c.Request.Context()
	header := c.GetHeader("Authorization")
	if header == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization token required"})
		c.Abort()
		return
	}

	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
		c.Abort()
		return
	}
	token := parts[1]

	claims, err := auth.ValidateToken(token, types.TokenAccess, true)
	if err != nil {
		logger.Warn(ctx, "invalid jwt token", zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		c.Abort()
		return
	}

	if time.Now().After(claims.ExpiresAt.Time) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
		c.Abort()
		return
	}

	t, err := storage.GetJwtToken(token)
	if err != nil {
		logger.Error(ctx, "failed to get token from db", zap.Error(err))
		c.JSON(utils.FormInternalErrResponse())
		c.Abort()
		return
	}

	if !t.Active {
		logger.Warn(ctx, "inactive jwt token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "inactive token"})
		c.Abort()
		return
	}

	// Attach claims to context
	c.Set("claims", claims)

	c.Next()
}

func RoleAuthMiddleware(role types.Role) func(c *gin.Context) {
	return func(c *gin.Context) {
		claims, ok := c.Get("claims")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "no claims in context"})
			c.Abort()
			return
		}

		userClaims, ok := claims.(*auth.Claims)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid claims type"})
			c.Abort()
			return
		}

		if types.CompareRoles(types.Role(userClaims.Role), role) == -1 {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden: insufficient rights"})
			c.Abort()
			return
		}
		c.Next()
	}
}
