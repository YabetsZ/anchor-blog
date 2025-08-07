package usersvc

import (
	"anchor-blog/internal/domain/entities"
	errorr "anchor-blog/internal/errors"
	"context"
	"testing"
	"time"
)

// Mock repository for testing
type mockUserRepository struct {
	users map[string]*entities.User
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users: make(map[string]*entities.User),
	}
}

func (m *mockUserRepository) GetUserByID(ctx context.Context, id string) (*entities.User, error) {
	if user, exists := m.users[id]; exists {
		return user, nil
	}
	return nil, errorr.ErrUserNotFound
}

func (m *mockUserRepository) EditUserByID(ctx context.Context, id string, user *entities.User) error {
	if existingUser, exists := m.users[id]; exists {
		// Update only non-empty fields
		if user.Profile.Bio != "" {
			existingUser.Profile.Bio = user.Profile.Bio
		}
		if user.Profile.PictureURL != "" {
			existingUser.Profile.PictureURL = user.Profile.PictureURL
		}
		if len(user.Profile.SocialLinks) > 0 {
			existingUser.Profile.SocialLinks = user.Profile.SocialLinks
		}
		existingUser.UpdatedAt = user.UpdatedAt
		return nil
	}
	return errorr.ErrUserNotFound
}

// Implement other required methods (not used in tests)
func (m *mockUserRepository) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) { return nil, nil }
func (m *mockUserRepository) GetUserByUsername(ctx context.Context, username string) (*entities.User, error) { return nil, nil }
func (m *mockUserRepository) GetUsers(ctx context.Context, limit, offset int64) ([]*entities.User, error) { return nil, nil }
func (m *mockUserRepository) CountUsersByRole(ctx context.Context, role string) (int64, error) { return 0, nil }
func (m *mockUserRepository) CountAllUsers(ctx context.Context) (int64, error) { return 0, nil }
func (m *mockUserRepository) CountActiveUsers(ctx context.Context) (int64, error) { return 0, nil }
func (m *mockUserRepository) CountInactiveUsers(ctx context.Context) (int64, error) { return 0, nil }
func (m *mockUserRepository) GetUserRoleByID(ctx context.Context, userID string) (string, error) { return "", nil }
func (m *mockUserRepository) CreateUser(ctx context.Context, user *entities.User) (string, error) { return "", nil }
func (m *mockUserRepository) DeleteUserByID(ctx context.Context, id string) error { return nil }
func (m *mockUserRepository) SetLastSeen(ctx context.Context, id string, timestamp time.Time) error { return nil }
func (m *mockUserRepository) CheckEmail(ctx context.Context, email string) (bool, error) { return false, nil }
func (m *mockUserRepository) CheckUsername(ctx context.Context, username string) (bool, error) { return false, nil }
func (m *mockUserRepository) ChangePassword(ctx context.Context, id string, newHashedPassword string) error { return nil }
func (m *mockUserRepository) ChangeEmail(ctx context.Context, email string, newEmail string) error { return nil }
func (m *mockUserRepository) SetRole(ctx context.Context, id string, role string) error { return nil }
func (m *mockUserRepository) ActivateUserByID(ctx context.Context, id string) error { return nil }
func (m *mockUserRepository) DeactivateUserByID(ctx context.Context, id string) error { return nil }

func TestProfileService_GetUserProfile(t *testing.T) {
	mockRepo := newMockUserRepository()
	profileService := NewProfileService(mockRepo)

	// Setup test user
	userID := "test-user-id"
	testUser := &entities.User{
		ID: userID,
		Profile: entities.UserProfile{
			Bio:        "Test bio",
			PictureURL: "https://example.com/pic.jpg",
			SocialLinks: []entities.SocialLink{
				{Platform: "twitter", URL: "https://twitter.com/test"},
				{Platform: "github", URL: "https://github.com/test"},
			},
		},
	}
	mockRepo.users[userID] = testUser

	// Test GetUserProfile
	profile, err := profileService.GetUserProfile(context.Background(), userID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if profile.Bio != "Test bio" {
		t.Errorf("Expected bio 'Test bio', got '%s'", profile.Bio)
	}

	if profile.PictureURL != "https://example.com/pic.jpg" {
		t.Errorf("Expected picture URL 'https://example.com/pic.jpg', got '%s'", profile.PictureURL)
	}

	if len(profile.SocialLinks) != 2 {
		t.Errorf("Expected 2 social links, got %d", len(profile.SocialLinks))
	}
}

func TestProfileService_UpdateUserProfile(t *testing.T) {
	mockRepo := newMockUserRepository()
	profileService := NewProfileService(mockRepo)

	// Setup test user
	userID := "test-user-id"
	testUser := &entities.User{
		ID: userID,
		Profile: entities.UserProfile{
			Bio:        "Old bio",
			PictureURL: "https://example.com/old.jpg",
			SocialLinks: []entities.SocialLink{
				{Platform: "twitter", URL: "https://twitter.com/old"},
			},
		},
	}
	mockRepo.users[userID] = testUser

	// Test UpdateUserProfile
	newBio := "Updated bio"
	updateReq := &UpdateProfileRequest{
		Bio: &newBio,
	}

	profile, err := profileService.UpdateUserProfile(context.Background(), userID, updateReq)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if profile.Bio != "Updated bio" {
		t.Errorf("Expected bio 'Updated bio', got '%s'", profile.Bio)
	}

	// Verify other fields remain unchanged
	if profile.PictureURL != "https://example.com/old.jpg" {
		t.Errorf("Expected picture URL to remain unchanged, got '%s'", profile.PictureURL)
	}
}