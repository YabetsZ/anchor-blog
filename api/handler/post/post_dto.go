package post

import (
	"anchor-blog/internal/domain/entities"
	"time"
)

type PostDTO struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	AuthorID  string    `json:"author_id"`
	Tags      []string  `json:"tags"`
	ViewCount int       `json:"view_count"`
	Likes     []string  `json:"likes"`
	Dislikes  []string  `json:"dislikes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func MapPostToDTO(post *entities.Post) *PostDTO {
	return &PostDTO{
		ID:        post.ID,
		Title:     post.Title,
		Content:   post.Content,
		AuthorID:  post.AuthorID,
		Tags:      post.Tags,
		ViewCount: post.ViewCount,
		Likes:     post.Likes,
		Dislikes:  post.Dislikes,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}
}

func MapDTOToPost(dto *PostDTO) *entities.Post {
	return &entities.Post{
		ID:        dto.ID,
		Title:     dto.Title,
		Content:   dto.Content,
		AuthorID:  dto.AuthorID,
		Tags:      dto.Tags,
		ViewCount: dto.ViewCount,
		Likes:     dto.Likes,
		Dislikes:  dto.Dislikes,
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
	}
}
