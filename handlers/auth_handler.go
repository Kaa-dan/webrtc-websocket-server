package handlers

import (
	"net/http"

	"github.com/Kaa-dan/webrtc-websocket-server.git/managers"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	groupName   string
	authManager *managers.AuthManager
}

func NewAuthHandler(auth_Manager *managers.AuthManager) *AuthHandler {
	return &AuthHandler{
		groupName:   "/api/auth",
		authManager: auth_Manager,
	}
}

func (h *AuthHandler) RegisterAuthApis(r *gin.Engine) {
	authGroup := r.Group(h.groupName)
	{
		authGroup.POST("/sign-up", h.SignUp)
	}
}

func (h *AuthHandler) SignUp(ctx *gin.Context) {

	h.handleSuccess(ctx, http.StatusCreated, "User created successfully")
}

func (h *AuthHandler) handleSuccess(ctx *gin.Context, statusCode int, message string) {
	ctx.JSON(statusCode, gin.H{
		"success": true,
		"message": message,
		// "data":    data,
	})
}
