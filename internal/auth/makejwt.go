package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	now := jwt.NewNumericDate(time.Now().UTC())
	expiration := jwt.NewNumericDate(time.Now().UTC().Add(expiresIn))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Issuer:    "chirpy",
			IssuedAt:  now,
			ExpiresAt: expiration,
			Subject:   userID.String()})

	return token.SignedString([]byte(tokenSecret))
}
