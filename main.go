//go:build !embed

package main

// Author: deepseek-v4-pro / opencode

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"mcp_plat-console/config"
	"mcp_plat-console/database"
	"mcp_plat-console/logging"
	"mcp_plat-console/router"
	"mcp_plat-console/service"
)

func main() {
	port := flag.String("port", "", "server listen port (overrides config.yaml)")
	ver := flag.Bool("version", false, "print version and exit")
	flag.Parse()

	if *ver {
		fmt.Printf("mcp_plat-console version %s (commit %s), built at %s\n", Version, GitCommit, BuildTime)
		os.Exit(0)
	}

	config.Load()
	config.SetServerPort(*port)

	fmt.Printf("mcp_plat-console version %s (commit %s), built at %s\n", Version, GitCommit, BuildTime)

	database.Init()
	logging.Setup()

	service.StartScheduler()

	r := router.Setup()

	srv := &http.Server{
		Addr:              ":" + config.AppConfig.ServerPort,
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1MiB
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
