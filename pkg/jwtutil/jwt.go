package jwtutil

import (
	"anchor-blog/internal/domain/entities"
	"anchor-blog/internal/errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	RefreshTokenDuration = time.Hour * 24 * 7
	AccessTokenDuration  = time.Hour * 1
)

func GenerateAccessToken(user *entities.User, secret string) (string, error) {
	claims := entities.CustomClaims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Printf("ERROR: Failed to generate access JWT for user '%s': %v", user.Username, err)
		return "", errors.ErrInternalServer
	}
	return signedToken, nil
}

func GenerateRefreshToken(user *entities.User, secret string) (string, error) {
	claims := entities.CustomClaims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(RefreshTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Printf("ERROR: Failed to generate refresh JWT for user '%s': %v", user.Username, err)
		return "", errors.ErrInternalServer
	}
	return signedToken, nil
}

func ValidateToken(tokenString string, secret string) (*entities.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &entities.CustomClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(secret), nil
	})
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	claims, ok := token.Claims.(*entities.CustomClaims)
	if !ok || !token.Valid {
		return nil, errors.ErrInvalidToken
	}

	return claims, nil
}
