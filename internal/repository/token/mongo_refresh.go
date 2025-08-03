package tokenrepo

import (
	"anchor-blog/internal/domain/entities"
	"anchor-blog/internal/errors"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoTokenRepository struct {
	collection *mongo.Collection
}

func NewMongoTokenRepository(collection *mongo.Collection) entities.ITokenRepository {
	ctx := context.Background()
	if err := ensureIndexes(ctx, collection); err != nil {
		log.Fatalf("failed to create indeces on token_hash and expires_at: %v", err)
	}
	return &mongoTokenRepository{collection}
}

// creates indexes for refresh tokens incase the don't exist
func ensureIndexes(ctx context.Context, col *mongo.Collection) error {
	_, err := col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "token_hash", Value: 1}},
			Options: options.Index().
				SetName("idx_token_hash").
				SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "expires_at", Value: 1}},
			Options: options.Index().
				SetExpireAfterSeconds(0).
				SetName("idx_token_expiry"),
		},
	})
	return err
}

func (mt *mongoTokenRepository) StoreRefreshToken(ctx context.Context, token *entities.RefreshToken) error {
	mToken, err := FromDomainToken(token)
	if err != nil {
		return err
	}

	_, err = mt.collection.InsertOne(ctx, mToken)
	if err != nil {
		log.Println("unexpected error during Insertion: ", err)
		return errors.ErrInternalServer
	}

	return nil
}

func (mt *mongoTokenRepository) FindByHash(ctx context.Context, hash string) (*entities.RefreshToken, error) {
	token := mongoRefreshToken{}
	err := mt.collection.FindOne(ctx, bson.M{"token_hash": hash}).Decode(&token)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.ErrNotFound
		}
		log.Println("error while finding hash: ", err)
		return nil, errors.ErrInternalServer
	}
	return ToDomainToken(&token), nil
}

func (mt *mongoTokenRepository) DeleteByHash(ctx context.Context, tokenHash string) error {
	_, err := mt.collection.DeleteOne(ctx, bson.M{"token_hash": tokenHash})
	if err != nil {
		log.Printf("failed to delete token with hash %s: %v", tokenHash, err)
		return errors.ErrInternalServer
	}
	return nil
}

func (mt *mongoTokenRepository) DeleteAllByUserID(ctx context.Context, userID string) error {
	ID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println("invalid user id:", userID)
		return errors.ErrInvalidUserID
	}
	_, err = mt.collection.DeleteMany(ctx, bson.M{"user_id": ID})
	if err != nil {
		log.Printf("failed to delete token with user id %s: %v", userID, err)
		return errors.ErrInternalServer
	}
	return nil
}
