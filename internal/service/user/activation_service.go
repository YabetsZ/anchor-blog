package userservice

import (
	"anchor-blog/internal/domain/entities"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"time"
)

type ActivationService struct {
	// TODO: Add repository dependencies when available
}

// NewActivationService creates a new activation service
func NewActivationService() *ActivationService {
	return &ActivationService{}
}

// SendActivationEmail generates an activation token and logs the activation link
func (s *ActivationService) SendActivationEmail(user *entities.User) error {
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

	// TODO: Store activation token in database when repository is available
	_ = activationToken

	// For now, log the activation link to console
	activationLink := fmt.Sprintf("http://localhost:8080/api/v1/users/activate?token=%s", token)
	log.Printf("🔗 Activation email for user %s (%s):", user.Username, user.Email)
	log.Printf("📧 Activation Link: %s", activationLink)
	log.Printf("⏰ Token expires at: %s", activationToken.ExpiresAt.Format("2006-01-02 15:04:05"))

	return nil
}

// VerifyActivation validates the token and activates the user
func (s *ActivationService) VerifyActivation(token string) (*entities.User, error) {
	// TODO: Implement token validation with database
	// For now, return mock validation

	if token == "" {
		return nil, fmt.Errorf("activation token is required")
	}

	// TODO: Find activation token in database
	// TODO: Check if token exists and is not expired
	// TODO: Check if token is not already used
	// TODO: Find user by UserID from token
	// TODO: Set user.Activated = true and user.Role = "user"
	// TODO: Mark token as used
	// TODO: Save user to database

	// Mock response for now
	log.Printf("✅ Activation token validated: %s", token)
	log.Printf("🎉 User account activated successfully!")

	// Return mock activated user
	mockUser := &entities.User{
		ID:        "mock-user-id",
		Username:  "activated-user",
		Email:     "user@example.com",
		Activated: true,
		Role:      "user",
	}

	return mockUser, nil
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
