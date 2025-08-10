package usersvc

import (
	"context"
	"testing"
	"time"

	"anchor-blog/config"
	"anchor-blog/internal/domain/entities"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Minimal mock for testing registration logic
type MockUserRepoForRegistration struct {
	mock.Mock
}

func (m *MockUserRepoForRegistration) CheckUsername(ctx context.Context, username string) (bool, error) {
	args := m.Called(ctx, username)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepoForRegistration) CheckEmail(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepoForRegistration) CountAllUsers(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserRepoForRegistration) CreateUser(ctx context.Context, user *entities.User) (string, error) {
	args := m.Called(ctx, user)
	return args.String(0), args.Error(1)
}

func (m *MockUserRepoForRegistration) GetUserByID(ctx context.Context, id string) (*entities.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entities.User), args.Error(1)
}

// Implement other required interface methods as no-ops for this test
func (m *MockUserRepoForRegistration) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	return nil, nil
}
func (m *MockUserRepoForRegistration) GetUserByUsername(ctx context.Context, username string) (*entities.User, error) {
	return nil, nil
}
func (m *MockUserRepoForRegistration) GetUsers(ctx context.Context, limit, offset int64) ([]*entities.User, error) {
	return nil, nil
}
func (m *MockUserRepoForRegistration) CountUsersByRole(ctx context.Context, role string) (int64, error) {
	return 0, nil
}
func (m *MockUserRepoForRegistration) CountActiveUsers(ctx context.Context) (int64, error) {
	return 0, nil
}
func (m *MockUserRepoForRegistration) CountInactiveUsers(ctx context.Context) (int64, error) {
	return 0, nil
}
func (m *MockUserRepoForRegistration) GetUserRoleByID(ctx context.Context, userID string) (string, error) {
	return "", nil
}
func (m *MockUserRepoForRegistration) EditUserByID(ctx context.Context, id string, user *entities.User) error {
	return nil
}
func (m *MockUserRepoForRegistration) DeleteUserByID(ctx context.Context, id string) error {
	return nil
}
func (m *MockUserRepoForRegistration) SetLastSeen(ctx context.Context, id string, timestamp time.Time) error {
	return nil
}
func (m *MockUserRepoForRegistration) UpdateUserRole(ctx context.Context, adminID, targetID, role string) error {
	return nil
}
func (m *MockUserRepoForRegistration) ChangePassword(ctx context.Context, id string, newHashedPassword string) error {
	return nil
}
func (m *MockUserRepoForRegistration) ChangeEmail(ctx context.Context, email string, newEmail string) error {
	return nil
}
func (m *MockUserRepoForRegistration) SetRole(ctx context.Context, id string, role string) error {
	return nil
}
func (m *MockUserRepoForRegistration) ActivateUserByID(ctx context.Context, id string) error {
	return nil
}
func (m *MockUserRepoForRegistration) DeactivateUserByID(ctx context.Context, id string) error {
	return nil
}

type MockTokenRepoForRegistration struct {
	mock.Mock
}

func (m *MockTokenRepoForRegistration) StoreRefreshToken(ctx context.Context, token *entities.RefreshToken) error {
	return nil
}
func (m *MockTokenRepoForRegistration) FindByHash(ctx context.Context, hash string) (*entities.RefreshToken, error) {
	return nil, nil
}
func (m *MockTokenRepoForRegistration) DeleteByHash(ctx context.Context, hash string) error {
	return nil
}
func (m *MockTokenRepoForRegistration) DeleteAllByUserID(ctx context.Context, userID string) error {
	return nil
}

func TestRegistration_FirstUser_BecomesSuperAdmin(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepoForRegistration)
	mockTokenRepo := new(MockTokenRepoForRegistration)
	cfg := &config.Config{}

	// Create UserServices with mocks
	userServices := &UserServices{
		userRepo:  mockUserRepo,
		tokenRepo: mockTokenRepo,
		cfg:       cfg,
	}

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
		return user.Role == "superadmin" && user.Username == "admin" && !user.Activated
	})).Return(expectedUserID, nil)

	// Execute
	result, err := userServices.Register(context.Background(), userDTO)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedUserID, result)

	// Verify mocks were called
	mockUserRepo.AssertExpectations(t)
}

func TestRegistration_RegularUser_BecomesUnverified(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepoForRegistration)
	mockTokenRepo := new(MockTokenRepoForRegistration)
	cfg := &config.Config{}

	userServices := &UserServices{
		userRepo:  mockUserRepo,
		tokenRepo: mockTokenRepo,
		cfg:       cfg,
	}

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
	result, err := userServices.Register(context.Background(), userDTO)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedUserID, result)

	// Verify mocks were called
	mockUserRepo.AssertExpectations(t)
}

func TestRegistration_GetUserByID_Success(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepoForRegistration)
	mockTokenRepo := new(MockTokenRepoForRegistration)
	cfg := &config.Config{}

	userServices := &UserServices{
		userRepo:  mockUserRepo,
		tokenRepo: mockTokenRepo,
		cfg:       cfg,
	}

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
	result, err := userServices.GetUserByID(context.Background(), userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, userID, result.ID)
	assert.Equal(t, "testuser", result.Username)
	assert.Equal(t, "test@example.com", result.Email)

	// Verify mock was called
	mockUserRepo.AssertExpectations(t)
}

func TestRegistration_ValidationErrors(t *testing.T) {
	// Setup
	cfg := &config.Config{}

	// Test cases for validation errors
	testCases := []struct {
		name    string
		userDTO *UserDTO
		wantErr bool
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
			wantErr: true,
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
			wantErr: true,
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
			wantErr: true,
		},
		{
			name: "Valid input",
			userDTO: &UserDTO{
				Username:  "testuser",
				Email:     "test@example.com",
				Password:  "password123",
				FirstName: "Test",
				LastName:  "User",
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create fresh mocks for each test case
			freshMockUserRepo := new(MockUserRepoForRegistration)
			freshMockTokenRepo := new(MockTokenRepoForRegistration)
			
			freshUserServices := &UserServices{
				userRepo:  freshMockUserRepo,
				tokenRepo: freshMockTokenRepo,
				cfg:       cfg,
			}

			if !tc.wantErr {
				// Setup mocks for valid case
				freshMockUserRepo.On("CheckUsername", mock.Anything, tc.userDTO.Username).Return(false, nil)
				freshMockUserRepo.On("CheckEmail", mock.Anything, tc.userDTO.Email).Return(false, nil)
				freshMockUserRepo.On("CountAllUsers", mock.Anything).Return(int64(1), nil)
				freshMockUserRepo.On("CreateUser", mock.Anything, mock.AnythingOfType("*entities.User")).Return("test-id", nil)
			} else {
				// For error cases, we might still need to mock some calls depending on where validation fails
				if tc.userDTO.Username != "" && tc.userDTO.Email != "" {
					// These will be called before validation fails
					freshMockUserRepo.On("CheckUsername", mock.Anything, tc.userDTO.Username).Return(false, nil)
					freshMockUserRepo.On("CheckEmail", mock.Anything, tc.userDTO.Email).Return(false, nil)
				}
			}

			// Execute
			result, err := freshUserServices.Register(context.Background(), tc.userDTO)

			// Assert
			if tc.wantErr {
				assert.Error(t, err)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, result)
			}
		})
	}
}