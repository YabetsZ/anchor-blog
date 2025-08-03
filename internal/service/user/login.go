package usersvc

import (
	"anchor-blog/internal/domain/entities"
	"anchor-blog/internal/errors"
	"anchor-blog/pkg/hashutil"
	"anchor-blog/pkg/jwtutil"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// This service will be merged into one file with others when other user-services are completed.
type userServices struct {
	userRepo           entities.IUserRepository
	tokenRepo          entities.ITokenRepository
	accessTokenSecret  string // from viper struct
	refreshTokenSecret string
}

func NewUserServices(userRepo entities.IUserRepository, tokenRepo entities.ITokenRepository, accessTokenSecret, refreshTokenSecret string) *userServices {
	return &userServices{userRepo, tokenRepo, accessTokenSecret, refreshTokenSecret}
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (us *userServices) Login(username, password string) (*LoginResponse, error) {
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

	accessToken, err := jwtutil.GenerateAccessToken(user, us.accessTokenSecret)
	if err != nil {
		log.Printf("failed to produce access token: %v", err)
		return nil, errors.ErrInternalServer
	}
	refreshToken, err := jwtutil.GenerateRefreshToken(user, us.refreshTokenSecret)
	if err != nil {
		log.Printf("failed to produce refresh token: %v", err)
		return nil, errors.ErrInternalServer
	}

	// Persist refresh token
	tokenHash, err := hashutil.HashPassword(refreshToken)
	if err != nil {
		log.Printf("failed to hash refresh token: %v", err)
		return nil, errors.ErrInternalServer
	}
	err = us.tokenRepo.StoreRefreshToken(ctx, &entities.RefreshToken{
		ID:        primitive.NewObjectID().Hex(),
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(jwtutil.RefreshTokenDuration),
	})
	if err != nil {
		log.Println("failed to store refresh token")
		return nil, errors.ErrInternalServer
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
