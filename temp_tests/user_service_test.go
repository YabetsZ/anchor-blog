package usersvc

import (
	"context"
	"testing"

	"anchor-blog/config"
	"anchor-blog/internal/domain/entities"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock interfaces for testing
type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) GetUserByID(ctx context.Context, id string) (*entities.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserRepo) GetUserByUsername(ctx context.Context, username string) (*entities.User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserRepo) CreateUser(ctx context.Context, user *entities.User) (string, error) {
	args := m.Called(ctx, user)
	return args.String(0), args.Error(1)
}

func (m *MockUserRepo) CheckUsername(ctx context.Context, username string) (bool, error) {
	args := m.Called(ctx, username)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepo) CheckEmail(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepo) CountAllUsers(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserRepo) SetLastSeen(ctx context.Context, userID string, lastSeen interface{}) error {
	args := m.Called(ctx, userID, lastSeen)
	return args.Error(0)
}

type MockTokenRepo struct {
	mock.Mock
}

func (m *MockTokenRepo) StoreRefreshToken(ctx context.Context, token *entities.RefreshToken) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockTokenRepo) DeleteAllByUserID(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func TestUserServices_GetUserByID_Success(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepo)
	mockTokenRepo := new(MockTokenRepo)
	cfg := &config.Config{}

	service := NewUserServices(mockUserRepo, mockTokenRepo, cfg)

	// Test data
	userID := "user-123"
	expectedUser := &entities.User{
		ID:       userID,
		Username: "testuser",
		Email:    "test@example.com",
	}

	// Mock expectations
	mockUserRepo.On("GetUserByID", mock.Anything, userID).Return(expectedUser, nil)

	// Execute
	result, err := service.GetUserByID(context.Background(), userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, userID, result.ID)
	assert.Equal(t, "testuser", result.Username)
	assert.Equal(t, "test@example.com", result.Email)

	// Verify mock was called
	mockUserRepo.AssertExpectations(t)
}

func TestUserServices_GetUserByID_NotFound(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepo)
	mockTokenRepo := new(MockTokenRepo)
	cfg := &config.Config{}

	service := NewUserServices(mockUserRepo, mockTokenRepo, cfg)

	// Test data
	userID := "nonexistent-user"

	// Mock expectations - user not found
	mockUserRepo.On("GetUserByID", mock.Anything, userID).Return((*entities.User)(nil), assert.AnError)

	// Execute
	result, err := service.GetUserByID(context.Background(), userID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)

	// Verify mock was called
	mockUserRepo.AssertExpectations(t)
}

func TestUserServices_Register_FirstUser_SuperAdmin(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepo)
	mockTokenRepo := new(MockTokenRepo)
	cfg := &config.Config{}

	service := NewUserServices(mockUserRepo, mockTokenRepo, cfg)

	// Test data
	userDTO := &UserDTO{
		Username:  "admin",
		Email:     "admin@example.com",
		Password:  "password123",
		FirstName: "Admin",
		LastName:  "User",
	}

	expectedUserID := "admin-user-123"

	// Mock expectations - first user (count = 0)
	mockUserRepo.On("CheckUsername", mock.Anything, "admin").Return(false, nil)
	mockUserRepo.On("CheckEmail", mock.Anything, "admin@example.com").Return(false, nil)
	mockUserRepo.On("CountAllUsers", mock.Anything).Return(int64(0), nil)
	mockUserRepo.On("CreateUser", mock.Anything, mock.MatchedBy(func(user *entities.User) bool {
		return user.Role == "superadmin" && user.Username == "admin"
	})).Return(expectedUserID, nil)

	// Execute
	result, err := service.Register(context.Background(), userDTO)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedUserID, result)

	// Verify mocks were called
	mockUserRepo.AssertExpectations(t)
}

