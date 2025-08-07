package postrepo

import (
	"context"
	"errors"
	"fmt"
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
	return &mongoPostRepository{collection}
}

func (r *mongoPostRepository) Create(ctx context.Context, dPost *entities.Post) (*entities.Post, error) {
	dPost.ID = primitive.NewObjectID().Hex()
	post, err := FromDomainPost(dPost)
	if err != nil {
		return nil, err
	}
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
	var post Post
	err = r.collection.FindOne(ctx, bson.M{"_id": objId}).Decode(&post)
	if err != nil {
		return nil, err
	}
	return ToDomainPost(&post), nil
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

	var posts []Post
	if err := cursor.All(ctx, &posts); err != nil {
		return nil, err
	}

	result := make([]*entities.Post, len(posts))
	for idx, post := range posts {
		result[idx] = ToDomainPost(&post)
	}
	return result, nil
}

func (r *mongoPostRepository) DeleteByID(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("invalid post ID format: %s, error: %v", id, err)
		return fmt.Errorf("%w: invalid post ID", AppError.ErrBadRequest)
	}

	filter := bson.M{"_id": objID}
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Printf("failed to delete post %s: %v", id, err)
		return fmt.Errorf("%w: failed to delete post", AppError.ErrInternalServer)
	}

	if result.DeletedCount == 0 {
		log.Printf("post not found for deletion: %s", id)
		return fmt.Errorf("%w: post not found", AppError.ErrNotFound)
	}

	return nil
}

func (r *mongoPostRepository) Creator(ctx context.Context, id string) (string, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("invalid post ID format: %s, error: %v", id, err)
		return "", fmt.Errorf("%w: invalid post ID", AppError.ErrBadRequest)
	}

	// Fixed: Changed "id" to "_id" in the filter
	filter := bson.M{"_id": objID}
	var post entities.Post
	err = r.collection.FindOne(ctx, filter).Decode(&post)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("post not found: %s", id)
			return "", fmt.Errorf("%w: post not found", AppError.ErrNotFound)
		}
		log.Printf("failed to fetch post %s: %v", id, err)
		return "", fmt.Errorf("%w: failed to fetch post", AppError.ErrInternalServer)
	}

	if post.AuthorID == "" {
		log.Printf("post has no author ID: %s", id)
		return "", fmt.Errorf("%w: post author missing", AppError.ErrInternalServer)
	}

	return post.AuthorID, nil
}

func (r *mongoPostRepository) UpdateByID(ctx context.Context, id string, post *entities.Post) error {
	if post == nil {
		log.Println("nil post provided for update")
		return fmt.Errorf("%w: nil post", AppError.ErrBadRequest)
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("invalid post ID format: %s, error: %v", id, err)
		return fmt.Errorf("%w: invalid post ID", AppError.ErrBadRequest)
	}

	postDoc, err := FromDomainPost(post)
	if err != nil {
		log.Printf("failed to convert post to document: %v", err)
		return fmt.Errorf("%w: invalid post data", AppError.ErrBadRequest)
	}

	// Build update document dynamically based on what fields are provided
	update := bson.M{}
	if postDoc.Content != "" {
		update["content"] = postDoc.Content
	}
	if postDoc.Title != "" {
		update["title"] = postDoc.Title
	}
	if len(postDoc.Tags) > 0 {
		update["tags"] = postDoc.Tags
	}

	if len(update) == 0 {
		log.Println("no valid fields provided for update")
		return fmt.Errorf("%w: no fields to update", AppError.ErrBadRequest)
	}

	filter := bson.M{"_id": objID}
	result, err := r.collection.UpdateOne(
		ctx,
		filter,
		bson.M{"$set": update},
	)

	if err != nil {
		log.Printf("failed to update post %s: %v", id, err)
		return fmt.Errorf("%w: failed to update post", AppError.ErrInternalServer)
	}

	if result.MatchedCount == 0 {
		log.Printf("post not found for update: %s", id)
		return fmt.Errorf("%w: post not found", AppError.ErrNotFound)
	}

	return nil
}

func (r *mongoPostRepository) CountViews(ctx context.Context, id string) error {
	return nil
}
func (r *mongoPostRepository) LikePost(ctx context.Context, postID, userID string) (bool, error) {
	objID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return false, fmt.Errorf("invalid post ID: %w", err)
	}

	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return false, fmt.Errorf("invalid user ID: %w", err)
	}

	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{
			"_id":   objID,
			"likes": bson.M{"$ne": userObjID},
		},
		bson.M{
			"$addToSet": bson.M{"likes": userObjID},
			"$pull":     bson.M{"dislikes": userObjID},
		},
	)

	if err != nil {
		return false, fmt.Errorf("database error: %w", err)
	}

	if result.ModifiedCount == 0 {
		var post Post
		err := r.collection.FindOne(ctx, bson.M{"_id": objID, "likes": userObjID}).Decode(&post)
		if err == nil {
			return false, errors.New("already liked")
		}
		return false, errors.New("post not found or no changes made")
	}

	return true, nil
}

func (r *mongoPostRepository) DislikePost(ctx context.Context, postID, userID string) (bool, error) {
	objID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return false, fmt.Errorf("invalid post ID: %w", err)
	}

	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return false, fmt.Errorf("invalid user ID: %w", err)
	}

	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{
			"_id":      objID,
			"dislikes": bson.M{"$ne": userObjID},
		},
		bson.M{
			"$addToSet": bson.M{"dislikes": userObjID},
			"$pull":     bson.M{"likes": userObjID},
		},
	)

	if err != nil {
		return false, fmt.Errorf("database error: %w", err)
	}

	if result.ModifiedCount == 0 {
		var post Post
		err := r.collection.FindOne(ctx, bson.M{"_id": objID, "dislikes": userObjID}).Decode(&post)
		if err == nil {
			return false, errors.New("already disliked")
		}
		return false, errors.New("post not found or no changes made")
	}

	return true, nil
}
