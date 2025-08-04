package usersvc

import (
	"anchor-blog/internal/domain/entities"
	AppError "anchor-blog/internal/errors"
	"anchor-blog/pkg/hashutil"
	"anchor-blog/pkg/jwtutil"
	"context"
	"errors"
	"time"
)

func (us *UserServices) Refresh(refreshToken string) (*LoginResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	claim, err := jwtutil.ValidateToken(refreshToken, us.cfg.JWT.RefreshTokenSecret)
	if err != nil {
		return nil, err
	}
	if claim.ExpiresAt.Time.Before(time.Now()) {
		return nil, AppError.ErrInvalidToken
	}

	tokenHash := hashutil.HashToken(refreshToken, us.cfg.HMAC.Secret)
	_, err = us.tokenRepo.FindByHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, AppError.ErrNotFound) {
			return nil, AppError.ErrInvalidToken
		}
		return nil, err
	}

	err = us.tokenRepo.DeleteByHash(ctx, tokenHash)
	if err != nil {
		return nil, err
	}

	user := &entities.User{
		ID:       claim.UserID,
		Username: claim.Username,
		Role:     claim.Username,
	}
	newAccessToken, err := jwtutil.GenerateAccessToken(user, us.cfg.JWT.AccessTokenSecret)
	if err != nil {
		return nil, err
	}
	newRefreshToken, err := jwtutil.GenerateRefreshToken(user, us.cfg.JWT.RefreshTokenSecret)
	if err != nil {
		return nil, err
	}

	err = us.tokenRepo.StoreRefreshToken(ctx, &entities.RefreshToken{
		UserID:    claim.ID,
		TokenHash: hashutil.HashToken(newRefreshToken, us.cfg.HMAC.Secret),
		ExpiresAt: time.Now().Add(jwtutil.RefreshTokenDuration),
	})

	return &LoginResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
