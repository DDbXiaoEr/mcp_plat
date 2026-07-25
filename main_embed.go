//go:build embed

package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"mcp_plat-console/config"
	"mcp_plat-console/database"
	"mcp_plat-console/logging"
	"mcp_plat-console/router"

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

	r := router.Setup()

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

	srv := &http.Server{
		Addr:              ":" + config.AppConfig.ServerPort,
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	go func() {
		log.Printf("server starting on :%s", config.AppConfig.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("server shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}
	log.Println("server exited")
}
