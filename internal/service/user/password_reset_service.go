package userservice

import (
	"anchor-blog/internal/domain/entities"
	"anchor-blog/pkg/hashutil"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"time"
)

type PasswordResetService struct {
	// TODO: Add repository dependencies when available
}

// NewPasswordResetService creates a new password reset service
func NewPasswordResetService() *PasswordResetService {
	return &PasswordResetService{}
}

// ForgotPassword generates a reset token and sends reset email
func (s *PasswordResetService) ForgotPassword(email string) error {
	// TODO: Find user by email in database
	// For now, validate email format and proceed with mock user
	if email == "" {
		return fmt.Errorf("email is required")
	}

	// Generate a unique reset token
	token, err := s.generateResetToken()
	if err != nil {
		return fmt.Errorf("failed to generate reset token: %w", err)
	}

	// Create password reset token entity
	resetToken := &entities.PasswordResetToken{
		ID:        s.generateTokenID(),
		UserID:    "mock-user-id", // TODO: Get from database lookup
		Token:     token,
		ExpiresAt: time.Now().Add(1 * time.Hour), // Token expires in 1 hour
		Used:      false,
		CreatedAt: time.Now(),
	}

	// TODO: Store reset token in database when repository is available
	_ = resetToken

	// For now, log the reset link to console
	resetLink := fmt.Sprintf("http://localhost:8080/api/v1/users/reset-password?token=%s", token)
	log.Printf("🔐 Password reset email for: %s", email)
	log.Printf("📧 Reset Link: %s", resetLink)
	log.Printf("⏰ Token expires at: %s", resetToken.ExpiresAt.Format("2006-01-02 15:04:05"))

	return nil
}

// ResetPassword validates the token and updates the user's password
func (s *PasswordResetService) ResetPassword(token, newPassword string) (*entities.User, error) {
	// Validate inputs
	if token == "" {
		return nil, fmt.Errorf("reset token is required")
	}
	if newPassword == "" {
		return nil, fmt.Errorf("new password is required")
	}
	if len(newPassword) < 6 {
		return nil, fmt.Errorf("password must be at least 6 characters long")
	}

	// TODO: Find reset token in database
	// TODO: Check if token exists and is not expired
	// TODO: Check if token is not already used
	// TODO: Find user by UserID from token

	// Hash the new password
	hashedPassword, err := hashutil.HashPassword(newPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// TODO: Update user password in database
	// TODO: Mark token as used
	// TODO: Save user to database

	// Mock response for now
	log.Printf("✅ Password reset token validated: %s", token)
	log.Printf("🔒 Password updated successfully!")

	// Return mock updated user
	mockUser := &entities.User{
		ID:           "mock-user-id",
		Username:     "reset-user",
		Email:        "user@example.com",
		PasswordHash: hashedPassword,
		Activated:    true,
		Role:         "user",
		UpdatedAt:    time.Now(),
	}

	return mockUser, nil
}

// generateResetToken creates a random hex token
func (s *PasswordResetService) generateResetToken() (string, error) {
	bytes := make([]byte, 32) // 32 bytes = 64 hex characters
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// generateTokenID creates a unique ID for the reset token
func (s *PasswordResetService) generateTokenID() string {
	bytes := make([]byte, 12)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}