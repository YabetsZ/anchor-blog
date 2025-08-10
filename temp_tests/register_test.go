package user

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"anchor-blog/internal/domain/entities"
	usersvc "anchor-blog/internal/service/user"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock UserServices
type MockUserServices struct {
	mock.Mock
}

func (m *MockUserServices) Register(ctx context.Context, userDto *usersvc.UserDTO) (string, error) {
	args := m.Called(ctx, userDto)
	return args.String(0), args.Error(1)
}

func (m *MockUserServices) GetUserByID(ctx context.Context, userID string) (*entities.User, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*entities.User), args.Error(1)
}

// Mock ActivationService
type MockActivationService struct {
	mock.Mock
}

func (m *MockActivationService) SendActivationEmail(ctx context.Context, user *entities.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockActivationService) VerifyActivation(ctx context.Context, token string) (*entities.User, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(*entities.User), args.Error(1)
}

func TestUserHandler_Register_Success(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	
	mockUserService := new(MockUserServices)
	mockActivationService := new(MockActivationService)
	
	handler := &UserHandler{
		UserService:       mockUserService,
		ActivationService: mockActivationService,
	}

	// Test data
	userID := "test-user-id-123"
	testUser := &entities.User{
		ID:        userID,
		Username:  "testuser",
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
		Role:      "unverified",
		Activated: false,
	}

	registerRequest := usersvc.UserDTO{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
	}

	// Mock expectations
	mockUserService.On("Register", mock.Anything, &registerRequest).Return(userID, nil)
	mockUserService.On("GetUserByID", mock.Anything, userID).Return(testUser, nil)
	mockActivationService.On("SendActivationEmail", mock.Anything, testUser).Return(nil)

	// Create request
	requestBody, _ := json.Marshal(registerRequest)
	req, _ := http.NewRequest("POST", "/api/v1/user/register", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Execute
	handler.Register(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response registerResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, userID, response.ID)

	// Verify mocks were called
	mockUserService.AssertExpectations(t)
	mockActivationService.AssertExpectations(t)
}

func TestUserHandler_Register_SuperAdmin_NoActivation(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	
	mockUserService := new(MockUserServices)
	mockActivationService := new(MockActivationService)
	
	handler := &UserHandler{
		UserService:       mockUserService,
		ActivationService: mockActivationService,
	}

	// Test data - superadmin user
	userID := "superadmin-id-123"
	testUser := &entities.User{
		ID:        userID,
		Username:  "admin",
		Email:     "admin@example.com",
		FirstName: "Super",
		LastName:  "Admin",
		Role:      "superadmin", // This should skip activation
		Activated: true,
	}

	registerRequest := usersvc.UserDTO{
		Username:  "admin",
		Email:     "admin@example.com",
		Password:  "password123",
		FirstName: "Super",
		LastName:  "Admin",
	}

	// Mock expectations - no activation service call expected for superadmin
	mockUserService.On("Register", mock.Anything, &registerRequest).Return(userID, nil)
	mockUserService.On("GetUserByID", mock.Anything, userID).Return(testUser, nil)
	// Note: No expectation for SendActivationEmail since it shouldn't be called

	// Create request
	requestBody, _ := json.Marshal(registerRequest)
	req, _ := http.NewRequest("POST", "/api/v1/user/register", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Execute
	handler.Register(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response registerResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, userID, response.ID)

	// Verify mocks were called (activation service should NOT be called)
	mockUserService.AssertExpectations(t)
	mockActivationService.AssertNotCalled(t, "SendActivationEmail")
}

func TestUserHandler_Register_InvalidInput(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	
	mockUserService := new(MockUserServices)
	mockActivationService := new(MockActivationService)
	
	handler := &UserHandler{
		UserService:       mockUserService,
		ActivationService: mockActivationService,
	}

	// Create invalid request (missing required fields)
	invalidRequest := map[string]interface{}{
		"username": "testuser",
		// Missing email, password, first_name, last_name
	}

	requestBody, _ := json.Marshal(invalidRequest)
	req, _ := http.NewRequest("POST", "/api/v1/user/register", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Execute
	handler.Register(c)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Verify no service methods were called
	mockUserService.AssertNotCalled(t, "Register")
	mockActivationService.AssertNotCalled(t, "SendActivationEmail")
}

func TestUserHandler_Register_ActivationFailure_StillSucceeds(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	
	mockUserService := new(MockUserServices)
	mockActivationService := new(MockActivationService)
	
	handler := &UserHandler{
		UserService:       mockUserService,
		ActivationService: mockActivationService,
	}

	// Test data
	userID := "test-user-id-123"
	testUser := &entities.User{
		ID:        userID,
		Username:  "testuser",
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
		Role:      "unverified",
		Activated: false,
	}

	registerRequest := usersvc.UserDTO{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
	}

	// Mock expectations - activation fails but registration should still succeed
	mockUserService.On("Register", mock.Anything, &registerRequest).Return(userID, nil)
	mockUserService.On("GetUserByID", mock.Anything, userID).Return(testUser, nil)
	mockActivationService.On("SendActivationEmail", mock.Anything, testUser).Return(assert.AnError)

	// Create request
	requestBody, _ := json.Marshal(registerRequest)
	req, _ := http.NewRequest("POST", "/api/v1/user/register", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Execute
	handler.Register(c)

	// Assert - registration should still succeed even if activation fails
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Failed to send activation email", w.Header().Get("X-Activation-Warning"))
	
	var response registerResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, userID, response.ID)

	// Verify mocks were called
	mockUserService.AssertExpectations(t)
	mockActivationService.AssertExpectations(t)
}