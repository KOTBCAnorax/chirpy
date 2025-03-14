package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/KOTBCAnorax/chirpy/internal/auth"
	"github.com/KOTBCAnorax/chirpy/internal/database"
)

type RefreshTokenResponse struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
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

	expiresIn := params.ExpiresInSeconds
	if expiresIn <= 0 || expiresIn > 3600 {
		expiresIn = 3600
	}

	token, err := auth.MakeJWT(user.ID, cfg.secret, time.Duration(expiresIn)*time.Second)
	if err != nil {
		log.Printf("Token generation failed: %v\n", err)
		generateErrorResponse(w, 500, "Server error")
		return
	}

	refreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		log.Printf("Refresh token generation failed: %v\n", err)
		generateErrorResponse(w, 500, "Server error")
		return
	}

	newRefreshTokenDb := database.GenerateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 60),
	}

	_, err = cfg.db.GenerateRefreshToken(r.Context(), newRefreshTokenDb)
	if err != nil {
		log.Printf("Failed to insert new refresh token into database: %v\n", err)
		generateErrorResponse(w, 500, "Server error")
		return
	}

	responseBody := User{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        token,
		RefreshToken: refreshToken,
		IsChirpyRed:  user.IsChirpyRed,
	}

	data, err := json.Marshal(responseBody)
	if err != nil {
		log.Printf("Error encoding response: %v\n", err)
		generateErrorResponse(w, 500)
		return
	}

	w.WriteHeader(200)
	w.Write(data)
}

func (cfg *apiConfig) handleRefresh(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	requestToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error with provided refresh token: %v\n", err)
		generateErrorResponse(w, 500)
		return
	}

	dbToken, err := cfg.db.FindRefreshToken(r.Context(), requestToken)
	if err == sql.ErrNoRows {
		log.Printf("No such refresh token in database: %v\n", err)
		generateErrorResponse(w, 401)
		return
	}
	if err != nil {
		log.Printf("Failed to extract refresh token from database: %v\n", err)
		generateErrorResponse(w, 500)
		return
	}
	if dbToken.ExpiresAt.Before(time.Now()) {
		log.Printf("Refresh token has expired\n")
		generateErrorResponse(w, 401)
		return
	}

	searchParams := database.FindUserByRefreshTokenParams{
		Token:     dbToken.Token,
		ExpiresAt: time.Now(),
	}
	user, err := cfg.db.FindUserByRefreshToken(r.Context(), searchParams)
	if err != nil {
		log.Printf("Error finding user by refresh token: %v\n", err)
		generateErrorResponse(w, 401)
		return
	}

	newAccessToken, err := auth.MakeJWT(user.ID, cfg.secret, time.Hour)
	if err != nil {
		log.Printf("Error finding user by refresh token: %v\n", err)
		generateErrorResponse(w, 500)
		return
	}

	responseBody := RefreshTokenResponse{
		Token: newAccessToken,
	}

	data, err := json.Marshal(responseBody)
	if err != nil {
		log.Printf("Error encoding refresh token response: %v\n", err)
		generateErrorResponse(w, 500)
		return
	}

	w.WriteHeader(200)
	w.Write(data)
}

func (cfg *apiConfig) handleRevoke(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	requestToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error with provided refresh token: %v\n", err)
		generateErrorResponse(w, 500)
		return
	}

	dbRefreshToken, err := cfg.db.FindRefreshToken(r.Context(), requestToken)
	if err == sql.ErrNoRows {
		log.Printf("Could not find refresh token in database: %v\n", err)
		generateErrorResponse(w, 401)
		return
	}
	if err != nil {
		log.Printf("Error finding refresh token: %v\n", err)
		generateErrorResponse(w, 500)
		return
	}
	if dbRefreshToken.ExpiresAt.Before(time.Now()) {
		log.Printf("Refresh token expired\n")
		generateErrorResponse(w, 401)
		return
	}

	_, err = cfg.db.RevokeRefreshToken(r.Context(), requestToken)
	if err != nil {
		log.Printf("Error revoking refresh token: %v\n", err)
		generateErrorResponse(w, 500)
		return
	}

	w.WriteHeader(204)
}
