package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type UserEmail struct {
	Email string `json:"email"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handleUserCreation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	decoder := json.NewDecoder(r.Body)
	params := UserEmail{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error encoding response: %s", err)
		generateErrorResponse(w, 500)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), params.Email)
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
		log.Printf("Error encoding response: %s", err)
		generateErrorResponse(w, 500)
		return
	}

	w.WriteHeader(201)
	w.Write(data)
}
