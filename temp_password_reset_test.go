package usersvc

import (
	"context"
	"testing"
	"time"

	"anchor-blog/internal/domain/entities"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock for password reset tests
type MockUserRepoForPasswordReset struct {
	mock.Mock
}

func (m *MockUserRepoForPasswordReset) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserRepoForPasswordReset) GetUserByID(ctx context.Context, id string) (*entities.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserRepoForPasswordReset) ChangePassword(ctx context.Context, id string, newHashedPassword string) error {
	args := m.Called(ctx, id, newHashedPassword)
	return args.Error(0)
}

// Implement other required interface methods as no-ops
func (m *MockUserRepoForPasswordReset) GetUserByUsername(ctx context.Context, username string) (*entities.User, error) { return nil, nil }
func (m *MockUserRepoForPasswordReset) GetUsers(ctx context.Context, limit, offset int64) ([]*entities.User, error) { return nil, nil }
func (m *MockUserRepoForPasswordReset) CountUsersByRole(ctx context.Context, role string) (int64, error) { return 0, nil }
func (m *MockUserRepoForPasswordReset) CountAllUsers(ctx context.Context) (int64, error) { return 0, nil }
func (m *MockUserRepoForPasswordReset) CountActiveUsers(ctx context.Context) (int64, error) { return 0, nil }
func (m *MockUserRepoForPasswordReset) CountInactiveUsers(ctx context.Context) (int64, error) { return 0, nil }
func (m *MockUserRepoForPasswordReset) GetUserRoleByID(ctx context.Context, userID string) (string, error) { return "", nil }
func (m *MockUserRepoForPasswordReset) CreateUser(ctx context.Context, user *entities.User) (string, error) { return "", nil }
func (m *MockUserRepoForPasswordReset) EditUserByID(ctx context.Context, id string, user *entities.User) error { return nil }
func (m *MockUserRepoForPasswordReset) DeleteUserByID(ctx context.Context, id string) error { return nil }
func (m *MockUserRepoForPasswordReset) SetLastSeen(ctx context.Context, id string, timestamp time.Time) error { return nil }
func (m *MockUserRepoForPasswordReset) UpdateUserRole(ctx context.Context, adminID, targetID, role string) error { return nil }
func (m *MockUserRepoForPasswordReset) CheckEmail(ctx context.Context, email string) (bool, error) { return false, nil }
func (m *MockUserRepoForPasswordReset) CheckUsername(ctx context.Context, username string) (bool, error) { return false, nil }
func (m *MockUserRepoForPasswordReset) ChangeEmail(ctx context.Context, email string, newEmail string) error { return nil }
func (m *MockUserRepoForPasswordReset) SetRole(ctx context.Context, id string, role string) error { return nil }
func (m *MockUserRepoForPasswordReset) ActivateUserByID(ctx context.Context, id string) error { return nil }
func (m *MockUserRepoForPasswordReset) DeactivateUserByID(ctx context.Context, id string) error { return nil }

type MockPasswordResetTokenRepo struct {
	mock.Mock
}

func (m *MockPasswordResetTokenRepo) StorePasswordResetToken(ctx context.Context, token *entities.PasswordResetToken) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockPasswordResetTokenRepo) FindPasswordResetToken(ctx context.Context, token string) (*entities.PasswordResetToken, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(*entities.PasswordResetToken), args.Error(1)
}

func (m *MockPasswordResetTokenRepo) IsTokenValid(ctx context.Context, token string) (bool, error) {
	args := m.Called(ctx, token)
	return args.Bool(0), args.Error(1)
}

func (m *MockPasswordResetTokenRepo) MarkTokenAsUsed(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func TestPasswordResetService_SendPasswordResetEmail_Success(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepoForPasswordReset)
	mockTokenRepo := new(MockPasswordResetTokenRepo)
	
	service := &PasswordResetService{
		userRepo:                mockUserRepo,
		passwordResetTokenRepo: mockTokenRepo,
	}

	// Test data
	email := "test@example.com"
	testUser := &entities.User{
		ID:       "user-123",
		Username: "testuser",
		Email:    email,
	}

	// Mock expectations
	mockUserRepo.On("GetUserByEmail", mock.Anything, email).Return(testUser, nil)
	mockTokenRepo.On("StorePasswordResetToken", mock.Anything, mock.AnythingOfType("*entities.PasswordResetToken")).Return(nil)

	// Execute
	err := service.SendPasswordResetEmail(context.Background(), email)

	// Assert
	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
	mockTokenRepo.AssertExpectations(t)
}

