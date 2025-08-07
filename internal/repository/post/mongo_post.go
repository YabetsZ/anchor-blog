package postrepo

import (
	"context"
	"errors"
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
		return AppError.ErrInternalServer
	}
	filter := bson.M{"id": objID}
	_, err = r.collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Printf("error while delete post %v", err.Error())
		return AppError.ErrInternalServer
	}
	return nil
}

func (r *mongoPostRepository) Creator(ctx context.Context, id string) (string, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return "", AppError.ErrInternalServer
	}
	filter := bson.M{"id": objID}
	var post entities.Post
	err = r.collection.FindOne(ctx, filter).Decode(&post)

	if err != nil {
		log.Printf("error while fetch post %v", err.Error())
		return "", AppError.ErrInternalServer
	}
	return post.AuthorID, nil
}
func (r *mongoPostRepository) UpdateByID(ctx context.Context, id string, post *entities.Post) error {
	postDoc, err := FromDomainPost(post)
	if err != nil {
		return AppError.ErrInternalServer
	}
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return AppError.ErrInternalServer
	}
	var foundPost Post

	filter := bson.M{"id": objID}
	err = r.collection.FindOne(ctx, filter).Decode(&foundPost)
	if err != nil {
		return AppError.ErrInternalServer
	}

	if postDoc.Content != "" {
		foundPost.Content = postDoc.Content
	}
	if postDoc.Title != "" {
		foundPost.Title = postDoc.Title
	}
	if len(postDoc.Tags) > 0 {
		foundPost.Tags = postDoc.Tags
	}
	_, err = r.collection.UpdateOne(ctx, filter, foundPost)

	if err != nil {
		return AppError.ErrInternalServer
	}
	return nil
}

func (r *mongoPostRepository) CountViews(ctx context.Context, id string) error {
	return nil
}
func (r *mongoPostRepository) LikePost(ctx context.Context, postID, userID string) (bool, error) {
	objID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return false, AppError.ErrInternalServer
	}

	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return false, AppError.ErrInternalServer
	}

	filter := bson.M{"id": objID}
	var post Post
	err = r.collection.FindOne(ctx, filter).Decode(&post)
	if err != nil {
		return false, AppError.ErrInternalServer
	}

	dislikes := post.Dislikes

	for index := 0; index < len(dislikes); index++ {
		if dislikes[index] == userObjID {
			dislikes = append(dislikes[:index], dislikes[index+1:]...)
			break
		}
	}

	likes := post.Likes
	found := false
	for index := 0; index < len(likes); index++ {
		if likes[index] == userObjID {
			found = true
		}
	}
	if !found {
		likes = append(likes, userObjID)
	} else {
		return false, errors.New("already liked")
	}

	post.Dislikes = dislikes
	post.Likes = likes

	_, err = r.collection.UpdateOne(ctx, filter, post)
	if err != nil {
		return false, AppError.ErrInternalServer
	}

	return true, nil
}

func (r *mongoPostRepository) DislikePost(ctx context.Context, postID, userID string) (bool, error) {

	objID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return false, AppError.ErrInternalServer
	}

	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return false, AppError.ErrInternalServer
	}

	filter := bson.M{"id": objID}
	var post Post
	err = r.collection.FindOne(ctx, filter).Decode(&post)

	if err != nil {
		return false, AppError.ErrInternalServer
	}

	likes := post.Likes

	for index := 0; index < len(likes); index++ {
		if likes[index] == userObjID {
			likes = append(likes[:index], likes[index+1:]...)
			break
		}
	}

	dislikes := post.Dislikes
	found := false
	for index := 0; index < len(dislikes); index++ {
		if dislikes[index] == userObjID {
			found = true
		}
	}
	if !found {
		dislikes = append(dislikes, userObjID)
	} else {
		return false, errors.New("already disliked")
	}

	post.Dislikes = dislikes
	post.Likes = likes

	_, err = r.collection.UpdateOne(ctx, filter, post)
	if err != nil {
		return false, AppError.ErrInternalServer
	}

	return true, nil
}
