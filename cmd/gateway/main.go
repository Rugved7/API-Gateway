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
	"github.com/Rugved7/api-gateway/internal/middleware/ratelimit"
	"github.com/Rugved7/api-gateway/internal/observability"
	"github.com/Rugved7/api-gateway/internal/proxy"
	"github.com/Rugved7/api-gateway/internal/router"
	"github.com/Rugved7/api-gateway/internal/server"
)

func main() {

	// Config loading
	cfgPath := os.Getenv("GATEWAY_CONFIG")
	if cfgPath == "" {
		cfgPath = "configs/gateway.yaml"
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Router setup
	routes := make([]router.Route, 0, len(cfg.Routes))
	for _, r := range cfg.Routes {
		routes = append(routes, router.Route{
			Prefix:   r.Prefix,
			Upstream: r.Upstream,
		})
	}

	rt := router.New(routes)

	// Proxy setup
	px := proxy.New(5 * time.Second)

	// Core handler (routing + proxy)
	baseHandler := http.NewServeMux()
	baseHandler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		route, ok := rt.Match(r)
		if !ok {
			http.NotFound(w, r)
			return
		}

		h, err := px.Handler(route.Upstream)
		if err != nil {
			http.Error(w, "bad gateway", http.StatusBadGateway)
			return
		}

		h.ServeHTTP(w, r)
	})

	// Middleware chain (ORDER IS CRITICAL)
	validator := auth.NewValidator(cfg.Auth.JWTSecret)
	limiter := ratelimit.NewLimiter(
		cfg.RateLimit.Capacity,
		cfg.RateLimit.RefillRate,
	)

	handler := middleware.Chain(
		baseHandler,
		observability.LoggingMiddleware,
		auth.Middleware(validator),
		ratelimit.Middleware(limiter),
	)

	// HTTP server
	srv := server.New(cfg.Server.Address, handler)

	go func() {
		log.Printf("api gateway listening on %s", cfg.Server.Address)
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("shutting down gateway")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.ShutDown(ctx); err != nil {
		log.Printf("shutdown error: %v", err)
	}
}
