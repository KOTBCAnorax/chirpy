package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/KOTBCAnorax/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handleGetChirp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	id, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		log.Printf("Error while parsing given ID: %s\n", err)
		generateErrorResponse(w, 500)
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), id)
	if err == sql.ErrNoRows {
		log.Printf("Chirp with the given ID doesn't exist: ErrNoRows")
		generateErrorResponse(w, 404)
		return
	}
	if err != nil {
		log.Printf("Error retrieving the chirp: %s\n", err)
		generateErrorResponse(w, 500)
		return
	}

	responseChirp := ChirpResponse{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	data, err := json.Marshal(responseChirp)
	if err != nil {
		log.Printf("Failed to marshal chirp: %v\n", err)
		generateErrorResponse(w, 500)
		return
	}
	w.WriteHeader(200)
	w.Write(data)
}

func (cfg *apiConfig) handleListChirps(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	userIDStr := r.URL.Query().Get("author_id")

	var chirpsList []database.Chirp
	var err error

	if userIDStr == "" {
		chirpsList, err = cfg.db.ListChirps(r.Context())
	} else {
		var userID uuid.UUID
		userID, err = uuid.Parse(userIDStr)
		if err == nil {
			chirpsList, err = cfg.db.ListUserChirps(r.Context(), userID)
		}
	}

	if err != nil {
		log.Printf("Error retrieving chirps: %s\n", err)
		generateErrorResponse(w, 500)
		return
	}

	responseList := make([]ChirpResponse, len(chirpsList))
	for i, chirp := range chirpsList {
		responseList[i] = ChirpResponse{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
	}

	data, err := json.Marshal(responseList)
	if err != nil {
		log.Printf("Failed to marshal chirp: %v\n", err)
		generateErrorResponse(w, 500)
		return
	}
	w.WriteHeader(200)
	w.Write(data)
}
