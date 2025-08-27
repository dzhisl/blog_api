package router

import (
	"example.com/m/internal/api/handlers"
	"example.com/m/internal/api/handlers/admin"
	"example.com/m/internal/api/handlers/users"
	"example.com/m/internal/api/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {

	r := gin.New()

	r.Use(middleware.RequestIDMiddleware)
	r.Use(middleware.CorsMiddleware())
	RouterGroup := r.Group("/api")

	registerPublicRoutes(*RouterGroup)
	registerUserRoutes(*RouterGroup)
	registerAdminRoutes(*RouterGroup)
	return r
}

func registerPublicRoutes(r gin.RouterGroup) {
	r.GET("ping", handlers.PingHandler)

	r.POST("sign-up", users.SignUpHandler)
	r.POST("sign-in", users.SignInHandler)
	r.POST("refresh-token", users.RefreshTokenHandler)

}

func registerUserRoutes(r gin.RouterGroup) {
	r.Use(middleware.UserAuthMiddleware)
	r.GET("/user", users.GetUserHandler)
	r.GET("ping_user", handlers.PingHandler)
}

func registerAdminRoutes(r gin.RouterGroup) {
	// Apply UserAuthMiddleware first, then AdminAuthMiddleware
	r.Use(middleware.UserAuthMiddleware, middleware.AdminAuthMiddleware)

	r.GET("admin/ping", handlers.PingHandler)

	r.POST("admin/token/invalidate", admin.InvalidateTokenHandler)
	r.POST("admin/token/invalidate-all", admin.InvalidateAllTokensHandler)
}
