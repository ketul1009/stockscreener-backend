package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/ketul1009/stockscreener-backend/db"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	now := pgtype.Timestamp{Time: time.Now().UTC(), Valid: true}
	user, err := cfg.DB.CreateUser(r.Context(), db.CreateUserParams{
		ID:        pgtype.UUID{Bytes: uuid.New(), Valid: true},
		CreatedAt: now,
		UpdatedAt: now,
		Name:      params.Name,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error creating user: %v", err))
		return
	}

	respondWithJSON(w, 201, user)
}
