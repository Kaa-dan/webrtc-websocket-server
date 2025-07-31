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

// HandleSuccessWithData sends a success response with data (200 OK)
func HandleSuccessWithData(ctx *gin.Context, message string, data interface{}) {
	response := successResponse(message, data)
	ctx.JSON(http.StatusOK, response)
}

// HandleSuccessMessage sends a success response with just a message (200 OK)
func HandleSuccessMessage(ctx *gin.Context, message string) {
	response := successResponse(message, nil)
	ctx.JSON(http.StatusOK, response)
}

// HandleError sends a standardized error response
func HandleError(ctx *gin.Context, statusCode int, message string, err error) {
	response := errorResponse(statusCode, message, err)
	ctx.JSON(statusCode, response)
}

// HandleBadRequest sends a 400 Bad Request error response
func HandleBadRequest(ctx *gin.Context, message string, err error) {
	response := errorResponse(http.StatusBadRequest, message, err)
	ctx.JSON(http.StatusBadRequest, response)
}

// HandleUnauthorized sends a 401 Unauthorized error response
func HandleUnauthorized(ctx *gin.Context, message string, err error) {
	response := errorResponse(http.StatusUnauthorized, message, err)
	ctx.JSON(http.StatusUnauthorized, response)
}

// HandleForbidden sends a 403 Forbidden error response
func HandleForbidden(ctx *gin.Context, message string, err error) {
	response := errorResponse(http.StatusForbidden, message, err)
	ctx.JSON(http.StatusForbidden, response)
}

// HandleNotFound sends a 404 Not Found error response
func HandleNotFound(ctx *gin.Context, message string, err error) {
	response := errorResponse(http.StatusNotFound, message, err)
	ctx.JSON(http.StatusNotFound, response)
}

// HandleInternalServerError sends a 500 Internal Server Error response
func HandleInternalServerError(ctx *gin.Context, message string, err error) {
	response := errorResponse(http.StatusInternalServerError, message, err)
	ctx.JSON(http.StatusInternalServerError, response)
}

// HandleCustomError handles errors with custom message mapping
func HandleCustomError(ctx *gin.Context, err error, customMessage string) {
	// Determine status code based on error type
	statusCode := http.StatusBadRequest

	switch err {
	case ErrInvalidInput, ErrMissingUsername, ErrMissingEmail, ErrInvalidEmail, ErrInvalidUserID:
		statusCode = http.StatusBadRequest
	default:
		statusCode = http.StatusInternalServerError
	}

	message := customMessage
	if message == "" {
		// Fallback to error message if no custom message provided
		message = err.Error()
	}

	HandleError(ctx, statusCode, message, err)
}
