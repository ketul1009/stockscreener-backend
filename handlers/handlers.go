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
	engine "github.com/ketul1009/stockscreener-backend/stock-engine"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type ApiConfig struct {
	DB              *db.Queries
	AuthService     *service.AuthService
	ScreenerService *service.ScreenerService
	RedisClient     *redis.Client
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

func (cfg *ApiConfig) HandlerDeleteScreener(w http.ResponseWriter, r *http.Request) {
	screenerHandler := ScreenerHandler{ScreenerService: cfg.ScreenerService}
	screenerHandler.DeleteScreener(w, r)
}

func (cfg *ApiConfig) HandlerGetUserFromToken(w http.ResponseWriter, r *http.Request) {
	authHandler := AuthHandler{AuthService: cfg.AuthService}
	authHandler.HandlerGetUserFromToken(w, r)
}

func (cfg *ApiConfig) HandlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	authHandler := AuthHandler{AuthService: cfg.AuthService}
	authHandler.HandlerUpdateUser(w, r)
}

// Handler to produce a new screener job
func (cfg *ApiConfig) HandlerCreateJob(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Rules  []engine.Rule `json:"rules"`
		UserID string        `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if params.UserID == "" {
		respondWithError(w, http.StatusBadRequest, "User ID is required", http.StatusBadRequest)
		return
	}

	jobID := uuid.New().String()
	jobTracker, err := cfg.DB.GetJobTrackerByUserID(r.Context(), pgtype.UUID{Bytes: uuid.MustParse(params.UserID), Valid: true})

	if err != nil && err.Error() != "sql: no rows in result set" {
		jobTracker, err = cfg.DB.CreateJobTracker(r.Context(), db.CreateJobTrackerParams{
			JobID:        pgtype.UUID{Bytes: uuid.MustParse(jobID), Valid: true},
			UserID:       pgtype.UUID{Bytes: uuid.MustParse(params.UserID), Valid: true},
			JobStatus:    "pending",
			JobCreatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
			JobUpdatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
		})

		if err != nil {
			logger.Error("Failed to create job tracker", zap.Error(err))
			respondWithError(w, http.StatusInternalServerError, "Failed to create job tracker", http.StatusInternalServerError)
			return
		}
	} else if jobTracker.JobStatus == "running" || jobTracker.JobStatus == "pending" {
		respondWithError(w, http.StatusBadRequest, "User already has a running or pending job", http.StatusBadRequest)
		return
	}

	jobTracker, err = cfg.DB.UpdateJobTrackerForNewJob(r.Context(), db.UpdateJobTrackerForNewJobParams{
		JobID:        pgtype.UUID{Bytes: uuid.MustParse(jobID), Valid: true},
		UserID:       pgtype.UUID{Bytes: uuid.MustParse(params.UserID), Valid: true},
		JobStatus:    "pending",
		JobCreatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
		JobUpdatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
	})

	if err != nil {
		logger.Error("Failed to update job tracker", zap.Error(err))
		respondWithError(w, http.StatusInternalServerError, "Failed to update job tracker [Error: UpdateJobTrackerForNewJob]", http.StatusInternalServerError)
		return
	}

	job := engine.ScreenerJob{
		JobID:  jobID,
		Rules:  params.Rules,
		UserID: params.UserID,
	}
	jobJSON, err := json.Marshal(job)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to marshal job", http.StatusInternalServerError)
		return
	}

	// Push job to Redis queue
	err = cfg.RedisClient.LPush(r.Context(), "screener_jobs", jobJSON).Err()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to enqueue job", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusAccepted, map[string]string{"job_id": jobID})
}

// Handler to fetch job result
func (cfg *ApiConfig) HandlerGetJobResult(w http.ResponseWriter, r *http.Request) {
	jobID := r.URL.Query().Get("job_id")
	if jobID == "" {
		respondWithError(w, http.StatusBadRequest, "job_id is required", http.StatusBadRequest)
		return
	}

	jobTracker, err := cfg.DB.GetJobTrackerByJobID(r.Context(), pgtype.UUID{Bytes: uuid.MustParse(jobID), Valid: true})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get job tracker", http.StatusInternalServerError)
		return
	}

	if jobTracker.JobStatus == "pending" || jobTracker.JobStatus == "running" {
		respondWithError(w, http.StatusAccepted, "Job is still pending", http.StatusAccepted)
		return
	}

	result, err := cfg.RedisClient.Get(r.Context(), "screener_result:"+jobID).Result()
	if err == redis.Nil {
		respondWithError(w, http.StatusNotFound, "Result not ready", http.StatusNotFound)
		return
	} else if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch result", http.StatusInternalServerError)
		logger.Error("Failed to fetch result", zap.Error(err))
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": result, "job_id": jobTracker.JobID.String(), "job_status": jobTracker.JobStatus})
}

func (cfg *ApiConfig) HandlerGetJobId(w http.ResponseWriter, r *http.Request) {
	screenerHandler := ScreenerHandler{ScreenerService: cfg.ScreenerService}
	screenerHandler.GetJobId(w, r)
}
