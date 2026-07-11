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
	}

	return r
}
