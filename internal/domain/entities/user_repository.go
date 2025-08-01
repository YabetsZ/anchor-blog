package entities

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// For Read
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

// For Write
type IUserWriterRepository interface {
	CreateUser(ctx context.Context, user *User) (primitive.ObjectID, error)
	UpdateUserByID(ctx context.Context, id primitive.ObjectID, update map[string]interface{}) error
	DeleteUserByID(ctx context.Context, id primitive.ObjectID) error
	SetLastSeen(ctx context.Context, id primitive.ObjectID, timestamp int64) error
}

// For Auth
type IUserAuthRepository interface {
	CheckEmail(ctx context.Context, email string) (bool, error)
	CheckUsername(ctx context.Context, username string) (bool, error)
	UpdatePassword(ctx context.Context, id primitive.ObjectID, newHashedPassword string) error
}

// Only for ADMIN
type IUserAdminRepository interface {
	SetRole(ctx context.Context, id primitive.ObjectID, role string) error
	ActivateUserByID(ctx context.Context, id primitive.ObjectID) error
	DeactivateUserByID(ctx context.Context, id primitive.ObjectID) error
}

// User Repository
type IUserRepository interface {
	IUserReaderRepository
	IUserWriterRepository
	IUserAuthRepository
	IUserAdminRepository
}
