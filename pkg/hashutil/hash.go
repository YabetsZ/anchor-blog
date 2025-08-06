package hashutil

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes the plain password using bcrypt.
func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// ComparePassword compares a hashed password with a plain password.
//
// Returns nil if they match.
func ComparePassword(hashedPassword, plainPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	if err != nil {
		return err
	}
	return nil
}

// HashToken hashes the token string using HMAC-SHA256
func HashToken(token string, hmacSecret string) string {
	mac := hmac.New(sha256.New, []byte(hmacSecret))
	mac.Write([]byte(token))
	hashed := mac.Sum(nil)
	return base64.URLEncoding.EncodeToString(hashed)
}

// Compare HMAC hashed tokens
func CompareTokens(storedHash, token string, hmacSecret string) bool {
	// Hash the token presented by the user
	hashOfToken := HashToken(token, hmacSecret)

	return hmac.Equal([]byte(storedHash), []byte(hashOfToken))
}
