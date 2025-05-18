package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ketul1009/stockscreener-backend/service"
)

type WatchlistHandler struct {
	WatchlistService *service.WatchlistService
}

func (h *WatchlistHandler) CreateWatchlist(w http.ResponseWriter, r *http.Request) {
	var watchlist service.Watchlist
	if err := json.NewDecoder(r.Body).Decode(&watchlist); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), 400)
		return
	}
	createdWatchlist, err := h.WatchlistService.CreateWatchlist(r.Context(), &watchlist)
	if err != nil {
		if err.Error() == "watchlist name already exists" {
			respondWithError(w, http.StatusBadRequest, err.Error(), 400)
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error(), 500)
		}
		return
	}
	respondWithJSON(w, http.StatusCreated, createdWatchlist)
}

func (h *WatchlistHandler) GetWatchlist(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	idInt, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID format", 400)
		return
	}
	watchlist, err := h.WatchlistService.GetWatchlist(r.Context(), int32(idInt))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), 500)
	}
	respondWithJSON(w, http.StatusOK, watchlist)
}

func (h *WatchlistHandler) GetAllWatchlists(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	watchlists, err := h.WatchlistService.GetAllWatchlists(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), 500)
	}
	respondWithJSON(w, http.StatusOK, watchlists)
}

func (h *WatchlistHandler) UpdateWatchlist(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	idInt, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID format", 400)
		return
	}
	var watchlist service.Watchlist
	if err := json.NewDecoder(r.Body).Decode(&watchlist); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), 400)
		return
	}
	updatedWatchlist, err := h.WatchlistService.UpdateWatchlist(r.Context(), int32(idInt), &watchlist)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), 500)
		return
	}
	respondWithJSON(w, http.StatusOK, updatedWatchlist)
}

func (h *WatchlistHandler) DeleteWatchlist(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	idInt, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID format", 400)
		return
	}
	err = h.WatchlistService.DeleteWatchlist(r.Context(), int32(idInt))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), 500)
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Watchlist deleted successfully"})
}
