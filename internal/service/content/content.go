package contentsvc

import (
	"anchor-blog/internal/domain/entities"
	AppError "anchor-blog/internal/errors"
	"context"
)

type ContentUsecase interface {
	GenerateContent(ctx context.Context, req entities.ContentRequest) (*entities.ContentResponse, error)
}

type contentUsecase struct {
	repo ContentRepository
}

func NewContentUsecase(r ContentRepository) ContentUsecase {
	return &contentUsecase{repo: r}
}

func (uc *contentUsecase) GenerateContent(ctx context.Context, req entities.ContentRequest) (*entities.ContentResponse, error) {

	resp, err := uc.repo.Generate(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.SafetyReport.Blocked {
		return nil, AppError.ErrContentBlocked
	}

	return resp, nil
}

type ContentRepository interface {
	Generate(ctx context.Context, req entities.ContentRequest) (*entities.ContentResponse, error)
}
