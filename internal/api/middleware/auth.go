package middleware

import (
	"net/http"
	"strings"
	"time"

	"example.com/m/internal/api/auth"
	"example.com/m/internal/types"
	"example.com/m/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func UserAuthMiddleware(c *gin.Context) {
	ctx := c.Request.Context()
	header := c.GetHeader("Authorization")
	if header == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
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

	claims, err := auth.ValidateToken(token, types.TokenAccess)
	if err != nil {
		logger.Warn(ctx, "invalid jwt token", zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		c.Abort()
		return
	}

	if time.Now().After(claims.ExpiresAt.Time) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token is expired"})
		c.Abort()
		return
	}

	// Attach claims to context
	c.Set("claims", claims)

	// Continue to next middleware/handler
	c.Next()
}

func AdminAuthMiddleware(c *gin.Context) {
	// Get claims from context (should be set by UserAuthMiddleware)
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

	// Check if user has admin role
	if userClaims.Role != string(types.RoleAdmin) {
		c.JSON(http.StatusForbidden, gin.H{"error": "not enough rights"})
		c.Abort()
		return
	}

	// Continue to next middleware/handler
	c.Next()
}
