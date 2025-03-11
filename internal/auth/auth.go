package auth

import (
	"fmt"
	"log"

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
