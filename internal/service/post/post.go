package postsvc

import (
	"context"

	"anchor-blog/internal/domain/entities"
)

type PostService struct {
	postRepo entities.IPostRepository
}

// NewPostService creates a new post service.
func NewPostService(repo entities.IPostRepository) *PostService {
	return &PostService{postRepo: repo}
}

func (s *PostService) CreatePost(ctx context.Context, title, content string, authorID string, tags []string) (*entities.Post, error) {
	post := &entities.Post{
		Title:    title,
		Content:  content,
		AuthorID: authorID,
		Tags:     tags,
	}

	return s.postRepo.Create(ctx, post)
}

func (s *PostService) GetPostByID(ctx context.Context, id string) (*entities.Post, error) {
	return s.postRepo.FindByID(ctx, id)
}

func (s *PostService) ListPosts(ctx context.Context, page, limit int64) ([]*entities.Post, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10 // Default limit
	}

	opts := entities.PaginationOptions{
		Page:  page,
		Limit: limit,
	}

	return s.postRepo.FindAll(ctx, opts)
}

// UpdatePost updates an existing post
func (s *PostService) UpdatePost(ctx context.Context, id string, title, content string, tags []string) (*entities.Post, error) {
	post := &entities.Post{
		Title:   title,
		Content: content,
		Tags:    tags,
	}

	return s.postRepo.Update(ctx, id, post)
}

// DeletePost deletes a post by ID
func (s *PostService) DeletePost(ctx context.Context, id string) error {
	return s.postRepo.Delete(ctx, id)
}

// SearchPosts searches posts by title or author
func (s *PostService) SearchPosts(ctx context.Context, query, searchType string, page, limit int64) ([]*entities.Post, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	opts := entities.PaginationOptions{
		Page:  page,
		Limit: limit,
	}

	switch searchType {
	case "title":
		return s.postRepo.SearchByTitle(ctx, query, opts)
	case "author":
		return s.postRepo.SearchByAuthor(ctx, query, opts)
	default:
		// Default to title search
		return s.postRepo.SearchByTitle(ctx, query, opts)
	}
}

// FilterPosts filters posts by various criteria
func (s *PostService) FilterPostsByTags(ctx context.Context, tags []string, page, limit int64) ([]*entities.Post, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	opts := entities.PaginationOptions{
		Page:  page,
		Limit: limit,
	}

	return s.postRepo.FilterByTags(ctx, tags, opts)
}

// FilterPostsByDateRange filters posts by date range
func (s *PostService) FilterPostsByDateRange(ctx context.Context, startDate, endDate string, page, limit int64) ([]*entities.Post, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	opts := entities.PaginationOptions{
		Page:  page,
		Limit: limit,
	}

	return s.postRepo.FilterByDateRange(ctx, startDate, endDate, opts)
}

// LikePost adds a like to a post
func (s *PostService) LikePost(ctx context.Context, postID, userID string) error {
	return s.postRepo.AddLike(ctx, postID, userID)
}

// UnlikePost removes a like from a post
func (s *PostService) UnlikePost(ctx context.Context, postID, userID string) error {
	return s.postRepo.RemoveLike(ctx, postID, userID)
}

// DislikePost adds a dislike to a post
func (s *PostService) DislikePost(ctx context.Context, postID, userID string) error {
	return s.postRepo.AddDislike(ctx, postID, userID)
}

// UndislikePost removes a dislike from a post
func (s *PostService) UndislikePost(ctx context.Context, postID, userID string) error {
	return s.postRepo.RemoveDislike(ctx, postID, userID)
}

// GetPostLikeStatus gets the like/dislike status for a user
func (s *PostService) GetPostLikeStatus(ctx context.Context, postID, userID string) (liked bool, disliked bool, err error) {
	return s.postRepo.GetLikeStatus(ctx, postID, userID)
}
