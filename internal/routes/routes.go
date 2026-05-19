package routes

import (
	"backend-assignment/internal/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	router.POST("/request", handlers.HandleRequest)
	router.GET("/stats", handlers.GetStats)
}