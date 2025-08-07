package postsvc

import (
	"context"
	"errors"

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

func (s *PostService) LikePost(ctx context.Context, postID, userID string) (bool, error) {
	return s.postRepo.LikePost(ctx, postID, userID)
}

func (s *PostService) DislikePost(ctx context.Context, postID, userID string) (bool, error) {
	return s.postRepo.DislikePost(ctx, postID, userID)
}

func (s *PostService) Delete(ctx context.Context, postID, userID string) error {
	authorID, err := s.postRepo.Creator(ctx, postID)
	if authorID != userID || err != nil {
		return errors.New("cannot edit the post")
	}
	return s.postRepo.DeleteByID(ctx, postID)
}
