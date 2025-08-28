package router

import (
	"example.com/m/internal/api/handlers"
	"example.com/m/internal/api/handlers/admin"
	"example.com/m/internal/api/handlers/user"
	"example.com/m/internal/api/middleware"
	"example.com/m/internal/types"
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

	r.POST("sign-up", user.SignUpHandler)
	r.POST("sign-in", user.SignInHandler)
	r.POST("refresh-token", user.RefreshTokenHandler)

}

func registerUserRoutes(r gin.RouterGroup) {
	r.Use(middleware.UserAuthMiddleware)
	r.GET("/user", user.GetUserHandler)
	r.POST("/profile/update", user.ChangeProfileHandler)
}

func registerAdminRoutes(r gin.RouterGroup) {
	r.Use(middleware.UserAuthMiddleware, middleware.RoleAuthMiddleware(types.RoleAdmin))
	r.POST("admin/token/invalidate", admin.InvalidateTokenHandler)
	r.POST("admin/token/invalidate-all", admin.InvalidateAllTokensHandler)
	r.POST("admin/user/:user_id/ban", admin.ChangeUserStatushandler(types.StatusBanned))
	r.POST("admin/user/:user_id/unban", admin.ChangeUserStatushandler(types.StatusOk))
	r.POST("admin/user/:user_id/change-role", middleware.RoleAuthMiddleware(types.RoleOwner), admin.ChangeUserRoleHandler)
	r.GET("admin/users", middleware.RoleAuthMiddleware(types.RoleSuperAdmin), admin.GetAllUsersHandler)
}
