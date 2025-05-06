package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

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
	if username == "" {
		respondWithError(w, http.StatusBadRequest, "Username is required", 400)
		return
	}

	screeners, err := h.ScreenerService.GetScreeners(r.Context(), username)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get screeners", 500)
		return
	}
	respondWithJSON(w, http.StatusOK, screeners)
}
