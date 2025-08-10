package usersvc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Simple test to verify testing framework is working
func TestSimple_TestingFramework(t *testing.T) {
	// Test basic assertions
	assert.True(t, true)
	assert.Equal(t, "test", "test")
	assert.NotNil(t, &UserDTO{})
	
	// Test UserDTO creation
	dto := &UserDTO{
		Username:  "testuser",
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
	}
	
	assert.Equal(t, "testuser", dto.Username)
	assert.Equal(t, "test@example.com", dto.Email)
	assert.Equal(t, "Test", dto.FirstName)
	assert.Equal(t, "User", dto.LastName)
}

func TestDTOToEntity_Conversion(t *testing.T) {
	// Test DTO to Entity conversion
	dto := UserDTO{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
	}
	
	entity := DTOToEntity(dto)
	
	assert.Equal(t, dto.Username, entity.Username)
	assert.Equal(t, dto.Email, entity.Email)
	assert.Equal(t, dto.FirstName, entity.FirstName)
	assert.Equal(t, dto.LastName, entity.LastName)
	// DTOToEntity copies the password to PasswordHash field (though it should be hashed in Register method)
	assert.Equal(t, dto.Password, entity.PasswordHash)
	// Note: This is the current behavior, though ideally password should be hashed before storing
}