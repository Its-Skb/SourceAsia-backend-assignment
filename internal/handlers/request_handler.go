package handlers

import (
	"net/http"

	"backend-assignment/internal/limiter"
	"backend-assignment/internal/models"

	"github.com/gin-gonic/gin"
)

func HandleRequest(c *gin.Context) {
	var req models.RequestPayload

	// Validate JSON body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "invalid JSON body",
		})
		return
	}

	// Validate user_id
	if req.UserID == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "user_id is required",
		})
		return
	}

	// Validate payload
	if req.Payload == nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "payload is required",
		})
		return
	}

	// Apply rate limiting
	allowed := limiter.AllowRequest(req.UserID)

	if !allowed {
		c.JSON(http.StatusTooManyRequests, models.ErrorResponse{
			Success: false,
			Error:   "rate limit exceeded",
		})
		return
	}

	// Success response
	c.JSON(http.StatusCreated, models.SuccessResponse{
		Success: true,
		Message: "request accepted",
	})
}