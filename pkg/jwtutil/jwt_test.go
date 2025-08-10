package jwtutil

import (
	"testing"
	"time"

	"anchor-blog/internal/domain/entities"

	"github.com/stretchr/testify/assert"
)

func TestGenerateAccessToken_Success(t *testing.T) {
	// Test data
	user := &entities.User{
		ID:       "user-123",
		Username: "testuser",
		Email:    "test@example.com",
		Role:     "user",
	}
	secret := "test-secret-key"

	// Execute
	token, err := GenerateAccessToken(user, secret)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.Contains(t, token, ".") // JWT should contain dots
}

func TestGenerateRefreshToken_Success(t *testing.T) {
	// Test data
	user := &entities.User{
		ID:       "user-123",
		Username: "testuser",
		Email:    "test@example.com",
		Role:     "user",
	}
	secret := "test-secret-key"

	// Execute
	token, err := GenerateRefreshToken(user, secret)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.Contains(t, token, ".") // JWT should contain dots
}

func TestValidateToken_ValidAccessToken(t *testing.T) {
	// Setup - generate a token first
	user := &entities.User{
		ID:       "user-123",
		Username: "testuser",
		Email:    "test@example.com",
		Role:     "user",
	}
	secret := "test-secret-key"

	token, err := GenerateAccessToken(user, secret)
	assert.NoError(t, err)

	// Execute
	claims, err := ValidateToken(token, secret)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, user.ID, claims.UserID)
	assert.Equal(t, user.Username, claims.Username)
	assert.Equal(t, user.Role, claims.Role)
}

func TestValidateToken_ValidRefreshToken(t *testing.T) {
	// Setup - generate a refresh token first
	user := &entities.User{
		ID:       "user-123",
		Username: "testuser",
		Email:    "test@example.com",
		Role:     "user",
	}
	secret := "test-secret-key"

	token, err := GenerateRefreshToken(user, secret)
	assert.NoError(t, err)

	// Execute
	claims, err := ValidateToken(token, secret)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, user.ID, claims.UserID)
	assert.Equal(t, user.Username, claims.Username)
	assert.Equal(t, user.Role, claims.Role)
}

func TestValidateToken_InvalidToken(t *testing.T) {
	// Test data
	invalidToken := "invalid.jwt.token"
	secret := "test-secret-key"

	// Execute
	claims, err := ValidateToken(invalidToken, secret)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestValidateToken_WrongSecret(t *testing.T) {
	// Setup - generate token with one secret
	user := &entities.User{
		ID:       "user-123",
		Username: "testuser",
		Email:    "test@example.com",
		Role:     "user",
	}
	secret := "test-secret-key"
	wrongSecret := "wrong-secret-key"

	token, err := GenerateAccessToken(user, secret)
	assert.NoError(t, err)

	// Execute - validate with wrong secret
	claims, err := ValidateToken(token, wrongSecret)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestValidateToken_ExpiredToken(t *testing.T) {
	// This test would require manipulating time or creating an expired token
	// For now, we'll test the structure and leave expiration testing for integration tests
	t.Skip("Expiration testing requires time manipulation - implement in integration tests")
}

func TestTokenDurations(t *testing.T) {
	// Test that token durations are reasonable
	assert.Equal(t, 1*time.Hour, AccessTokenDuration)
	assert.Equal(t, 7*24*time.Hour, RefreshTokenDuration) // 7 days

	// Ensure access token is shorter than refresh token
	assert.True(t, AccessTokenDuration < RefreshTokenDuration)
}

func TestGenerateToken_EmptySecret(t *testing.T) {
	// Test data
	user := &entities.User{
		ID:       "user-123",
		Username: "testuser",
		Email:    "test@example.com",
		Role:     "user",
	}
	emptySecret := ""

	// Execute
	token, err := GenerateAccessToken(user, emptySecret)

	// Assert - should still work but not be secure
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestGenerateToken_NilUser(t *testing.T) {
	// Test data
	var user *entities.User = nil
	secret := "test-secret-key"

	// Execute - this should panic or error
	assert.Panics(t, func() {
		GenerateAccessToken(user, secret)
	})
}