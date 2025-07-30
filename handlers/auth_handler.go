package handlers

import "github.com/Kaa-dan/webrtc-websocket-server.git/managers"

type AuthHandler struct {
	group       string
	authManager *managers.AuthManager
}

func NewAuthHandler(auth_Manager *managers.AuthManager) *AuthHandler {
	return &AuthHandler{
		group:       "/api/users",
		authManager: auth_Manager,
	}
}
