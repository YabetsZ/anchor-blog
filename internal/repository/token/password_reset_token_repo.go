package tokenrepo

import (
	"anchor-blog/internal/domain/entities"
	"anchor-blog/internal/errors"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PasswordResetTokenRepository struct {
	collection *mongo.Collection
}

func NewPasswordResetTokenRepository(collection *mongo.Collection) *PasswordResetTokenRepository {
	ctx := context.Background()
	if err := ensurePasswordResetTokenIndexes(ctx, collection); err != nil {
		log.Printf("failed to create indexes on password reset tokens: %v", err)
	}
	return &PasswordResetTokenRepository{collection}
}

func ensurePasswordResetTokenIndexes(ctx context.Context, col *mongo.Collection) error {
	_, err := col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "token", Value: 1}},
			Options: options.Index().
				SetName("idx_password_reset_token").
				SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "expires_at", Value: 1}},
			Options: options.Index().
				SetExpireAfterSeconds(0).
				SetName("idx_password_reset_token_expiry"),
		},
	})
	return err
}

func (r *PasswordResetTokenRepository) StorePasswordResetToken(ctx context.Context, token *entities.PasswordResetToken) error {
	doc := bson.M{
		"_id":        primitive.NewObjectID(),
		"user_id":    token.UserID,
		"token":      token.Token,
		"expires_at": token.ExpiresAt,
		"used":       token.Used,
		"created_at": token.CreatedAt,
	}

	_, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		log.Printf("error storing password reset token: %v", err)
		return errors.ErrInternalServer
	}
	return nil
}

func (r *PasswordResetTokenRepository) FindPasswordResetToken(ctx context.Context, token string) (*entities.PasswordResetToken, error) {
	var result bson.M
	err := r.collection.FindOne(ctx, bson.M{"token": token}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.ErrNotFound
		}
		log.Printf("error finding password reset token: %v", err)
		return nil, errors.ErrInternalServer
	}

	resetToken := &entities.PasswordResetToken{
		ID:        result["_id"].(primitive.ObjectID).Hex(),
		UserID:    result["user_id"].(string),
		Token:     result["token"].(string),
		ExpiresAt: result["expires_at"].(primitive.DateTime).Time(),
		Used:      result["used"].(bool),
		CreatedAt: result["created_at"].(primitive.DateTime).Time(),
	}

	return resetToken, nil
}

func (r *PasswordResetTokenRepository) MarkTokenAsUsed(ctx context.Context, token string) error {
	filter := bson.M{"token": token}
	update := bson.M{"$set": bson.M{"used": true}}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("error marking password reset token as used: %v", err)
		return errors.ErrInternalServer
	}
	return nil
}

func (r *PasswordResetTokenRepository) IsTokenValid(ctx context.Context, token string) (bool, error) {
	resetToken, err := r.FindPasswordResetToken(ctx, token)
	if err != nil {
		return false, err
	}

	// Check if token is expired
	if time.Now().After(resetToken.ExpiresAt) {
		return false, fmt.Errorf("password reset token has expired")
	}

	// Check if token is already used
	if resetToken.Used {
		return false, fmt.Errorf("password reset token has already been used")
	}

	return true, nil
}