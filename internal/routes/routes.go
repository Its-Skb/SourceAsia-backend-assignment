package routes

import (
	"backend-assignment/internal/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"success": true,
			"message": "Backend Assignment API is running",
		})
	})

	// Part 1 routes
	router.POST("/request", handlers.HandleRequest)
	router.GET("/stats", handlers.GetStats)

	// Part 2 routes
	router.POST("/products", handlers.CreateProduct)
	router.GET("/products", handlers.GetProducts)
	router.GET("/products/:id", handlers.GetProductByID)
	router.POST("/products/:id/media", handlers.AddProductMedia)
}
