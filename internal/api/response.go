package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response represents the standard JSON response structure
type Response struct {
	Status  string      `json:"status"`            // "success" or "error"
	Message string      `json:"message,omitempty"` // Success or error message
	Data    interface{} `json:"data,omitempty"`    // Payload for success responses
	Error   string      `json:"error,omitempty"`   // Error details
}

// SendSuccess sends a JSON success response
func SendSuccess(c *gin.Context, message string, data interface{}) {
	response := Response{
		Status:  "success",
		Message: message,
		Data:    data,
	}
	c.JSON(http.StatusOK, response)
}

// SendError sends a JSON error response
func SendError(c *gin.Context, message string, errDetail string, statusCode int) {
	response := Response{
		Status:  "error",
		Message: message,
		Error:   errDetail,
	}
	c.JSON(statusCode, response)
}
