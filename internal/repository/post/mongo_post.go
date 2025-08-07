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

// IncrementViewCount increments the view count for a specific post
func (r *mongoPostRepository) IncrementViewCount(ctx context.Context, postID string) error {
	objId, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		log.Println("unable to convert id to object id", postID)
		return AppError.ErrInvalidPostID
	}

	filter := bson.M{"_id": objId}
	update := bson.M{"$inc": bson.M{"view_count": 1}}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Error incrementing view count for post %s: %v", postID, err)
		return AppError.ErrInternalServer
	}

	if result.MatchedCount == 0 {
		return AppError.ErrNotFound
	}

	return nil
}

// GetViewCount retrieves the current view count for a specific post
func (r *mongoPostRepository) GetViewCount(ctx context.Context, postID string) (int, error) {
	objId, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		log.Println("unable to convert id to object id", postID)
		return 0, AppError.ErrInvalidPostID
	}

	var post Post
	filter := bson.M{"_id": objId}
	projection := bson.M{"view_count": 1}

	err = r.collection.FindOne(ctx, filter, options.FindOne().SetProjection(projection)).Decode(&post)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 0, AppError.ErrNotFound
		}
		return 0, AppError.ErrInternalServer
	}

	return post.ViewCount, nil
}

// GetTotalViews calculates the total view count across all posts
func (r *mongoPostRepository) GetTotalViews(ctx context.Context) (int64, error) {
	pipeline := []bson.M{
		{
			"$group": bson.M{
				"_id":         nil,
				"total_views": bson.M{"$sum": "$view_count"},
			},
		},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, AppError.ErrInternalServer
	}
	defer cursor.Close(ctx)

	var result []bson.M
	if err := cursor.All(ctx, &result); err != nil {
		return 0, AppError.ErrInternalServer
	}

	if len(result) == 0 {
		return 0, nil
	}

	totalViews, ok := result[0]["total_views"].(int32)
	if !ok {
		return 0, nil
	}

	return int64(totalViews), nil
}

// GetPostsByViewCount retrieves posts ordered by view count (most viewed first)
func (r *mongoPostRepository) GetPostsByViewCount(ctx context.Context, limit int) ([]*entities.Post, error) {
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "view_count", Value: -1}}) // Sort by most viewed
	findOptions.SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return nil, AppError.ErrInternalServer
	}
	defer cursor.Close(ctx)

	var posts []Post
	if err := cursor.All(ctx, &posts); err != nil {
		return nil, AppError.ErrInternalServer
	}

	result := make([]*entities.Post, len(posts))
	for idx, post := range posts {
		result[idx] = ToDomainPost(&post)
	}
	return result, nil
}

// ResetViewCount resets the view count for a specific post to zero
func (r *mongoPostRepository) ResetViewCount(ctx context.Context, postID string) error {
	objId, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		log.Println("unable to convert id to object id", postID)
		return AppError.ErrInvalidPostID
	}

	filter := bson.M{"_id": objId}
	update := bson.M{"$set": bson.M{"view_count": 0}}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Error resetting view count for post %s: %v", postID, err)
		return AppError.ErrInternalServer
	}

	if result.MatchedCount == 0 {
		return AppError.ErrNotFound
	}

	return nil
}
// Update updates an existing post
func (r *mongoPostRepository) Update(ctx context.Context, id string, post *entities.Post) (*entities.Post, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("unable to convert id to object id", id)
		return nil, AppError.ErrInvalidPostID
	}

	// Convert entity to model
	postModel, err := FromDomainPost(post)
	if err != nil {
		return nil, AppError.ErrInternalServer
	}

	postModel.UpdatedAt = time.Now()
	
	filter := bson.M{"_id": objId}
	update := bson.M{"$set": bson.M{
		"title":      postModel.Title,
		"content":    postModel.Content,
		"tags":       postModel.Tags,
		"updated_at": postModel.UpdatedAt,
	}}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Error updating post %s: %v", id, err)
		return nil, AppError.ErrInternalServer
	}

	if result.MatchedCount == 0 {
		return nil, AppError.ErrNotFound
	}

	// Return updated post
	return r.FindByID(ctx, id)
}

// Delete removes a post by ID
func (r *mongoPostRepository) Delete(ctx context.Context, id string) error {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("unable to convert id to object id", id)
		return AppError.ErrInvalidPostID
	}

	filter := bson.M{"_id": objId}
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Printf("Error deleting post %s: %v", id, err)
		return AppError.ErrInternalServer
	}

	if result.DeletedCount == 0 {
		return AppError.ErrNotFound
	}

	return nil
}

// SearchByTitle searches posts by title
func (r *mongoPostRepository) SearchByTitle(ctx context.Context, query string, opts entities.PaginationOptions) ([]*entities.Post, error) {
	filter := bson.M{
		"title": bson.M{
			"$regex":   query,
			"$options": "i", // case insensitive
		},
	}

	return r.findWithFilter(ctx, filter, opts)
}

