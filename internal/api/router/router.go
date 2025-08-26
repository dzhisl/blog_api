package router

import (
	"example.com/m/internal/api/handlers"
	"example.com/m/internal/api/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {

	r := gin.New()
	r.Use(middleware.RequestIDMiddleware)
	RouterGroup := r.Group("/api")

	registerPublicRoutes(*RouterGroup)
	registerUserRoutes(*RouterGroup)
	registerAdminRoutes(*RouterGroup)
	return r
}

func registerPublicRoutes(r gin.RouterGroup) {
	r.GET("ping", handlers.PingHandler)
	r.POST("sign-up", handlers.SignUpHandler)
	r.POST("sign-in", handlers.SignInHandler)
	r.POST("refresh-token", handlers.RefreshTokenHandler)
}

func registerUserRoutes(r gin.RouterGroup) {
	r.Use(middleware.UserAuthMiddleware)
	r.GET("ping_user", handlers.PingHandler)
}

func registerAdminRoutes(r gin.RouterGroup) {
	// Apply UserAuthMiddleware first, then AdminAuthMiddleware
	r.Use(middleware.UserAuthMiddleware, middleware.AdminAuthMiddleware)
	r.GET("ping_admin", handlers.PingHandler)
}
