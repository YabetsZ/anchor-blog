package entities

import (
	"context"
)

// PaginationOptions holds the parameters for pagination.
type PaginationOptions struct {
	Page  int64
	Limit int64
}

// PostRepository defines the interface for post data operations.
type IPostRepository interface {
	Create(ctx context.Context, post *Post) (*Post, error)
	FindByID(ctx context.Context, id string) (*Post, error)
	FindAll(ctx context.Context, opts PaginationOptions) ([]*Post, error)
	
	// CRUD operations
	Update(ctx context.Context, id string, post *Post) (*Post, error)
	Delete(ctx context.Context, id string) error
	
	// Search and filter operations
	SearchByTitle(ctx context.Context, query string, opts PaginationOptions) ([]*Post, error)
	SearchByAuthor(ctx context.Context, authorID string, opts PaginationOptions) ([]*Post, error)
	FilterByTags(ctx context.Context, tags []string, opts PaginationOptions) ([]*Post, error)
	FilterByDateRange(ctx context.Context, startDate, endDate string, opts PaginationOptions) ([]*Post, error)
	
	// Like/Dislike operations
	AddLike(ctx context.Context, postID, userID string) error
	RemoveLike(ctx context.Context, postID, userID string) error
	AddDislike(ctx context.Context, postID, userID string) error
	RemoveDislike(ctx context.Context, postID, userID string) error
	GetLikeStatus(ctx context.Context, postID, userID string) (liked bool, disliked bool, error error)
	
	// View tracking methods
	IncrementViewCount(ctx context.Context, postID string) error
	GetViewCount(ctx context.Context, postID string) (int, error)
	GetTotalViews(ctx context.Context) (int64, error)
	GetPostsByViewCount(ctx context.Context, limit int) ([]*Post, error)
	ResetViewCount(ctx context.Context, postID string) error
}
