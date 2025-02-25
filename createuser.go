package main

import (
	"database/sql"
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
		w.WriteHeader(500)
		w.Write(generateErrorResponse())
		return
	}

	email := sql.NullString{
		String: params.Email,
		Valid:  true,
	}
	user, err := cfg.db.CreateUser(r.Context(), email)
	if err != nil {
		fmt.Printf("Database error: %v\n", err)
		w.WriteHeader(500)
		w.Write(generateErrorResponse("--->Couldn't create user\n"))
		return
	}

	responseBody := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
		Email:     user.Email.String,
	}

	data, err := json.Marshal(responseBody)
	if err != nil {
		log.Printf("Error encoding response: %s", err)
		w.WriteHeader(500)
		w.Write(generateErrorResponse())
		return
	}

	w.WriteHeader(201)
	w.Write(data)
}
