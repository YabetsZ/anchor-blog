package usersvc

import (
	"context"
	"testing"
	"time"

	"anchor-blog/config"
	"anchor-blog/internal/domain/entities"
	"anchor-blog/internal/errors"
	"anchor-blog/pkg/hashutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock for login tests
type MockUserRepoForLogin struct {
	mock.Mock
}

func (m *MockUserRepoForLogin) GetUserByUsername(ctx context.Context, username string) (*entities.User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(*entities.User), args.Error(1)
}

// Implement other required interface methods as no-ops
func (m *MockUserRepoForLogin) GetUserByID(ctx context.Context, id string) (*entities.User, error) { return nil, nil }
func (m *MockUserRepoForLogin) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) { return nil, nil }
func (m *MockUserRepoForLogin) GetUsers(ctx context.Context, limit, offset int64) ([]*entities.User, error) { return nil, nil }
func (m *MockUserRepoForLogin) CountUsersByRole(ctx context.Context, role string) (int64, error) { return 0, nil }
func (m *MockUserRepoForLogin) CountAllUsers(ctx context.Context) (int64, error) { return 0, nil }
func (m *MockUserRepoForLogin) CountActiveUsers(ctx context.Context) (int64, error) { return 0, nil }
func (m *MockUserRepoForLogin) CountInactiveUsers(ctx context.Context) (int64, error) { return 0, nil }
func (m *MockUserRepoForLogin) GetUserRoleByID(ctx context.Context, userID string) (string, error) { return "", nil }
func (m *MockUserRepoForLogin) CreateUser(ctx context.Context, user *entities.User) (string, error) { return "", nil }
func (m *MockUserRepoForLogin) EditUserByID(ctx context.Context, id string, user *entities.User) error { return nil }
func (m *MockUserRepoForLogin) DeleteUserByID(ctx context.Context, id string) error { return nil }
func (m *MockUserRepoForLogin) SetLastSeen(ctx context.Context, id string, timestamp time.Time) error { return nil }
func (m *MockUserRepoForLogin) UpdateUserRole(ctx context.Context, adminID, targetID, role string) error { return nil }
func (m *MockUserRepoForLogin) CheckEmail(ctx context.Context, email string) (bool, error) { return false, nil }
func (m *MockUserRepoForLogin) CheckUsername(ctx context.Context, username string) (bool, error) { return false, nil }
func (m *MockUserRepoForLogin) ChangePassword(ctx context.Context, id string, newHashedPassword string) error { return nil }
func (m *MockUserRepoForLogin) ChangeEmail(ctx context.Context, email string, newEmail string) error { return nil }
func (m *MockUserRepoForLogin) SetRole(ctx context.Context, id string, role string) error { return nil }
func (m *MockUserRepoForLogin) ActivateUserByID(ctx context.Context, id string) error { return nil }
func (m *MockUserRepoForLogin) DeactivateUserByID(ctx context.Context, id string) error { return nil }

type MockTokenRepoForLogin struct {
	mock.Mock
}

func (m *MockTokenRepoForLogin) StoreRefreshToken(ctx context.Context, token *entities.RefreshToken) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockTokenRepoForLogin) FindByHash(ctx context.Context, hash string) (*entities.RefreshToken, error) { return nil, nil }
func (m *MockTokenRepoForLogin) DeleteByHash(ctx context.Context, hash string) error { return nil }
func (m *MockTokenRepoForLogin) DeleteAllByUserID(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func TestLogin_Success(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepoForLogin)
	mockTokenRepo := new(MockTokenRepoForLogin)
	cfg := &config.Config{
		JWT: struct {
			AccessTokenSecret  string `mapstructure:"access_token_secret"`
			RefreshTokenSecret string `mapstructure:"refresh_token_secret"`
		}{
			AccessTokenSecret:  "test-access-secret",
			RefreshTokenSecret: "test-refresh-secret",
		},
		HMAC: struct {
			Secret string `mapstructure:"hmac_secret"`
		}{
			Secret: "test-hmac-secret",
		},
	}

	userService := &UserServices{
		userRepo:  mockUserRepo,
		tokenRepo: mockTokenRepo,
		cfg:       cfg,
	}

	// Test data - user with hashed password "password123"
	// Generate a proper bcrypt hash for testing
	hashedPassword, _ := hashutil.HashPassword("password123")
	testUser := &entities.User{
		ID:           "user-123",
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: hashedPassword,
		Role:         "user",
		Activated:    true,
	}

	// Mock expectations
	mockUserRepo.On("GetUserByUsername", mock.Anything, "testuser").Return(testUser, nil)
	mockTokenRepo.On("StoreRefreshToken", mock.Anything, mock.AnythingOfType("*entities.RefreshToken")).Return(nil)

	// Execute
	result, err := userService.Login(context.Background(), "testuser", "password123")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.AccessToken)
	assert.NotEmpty(t, result.RefreshToken)

	// Verify mocks were called
	mockUserRepo.AssertExpectations(t)
	mockTokenRepo.AssertExpectations(t)
}

