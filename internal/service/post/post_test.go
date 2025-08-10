package postsvc

import (
	"context"
	"testing"

	"anchor-blog/internal/domain/entities"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock PostRepository for testing
type MockPostRepository struct {
	mock.Mock
}

func (m *MockPostRepository) Create(ctx context.Context, post *entities.Post) (*entities.Post, error) {
	args := m.Called(ctx, post)
	return args.Get(0).(*entities.Post), args.Error(1)
}

func (m *MockPostRepository) FindByID(ctx context.Context, id string) (*entities.Post, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entities.Post), args.Error(1)
}

func (m *MockPostRepository) FindAll(ctx context.Context, opts entities.PaginationOptions) ([]*entities.Post, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]*entities.Post), args.Error(1)
}

func (m *MockPostRepository) Update(ctx context.Context, id string, post *entities.Post) (*entities.Post, error) {
	args := m.Called(ctx, id, post)
	return args.Get(0).(*entities.Post), args.Error(1)
}

func (m *MockPostRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPostRepository) SearchByTitle(ctx context.Context, query string, opts entities.PaginationOptions) ([]*entities.Post, error) {
	args := m.Called(ctx, query, opts)
	return args.Get(0).([]*entities.Post), args.Error(1)
}

func (m *MockPostRepository) SearchByAuthor(ctx context.Context, query string, opts entities.PaginationOptions) ([]*entities.Post, error) {
	args := m.Called(ctx, query, opts)
	return args.Get(0).([]*entities.Post), args.Error(1)
}

func (m *MockPostRepository) FilterByTags(ctx context.Context, tags []string, opts entities.PaginationOptions) ([]*entities.Post, error) {
	args := m.Called(ctx, tags, opts)
	return args.Get(0).([]*entities.Post), args.Error(1)
}

func (m *MockPostRepository) FilterByDateRange(ctx context.Context, startDate, endDate string, opts entities.PaginationOptions) ([]*entities.Post, error) {
	args := m.Called(ctx, startDate, endDate, opts)
	return args.Get(0).([]*entities.Post), args.Error(1)
}

func (m *MockPostRepository) AddLike(ctx context.Context, postID, userID string) error {
	args := m.Called(ctx, postID, userID)
	return args.Error(0)
}

func (m *MockPostRepository) RemoveLike(ctx context.Context, postID, userID string) error {
	args := m.Called(ctx, postID, userID)
	return args.Error(0)
}

func (m *MockPostRepository) AddDislike(ctx context.Context, postID, userID string) error {
	args := m.Called(ctx, postID, userID)
	return args.Error(0)
}

func (m *MockPostRepository) RemoveDislike(ctx context.Context, postID, userID string) error {
	args := m.Called(ctx, postID, userID)
	return args.Error(0)
}

func (m *MockPostRepository) GetLikeStatus(ctx context.Context, postID, userID string) (bool, bool, error) {
	args := m.Called(ctx, postID, userID)
	return args.Bool(0), args.Bool(1), args.Error(2)
}

func (m *MockPostRepository) IncrementViewCount(ctx context.Context, postID string) error {
	args := m.Called(ctx, postID)
	return args.Error(0)
}

func (m *MockPostRepository) GetViewCount(ctx context.Context, postID string) (int, error) {
	args := m.Called(ctx, postID)
	return args.Int(0), args.Error(1)
}

func (m *MockPostRepository) GetTotalViews(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockPostRepository) GetPostsByViewCount(ctx context.Context, limit int) ([]*entities.Post, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]*entities.Post), args.Error(1)
}

func (m *MockPostRepository) ResetViewCount(ctx context.Context, postID string) error {
	args := m.Called(ctx, postID)
	return args.Error(0)
}

func TestPostService_CreatePost_Success(t *testing.T) {
	// Setup
	mockRepo := new(MockPostRepository)
	service := NewPostService(mockRepo)

	// Test data
	title := "Test Post"
	content := "This is a test post content"
	authorID := "author-123"
	tags := []string{"test", "golang"}

	expectedPost := &entities.Post{
		ID:       "post-123",
		Title:    title,
		Content:  content,
		AuthorID: authorID,
		Tags:     tags,
	}

	// Mock expectations
	mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(post *entities.Post) bool {
		return post.Title == title && post.Content == content && post.AuthorID == authorID
	})).Return(expectedPost, nil)

	// Execute
	result, err := service.CreatePost(context.Background(), title, content, authorID, tags)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedPost.ID, result.ID)
	assert.Equal(t, title, result.Title)
	assert.Equal(t, content, result.Content)
	assert.Equal(t, authorID, result.AuthorID)
	assert.Equal(t, tags, result.Tags)

	mockRepo.AssertExpectations(t)
}

func TestPostService_GetPostByID_Success(t *testing.T) {
	// Setup
	mockRepo := new(MockPostRepository)
	service := NewPostService(mockRepo)

	// Test data
	postID := "post-123"
	expectedPost := &entities.Post{
		ID:       postID,
		Title:    "Test Post",
		Content:  "Test content",
		AuthorID: "author-123",
	}

	// Mock expectations
	mockRepo.On("FindByID", mock.Anything, postID).Return(expectedPost, nil)

	// Execute
	result, err := service.GetPostByID(context.Background(), postID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedPost.ID, result.ID)
	assert.Equal(t, expectedPost.Title, result.Title)

	mockRepo.AssertExpectations(t)
}

