package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (cfg *apiConfig) handleListChirps(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	chirpsList, err := cfg.db.ListChirps(r.Context())
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
