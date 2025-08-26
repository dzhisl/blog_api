package auth

import (
	"crypto/sha256"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func HashToken(input string) string {
	hash := sha256.Sum256([]byte(input)) // returns [32]byte
	return hex.EncodeToString(hash[:])   // convert to hex string
}

func CheckTokenHash(plain, hash string) bool {
	expectedHash := HashToken(plain)
	return expectedHash == hash
}
