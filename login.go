package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/KOTBCAnorax/chirpy/internal/auth"
)

func (cfg *apiConfig) HandleLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	decoder := json.NewDecoder(r.Body)
	params := ReqUser{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding request: %s\n", err)
		generateErrorResponse(w, 500)
		return
	}

	user, err := cfg.db.FindUserByEmail(r.Context(), params.Email)
	if err == sql.ErrNoRows {
		fmt.Printf("Database error: %v\n", err)
		generateErrorResponse(w, 401, "Incorrect email or password")
		return
	}
	if err != nil {
		fmt.Printf("Database error: %v\n", err)
		generateErrorResponse(w, 500, "Couldn't create user")
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		fmt.Printf("Database error: %v\n", err)
		generateErrorResponse(w, 401, "Incorrect email or password")
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

	w.WriteHeader(200)
	w.Write(data)
}
