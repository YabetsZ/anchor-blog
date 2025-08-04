package usersvc

import (
	"anchor-blog/config"
	"anchor-blog/internal/domain/entities"
	"anchor-blog/internal/errors"
	"anchor-blog/pkg/hashutil"
	"anchor-blog/pkg/jwtutil"
	"context"
	"log"
	"time"
)

// This service will be merged into one file with others when other user-services are completed.
type UserServices struct {
	userRepo  entities.IUserRepository
	tokenRepo entities.ITokenRepository
	cfg       *config.Config
}

func NewUserServices(userRepo entities.IUserRepository, tokenRepo entities.ITokenRepository, cfg *config.Config) *UserServices {
	return &UserServices{userRepo, tokenRepo, cfg}
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (us *UserServices) Login(username, password string) (*LoginResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	user, err := us.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	err = hashutil.ComparePassword(user.PasswordHash, password)
	if err != nil {
		log.Printf("login failed for username '%s': invalid password", username)
		return nil, errors.ErrInvalidCredentials
	}

	accessToken, err := jwtutil.GenerateAccessToken(user, us.cfg.JWT.AccessTokenSecret)
	if err != nil {
		log.Printf("failed to produce access token: %v", err)
		return nil, err
	}
	refreshToken, err := jwtutil.GenerateRefreshToken(user, us.cfg.JWT.RefreshTokenSecret)
	if err != nil {
		log.Printf("failed to produce refresh token: %v", err)
		return nil, err
	}

	// Persist refresh token
	tokenHash := hashutil.HashToken(refreshToken, us.cfg.HMAC.Secret)

	err = us.tokenRepo.StoreRefreshToken(ctx, &entities.RefreshToken{
		// ID:        primitive.NewObjectID().Hex(),
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(jwtutil.RefreshTokenDuration),
	})
	if err != nil {
		log.Println("failed to store refresh token")
		return nil, err
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
