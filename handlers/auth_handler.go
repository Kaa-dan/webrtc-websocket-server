package handlers

import (
	"net/http"

	"github.com/Kaa-dan/webrtc-websocket-server.git/commons"
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
	userData := commons.NewSignupInput()
	if err := ctx.ShouldBindJSON(userData); err != nil {
		commons.HandleError(ctx, http.StatusBadRequest, "Invalid JSON data", err)
		return
	}

	// Validate input data
	if err := commons.HandleValidationError(userData); err != nil {
		commons.HandleCustomError(ctx, err, "Validation failed")
		return
	}

	newUser, err := h.authManager.SignUp(userData)
	if err != nil {
		commons.HandleCustomError(ctx, err, "Failed to create user")
		return
	}

	commons.HandleSuccess(ctx, http.StatusCreated, "User created successfully", newUser)
}
