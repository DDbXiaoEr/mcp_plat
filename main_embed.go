//go:build embed

package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"strings"

	"mcp_plat-console/config"
	"mcp_plat-console/database"
	"mcp_plat-console/handler"
	"mcp_plat-console/logging"
	"mcp_plat-console/middleware"

	"github.com/gin-gonic/gin"
)

//go:embed web/dist
var dist embed.FS

func main() {
	config.Load()

	fmt.Printf("mcp_plat-console version %s (commit %s), built at %s\n", Version, GitCommit, BuildTime)

	database.Init()
	logging.Setup()

	staticFS, err := fs.Sub(dist, "web/dist")
	if err != nil {
		log.Fatalf("failed to init embedded frontend: %v", err)
	}

	r := gin.Default()

	authHandler := handler.NewAuthHandler()
	accessKeyHandler := handler.NewAccessKeyHandler()
	historyHandler := handler.NewHistoryHandler()
	serverHandler := handler.NewMCPServerHandler()
	roleHandler := handler.NewRoleHandler()
	userHandler := handler.NewRBACUserHandler()
	settingHandler := handler.NewSettingHandler()

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
		auth.POST("/servers/fetch-tools", serverHandler.FetchTools)
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

		admin.GET("/settings", settingHandler.Get)
		admin.PUT("/settings/:key", settingHandler.Save)
	}

	r.GET("/favicon.svg", func(c *gin.Context) {
		c.FileFromFS("/favicon.svg", http.FS(staticFS))
	})

	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		if path != "/" {
			if f, err := staticFS.Open(strings.TrimPrefix(path, "/")); err == nil {
				if info, err := f.Stat(); err == nil && !info.IsDir() {
					f.Close()
					c.FileFromFS(path, http.FS(staticFS))
					return
				}
				f.Close()
			}
		}

		data, err := fs.ReadFile(staticFS, "index.html")
		if err != nil {
			c.String(http.StatusNotFound, "not found")
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", data)
	})

	log.Printf("server starting on :%s", config.AppConfig.ServerPort)
	if err := r.Run(":" + config.AppConfig.ServerPort); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
