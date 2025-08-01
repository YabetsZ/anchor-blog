package entities

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RefreshToken struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	TokenHash string             `bson:"token_hash"`
	UserID    string             `bson:"user_id"`
	IssuedAt  time.Time          `bson:"issued_at"`
	ExpiresAt time.Time          `bson:"expires_at"`
	Revoked   bool               `bson:"revoked"`
}

type CustomClaims struct {
	UserID   string
	Username string
	Role     string
	// may be activated?
	jwt.RegisteredClaims
}
