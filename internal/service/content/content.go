package contentsvc

import (
	"anchor-blog/internal/domain/entities"
	"context"
)

type ContentUsecase interface {
	GenerateContent(ctx context.Context, req entities.ContentRequest) (string, string, error)
}

type contentUsecase struct {
	repo ContentRepository
}

func NewContentUsecase(r ContentRepository) ContentUsecase {
	return &contentUsecase{repo: r}
}

func (uc *contentUsecase) GenerateContent(ctx context.Context, req entities.ContentRequest) (string, string, error) {

	title, content, err := uc.repo.Generate(ctx, req)
	if err != nil {
		return "", "", err
	}

	return title, content, nil
}

type ContentRepository interface {
	Generate(ctx context.Context, req entities.ContentRequest) (string, string, error)
}
