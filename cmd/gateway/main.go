// Package main is the entry point to the application
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Rugved7/api-gateway/internal/config"
	"github.com/Rugved7/api-gateway/internal/middleware"
	"github.com/Rugved7/api-gateway/internal/middleware/auth"
	"github.com/Rugved7/api-gateway/internal/observability"
	"github.com/Rugved7/api-gateway/internal/server"
)

func main() {
	cfgPath := os.Getenv("GATEWAY_CONFIG")
	if cfgPath == "" {
		cfgPath = "configs/gateway.yaml"
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	handler := http.NewServeMux()
	handler.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	validator := auth.NewValidator(cfg.Auth.JWTSecret)

	wrapperHandler := middleware.Chain(
		handler,
		observability.LoggingMiddleware,
		auth.Middleware(validator),
	)
	srv := server.New(cfg.Server.Address, wrapperHandler)

	go func() {
		log.Printf("gateway listening on %s", cfg.Server.Address)
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("error listening to server : %v", err)
		}
	}()

	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, syscall.SIGINT, syscall.SIGTERM)
	<-shutdownCh

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("shutting down api-gateway")
	if err := srv.ShutDown(ctx); err != nil {
		log.Printf("graceful shutting down failed: %v", err)
	}
}
