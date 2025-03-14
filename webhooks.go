package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/KOTBCAnorax/chirpy/internal/auth"
	"github.com/google/uuid"
)

type UpgradeRequest struct {
	Event string            `json:"event"`
	Data  map[string]string `json:"data"`
}

type UpgradeResponse struct {
	ID          uuid.UUID `json:"user_id"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (cfg *apiConfig) handleUpgrade(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	requestApiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		log.Printf("Error getting api key: %v\n", err)
		generateErrorResponse(w, 401)
		return
	}

	if requestApiKey != cfg.polkaKey {
		log.Printf("Api key mismatch\n")
		generateErrorResponse(w, 401)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := UpgradeRequest{}
	err = decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding upgrade request: %v\n", err)
		generateErrorResponse(w, 500)
		return
	}

	if params.Event != "user.upgraded" {
		log.Printf("Invalid event\n")
		w.WriteHeader(204)
		return
	}

	if _, ok := params.Data["user_id"]; !ok {
		log.Printf("No user id provided\n")
		w.WriteHeader(204)
		return
	}

	userID, err := uuid.Parse(params.Data["user_id"])
	if err != nil {
		log.Printf("Error parsing user id: %v\n", err)
		generateErrorResponse(w, 500)
		return
	}

	_, err = cfg.db.UpgradeToRed(r.Context(), userID)
	if err == sql.ErrNoRows {
		log.Printf("No user with given id\n")
		generateErrorResponse(w, 404)
		return
	}
	if err != nil {
		log.Printf("Database error: %v\n", err)
		generateErrorResponse(w, 500)
		return
	}

	w.WriteHeader(204)
}
