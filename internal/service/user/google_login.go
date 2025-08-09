package usersvc

import (
	"anchor-blog/internal/domain/entities"
	AppError "anchor-blog/internal/errors"
	"anchor-blog/pkg/hashutil"
	"anchor-blog/pkg/jwtutil"
	"context"
	"encoding/json"
	"log"
	"time"
)

func (us *UserServices) HandleGoogleLogin(ctx context.Context, googleUserInfoData []byte) (*LoginResponse, error) {
	var userInfo GoogleUserInfo
	if err := json.Unmarshal(googleUserInfoData, &userInfo); err != nil {
		log.Printf("failed to parse user info. \n%v", err)
		return nil, AppError.ErrFailedToParse
	}

	// Check if a user with this email already exists.
	user, err := us.userRepo.GetUserByEmail(ctx, userInfo.Email)
	if err != nil {
		// If user not found, create a new one.
		user = &entities.User{
			Email:     userInfo.Email,
			FirstName: userInfo.FirstName,
			LastName:  userInfo.LastName,
			Username:  userInfo.Email,
			Role:      "user",
			Activated: true,
			Profile: entities.UserProfile{
				PictureURL: userInfo.ProfilePicture,
			},
			UpdatedAt: time.Now(),
			CreatedAt: time.Now(),
			LastSeen:  time.Now(),
		}
		id, err := us.userRepo.CreateUser(ctx, user)
		if err != nil {
			log.Println("failed to create use from google info.\n", err)
			return nil, err
		}
		user.ID = id
	}
	// TODO: If user exists, optionally update their details.like first and second name

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
		log.Println("failed to store refresh token: ", err.Error())
		return nil, err
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