func TestPasswordResetService_SendPasswordResetEmail_UserNotFound(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepoForPasswordReset)
	mockTokenRepo := new(MockPasswordResetTokenRepo)
	
	service := &PasswordResetService{
		userRepo:                mockUserRepo,
		passwordResetTokenRepo: mockTokenRepo,
	}

	// Test data
	email := "nonexistent@example.com"

	// Mock expectations - user not found
	mockUserRepo.On("GetUserByEmail", mock.Anything, email).Return((*entities.User)(nil), assert.AnError)

	// Execute
	err := service.SendPasswordResetEmail(context.Background(), email)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to find user")
	mockUserRepo.AssertExpectations(t)
	mockTokenRepo.AssertNotCalled(t, "StorePasswordResetToken")
}

func TestPasswordResetService_ResetPassword_Success(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepoForPasswordReset)
	mockTokenRepo := new(MockPasswordResetTokenRepo)
	
	service := &PasswordResetService{
		userRepo:                mockUserRepo,
		passwordResetTokenRepo: mockTokenRepo,
	}

	// Test data
	token := "valid-reset-token"
	newPassword := "newpassword123"
	userID := "user-123"
	
	resetToken := &entities.PasswordResetToken{
		ID:        "token-id",
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		Used:      false,
	}

	testUser := &entities.User{
		ID:       userID,
		Username: "testuser",
		Email:    "test@example.com",
	}

	updatedUser := &entities.User{
		ID:       userID,
		Username: "testuser",
		Email:    "test@example.com",
	}

	// Mock expectations
	mockTokenRepo.On("IsTokenValid", mock.Anything, token).Return(true, nil)
	mockTokenRepo.On("FindPasswordResetToken", mock.Anything, token).Return(resetToken, nil)
	mockUserRepo.On("GetUserByID", mock.Anything, userID).Return(testUser, nil).Once()
	mockUserRepo.On("ChangePassword", mock.Anything, userID, mock.AnythingOfType("string")).Return(nil)
	mockTokenRepo.On("MarkTokenAsUsed", mock.Anything, token).Return(nil)
	mockUserRepo.On("GetUserByID", mock.Anything, userID).Return(updatedUser, nil).Once()

	// Execute
	result, err := service.ResetPassword(context.Background(), token, newPassword)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, userID, result.ID)

	// Verify all mocks were called
	mockUserRepo.AssertExpectations(t)
	mockTokenRepo.AssertExpectations(t)
}

func TestPasswordResetService_ResetPassword_InvalidToken(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepoForPasswordReset)
	mockTokenRepo := new(MockPasswordResetTokenRepo)
	
	service := &PasswordResetService{
		userRepo:                mockUserRepo,
		passwordResetTokenRepo: mockTokenRepo,
	}

	// Test data
	token := "invalid-token"
	newPassword := "newpassword123"

	// Mock expectations - token is invalid
	mockTokenRepo.On("IsTokenValid", mock.Anything, token).Return(false, nil)

	// Execute
	result, err := service.ResetPassword(context.Background(), token, newPassword)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid or expired")

	// Verify only token validation was called
	mockTokenRepo.AssertExpectations(t)
	mockUserRepo.AssertNotCalled(t, "GetUserByID")
	mockUserRepo.AssertNotCalled(t, "ChangePassword")
}

func TestPasswordResetService_ResetPassword_EmptyToken(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepoForPasswordReset)
	mockTokenRepo := new(MockPasswordResetTokenRepo)
	
	service := &PasswordResetService{
		userRepo:                mockUserRepo,
		passwordResetTokenRepo: mockTokenRepo,
	}

	// Execute with empty token
	result, err := service.ResetPassword(context.Background(), "", "newpassword123")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "token is required")

	// Verify no methods were called
	mockTokenRepo.AssertNotCalled(t, "IsTokenValid")
	mockUserRepo.AssertNotCalled(t, "GetUserByID")
}

func TestPasswordResetService_ResetPassword_EmptyPassword(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepoForPasswordReset)
	mockTokenRepo := new(MockPasswordResetTokenRepo)
	
	service := &PasswordResetService{
		userRepo:                mockUserRepo,
		passwordResetTokenRepo: mockTokenRepo,
	}

	// Execute with empty password
	result, err := service.ResetPassword(context.Background(), "valid-token", "")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "password is required")

	// Verify no methods were called
	mockTokenRepo.AssertNotCalled(t, "IsTokenValid")
	mockUserRepo.AssertNotCalled(t, "GetUserByID")
}