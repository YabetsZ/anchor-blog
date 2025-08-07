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
	
	// View tracking methods
	IncrementViewCount(ctx context.Context, postID string) error
	GetViewCount(ctx context.Context, postID string) (int, error)
	GetTotalViews(ctx context.Context) (int64, error)
	GetPostsByViewCount(ctx context.Context, limit int) ([]*Post, error)
	ResetViewCount(ctx context.Context, postID string) error
	
	// add Update and Delete
}
