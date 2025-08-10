package usersvc

import (
	"context"
	"testing"
	"time"

	"anchor-blog/internal/domain/entities"
	tokenrepo "anchor-blog/internal/repository/token"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, id string) (*entities.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, id string, user *entities.User) (*entities.User, error) {
	args := m.Called(ctx, id, user)
	return args.Get(0).(*entities.User), args.Error(1)
}

// Mock ActivationTokenRepository
type MockActivationTokenRepository struct {
	mock.Mock
}

func (m *MockActivationTokenRepository) StoreActivationToken(ctx context.Context, token *entities.ActivationToken) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockActivationTokenRepository) FindActivationToken(ctx context.Context, token string) (*entities.ActivationToken, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(*entities.ActivationToken), args.Error(1)
}

func (m *MockActivationTokenRepository) IsTokenValid(ctx context.Context, token string) (bool, error) {
	args := m.Called(ctx, token)
	return args.Bool(0), args.Error(1)
}

func (m *MockActivationTokenRepository) MarkTokenAsUsed(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func TestActivationService_SendActivationEmail_Success(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepository)
	mockTokenRepo := new(MockActivationTokenRepository)
	
	service := &ActivationService{
		userRepo:            mockUserRepo,
		activationTokenRepo: mockTokenRepo,
	}

	// Test data
	testUser := &entities.User{
		ID:       "user-123",
		Username: "testuser",
		Email:    "test@example.com",
	}

	// Mock expectations
	mockTokenRepo.On("StoreActivationToken", mock.Anything, mock.AnythingOfType("*entities.ActivationToken")).Return(nil)

	// Execute
	err := service.SendActivationEmail(context.Background(), testUser)

	// Assert
	assert.NoError(t, err)
	mockTokenRepo.AssertExpectations(t)
}

func TestActivationService_SendActivationEmail_StoreFailure(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepository)
	mockTokenRepo := new(MockActivationTokenRepository)
	
	service := &ActivationService{
		userRepo:            mockUserRepo,
		activationTokenRepo: mockTokenRepo,
	}

	// Test data
	testUser := &entities.User{
		ID:       "user-123",
		Username: "testuser",
		Email:    "test@example.com",
	}

	// Mock expectations - store fails
	mockTokenRepo.On("StoreActivationToken", mock.Anything, mock.AnythingOfType("*entities.ActivationToken")).Return(assert.AnError)

	// Execute
	err := service.SendActivationEmail(context.Background(), testUser)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to store activation token")
	mockTokenRepo.AssertExpectations(t)
}

func TestActivationService_VerifyActivation_Success(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepository)
	mockTokenRepo := new(MockActivationTokenRepository)
	
	service := &ActivationService{
		userRepo:            mockUserRepo,
		activationTokenRepo: mockTokenRepo,
	}

	// Test data
	token := "valid-token-123"
	userID := "user-123"
	
	activationToken := &entities.ActivationToken{
		ID:        "token-id-123",
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		Used:      false,
	}

	testUser := &entities.User{
		ID:        userID,
		Username:  "testuser",
		Email:     "test@example.com",
		Activated: false,
		Role:      "unverified",
	}

	updatedUser := &entities.User{
		ID:        userID,
		Username:  "testuser",
		Email:     "test@example.com",
		Activated: true,
		Role:      "user",
	}

	// Mock expectations
	mockTokenRepo.On("IsTokenValid", mock.Anything, token).Return(true, nil)
	mockTokenRepo.On("FindActivationToken", mock.Anything, token).Return(activationToken, nil)
	mockUserRepo.On("GetUserByID", mock.Anything, userID).Return(testUser, nil)
	mockUserRepo.On("UpdateUser", mock.Anything, userID, mock.AnythingOfType("*entities.User")).Return(updatedUser, nil)
	mockTokenRepo.On("MarkTokenAsUsed", mock.Anything, token).Return(nil)

	// Execute
	result, err := service.VerifyActivation(context.Background(), token)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, userID, result.ID)
	assert.True(t, result.Activated)
	assert.Equal(t, "user", result.Role)

	// Verify all mocks were called
	mockTokenRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestActivationService_VerifyActivation_EmptyToken(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepository)
	mockTokenRepo := new(MockActivationTokenRepository)
	
	service := &ActivationService{
		userRepo:            mockUserRepo,
		activationTokenRepo: mockTokenRepo,
	}

	// Execute with empty token
	result, err := service.VerifyActivation(context.Background(), "")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "activation token is required")

	// Verify no mocks were called
	mockTokenRepo.AssertNotCalled(t, "IsTokenValid")
	mockUserRepo.AssertNotCalled(t, "GetUserByID")
}

func TestActivationService_VerifyActivation_InvalidToken(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepository)
	mockTokenRepo := new(MockActivationTokenRepository)
	
	service := &ActivationService{
		userRepo:            mockUserRepo,
		activationTokenRepo: mockTokenRepo,
	}

	token := "invalid-token-123"

	// Mock expectations - token is invalid
	mockTokenRepo.On("IsTokenValid", mock.Anything, token).Return(false, nil)

	// Execute
	result, err := service.VerifyActivation(context.Background(), token)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid or expired activation token")

	// Verify only token validation was called
	mockTokenRepo.AssertExpectations(t)
	mockUserRepo.AssertNotCalled(t, "GetUserByID")
}

func TestActivationService_VerifyActivation_UserNotFound(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepository)
	mockTokenRepo := new(MockActivationTokenRepository)
	
	service := &ActivationService{
		userRepo:            mockUserRepo,
		activationTokenRepo: mockTokenRepo,
	}

	token := "valid-token-123"
	userID := "nonexistent-user"
	
	activationToken := &entities.ActivationToken{
		ID:        "token-id-123",
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		Used:      false,
	}

	// Mock expectations - user not found
	mockTokenRepo.On("IsTokenValid", mock.Anything, token).Return(true, nil)
	mockTokenRepo.On("FindActivationToken", mock.Anything, token).Return(activationToken, nil)
	mockUserRepo.On("GetUserByID", mock.Anything, userID).Return((*entities.User)(nil), assert.AnError)

	// Execute
	result, err := service.VerifyActivation(context.Background(), token)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to find user")

	// Verify mocks were called up to the failure point
	mockTokenRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestActivationService_VerifyActivation_AlreadyActivated(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepository)
	mockTokenRepo := new(MockActivationTokenRepository)
	
	service := &ActivationService{
		userRepo:            mockUserRepo,
		activationTokenRepo: mockTokenRepo,
	}

	token := "valid-token-123"
	userID := "user-123"
	
	activationToken := &entities.ActivationToken{
		ID:        "token-id-123",
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		Used:      false,
	}

	// User is already activated
	testUser := &entities.User{
		ID:        userID,
		Username:  "testuser",
		Email:     "test@example.com",
		Activated: true, // Already activated
		Role:      "user",
	}

	// Mock expectations
	mockTokenRepo.On("IsTokenValid", mock.Anything, token).Return(true, nil)
	mockTokenRepo.On("FindActivationToken", mock.Anything, token).Return(activationToken, nil)
	mockUserRepo.On("GetUserByID", mock.Anything, userID).Return(testUser, nil)

	// Execute
	result, err := service.VerifyActivation(context.Background(), token)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "user account is already activated")

	// Verify mocks were called up to the check
	mockTokenRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
	
	// Update and mark token should not be called
	mockUserRepo.AssertNotCalled(t, "UpdateUser")
	mockTokenRepo.AssertNotCalled(t, "MarkTokenAsUsed")
}