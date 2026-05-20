package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Process request
		c.Next()

		duration := time.Since(startTime)

		log.Printf(
			"%s | %d | %s | %s %s",
			startTime.Format(time.RFC3339),
			c.Writer.Status(),
			duration,
			c.Request.Method,
			c.Request.URL.Path,
		)
	}
}
