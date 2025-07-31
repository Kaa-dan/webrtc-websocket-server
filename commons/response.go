package commons

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ResponseData represents the standard API response structure
type ResponseData struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Status  int         `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// successResponse creates a success response structure
func successResponse(message string, data interface{}) *ResponseData {
	return &ResponseData{
		Success: true,
		Message: message,
		Status:  http.StatusOK,
		Data:    data,
	}
}

// errorResponse creates an error response structure
func errorResponse(statusCode int, message string, err error) *ResponseData {
	response := &ResponseData{
		Success: false,
		Message: message,
		Status:  statusCode,
		Data:    nil,
	}

	// Only include error details in debug mode
	if gin.Mode() == gin.DebugMode && err != nil {
		response.Error = err.Error()
	}

	return response
}

// HandleSuccess sends a standardized success response
func HandleSuccess(ctx *gin.Context, statusCode int, message string, data interface{}) {
	response := successResponse(message, data)
	response.Status = statusCode
	ctx.JSON(statusCode, response)
}

// HandleBadRequest sends a 400 Bad Request error response
func HandleBadRequest(ctx *gin.Context, status int, message string, err error) {
	response := errorResponse(status, message, err)
	ctx.JSON(status, response)
}
