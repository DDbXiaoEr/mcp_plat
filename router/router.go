package router

import (
	"mcp_plat-console/handler"
	"mcp_plat-console/middleware"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	r := gin.Default()

	authHandler := handler.NewAuthHandler()
	accessKeyHandler := handler.NewAccessKeyHandler()
	historyHandler := handler.NewHistoryHandler()
	serverHandler := handler.NewMCPServerHandler()
	roleHandler := handler.NewRoleHandler()
	userHandler := handler.NewRBACUserHandler()

	r.POST("/api/auth/login", authHandler.Login)

	auth := r.Group("/api")
	auth.Use(middleware.AuthRequired())
	{
		auth.GET("/auth/profile", authHandler.Profile)

		auth.GET("/access-keys", accessKeyHandler.List)
		auth.POST("/access-keys", accessKeyHandler.Create)
		auth.PUT("/access-keys/:id", accessKeyHandler.Update)
		auth.DELETE("/access-keys/:id", accessKeyHandler.Delete)

		auth.GET("/history", historyHandler.List)

		auth.GET("/servers", serverHandler.List)
		auth.POST("/servers", serverHandler.Create)
		auth.PUT("/servers/:id", serverHandler.Update)
		auth.DELETE("/servers/:id", serverHandler.Delete)
	}

	admin := r.Group("/api")
	admin.Use(middleware.AuthRequired(), middleware.AdminRequired())
	{
		admin.GET("/roles", roleHandler.List)
		admin.POST("/roles", roleHandler.Create)
		admin.PUT("/roles/:id", roleHandler.Update)
		admin.DELETE("/roles/:id", roleHandler.Delete)
		admin.GET("/roles/:id/users", roleHandler.GetUsers)
		admin.PUT("/roles/:id/users", roleHandler.AssignUsers)

		admin.GET("/users", userHandler.List)
		admin.POST("/users", userHandler.Create)
		admin.PUT("/users/:id", userHandler.Update)
		admin.DELETE("/users/:id", userHandler.Delete)
	}

	return r
}
