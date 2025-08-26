package middleware

import (
	"bytes"
	"context"
	"io"
	"time"

	"example.com/m/internal/types"
	"example.com/m/pkg/logger"
	"github.com/gin-gonic/gin"
	uuid "github.com/google/uuid"
	"go.uber.org/zap"
)

func RequestIDMiddleware(c *gin.Context) {
	now := time.Now()
	reqID := uuid.New().String()

	// Store Request ID in context
	ctx := context.WithValue(c.Request.Context(), types.RequestIDKey, reqID)
	c.Request = c.Request.WithContext(ctx)

	// Read and restore the request body
	var bodyBytes []byte
	if c.Request.Body != nil {
		bodyBytes, _ = io.ReadAll(c.Request.Body)
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // Restore body for next handlers

	// Get query and path params
	queryParams := c.Request.URL.Query()
	pathParams := c.Params

	logger.Debug(ctx, "incoming request",
		zap.String("path", c.Request.URL.Path),
		zap.String("method", c.Request.Method),
		zap.ByteString("body", bodyBytes),
		zap.Any("query_params", queryParams),
		zap.Any("path_params", pathParams),
	)

	c.Next()

	if c.Writer.Status() == 200 {
		logger.Debug(ctx, "request completed",
			zap.Int("status", c.Writer.Status()),
			zap.Duration("duration", time.Since(now)),
		)
	} else {
		logger.Debug(ctx, "request completed with non-200 status",
			zap.Int("status", c.Writer.Status()),
			zap.Duration("duration", time.Since(now)),
		)
	}

}
