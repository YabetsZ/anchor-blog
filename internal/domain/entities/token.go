package entities

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type RefreshToken struct {
	ID        string
	TokenHash string
	UserID    string
	ExpiresAt time.Time
	// IssuedAt  time.Time          `bson:"issued_at"`
	// Revoked   bool               `bson:"revoked"`
}

type CustomClaims struct {
	UserID   string
	Username string
	Role     string
	// may be activated?
	jwt.RegisteredClaims
}
