package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ketul1009/stockscreener-backend/service"
)

type ScreenerHandler struct {
	ScreenerService *service.ScreenerService
}

func (h *ScreenerHandler) CreateScreener(w http.ResponseWriter, r *http.Request) {
	var screener service.Screener
	if err := json.NewDecoder(r.Body).Decode(&screener); err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createdScreener, err := h.ScreenerService.CreateScreener(r.Context(), &screener)
	if err != nil {
		fmt.Println(err.Error())
		if err.Error() == "ERROR: duplicate key value violates unique constraint \"unique_name_user_id\" (SQLSTATE 23505)" {
			respondWithJSON(w, http.StatusConflict, map[string]string{"error": "Screener with this name already exists"})
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to create screener", 500)
		return
	}

	respondWithJSON(w, http.StatusCreated, createdScreener)
}

func (h *ScreenerHandler) GetScreeners(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	id := r.URL.Query().Get("id")
	if username == "" && id == "" {
		respondWithError(w, http.StatusBadRequest, "Username or ID is required", 400)
		return
	}

	if username != "" {
		screeners, err := h.ScreenerService.GetScreeners(r.Context(), username)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to get screeners", 500)
			return
		}
		respondWithJSON(w, http.StatusOK, screeners)
		return
	}

	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID format", 400)
		return
	}

	if id != "" {
		screener, err := h.ScreenerService.GetScreener(r.Context(), idInt)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to get screener", 500)
			return
		}
		respondWithJSON(w, http.StatusOK, screener)
		return
	}

	respondWithError(w, http.StatusBadRequest, "Invalid request", 400)
}

func (h *ScreenerHandler) UpdateScreener(w http.ResponseWriter, r *http.Request) {
	var screener service.Screener
	err := json.NewDecoder(r.Body).Decode(&screener)
	if err != nil {
		fmt.Println(err.Error())
		respondWithError(w, http.StatusBadRequest, "Invalid request body", 400)
		return
	}

	updatedScreener, err := h.ScreenerService.UpdateScreener(r.Context(), screener.ID, &screener)
	if err != nil {
		fmt.Println(err.Error())
		respondWithError(w, http.StatusInternalServerError, "Failed to update screener", 500)
		return
	}

	respondWithJSON(w, http.StatusOK, updatedScreener)
}

func (h *ScreenerHandler) DeleteScreener(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "ID is required", 400)
		return
	}

	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID format", 400)
		return
	}

	_, err = h.ScreenerService.GetScreener(r.Context(), idInt)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Screener not found", 404)
		return
	}

	err = h.ScreenerService.DeleteScreener(r.Context(), idInt)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete screener", 500)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Screener deleted successfully"})
}

func (h *ScreenerHandler) GetJobId(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		respondWithError(w, http.StatusBadRequest, "User ID is required", 400)
		return
	}

	jobTracker, err := h.ScreenerService.GetJobId(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get job ID", 500)
		return
	}

	respondWithJSON(w, http.StatusOK, jobTracker)
}
