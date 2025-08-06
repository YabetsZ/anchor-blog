package usersvc

import (
	"anchor-blog/internal/domain/entities"
	tokenrepo "anchor-blog/internal/repository/token"
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"time"
)

type ActivationService struct {
	userRepo            entities.IUserRepository
	activationTokenRepo *tokenrepo.ActivationTokenRepository
}

// NewActivationService creates a new activation service
func NewActivationService(userRepo entities.IUserRepository, activationTokenRepo *tokenrepo.ActivationTokenRepository) *ActivationService {
	return &ActivationService{
		userRepo:            userRepo,
		activationTokenRepo: activationTokenRepo,
	}
}

// SendActivationEmail generates an activation token and logs the activation link
func (s *ActivationService) SendActivationEmail(ctx context.Context, user *entities.User) error {
	// Generate a unique activation token
	token, err := s.generateActivationToken()
	if err != nil {
		return fmt.Errorf("failed to generate activation token: %w", err)
	}

	// Create activation token entity
	activationToken := &entities.ActivationToken{
		ID:        s.generateTokenID(),
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour), // Token expires in 24 hours
		Used:      false,
		CreatedAt: time.Now(),
	}

	// Store activation token in database
	err = s.activationTokenRepo.StoreActivationToken(ctx, activationToken)
	if err != nil {
		return fmt.Errorf("failed to store activation token: %w", err)
	}

	// Log the activation link to console
	activationLink := fmt.Sprintf("http://localhost:8080/api/v1/users/activate?token=%s", token)
	log.Printf("🔗 Activation email for user %s (%s):", user.Username, user.Email)
	log.Printf("📧 Activation Link: %s", activationLink)
	log.Printf("⏰ Token expires at: %s", activationToken.ExpiresAt.Format("2006-01-02 15:04:05"))

	return nil
}

// VerifyActivation validates the token and activates the user
func (s *ActivationService) VerifyActivation(ctx context.Context, token string) (*entities.User, error) {
	if token == "" {
		return nil, fmt.Errorf("activation token is required")
	}

	// Validate token (check if exists, not expired, not used)
	isValid, err := s.activationTokenRepo.IsTokenValid(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("token validation failed: %w", err)
	}
	if !isValid {
		return nil, fmt.Errorf("invalid or expired activation token")
	}

	// Find activation token to get user ID
	activationToken, err := s.activationTokenRepo.FindActivationToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("failed to find activation token: %w", err)
	}

	// Find user by ID
	user, err := s.userRepo.GetUserByID(ctx, activationToken.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Activate user and set role
	err = s.userRepo.ActivateUserByID(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to activate user: %w", err)
	}

	// Set user role to "user" if not already set
	if user.Role == "unverified" {
		err = s.userRepo.SetRole(ctx, user.ID, "user")
		if err != nil {
			return nil, fmt.Errorf("failed to set user role: %w", err)
		}
	}

	// Mark token as used
	err = s.activationTokenRepo.MarkTokenAsUsed(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("failed to mark token as used: %w", err)
	}

	// Get updated user
	updatedUser, err := s.userRepo.GetUserByID(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated user: %w", err)
	}

	log.Printf("✅ Activation token validated: %s", token)
	log.Printf("🎉 User account activated successfully for: %s", updatedUser.Username)

	return updatedUser, nil
}

// generateActivationToken creates a random hex token
func (s *ActivationService) generateActivationToken() (string, error) {
	bytes := make([]byte, 32) // 32 bytes = 64 hex characters
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// generateTokenID creates a unique ID for the activation token
func (s *ActivationService) generateTokenID() string {
	bytes := make([]byte, 12)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
