package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"posixfy-cloud/backend/config"
	"posixfy-cloud/backend/database"
	"posixfy-cloud/backend/router"
	"posixfy-cloud/backend/service"
)

func main() {
	cfg := config.Load()

	db := database.Open(cfg.DBPath)
	defer db.Close()

	userService := service.NewUserService(db)
	userService.Bootstrap(cfg.AdminInitPassword)

	fsClient := service.NewFSClient(cfg.FilebridgeURL, cfg.APIKey)

	r := router.Setup(cfg, userService, fsClient)

	srv := &http.Server{
		Addr:    cfg.ListenAddr,
		Handler: r,
	}

	// Graceful shutdown
	go func() {
		log.Printf("starting server on %s", cfg.ListenAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}
	log.Println("server exited")
}
