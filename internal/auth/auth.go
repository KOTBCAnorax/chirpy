package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	if password == "" {
		err := fmt.Errorf("empty passwords are not allowed")
		return "", err
	}

	pwdByte := []byte(password)
	hashByte, err := bcrypt.GenerateFromPassword(pwdByte, bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error: %v", err)
		return "", err
	}

	return string(hashByte), nil
}

func CheckPasswordHash(password, hash string) error {
	pwdByte := []byte(password)
	hashByte := []byte(hash)
	return bcrypt.CompareHashAndPassword(hashByte, pwdByte)
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(tokenSecret), nil
		})
	if err != nil {
		log.Printf("could not parse token: %v\n", err)
		return uuid.Nil, err
	}

	userIDStr, err := token.Claims.GetSubject()
	if err != nil {
		log.Printf("could not extract user id: %v\n", err)
		return uuid.Nil, err
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		log.Printf("could not parse user id: %v\n", err)
		return uuid.Nil, err
	}

	return userID, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if len(authHeader) == 0 {
		return "", fmt.Errorf("authorization header missing")
	}

	if !strings.Contains(authHeader, "Bearer ") {
		return "", fmt.Errorf("'Bearer' before token missing")
	}

	token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))

	if len(token) == 0 {
		return "", fmt.Errorf("error: empty token")
	}

	return token, nil
}

func GenerateRefreshToken() (string, error) {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		log.Printf("error generating random bytes for refresh token: %v\n", err)
		return "", err
	}

	hexString := hex.EncodeToString(randomBytes)
	return hexString, nil
}

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if len(authHeader) == 0 {
		return "", fmt.Errorf("authorization header missing")
	}

	if !strings.Contains(authHeader, "ApiKey ") {
		return "", fmt.Errorf("'ApiKey' before token missing")
	}

	apiKey := strings.TrimSpace(strings.TrimPrefix(authHeader, "ApiKey "))

	if len(apiKey) == 0 {
		return "", fmt.Errorf("error: no api key")
	}

	return apiKey, nil
}
