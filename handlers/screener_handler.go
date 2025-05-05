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
		http.Error(w, "Failed to create screener", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusCreated, createdScreener)
}

func (h *ScreenerHandler) GetScreener(w http.ResponseWriter, r *http.Request) {

}
