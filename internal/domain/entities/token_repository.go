package entities

import (
	"context"
)

type ITokenRepository interface {
	StoreRefreshToken(ctx context.Context, token *RefreshToken) error
	FindByHash(ctx context.Context, hash string) (*RefreshToken, error)
	DeleteByHash(ctx context.Context, hash string) error
	DeleteAllByUserID(ctx context.Context, userID string) error
}