func TestPostService_GetPostByID_NotFound(t *testing.T) {
	// Setup
	mockRepo := new(MockPostRepository)
	service := NewPostService(mockRepo)

	// Test data
	postID := "nonexistent-post"

	// Mock expectations
	mockRepo.On("FindByID", mock.Anything, postID).Return((*entities.Post)(nil), assert.AnError)

	// Execute
	result, err := service.GetPostByID(context.Background(), postID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
}

func TestPostService_ListPosts_Success(t *testing.T) {
	// Setup
	mockRepo := new(MockPostRepository)
	service := NewPostService(mockRepo)

	// Test data
	page := int64(1)
	limit := int64(10)
	expectedPosts := []*entities.Post{
		{ID: "post-1", Title: "Post 1"},
		{ID: "post-2", Title: "Post 2"},
	}

	// Mock expectations
	mockRepo.On("FindAll", mock.Anything, entities.PaginationOptions{
		Page:  page,
		Limit: limit,
	}).Return(expectedPosts, nil)

	// Execute
	result, err := service.ListPosts(context.Background(), page, limit)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, "post-1", result[0].ID)
	assert.Equal(t, "post-2", result[1].ID)

	mockRepo.AssertExpectations(t)
}

func TestPostService_ListPosts_DefaultPagination(t *testing.T) {
	// Setup
	mockRepo := new(MockPostRepository)
	service := NewPostService(mockRepo)

	// Test data - invalid pagination values
	page := int64(0)  // Should default to 1
	limit := int64(0) // Should default to 10

	expectedPosts := []*entities.Post{
		{ID: "post-1", Title: "Post 1"},
	}

	// Mock expectations - should use default values
	mockRepo.On("FindAll", mock.Anything, entities.PaginationOptions{
		Page:  1,  // Default
		Limit: 10, // Default
	}).Return(expectedPosts, nil)

	// Execute
	result, err := service.ListPosts(context.Background(), page, limit)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)

	mockRepo.AssertExpectations(t)
}

func TestPostService_UpdatePost_Success(t *testing.T) {
	// Setup
	mockRepo := new(MockPostRepository)
	service := NewPostService(mockRepo)

	// Test data
	postID := "post-123"
	title := "Updated Title"
	content := "Updated content"
	tags := []string{"updated", "test"}

	updatedPost := &entities.Post{
		ID:      postID,
		Title:   title,
		Content: content,
		Tags:    tags,
	}

	// Mock expectations
	mockRepo.On("Update", mock.Anything, postID, mock.MatchedBy(func(post *entities.Post) bool {
		return post.Title == title && post.Content == content
	})).Return(updatedPost, nil)

	// Execute
	result, err := service.UpdatePost(context.Background(), postID, title, content, tags)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, postID, result.ID)
	assert.Equal(t, title, result.Title)
	assert.Equal(t, content, result.Content)
	assert.Equal(t, tags, result.Tags)

	mockRepo.AssertExpectations(t)
}

func TestPostService_DeletePost_Success(t *testing.T) {
	// Setup
	mockRepo := new(MockPostRepository)
	service := NewPostService(mockRepo)

	// Test data
	postID := "post-123"

	// Mock expectations
	mockRepo.On("Delete", mock.Anything, postID).Return(nil)

	// Execute
	err := service.DeletePost(context.Background(), postID)

	// Assert
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestPostService_SearchPosts_ByTitle(t *testing.T) {
	// Setup
	mockRepo := new(MockPostRepository)
	service := NewPostService(mockRepo)

	// Test data
	query := "golang"
	searchType := "title"
	page := int64(1)
	limit := int64(10)

	expectedPosts := []*entities.Post{
		{ID: "post-1", Title: "Golang Tutorial"},
	}

	// Mock expectations
	mockRepo.On("SearchByTitle", mock.Anything, query, entities.PaginationOptions{
		Page:  page,
		Limit: limit,
	}).Return(expectedPosts, nil)

	// Execute
	result, err := service.SearchPosts(context.Background(), query, searchType, page, limit)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 1)
	assert.Equal(t, "post-1", result[0].ID)

	mockRepo.AssertExpectations(t)
}

func TestPostService_SearchPosts_ByAuthor(t *testing.T) {
	// Setup
	mockRepo := new(MockPostRepository)
	service := NewPostService(mockRepo)

	// Test data
	query := "john"
	searchType := "author"
	page := int64(1)
	limit := int64(10)

	expectedPosts := []*entities.Post{
		{ID: "post-1", Title: "Post by John"},
	}

	// Mock expectations
	mockRepo.On("SearchByAuthor", mock.Anything, query, entities.PaginationOptions{
		Page:  page,
		Limit: limit,
	}).Return(expectedPosts, nil)

	// Execute
	result, err := service.SearchPosts(context.Background(), query, searchType, page, limit)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 1)

	mockRepo.AssertExpectations(t)
}

func TestPostService_LikePost_Success(t *testing.T) {
	// Setup
	mockRepo := new(MockPostRepository)
	service := NewPostService(mockRepo)

	// Test data
	postID := "post-123"
	userID := "user-123"

	// Mock expectations
	mockRepo.On("AddLike", mock.Anything, postID, userID).Return(nil)

	// Execute
	err := service.LikePost(context.Background(), postID, userID)

	// Assert
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestPostService_GetPostLikeStatus_Success(t *testing.T) {
	// Setup
	mockRepo := new(MockPostRepository)
	service := NewPostService(mockRepo)

	// Test data
	postID := "post-123"
	userID := "user-123"

	// Mock expectations
	mockRepo.On("GetLikeStatus", mock.Anything, postID, userID).Return(true, false, nil)

	// Execute
	liked, disliked, err := service.GetPostLikeStatus(context.Background(), postID, userID)

	// Assert
	assert.NoError(t, err)
	assert.True(t, liked)
	assert.False(t, disliked)

	mockRepo.AssertExpectations(t)
}