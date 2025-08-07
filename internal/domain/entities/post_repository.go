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
	// add Update and Delete

	DeleteByID(ctx context.Context, id string) error
	Creator(ctx context.Context, id string) (string, error)
	UpdateByID(ctx context.Context, id string, post *Post) error

	// track popularity
	CountViews(ctx context.Context, id string) error
	LikePost(ctx context.Context, postID, userID string) (bool, error)
	DislikePost(ctx context.Context, postID, userID string) (bool, error)
}
