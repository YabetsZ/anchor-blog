package Models

import (
	"context"
	"time"
)

type ITokenRepository interface {
	StoreRefreshToken(ctx context.Context, tokenHash, userID string, expiresAt time.Time) error
}
