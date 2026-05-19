package handlers

import (
	"net/http"

	"backend-assignment/internal/models"
	"backend-assignment/internal/storage"

	"github.com/gin-gonic/gin"
)

func GetStats(c *gin.Context) {
	storage.Mu.RLock()
	defer storage.Mu.RUnlock()

	response := make(map[string]interface{})

	totalAccepted := 0
	totalRejected := 0

	for userID, data := range storage.RequestStore {
		response[userID] = gin.H{
			"accepted_requests_current_window": len(data.AcceptedTimestamps),
			"rejected_requests_total":          data.RejectedCount,
		}

		totalAccepted += len(data.AcceptedTimestamps)
		totalRejected += data.RejectedCount
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "stats fetched successfully",
		Data: gin.H{
			"users": response,
			"global_totals": gin.H{
				"accepted_requests": totalAccepted,
				"rejected_requests": totalRejected,
			},
		},
	})
}