package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ketul1009/stockscreener-backend/pkg/logger"
	"github.com/ketul1009/stockscreener-backend/service"
	"go.uber.org/zap"
)

type AuthHandler struct {
	AuthService *service.AuthService
}

func (h *AuthHandler) HandlerLogin(w http.ResponseWriter, r *http.Request) {
	var loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response, err := h.AuthService.Login(r.Context(), loginRequest.Email, loginRequest.Password)
	if err != nil {
		logger.Error("Failed to login", zap.Error(err))
		respondWithJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
		return
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (h *AuthHandler) HandlerRegister(w http.ResponseWriter, r *http.Request) {
	var registerRequest struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&registerRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token, code, err := h.AuthService.Register(r.Context(), registerRequest.Username, registerRequest.Email, registerRequest.Password)
	if err != nil {
		logger.Error("Failed to register", zap.Error(err))
		respondWithJSON(w, code, map[string]string{"error": err.Error()})
		return
	}

	respondWithJSON(w, code, map[string]string{"token": token})
}

func (h *AuthHandler) HandlerGetUserFromToken(w http.ResponseWriter, r *http.Request) {
	// Get token from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		respondWithJSON(w, http.StatusUnauthorized, map[string]string{"error": "missing authorization header"})
		return
	}

	// Extract token from "Bearer <token>"
	token := authHeader
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		token = authHeader[7:]
	}

	user, err := h.AuthService.GetUserFromToken(r.Context(), token)
	if err != nil {
		logger.Error("Failed to get user from token", zap.Error(err))
		respondWithJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}

func (h *AuthHandler) HandlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	var updateRequest struct {
		ID       string `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.AuthService.UpdateUser(r.Context(), updateRequest.ID, updateRequest.Username, updateRequest.Email)
	if err != nil {
		logger.Error("Failed to update user", zap.Error(err))
		respondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}
