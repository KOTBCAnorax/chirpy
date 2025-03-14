package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/KOTBCAnorax/chirpy/internal/auth"
	"github.com/KOTBCAnorax/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handleChirpDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-length", "0")

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Token invalid or missing: %v\n", err)
		generateErrorResponse(w, 401)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		log.Printf("Error validating token: %v\n", err)
		generateErrorResponse(w, 401)
		return
	}

	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		log.Printf("Error while parsing given ID: %v\n", err)
		generateErrorResponse(w, 500)
		return
	}

	deleteParams := database.DeleteChirpParams{
		ID:     chirpID,
		UserID: userID,
	}

	isOwner, err := cfg.db.IsOwner(r.Context(), database.IsOwnerParams(deleteParams))
	if err != nil {
		log.Printf("Database error: %v\n", err)
		generateErrorResponse(w, 500)
		return
	}
	if !isOwner {
		log.Printf("Foreign chirp: %v\n", err)
		generateErrorResponse(w, 403)
		return
	}

	err = cfg.db.DeleteChirp(r.Context(), deleteParams)
	if err == sql.ErrNoRows {
		log.Printf("Chirp not found: %v\n", err)
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
