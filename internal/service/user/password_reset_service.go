package usersvc

import (
	"anchor-blog/internal/domain/entities"
	tokenrepo "anchor-blog/internal/repository/token"
	"anchor-blog/pkg/hashutil"
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"time"
)

type PasswordResetService struct {
	userRepo               entities.IUserRepository
	passwordResetTokenRepo *tokenrepo.PasswordResetTokenRepository
}

// NewPasswordResetService creates a new password reset service
func NewPasswordResetService(userRepo entities.IUserRepository, passwordResetTokenRepo *tokenrepo.PasswordResetTokenRepository) *PasswordResetService {
	return &PasswordResetService{
		userRepo:               userRepo,
		passwordResetTokenRepo: passwordResetTokenRepo,
	}
}

// ForgotPassword generates a reset token and sends reset email
func (s *PasswordResetService) ForgotPassword(ctx context.Context, email string) error {
	if email == "" {
		return fmt.Errorf("email is required")
	}

	// Find user by email
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("user with email %s not found", email)
	}

	// Generate a unique reset token
	token, err := s.generateResetToken()
	if err != nil {
		return fmt.Errorf("failed to generate reset token: %w", err)
	}

	// Create password reset token entity
	resetToken := &entities.PasswordResetToken{
		ID:        s.generateTokenID(),
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(1 * time.Hour), // Token expires in 1 hour
		Used:      false,
		CreatedAt: time.Now(),
	}

	// Store reset token in database
	err = s.passwordResetTokenRepo.StorePasswordResetToken(ctx, resetToken)
	if err != nil {
		return fmt.Errorf("failed to store reset token: %w", err)
	}

	// Log the reset link to console
	resetLink := fmt.Sprintf("http://localhost:8080/api/v1/users/reset-password?token=%s", token)
	log.Printf("🔐 Password reset email for: %s", email)
	log.Printf("📧 Reset Link: %s", resetLink)
	log.Printf("⏰ Token expires at: %s", resetToken.ExpiresAt.Format("2006-01-02 15:04:05"))

	return nil
}

// ResetPassword validates the token and updates the user's password
func (s *PasswordResetService) ResetPassword(ctx context.Context, token, newPassword string) (*entities.User, error) {
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

	// Validate token (check if exists, not expired, not used)
	isValid, err := s.passwordResetTokenRepo.IsTokenValid(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("token validation failed: %w", err)
	}
	if !isValid {
		return nil, fmt.Errorf("invalid or expired reset token")
	}

	// Find reset token to get user ID
	resetToken, err := s.passwordResetTokenRepo.FindPasswordResetToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("failed to find reset token: %w", err)
	}

	// Find user by ID
	user, err := s.userRepo.GetUserByID(ctx, resetToken.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Hash the new password
	hashedPassword, err := hashutil.HashPassword(newPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Update user password in database
	err = s.userRepo.ChangePassword(ctx, user.ID, hashedPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to update password: %w", err)
	}

	// Mark token as used
	err = s.passwordResetTokenRepo.MarkTokenAsUsed(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("failed to mark token as used: %w", err)
	}

	// Get updated user
	updatedUser, err := s.userRepo.GetUserByID(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated user: %w", err)
	}

	log.Printf("✅ Password reset token validated: %s", token)
	log.Printf("🔒 Password updated successfully for user: %s", updatedUser.Username)

	return updatedUser, nil
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