func TestUserServices_Register_RegularUser_Unverified(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepo)
	mockTokenRepo := new(MockTokenRepo)
	cfg := &config.Config{}

	service := NewUserServices(mockUserRepo, mockTokenRepo, cfg)

	// Test data
	userDTO := &UserDTO{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
	}

	expectedUserID := "test-user-123"

	// Mock expectations - not first user (count > 0)
	mockUserRepo.On("CheckUsername", mock.Anything, "testuser").Return(false, nil)
	mockUserRepo.On("CheckEmail", mock.Anything, "test@example.com").Return(false, nil)
	mockUserRepo.On("CountAllUsers", mock.Anything).Return(int64(5), nil) // Existing users
	mockUserRepo.On("CreateUser", mock.Anything, mock.MatchedBy(func(user *entities.User) bool {
		return user.Role == "unverified" && user.Username == "testuser" && !user.Activated
	})).Return(expectedUserID, nil)

	// Execute
	result, err := service.Register(context.Background(), userDTO)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedUserID, result)

	// Verify mocks were called
	mockUserRepo.AssertExpectations(t)
}

func TestUserServices_Register_UsernameExists(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepo)
	mockTokenRepo := new(MockTokenRepo)
	cfg := &config.Config{}

	service := NewUserServices(mockUserRepo, mockTokenRepo, cfg)

	// Test data
	userDTO := &UserDTO{
		Username:  "existinguser",
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
	}

	// Mock expectations - username already exists
	mockUserRepo.On("CheckUsername", mock.Anything, "existinguser").Return(true, nil)

	// Execute
	result, err := service.Register(context.Background(), userDTO)

	// Assert
	assert.Error(t, err)
	assert.Empty(t, result)

	// Verify only username check was called
	mockUserRepo.AssertExpectations(t)
	mockUserRepo.AssertNotCalled(t, "CheckEmail")
	mockUserRepo.AssertNotCalled(t, "CreateUser")
}

func TestUserServices_Register_EmailExists(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepo)
	mockTokenRepo := new(MockTokenRepo)
	cfg := &config.Config{}

	service := NewUserServices(mockUserRepo, mockTokenRepo, cfg)

	// Test data
	userDTO := &UserDTO{
		Username:  "testuser",
		Email:     "existing@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
	}

	// Mock expectations - email already exists
	mockUserRepo.On("CheckUsername", mock.Anything, "testuser").Return(false, nil)
	mockUserRepo.On("CheckEmail", mock.Anything, "existing@example.com").Return(true, nil)

	// Execute
	result, err := service.Register(context.Background(), userDTO)

	// Assert
	assert.Error(t, err)
	assert.Empty(t, result)

	// Verify checks were called but not create
	mockUserRepo.AssertExpectations(t)
	mockUserRepo.AssertNotCalled(t, "CreateUser")
}

func TestUserServices_Register_InvalidInput(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepo)
	mockTokenRepo := new(MockTokenRepo)
	cfg := &config.Config{}

	service := NewUserServices(mockUserRepo, mockTokenRepo, cfg)

	// Test cases for invalid input
	testCases := []struct {
		name    string
		userDTO *UserDTO
	}{
		{
			name: "Empty username",
			userDTO: &UserDTO{
				Username:  "",
				Email:     "test@example.com",
				Password:  "password123",
				FirstName: "Test",
				LastName:  "User",
			},
		},
		{
			name: "Empty email",
			userDTO: &UserDTO{
				Username:  "testuser",
				Email:     "",
				Password:  "password123",
				FirstName: "Test",
				LastName:  "User",
			},
		},
		{
			name: "Empty first name",
			userDTO: &UserDTO{
				Username:  "testuser",
				Email:     "test@example.com",
				Password:  "password123",
				FirstName: "",
				LastName:  "User",
			},
		},
		{
			name: "Empty last name",
			userDTO: &UserDTO{
				Username:  "testuser",
				Email:     "test@example.com",
				Password:  "password123",
				FirstName: "Test",
				LastName:  "",
			},
		},
		{
			name: "Short first name",
			userDTO: &UserDTO{
				Username:  "testuser",
				Email:     "test@example.com",
				Password:  "password123",
				FirstName: "Te", // Too short
				LastName:  "User",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Execute
			result, err := service.Register(context.Background(), tc.userDTO)

			// Assert
			assert.Error(t, err)
			assert.Empty(t, result)

			// Verify no repository methods were called for invalid input
			mockUserRepo.AssertNotCalled(t, "CheckUsername")
			mockUserRepo.AssertNotCalled(t, "CheckEmail")
			mockUserRepo.AssertNotCalled(t, "CreateUser")
		})
	}
}