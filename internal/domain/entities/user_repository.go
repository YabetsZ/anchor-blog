package entities

import (
	"context"
	"time"
)

// For Read
type IUserReaderRepository interface {
	GetUserByID(ctx context.Context, id string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	GetUsers(ctx context.Context, limit, offset int64) ([]*User, error)
	CountUsersByRole(ctx context.Context, role string) (int64, error)
	CountAllUsers(ctx context.Context) (int64, error)
	CountActiveUsers(ctx context.Context) (int64, error)
	CountInactiveUsers(ctx context.Context) (int64, error)
	// GetInactiveUsers(ctx context.Context) ([]*User, error)
	// SearchUsers(ctx context.Context, keyword string, limit int) ([]*User, error)
	// GetUserPosts(ctx context.Context, userID string) ([]*Post, error)
	GetUserRoleByID(ctx context.Context, userID string) (string, error)
}

// For Write
type IUserWriterRepository interface {
	CreateUser(ctx context.Context, user *User) (string, error)
	EditUserByID(ctx context.Context, id string, user *User) error
	DeleteUserByID(ctx context.Context, id string) error
	SetLastSeen(ctx context.Context, id string, timestamp time.Time) error
}

// For Auth
type IUserAuthRepository interface {
	CheckEmail(ctx context.Context, email string) (bool, error)
	CheckUsername(ctx context.Context, username string) (bool, error)
	ChangePassword(ctx context.Context, id string, newHashedPassword string) error
	ChangeEmail(ctx context.Context, email string, newEmail string) error
}

// Only for ADMIN
type IUserAdminRepository interface {
	SetRole(ctx context.Context, id string, role string) error
	ActivateUserByID(ctx context.Context, id string) error
	DeactivateUserByID(ctx context.Context, id string) error
}

// User Repository
type IUserRepository interface {
	IUserReaderRepository
	IUserWriterRepository
	IUserAuthRepository
	IUserAdminRepository
}