func TestLogin_UserNotFound(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepoForLogin)
	mockTokenRepo := new(MockTokenRepoForLogin)
	cfg := &config.Config{}

	userService := &UserServices{
		userRepo:  mockUserRepo,
		tokenRepo: mockTokenRepo,
		cfg:       cfg,
	}

	// Mock expectations - user not found
	mockUserRepo.On("GetUserByUsername", mock.Anything, "nonexistent").Return((*entities.User)(nil), errors.ErrUserNotFound)

	// Execute
	result, err := userService.Login(context.Background(), "nonexistent", "password123")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, errors.ErrUserNotFound, err)

	// Verify mocks were called
	mockUserRepo.AssertExpectations(t)
	mockTokenRepo.AssertNotCalled(t, "StoreRefreshToken")
}

func TestLogin_InvalidPassword(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepoForLogin)
	mockTokenRepo := new(MockTokenRepoForLogin)
	cfg := &config.Config{}

	userService := &UserServices{
		userRepo:  mockUserRepo,
		tokenRepo: mockTokenRepo,
		cfg:       cfg,
	}

	// Test data - user with different password
	hashedPassword, _ := hashutil.HashPassword("password123")
	testUser := &entities.User{
		ID:           "user-123",
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: hashedPassword,
		Role:         "user",
		Activated:    true,
	}

	// Mock expectations
	mockUserRepo.On("GetUserByUsername", mock.Anything, "testuser").Return(testUser, nil)

	// Execute with wrong password
	result, err := userService.Login(context.Background(), "testuser", "wrongpassword")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, errors.ErrInvalidCredentials, err)

	// Verify mocks were called
	mockUserRepo.AssertExpectations(t)
	mockTokenRepo.AssertNotCalled(t, "StoreRefreshToken")
}

func TestLogout_Success(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepoForLogin)
	mockTokenRepo := new(MockTokenRepoForLogin)
	cfg := &config.Config{}

	userService := &UserServices{
		userRepo:  mockUserRepo,
		tokenRepo: mockTokenRepo,
		cfg:       cfg,
	}

	userID := "user-123"

	// Mock expectations
	mockTokenRepo.On("DeleteAllByUserID", mock.Anything, userID).Return(nil)

	// Execute
	err := userService.Logout(context.Background(), userID)

	// Assert
	assert.NoError(t, err)

	// Verify mocks were called
	mockTokenRepo.AssertExpectations(t)
}

func TestLogout_TokenDeletionFails(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepoForLogin)
	mockTokenRepo := new(MockTokenRepoForLogin)
	cfg := &config.Config{}

	userService := &UserServices{
		userRepo:  mockUserRepo,
		tokenRepo: mockTokenRepo,
		cfg:       cfg,
	}

	userID := "user-123"

	// Mock expectations - deletion fails
	mockTokenRepo.On("DeleteAllByUserID", mock.Anything, userID).Return(assert.AnError)

	// Execute
	err := userService.Logout(context.Background(), userID)

	// Assert
	assert.Error(t, err)

	// Verify mocks were called
	mockTokenRepo.AssertExpectations(t)
}