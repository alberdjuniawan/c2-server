package utils

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword mengubah password biasa menjadi hash bcrypt
func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password:", err)
		return "", err
	}
	return string(hashed), nil
}

// VerifyPassword membandingkan password asli dengan hash dari database
func VerifyPassword(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}
