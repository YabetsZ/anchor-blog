package jwtutil

import (
	Models "anchor-blog/internal/domain/models"
	"anchor-blog/internal/errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateAccessToken(user *Models.User, secret string) (string, error) {
	claims := Models.CustomClaims{
		UserID:   user.ID.Hex(),
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 2)),
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

func GenerateRefreshToken(user *Models.User, secret string) (string, error) {
	claims := Models.CustomClaims{
		UserID:   user.ID.Hex(),
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
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

func ValidateToken(tokenString string, secret string) (*Models.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Models.CustomClaims{}, func(token *jwt.Token) (any, error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Models.CustomClaims)
	if !ok || !token.Valid {
		return nil, errors.ErrInvalidToken
	}

	return claims, nil
}
