package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ketul1009/stockscreener-backend/db"
	"github.com/ketul1009/stockscreener-backend/pkg/logger"
	"github.com/ketul1009/stockscreener-backend/service"
	"go.uber.org/zap"
)

type ApiConfig struct {
	DB              *db.Queries
	AuthService     *service.AuthService
	ScreenerService *service.ScreenerService
}

type errorResponse struct {
	Error      string `json:"error"`
	StatusCode int    `json:"status_code"`
}

type successResponse struct {
	Data       interface{} `json:"data"`
	StatusCode int         `json:"status_code"`
}

func ReadinessHandler(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, successResponse{
		Data: map[string]string{
			"status": "ok",
		},
		StatusCode: http.StatusOK,
	})
}

func ErrHandler(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, http.StatusInternalServerError, "Internal Server Error", http.StatusInternalServerError)
}

func (cfg *ApiConfig) HandlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name         string `json:"name"`
		Username     string `json:"username"`
		Email        string `json:"email"`
		PasswordHash string `json:"password_hash"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		logger.Error("Failed to decode request body", zap.Error(err))
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", http.StatusBadRequest)
		return
	}

	now := time.Now()
	user, err := cfg.DB.CreateUser(r.Context(), db.CreateUserParams{
		ID:           pgtype.UUID{Bytes: uuid.New(), Valid: true},
		CreatedAt:    pgtype.Timestamp{Time: now, Valid: true},
		UpdatedAt:    pgtype.Timestamp{Time: now, Valid: true},
		Username:     params.Username,
		Email:        params.Email,
		PasswordHash: params.PasswordHash,
	})
	if err != nil {
		logger.Error("Failed to create user", zap.Error(err))
		respondWithError(w, http.StatusInternalServerError, "Failed to create user", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusCreated, successResponse{
		Data: user,
	})
}

func (cfg *ApiConfig) HandlerGetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := cfg.DB.GetUsers(r.Context())
	if err != nil {
		logger.Error("Failed to fetch users", zap.Error(err))
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch users", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, successResponse{
		Data: users,
	})
}

func respondWithError(w http.ResponseWriter, code int, message string, statusCode int) {
	respondWithJSON(w, code, errorResponse{Error: message, StatusCode: statusCode})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		logger.Error("Failed to marshal response", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (cfg *ApiConfig) HandlerLogin(w http.ResponseWriter, r *http.Request) {
	authHandler := AuthHandler{AuthService: cfg.AuthService}
	authHandler.HandlerLogin(w, r)
}

func (cfg *ApiConfig) HandlerRegister(w http.ResponseWriter, r *http.Request) {
	authHandler := AuthHandler{AuthService: cfg.AuthService}
	authHandler.HandlerRegister(w, r)
}

func (cfg *ApiConfig) HandlerCreateScreener(w http.ResponseWriter, r *http.Request) {
	screenerHandler := ScreenerHandler{ScreenerService: cfg.ScreenerService}
	screenerHandler.CreateScreener(w, r)
}

func (cfg *ApiConfig) HandlerGetScreeners(w http.ResponseWriter, r *http.Request) {
	screenerHandler := ScreenerHandler{ScreenerService: cfg.ScreenerService}
	screenerHandler.GetScreeners(w, r)
}

func (cfg *ApiConfig) HandlerUpdateScreener(w http.ResponseWriter, r *http.Request) {
	screenerHandler := ScreenerHandler{ScreenerService: cfg.ScreenerService}
	screenerHandler.UpdateScreener(w, r)
}

func (cfg *ApiConfig) HandlerGetUserFromToken(w http.ResponseWriter, r *http.Request) {
	authHandler := AuthHandler{AuthService: cfg.AuthService}
	authHandler.HandlerGetUserFromToken(w, r)
}
