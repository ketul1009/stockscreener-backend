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

	token, err := h.AuthService.Login(r.Context(), loginRequest.Email, loginRequest.Password)
	if err != nil {
		logger.Error("Failed to login", zap.Error(err))
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	response := struct {
		Token string `json:"token"`
	}{
		Token: token,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
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

	token, err := h.AuthService.Register(r.Context(), registerRequest.Username, registerRequest.Email, registerRequest.Password)
	if err != nil {
		http.Error(w, "Failed to register", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(successResponse{Data: map[string]string{"token": token}})
}
