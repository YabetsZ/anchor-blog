package service

import (
	"context"
	"time"

	"anchor-blog/internal/domain/entities"
)

type PostService struct {
	postRepo entities.IPostRepository
}

// NewPostService creates a new post service.
func NewPostService(repo entities.IPostRepository) *PostService {
	return &PostService{postRepo: repo}
}

func (s *PostService) CreatePost(title, content string, authorID string, tags []string) (*entities.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	post := &entities.Post{
		Title:    title,
		Content:  content,
		AuthorID: authorID,
		Tags:     tags,
	}

	return s.postRepo.Create(ctx, post)
}

func (s *PostService) GetPostByID(id string) (*entities.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	return s.postRepo.FindByID(ctx, id)
}

func (s *PostService) ListPosts(page, limit int64) ([]*entities.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
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
