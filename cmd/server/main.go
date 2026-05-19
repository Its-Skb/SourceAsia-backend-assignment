package main

import (
	"backend-assignment/internal/middleware"
	"backend-assignment/internal/routes"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {

	// Use clean Gin engine
	router := gin.New()

	// Trusted proxies fix
	router.SetTrustedProxies(nil)

	// Middleware
	router.Use(gin.Recovery())
	router.Use(middleware.LoggerMiddleware())

	// Register routes
	routes.RegisterRoutes(router)

	// HTTP server configuration
	server := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    30 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	// Run server in goroutine
	go func() {
		log.Println("Server running on port 8080")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server startup failed: %v", err)
		}
	}()

	// Graceful shutdown listener
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	log.Println("Shutting down server...")

	// Shutdown timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Forced shutdown: %v", err)
	}

	log.Println("Server exited cleanly")
}