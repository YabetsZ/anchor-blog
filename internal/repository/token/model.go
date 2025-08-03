package tokenrepo

import (
	"anchor-blog/internal/domain/entities"
	"anchor-blog/internal/errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mongoRefreshToken struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	TokenHash string             `bson:"token_hash"`
	UserID    primitive.ObjectID `bson:"user_id"`
	ExpiresAt time.Time          `bson:"expires_at"`
	// IssuedAt  time.Time          `bson:"issued_at"`
	// Revoked   bool               `bson:"revoked"`
}

// :::::::: Mapping functions ::::::::

func FromDomainToken(token *entities.RefreshToken) (*mongoRefreshToken, error) {
	UserID, err := primitive.ObjectIDFromHex(token.UserID)
	if err != nil {
		log.Println("invalid user_id in token:", token.UserID)
		return nil, errors.ErrInvalidToken
	}
	ID, err := primitive.ObjectIDFromHex(token.ID)
	if err != nil {
		log.Println("invalid token id:", token.ID)
		return nil, errors.ErrInvalidToken
	}
	return &mongoRefreshToken{
		ID:        ID,
		TokenHash: token.TokenHash,
		UserID:    UserID,
		ExpiresAt: token.ExpiresAt,
	}, nil
}

func ToDomainToken(mToken *mongoRefreshToken) *entities.RefreshToken {
	return &entities.RefreshToken{
		ID:        mToken.ID.Hex(),
		TokenHash: mToken.TokenHash,
		UserID:    mToken.UserID.Hex(),
		ExpiresAt: mToken.ExpiresAt,
	}
}
