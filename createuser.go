package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/KOTBCAnorax/chirpy/internal/auth"
	"github.com/KOTBCAnorax/chirpy/internal/database"
	"github.com/google/uuid"
)

type ReqUser struct {
	Password         string `json:"password"`
	Email            string `json:"email"`
	ExpiresInSeconds int    `json:"expires_in_seconds"`
}

type User struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
}

func (cfg *apiConfig) handleUserCreation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	decoder := json.NewDecoder(r.Body)
	params := ReqUser{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding request: %s\n", err)
		generateErrorResponse(w, 500)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("Error hashing provided password: %s\n", err)
		generateErrorResponse(w, 500)
		return
	}

	newDbEntry := database.CreateUserParams{
		HashedPassword: hashedPassword,
		Email:          params.Email,
	}

	user, err := cfg.db.CreateUser(r.Context(), newDbEntry)
	if err != nil {
		fmt.Printf("Database error: %v\n", err)
		generateErrorResponse(w, 500, "Couldn't create user")
		return
	}

	responseBody := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	data, err := json.Marshal(responseBody)
	if err != nil {
		log.Printf("Error encoding response: %s\n", err)
		generateErrorResponse(w, 500)
		return
	}

	w.WriteHeader(201)
	w.Write(data)
}
