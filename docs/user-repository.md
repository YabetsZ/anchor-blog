# User Repository Interfaces

## Interface Overview

| Interface               | Responsibility                        |
| ----------------------- | ------------------------------------- |
| `IUserReaderRepository` | Read/query operations                 |
| `IUserWriterRepository` | Write/update/delete operations        |
| `IUserAuthRepository`   | Authentication-related operations     |
| `IUserAdminRepository`  | Admin-only operations (roles, status) |
| `IUserRepository`       | Aggregates all the above interfaces   |

---

## IUserReaderRepository

Responsible for **read-only** user data operations.

```go
type IUserReaderRepository interface {
	GetUserByID(ctx context.Context, id primitive.ObjectID) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	GetUser(ctx context.Context, filter map[string]interface{}) (*User, error)
	GetUsers(ctx context.Context, filter map[string]interface{}, limit, offset int) ([]*User, error)
	CountUsers(ctx context.Context, filter map[string]interface{}) (int64, error)
	GetInactiveUsers(ctx context.Context) ([]*User, error)
	SearchUsers(ctx context.Context, keyword string, limit int) ([]*User, error)
	GetUserPosts(ctx context.Context, userID primitive.ObjectID) ([]*Post, error)
	GetUserRoleByID(ctx context.Context, userID primitive.ObjectID) (string, error)
}
```

### Parameter Types:

* `context.Context`: Handles timeout, cancellation, and tracing
* `primitive.ObjectID`: MongoDB document IDs (user, post)
* `string`: Email, username, keyword
* `int`: Pagination (limit, offset)
* `map[string]interface{}`: Filter for flexible querying

---

## IUserWriterRepository

Handles **write** operations related to user persistence.

```go
type IUserWriterRepository interface {
	CreateUser(ctx context.Context, user *User) (primitive.ObjectID, error)
	UpdateUserByID(ctx context.Context, id primitive.ObjectID, update map[string]interface{}) error
	DeleteUserByID(ctx context.Context, id primitive.ObjectID) error
	SetLastSeen(ctx context.Context, id primitive.ObjectID, timestamp int64) error
}
```

### Parameter Types:

* `*User`: Full user model for creation
* `map[string]interface{}`: Partial update fields
* `int64`: UNIX timestamp for last seen tracking

---

## IUserAuthRepository

Handles user **authentication-related** operations.

```go
type IUserAuthRepository interface {
	CheckEmail(ctx context.Context, email string) (bool, error)
	CheckUsername(ctx context.Context, username string) (bool, error)
	UpdatePassword(ctx context.Context, id primitive.ObjectID, newHashedPassword string) error
}
```

### Parameter Types:

* `string`: Email, username, password hash
* `primitive.ObjectID`: Target user ID

---

## IUserAdminRepository

For **admin-level** user role and status control.

```go
type IUserAdminRepository interface {
	SetRole(ctx context.Context, id primitive.ObjectID, role string) error
	ActivateUserByID(ctx context.Context, id primitive.ObjectID) error
	DeactivateUserByID(ctx context.Context, id primitive.ObjectID) error
}
```

### Parameter Types:

* `string`: Role (e.g. `"user"`, `"admin"`)
* `primitive.ObjectID`: Target user ID

---

## IUserRepository

Composite interface combining all user-related operations.

```go
type IUserRepository interface {
	IUserReaderRepository
	IUserWriterRepository
	IUserAuthRepository
	IUserAdminRepository
}
```