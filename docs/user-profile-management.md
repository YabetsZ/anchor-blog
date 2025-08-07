# User Profile Management

This document describes the User Profile Management feature implementation.

## Overview

The User Profile Management feature allows authenticated users to view and update their profile information including:
- Bio
- Profile picture URL
- Social media links

## API Endpoints

### GET /api/v1/user/profile
Retrieves the current user's profile information.

**Authentication**: Required (Bearer token)

**Response**:
```json
{
  "success": true,
  "data": {
    "bio": "Software developer passionate about Go and web technologies",
    "picture_url": "https://example.com/profile.jpg",
    "social_links": [
      {
        "platform": "twitter",
        "url": "https://twitter.com/username"
      },
      {
        "platform": "github",
        "url": "https://github.com/username"
      }
    ]
  }
}
```

### PUT /api/v1/user/profile
Updates the current user's profile information.

**Authentication**: Required (Bearer token)

**Request Body**:
```json
{
  "bio": "Updated bio text",
  "picture_url": "https://example.com/new-profile.jpg",
  "social_links": [
    {
      "platform": "twitter",
      "url": "https://twitter.com/newusername"
    },
    {
      "platform": "linkedin",
      "url": "https://linkedin.com/in/username"
    }
  ]
}
```

**Response**:
```json
{
  "success": true,
  "message": "Profile updated successfully",
  "data": {
    "bio": "Updated bio text",
    "picture_url": "https://example.com/new-profile.jpg",
    "social_links": [
      {
        "platform": "twitter",
        "url": "https://twitter.com/newusername"
      },
      {
        "platform": "linkedin",
        "url": "https://linkedin.com/in/username"
      }
    ]
  }
}
```

## Implementation Details

### Model Layer
- **Entity**: `entities.User` with `Profile` field containing `UserProfile` struct
- **Repository Model**: MongoDB model with proper BSON tags
- **Mapping**: Conversion functions between entity and repository models

### Service Layer
- **ProfileService**: Handles business logic for profile operations
- **Methods**:
  - `GetUserProfile(ctx, userID)`: Retrieves user profile
  - `UpdateUserProfile(ctx, userID, request)`: Updates user profile with partial updates

### Handler Layer
- **UserHandler**: Extended with profile methods
- **Methods**:
  - `GetProfile(c *gin.Context)`: HTTP handler for GET profile
  - `UpdateProfile(c *gin.Context)`: HTTP handler for PUT profile

### Repository Layer
- **UserRepository**: Uses existing `EditUserByID` method
- **Partial Updates**: Only updates provided fields, preserves existing data

## Data Structure

### UserProfile Entity
```go
type UserProfile struct {
    Bio         string
    PictureURL  string
    SocialLinks []SocialLink
}

type SocialLink struct {
    Platform string
    URL      string
}
```

### Update Request
```go
type UpdateProfileRequest struct {
    Bio         *string          `json:"bio,omitempty"`
    PictureURL  *string          `json:"picture_url,omitempty"`
    SocialLinks *[]SocialLinkDTO `json:"social_links,omitempty"`
}
```

## Security Considerations

1. **Authentication**: All profile endpoints require valid JWT token
2. **Authorization**: Users can only access/modify their own profiles
3. **Input Validation**: Request data is validated before processing
4. **Partial Updates**: Only provided fields are updated, preventing accidental data loss

## Testing

Unit tests are provided for the ProfileService:
- `TestProfileService_GetUserProfile`: Tests profile retrieval
- `TestProfileService_UpdateUserProfile`: Tests profile updates

Run tests with:
```bash
go test ./internal/service/user -v
```

## Usage Examples

### Get Profile
```bash
curl -X GET http://localhost:8080/api/v1/user/profile \
  -H "Authorization: Bearer <your-jwt-token>"
```

### Update Profile
```bash
curl -X PUT http://localhost:8080/api/v1/user/profile \
  -H "Authorization: Bearer <your-jwt-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "bio": "Full-stack developer with 5+ years experience",
    "picture_url": "https://example.com/my-photo.jpg",
    "social_links": [
      {
        "platform": "github",
        "url": "https://github.com/myusername"
      }
    ]
  }'
```