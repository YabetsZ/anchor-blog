# Anchor Blog API Documentation

Complete API reference for testing with Postman. Base URL: `http://localhost:8080`

## 📋 Table of Contents
- [Authentication](#authentication)
- [Health Check](#health-check)
- [User Management](#user-management)
- [User Profile](#user-profile)
- [Account Activation](#account-activation)
- [Password Reset](#password-reset)
- [Blog Posts](#blog-posts)
- [AI Content Generation](#ai-content-generation)

---

## 🔐 Authentication

Most endpoints require JWT authentication. Include the token in the Authorization header:
```
Authorization: Bearer <your-jwt-token>
```

---

## ❤️ Health Check

### GET /health
Check if the API is running.

**Request:**
```http
GET /health
```

**Response:**
```json
{
  "status": "OK"
}
```

---

## 👤 User Management

### POST /api/v1/user/register
Register a new user account.

**Request:**
```http
POST /api/v1/user/register
Content-Type: application/json

{
  "username": "johndoe",
  "first_name": "John",
  "last_name": "Doe",
  "email": "john.doe@example.com",
  "password": "securepassword123",
  "role": "user",
  "profile": {
    "bio": "Software developer",
    "picture_url": "https://example.com/avatar.jpg",
    "social_links": [
      {
        "platform": "github",
        "url": "https://github.com/johndoe"
      }
    ]
  }
}
```

**Response:**
```json
{
  "id": "507f1f77bcf86cd799439011"
}
```

### POST /api/v1/user/login
Authenticate user and get access tokens.

**Request:**
```http
POST /api/v1/user/login
Content-Type: application/json

{
  "username": "johndoe",
  "password": "securepassword123"
}
```

**Response:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### POST /api/v1/refresh
Refresh access token using refresh token.

**Request:**
```http
POST /api/v1/refresh
Authorization: Bearer <refresh-token>
```

**Response:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

---

## 👤 User Profile

### GET /api/v1/user/profile
Get current user's profile information.

**Request:**
```http
GET /api/v1/user/profile
Authorization: Bearer <access-token>
```

**Response:**
```json
{
  "success": true,
  "data": {
    "bio": "Software developer passionate about Go and web technologies",
    "picture_url": "https://example.com/profile.jpg",
    "social_links": [
      {
        "platform": "github",
        "url": "https://github.com/johndoe"
      },
      {
        "platform": "twitter",
        "url": "https://twitter.com/johndoe"
      }
    ]
  }
}
```

### PUT /api/v1/user/profile
Update current user's profile information.

**Request:**
```http
PUT /api/v1/user/profile
Authorization: Bearer <access-token>
Content-Type: application/json

{
  "bio": "Updated bio text",
  "picture_url": "https://example.com/new-profile.jpg",
  "social_links": [
    {
      "platform": "github",
      "url": "https://github.com/johndoe"
    },
    {
      "platform": "linkedin",
      "url": "https://linkedin.com/in/johndoe"
    }
  ]
}
```

**Response:**
```json
{
  "success": true,
  "message": "Profile updated successfully",
  "data": {
    "bio": "Updated bio text",
    "picture_url": "https://example.com/new-profile.jpg",
    "social_links": [
      {
        "platform": "github",
        "url": "https://github.com/johndoe"
      },
      {
        "platform": "linkedin",
        "url": "https://linkedin.com/in/johndoe"
      }
    ]
  }
}
```

---

## ✅ Account Activation

### GET /api/v1/users/activate
Activate user account using activation token.

**Request:**
```http
GET /api/v1/users/activate?token=<activation-token>
```

**Response:**
```json
{
  "message": "Account activated successfully",
  "user": {
    "id": "507f1f77bcf86cd799439011",
    "username": "johndoe",
    "email": "john.doe@example.com",
    "activated": true,
    "role": "user"
  }
}
```

---

## 🔐 Password Reset

### POST /api/v1/users/forgot-password
Request password reset email.

**Request:**
```http
POST /api/v1/users/forgot-password
Content-Type: application/json

{
  "email": "john.doe@example.com"
}
```

**Response:**
```json
{
  "message": "Password reset email sent successfully",
  "email": "john.doe@example.com"
}
```

### POST /api/v1/users/reset-password
Reset password using reset token.

**Request:**
```http
POST /api/v1/users/reset-password
Content-Type: application/json

{
  "token": "<reset-token>",
  "new_password": "newSecurePassword123"
}
```

**Response:**
```json
{
  "message": "Password reset successfully",
  "user": {
    "id": "507f1f77bcf86cd799439011",
    "username": "johndoe",
    "email": "john.doe@example.com"
  }
}
```

---

## 📝 Blog Posts

### POST /api/v1/posts
Create a new blog post.

**Request:**
```http
POST /api/v1/posts
Authorization: Bearer <access-token>
Content-Type: application/json

{
  "title": "Getting Started with Go",
  "content": "Go is a powerful programming language...",
  "tags": ["golang", "programming", "tutorial"]
}
```

**Response:**
```json
{
  "id": "507f1f77bcf86cd799439012",
  "title": "Getting Started with Go",
  "content": "Go is a powerful programming language...",
  "author_id": "507f1f77bcf86cd799439011",
  "tags": ["golang", "programming", "tutorial"],
  "view_count": 0,
  "likes": [],
  "dislikes": [],
  "created_at": "2025-08-07T10:30:00Z",
  "updated_at": "2025-08-07T10:30:00Z"
}
```

### GET /api/v1/posts/:id
Get a specific blog post by ID.

**Request:**
```http
GET /api/v1/posts/507f1f77bcf86cd799439012
```

**Response:**
```json
{
  "id": "507f1f77bcf86cd799439012",
  "title": "Getting Started with Go",
  "content": "Go is a powerful programming language...",
  "author_id": "507f1f77bcf86cd799439011",
  "tags": ["golang", "programming", "tutorial"],
  "view_count": 1,
  "likes": [],
  "dislikes": [],
  "created_at": "2025-08-07T10:30:00Z",
  "updated_at": "2025-08-07T10:30:00Z"
}
```

### GET /api/v1/posts
Get list of blog posts with pagination.

**Request:**
```http
GET /api/v1/posts?page=1&limit=10
```

**Response:**
```json
[
  {
    "id": "507f1f77bcf86cd799439012",
    "title": "Getting Started with Go",
    "content": "Go is a powerful programming language...",
    "author_id": "507f1f77bcf86cd799439011",
    "tags": ["golang", "programming", "tutorial"],
    "view_count": 5,
    "likes": ["507f1f77bcf86cd799439013"],
    "dislikes": [],
    "created_at": "2025-08-07T10:30:00Z",
    "updated_at": "2025-08-07T10:30:00Z"
  }
]
```

### GET /api/v1/posts/popular
Get popular posts ordered by view count.

**Request:**
```http
GET /api/v1/posts/popular?limit=10
```

**Response:**
```json
{
  "posts": [
    {
      "id": "507f1f77bcf86cd799439012",
      "title": "Most Popular Post",
      "content": "This post has the most views...",
      "author_id": "507f1f77bcf86cd799439011",
      "tags": ["popular", "trending"],
      "view_count": 1250,
      "likes": ["507f1f77bcf86cd799439013"],
      "dislikes": [],
      "created_at": "2025-08-07T10:30:00Z",
      "updated_at": "2025-08-07T10:30:00Z"
    }
  ],
  "count": 10
}
```

### GET /api/v1/posts/:id/views
Get view count for a specific post.

**Request:**
```http
GET /api/v1/posts/507f1f77bcf86cd799439012/views
```

**Response:**
```json
{
  "post_id": "507f1f77bcf86cd799439012",
  "view_count": 1250
}
```

### GET /api/v1/stats/views
Get total view statistics across all posts.

**Request:**
```http
GET /api/v1/stats/views
```

**Response:**
```json
{
  "total_views": 15420
}
```

---

## 🤖 AI Content Generation

### POST /api/v1/ai/generate
Generate AI-powered content for blog posts.

**Request:**
```http
POST /api/v1/ai/generate
Authorization: Bearer <access-token>
Content-Type: application/json

{
  "topic": "Introduction to Machine Learning for Beginners",
  "tone": "professional",
  "word_count": 150,
  "audience": ["developers", "general"],
  "scope": "Cover basic concepts and practical applications"
}
```

**Response:**
```json
{
  "title": "Machine Learning for Beginners: A Comprehensive Guide",
  "meta_description": "Learn the fundamentals of machine learning with this beginner-friendly guide covering key concepts, algorithms, and practical applications.",
  "outline": [
    {
      "heading": "What is Machine Learning?",
      "key_points": [
        "Definition and core concepts",
        "Types of machine learning",
        "Real-world applications"
      ]
    },
    {
      "heading": "Getting Started with ML",
      "key_points": [
        "Essential tools and libraries",
        "Setting up your environment",
        "First ML project walkthrough"
      ]
    }
  ],
  "body": [
    {
      "heading": "What is Machine Learning?",
      "paragraphs": [
        "Machine learning is a subset of artificial intelligence that enables computers to learn and make decisions from data without being explicitly programmed.",
        "This revolutionary technology powers everything from recommendation systems to autonomous vehicles, transforming how we interact with technology."
      ],
      "h3s": [
        {
          "subheading": "Types of Machine Learning",
          "bullets": [
            "Supervised learning: Learning with labeled data",
            "Unsupervised learning: Finding patterns in unlabeled data",
            "Reinforcement learning: Learning through trial and error"
          ]
        }
      ]
    }
  ],
  "enhancements": {
    "seo_keywords": [
      "machine learning",
      "artificial intelligence",
      "data science",
      "ML algorithms",
      "beginner guide"
    ],
    "content_gaps": [
      "Consider adding more practical examples",
      "Include code snippets for better understanding"
    ],
    "audience_tips": [
      "Use simple analogies for complex concepts",
      "Provide step-by-step tutorials"
    ]
  },
  "safety_report": {
    "blocked": false,
    "block_reasons": [],
    "safe": true
  },
  "word_count": 150
}
```

---

## 🚨 Error Responses

All endpoints may return error responses in the following format:

### 400 Bad Request
```json
{
  "error": "Invalid request format",
  "details": "Field validation failed"
}
```

### 401 Unauthorized
```json
{
  "error": "Missing or malformed token"
}
```

### 404 Not Found
```json
{
  "error": "Resource not found"
}
```

### 500 Internal Server Error
```json
{
  "error": "Internal server error"
}
```

---

## 📚 Postman Collection Setup

### Environment Variables
Create a Postman environment with these variables:
- `base_url`: `http://localhost:8080`
- `access_token`: (will be set after login)
- `refresh_token`: (will be set after login)

### Pre-request Scripts
For authenticated endpoints, add this pre-request script:
```javascript
pm.request.headers.add({
    key: 'Authorization',
    value: 'Bearer ' + pm.environment.get('access_token')
});
```

### Test Scripts
For login endpoint, add this test script to save tokens:
```javascript
if (pm.response.code === 200) {
    const response = pm.response.json();
    pm.environment.set('access_token', response.access_token);
    pm.environment.set('refresh_token', response.refresh_token);
}
```

---

## 🔧 Testing Flow

1. **Health Check**: Test `/health` to ensure API is running
2. **Register**: Create a new user account
3. **Login**: Get access tokens
4. **Profile**: Test profile endpoints
5. **Posts**: Create and retrieve blog posts
6. **View Tracking**: Test view analytics endpoints
7. **AI Generation**: Test content generation
8. **Account Management**: Test activation and password reset

---

## 📝 Notes

- All timestamps are in ISO 8601 format (UTC)
- User IDs and Post IDs are MongoDB ObjectIDs (24-character hex strings)
- JWT tokens expire after a set time (check your config)
- Content generation requires valid API keys for Gemini AI
- Email functionality requires SMTP configuration
- **View Tracking**: Requires Redis server running on localhost:6379
- **IP Throttling**: Each IP can only increment view count once per 24 hours
- **Analytics**: View statistics update in real-time

## 🔧 Redis Setup for View Tracking

To test view tracking features, you need Redis running:

### Install Redis:
```bash
# Windows (Chocolatey)
choco install redis-64

# macOS (Homebrew)  
brew install redis

# Ubuntu/Debian
sudo apt-get install redis-server
```

### Start Redis:
```bash
redis-server
```

### Verify Redis:
```bash
redis-cli ping
# Should return: PONG
```

Happy testing! 🚀