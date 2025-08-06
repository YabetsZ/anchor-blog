package postrepo

import (
	"context"
	"log"
	"time"

	"anchor-blog/internal/domain/entities"
	AppError "anchor-blog/internal/errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoPostRepository struct {
	collection *mongo.Collection
}

// NewMongoPostRepository creates a new post repository with MongoDB implementation.
func NewMongoPostRepository(collection *mongo.Collection) entities.IPostRepository {
	return &mongoPostRepository{
		collection,
	}
}

func (r *mongoPostRepository) Create(ctx context.Context, dPost *entities.Post) (*entities.Post, error) {
	post, err := FromDomainPost(dPost)
	if err != nil {
		return nil, err
	}
	post.ID = primitive.NewObjectID()
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()
	post.Likes = []primitive.ObjectID{}
	post.Dislikes = []primitive.ObjectID{}
	post.ViewCount = 0

	_, err = r.collection.InsertOne(ctx, post)
	if err != nil {
		return nil, AppError.ErrInternalServer
	}

	return ToDomainPost(post), nil
}

func (r *mongoPostRepository) FindByID(ctx context.Context, id string) (*entities.Post, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("unable to convert id to object id", id)
		return nil, AppError.ErrInvalidPostID
	}
	var post entities.Post
	err = r.collection.FindOne(ctx, bson.M{"_id": objId}).Decode(&post)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *mongoPostRepository) FindAll(ctx context.Context, opts entities.PaginationOptions) ([]*entities.Post, error) {
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}}) // Sort by most recent
	findOptions.SetSkip((opts.Page - 1) * opts.Limit)
	findOptions.SetLimit(opts.Limit)

	cursor, err := r.collection.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var posts []*entities.Post
	if err := cursor.All(ctx, &posts); err != nil {
		return nil, err
	}

	return posts, nil
}
