package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/KOTBCAnorax/chirpy/internal/auth"
	"github.com/KOTBCAnorax/chirpy/internal/database"
	"github.com/google/uuid"
)

type ChirpRequest struct {
	Body   string `json:"body"`
	UserID string `json:"user_id"`
}

type ChirpResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handleChirp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	decoder := json.NewDecoder(r.Body)
	reqParams := ChirpRequest{}
	err := decoder.Decode(&reqParams)
	if err != nil {
		log.Printf("Error decoding parameters: %v\n", err)
		generateErrorResponse(w, 500)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error getting bearer token: %v\n", err)
		generateErrorResponse(w, 401, "Unauthorized")
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		log.Printf("Invalid token: %v\n", err)
		generateErrorResponse(w, 401, "Unauthorized")
		return
	}

	if len(reqParams.Body) > 140 {
		msg := "Chirp is too long"
		log.Println(msg)
		generateErrorResponse(w, 400, msg)
		return
	}

	reqParams.Body = FilterProfane(reqParams.Body)
	parsedParams := database.CreateChirpParams{
		Body:   reqParams.Body,
		UserID: userID,
	}

	chirp, err := cfg.db.CreateChirp(r.Context(), parsedParams)
	if err != nil {
		log.Printf("Database error: %v\n", err)
		generateErrorResponse(w, 500)
		return
	}

	chirpResponse := ChirpResponse{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	data, err := json.Marshal(chirpResponse)
	if err != nil {
		log.Printf("Failed to marshal chirp: %v\n", err)
		generateErrorResponse(w, 500)
		return
	}

	w.WriteHeader(201)
	w.Write(data)
}
