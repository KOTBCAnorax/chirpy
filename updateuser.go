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

type UpdateRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateResponse struct {
	ID        uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (cfg *apiConfig) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

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

	decoder := json.NewDecoder(r.Body)
	params := UpdateRequest{}
	err = decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding user update request: %v\n", err)
		generateErrorResponse(w, 500)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("Error hashing new password: %v\n", err)
		generateErrorResponse(w, 500)
		return
	}

	userUpdate := database.UpdateUserParams{
		ID:             userID,
		Email:          params.Email,
		HashedPassword: hashedPassword,
	}

	updatedUser, err := cfg.db.UpdateUser(r.Context(), userUpdate)
	if err != nil {
		log.Printf("Database error: %v\n", err)
		generateErrorResponse(w, 500)
		return
	}

	updateResponseBody := UpdateResponse{
		ID:        updatedUser.ID,
		Email:     updatedUser.Email,
		CreatedAt: updatedUser.CreatedAt,
		UpdatedAt: updatedUser.UpdatedAt,
	}

	data, err := json.Marshal(updateResponseBody)
	if err != nil {
		log.Printf("Error encoding user update response: %v\n", err)
		generateErrorResponse(w, 500)
		return
	}

	w.WriteHeader(200)
	w.Write(data)
}
