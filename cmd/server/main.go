package main

import (
	"backend-assignment/internal/middleware"
	"backend-assignment/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Custom logging middleware
	router.Use(middleware.LoggerMiddleware())

	// Register routes
	routes.RegisterRoutes(router)

	router.Run(":8080")
}