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

type ActivationTokenRepository struct {
	collection *mongo.Collection
}

func NewActivationTokenRepository(collection *mongo.Collection) *ActivationTokenRepository {
	ctx := context.Background()
	if err := ensureActivationTokenIndexes(ctx, collection); err != nil {
		log.Printf("failed to create indexes on activation tokens: %v", err)
	}
	return &ActivationTokenRepository{collection}
}

func ensureActivationTokenIndexes(ctx context.Context, col *mongo.Collection) error {
	_, err := col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "token", Value: 1}},
			Options: options.Index().
				SetName("idx_activation_token").
				SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "expires_at", Value: 1}},
			Options: options.Index().
				SetExpireAfterSeconds(0).
				SetName("idx_activation_token_expiry"),
		},
	})
	return err
}

func (r *ActivationTokenRepository) StoreActivationToken(ctx context.Context, token *entities.ActivationToken) error {
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
		log.Printf("error storing activation token: %v", err)
		return errors.ErrInternalServer
	}
	return nil
}

func (r *ActivationTokenRepository) FindActivationToken(ctx context.Context, token string) (*entities.ActivationToken, error) {
	var result bson.M
	err := r.collection.FindOne(ctx, bson.M{"token": token}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.ErrNotFound
		}
		log.Printf("error finding activation token: %v", err)
		return nil, errors.ErrInternalServer
	}

	activationToken := &entities.ActivationToken{
		ID:        result["_id"].(primitive.ObjectID).Hex(),
		UserID:    result["user_id"].(string),
		Token:     result["token"].(string),
		ExpiresAt: result["expires_at"].(primitive.DateTime).Time(),
		Used:      result["used"].(bool),
		CreatedAt: result["created_at"].(primitive.DateTime).Time(),
	}

	return activationToken, nil
}

func (r *ActivationTokenRepository) MarkTokenAsUsed(ctx context.Context, token string) error {
	filter := bson.M{"token": token}
	update := bson.M{"$set": bson.M{"used": true}}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("error marking activation token as used: %v", err)
		return errors.ErrInternalServer
	}
	return nil
}

func (r *ActivationTokenRepository) IsTokenValid(ctx context.Context, token string) (bool, error) {
	activationToken, err := r.FindActivationToken(ctx, token)
	if err != nil {
		return false, err
	}

	// Check if token is expired
	if time.Now().After(activationToken.ExpiresAt) {
		return false, fmt.Errorf("activation token has expired")
	}

	// Check if token is already used
	if activationToken.Used {
		return false, fmt.Errorf("activation token has already been used")
	}

	return true, nil
}