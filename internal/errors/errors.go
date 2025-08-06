package errors

import "errors"

var (
	ErrNotFound           = errors.New("not found") // broad sense: the resource in question doesn't exist
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidUserID      = errors.New("invalid user id")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrUsernameTaken      = errors.New("username already taken")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUnauthorized       = errors.New("unauthorized access")
	ErrForbidden          = errors.New("forbidden action")
	ErrValidationFailed   = errors.New("validation failed")
	ErrInternalServer     = errors.New("internal server error")
	ErrInvalidToken       = errors.New("invalid token")
	ErrNameCannotEmpty    = errors.New("name cannot be empty")
	ErrInvalidUsername    = errors.New("invalid username")
)
