package routes

import (
	"backend-assignment/internal/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {

	// Part 1 routes
	router.POST("/request", handlers.HandleRequest)
	router.GET("/stats", handlers.GetStats)

	// Part 2 routes
	router.POST("/products", handlers.CreateProduct)
}