// SearchByAuthor searches posts by author ID
func (r *mongoPostRepository) SearchByAuthor(ctx context.Context, authorID string, opts entities.PaginationOptions) ([]*entities.Post, error) {
	authorObjID, err := primitive.ObjectIDFromHex(authorID)
	if err != nil {
		return nil, AppError.ErrInvalidUserID
	}

	filter := bson.M{"author_id": authorObjID}
	return r.findWithFilter(ctx, filter, opts)
}

// FilterByTags filters posts by tags
func (r *mongoPostRepository) FilterByTags(ctx context.Context, tags []string, opts entities.PaginationOptions) ([]*entities.Post, error) {
	filter := bson.M{
		"tags": bson.M{
			"$in": tags,
		},
	}

	return r.findWithFilter(ctx, filter, opts)
}

// FilterByDateRange filters posts by date range
func (r *mongoPostRepository) FilterByDateRange(ctx context.Context, startDate, endDate string, opts entities.PaginationOptions) ([]*entities.Post, error) {
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, AppError.ErrValidationFailed
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, AppError.ErrValidationFailed
	}

	// Set end date to end of day
	end = end.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	filter := bson.M{
		"created_at": bson.M{
			"$gte": start,
			"$lte": end,
		},
	}

	return r.findWithFilter(ctx, filter, opts)
}

// AddLike adds a like to a post
func (r *mongoPostRepository) AddLike(ctx context.Context, postID, userID string) error {
	postObjID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return AppError.ErrInvalidPostID
	}

	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return AppError.ErrInvalidUserID
	}

	filter := bson.M{"_id": postObjID}
	update := bson.M{
		"$addToSet": bson.M{"likes": userObjID},
		"$pull":     bson.M{"dislikes": userObjID}, // Remove from dislikes if exists
	}

	_, err = r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Error adding like to post %s: %v", postID, err)
		return AppError.ErrInternalServer
	}

	return nil
}

// RemoveLike removes a like from a post
func (r *mongoPostRepository) RemoveLike(ctx context.Context, postID, userID string) error {
	postObjID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return AppError.ErrInvalidPostID
	}

	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return AppError.ErrInvalidUserID
	}

	filter := bson.M{"_id": postObjID}
	update := bson.M{
		"$pull": bson.M{"likes": userObjID},
	}

	_, err = r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Error removing like from post %s: %v", postID, err)
		return AppError.ErrInternalServer
	}

	return nil
}

// AddDislike adds a dislike to a post
func (r *mongoPostRepository) AddDislike(ctx context.Context, postID, userID string) error {
	postObjID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return AppError.ErrInvalidPostID
	}

	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return AppError.ErrInvalidUserID
	}

	filter := bson.M{"_id": postObjID}
	update := bson.M{
		"$addToSet": bson.M{"dislikes": userObjID},
		"$pull":     bson.M{"likes": userObjID}, // Remove from likes if exists
	}

	_, err = r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Error adding dislike to post %s: %v", postID, err)
		return AppError.ErrInternalServer
	}

	return nil
}

// RemoveDislike removes a dislike from a post
func (r *mongoPostRepository) RemoveDislike(ctx context.Context, postID, userID string) error {
	postObjID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return AppError.ErrInvalidPostID
	}

	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return AppError.ErrInvalidUserID
	}

	filter := bson.M{"_id": postObjID}
	update := bson.M{
		"$pull": bson.M{"dislikes": userObjID},
	}

	_, err = r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Error removing dislike from post %s: %v", postID, err)
		return AppError.ErrInternalServer
	}

	return nil
}

// GetLikeStatus checks if user has liked or disliked a post
func (r *mongoPostRepository) GetLikeStatus(ctx context.Context, postID, userID string) (bool, bool, error) {
	postObjID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return false, false, AppError.ErrInvalidPostID
	}

	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return false, false, AppError.ErrInvalidUserID
	}

	var post Post
	filter := bson.M{"_id": postObjID}
	projection := bson.M{"likes": 1, "dislikes": 1}

	err = r.collection.FindOne(ctx, filter, options.FindOne().SetProjection(projection)).Decode(&post)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, false, AppError.ErrNotFound
		}
		return false, false, AppError.ErrInternalServer
	}

	liked := false
	disliked := false

	for _, likeID := range post.Likes {
		if likeID == userObjID {
			liked = true
			break
		}
	}

	for _, dislikeID := range post.Dislikes {
		if dislikeID == userObjID {
			disliked = true
			break
		}
	}

	return liked, disliked, nil
}

// findWithFilter is a helper method for search and filter operations
func (r *mongoPostRepository) findWithFilter(ctx context.Context, filter bson.M, opts entities.PaginationOptions) ([]*entities.Post, error) {
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}}) // Sort by most recent
	findOptions.SetSkip((opts.Page - 1) * opts.Limit)
	findOptions.SetLimit(opts.Limit)

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, AppError.ErrInternalServer
	}
	defer cursor.Close(ctx)

	var posts []Post
	if err := cursor.All(ctx, &posts); err != nil {
		return nil, AppError.ErrInternalServer
	}

	result := make([]*entities.Post, len(posts))
	for idx, post := range posts {
		result[idx] = ToDomainPost(&post)
	}
	return result, nil
}