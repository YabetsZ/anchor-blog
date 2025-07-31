## `User` Domain Model Documentation

This document describes the structure and purpose of the `User` model defined in `models/user.go`. It represents the user entity in the system, structured for MongoDB using the Go MongoDB driver.

---

### Struct: `User`

```go
type User struct {
	ID           primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Username     string               `bson:"username" json:"username"`
	FirstName    string               `bson:"first_name" json:"first_name"`
	LastName     string               `bson:"last_name" json:"last_name"`
	Email        string               `bson:"email" json:"email"`
	PasswordHash string               `bson:"password_hash" json:"-"`
	Role         string               `bson:"role" json:"role"` // "user", "admin", "unverified"
	Activated    bool                 `bson:"activated" json:"activated"`
	LastSeen     time.Time            `bson:"last_seen" json:"last_seen"`
	Profile      UserProfile          `bson:"profile" json:"profile"`
	UpdatedBy    primitive.ObjectID   `bson:"updated_by" json:"updated_by"`
	CreatedAt    time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time            `bson:"updated_at" json:"updated_at"`
	UserPosts    []primitive.ObjectID `bson:"user_posts" json:"user_posts"`
}
```

---

### Field Descriptions

| Field          | Type                   | Description                                                   |
| -------------- | ---------------------- | ------------------------------------------------------------- |
| `ID`           | `primitive.ObjectID`   | Unique MongoDB identifier for the user (auto-generated).      |
| `Username`     | `string`               | Unique username used for login and display.                   |
| `FirstName`    | `string`               | User's first name.                                            |
| `LastName`     | `string`               | User's last name.                                             |
| `Email`        | `string`               | Unique email address of the user.                             |
| `PasswordHash` | `string`               | Hashed password (excluded from JSON responses).               |
| `Role`         | `string`               | Access level: `"user"`, `"admin"`, or `"unverified"`.         |
| `Activated`    | `bool`                 | Indicates if the user account is active.                      |
| `LastSeen`     | `time.Time`            | Timestamp of the user's last activity.                        |
| `Profile`      | `UserProfile`          | Embedded profile object (see below).                          |
| `UpdatedBy`    | `primitive.ObjectID`   | Reference to the user who last modified this record.          |
| `CreatedAt`    | `time.Time`            | Timestamp when the user was created.                          |
| `UpdatedAt`    | `time.Time`            | Timestamp when the user was last updated.                     |
| `UserPosts`    | `[]primitive.ObjectID` | List of post IDs created by this user.                        |

---

### Struct: `UserProfile`

```go
type UserProfile struct {
	Bio         string       `bson:"bio" json:"bio"`
	PictureURL  string       `bson:"picture_url" json:"picture_url"`
	SocialLinks []SocialLink `bson:"social_links" json:"social_links"`
}
```

| Field         | Type           | Description                              |
| ------------- | -------------- | ---------------------------------------- |
| `Bio`         | `string`       | Short biography of the user.             |
| `PictureURL`  | `string`       | URL to the user's profile picture.       |
| `SocialLinks` | `[]SocialLink` | List of social media platforms and URLs. |

---

### Struct: `SocialLink`

```go
type SocialLink struct {
	Platform string `bson:"platform" json:"platform"`
	URL      string `bson:"url" json:"url"`
}
```

| Field      | Type     | Description                                      |
| ---------- | -------- | ------------------------------------------------ |
| `Platform` | `string` | Platform name (e.g., `"twitter"`, `"linkedin"`). |
| `URL`      | `string` | URL to the user's profile on that platform.      |

---

### Security Notes

* `PasswordHash` is intentionally excluded from all JSON responses using `json:"-"` to prevent accidental exposure.
* Use strong hashing (e.g., `bcrypt`) and never store raw passwords.

---

### Example JSON Output (for API)

```json
{
  "id": "64fdf4bcab437ea77a6eb51e",
  "username": "Segni",
  "first_name": "Haile",
  "last_name": "Doe",
  "email": "segni@test.com",
  "role": "user",
  "activated": true,
  "last_seen": "2025-07-31T12:00:00Z",
  "profile": {
    "bio": "Full-stack developer.",
    "picture_url": "https://example.com/image.jpg",
    "social_links": [
      { "platform": "twitter", "url": "https://twitter.com/johndoe" }
    ]
  },
  "updated_by": "64fdf4b9ab437ea77a6eb51d",
  "created_at": "2025-07-01T10:00:00Z",
  "updated_at": "2025-07-31T11:45:00Z",
  "user_posts": ["64fe123fab437ea77a6eb111", "64fe1245ab437ea77a6eb112"]
}
```