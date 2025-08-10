package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"anchor-blog/internal/domain/entities"
	usersvc "anchor-blog/internal/service/user"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock ActivationService for handler tests
type MockActivationServiceForHandler struct {
	mock.Mock
}

func (m *MockActivationServiceForHandler) VerifyActivation(ctx interface{}, token string) (*entities.User, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockActivationServiceForHandler) SendActivationEmail(ctx interface{}, user *entities.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func TestActivationHandler_ActivateAccount_Success(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	
	mockActivationService := new(MockActivationServiceForHandler)
	handler := &ActivationHandler{
		activationService: (*usersvc.ActivationService)(mockActivationService),
	}

	// Test data
	token := "valid-activation-token-123"
	activatedUser := &entities.User{
		ID:        "user-123",
		Username:  "testuser",
		Email:     "test@example.com",
		Activated: true,
		Role:      "user",
	}

	// Mock expectations
	mockActivationService.On("VerifyActivation", mock.Anything, token).Return(activatedUser, nil)

	// Create request
	req, _ := http.NewRequest("GET", "/api/v1/users/activate?token="+token, nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Execute
	handler.ActivateAccount(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	
	// Check response contains success message and user data
	responseBody := w.Body.String()
	assert.Contains(t, responseBody, "Account activated successfully")
	assert.Contains(t, responseBody, "testuser")
	assert.Contains(t, responseBody, "test@example.com")
	assert.Contains(t, responseBody, "user-123")

	// Verify mock was called
	mockActivationService.AssertExpectations(t)
}

func TestActivationHandler_ActivateAccount_MissingToken(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	
	mockActivationService := new(MockActivationServiceForHandler)
	handler := &ActivationHandler{
		activationService: (*usersvc.ActivationService)(mockActivationService),
	}

	// Create request without token
	req, _ := http.NewRequest("GET", "/api/v1/users/activate", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Execute
	handler.ActivateAccount(c)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	// Check response contains error message
	responseBody := w.Body.String()
	assert.Contains(t, responseBody, "Activation token is required")

	// Verify no service methods were called
	mockActivationService.AssertNotCalled(t, "VerifyActivation")
}

func TestActivationHandler_ActivateAccount_InvalidToken(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	
	mockActivationService := new(MockActivationServiceForHandler)
	handler := &ActivationHandler{
		activationService: (*usersvc.ActivationService)(mockActivationService),
	}

	// Test data
	token := "invalid-token-123"

	// Mock expectations - verification fails
	mockActivationService.On("VerifyActivation", mock.Anything, token).Return((*entities.User)(nil), assert.AnError)

	// Create request
	req, _ := http.NewRequest("GET", "/api/v1/users/activate?token="+token, nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Execute
	handler.ActivateAccount(c)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	// Check response contains error message
	responseBody := w.Body.String()
	assert.Contains(t, responseBody, "Invalid or expired activation token")

	// Verify mock was called
	mockActivationService.AssertExpectations(t)
}

func TestActivationHandler_ActivateAccount_EmptyToken(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	
	mockActivationService := new(MockActivationServiceForHandler)
	handler := &ActivationHandler{
		activationService: (*usersvc.ActivationService)(mockActivationService),
	}

	// Create request with empty token
	req, _ := http.NewRequest("GET", "/api/v1/users/activate?token=", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Execute
	handler.ActivateAccount(c)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	// Check response contains error message
	responseBody := w.Body.String()
	assert.Contains(t, responseBody, "Activation token is required")

	// Verify no service methods were called
	mockActivationService.AssertNotCalled(t, "VerifyActivation")
}