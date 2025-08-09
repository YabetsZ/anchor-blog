package errors

import "errors"

var (
	ErrNotFound               = errors.New("not found") // broad sense: the resource in question doesn't exist
	ErrUserNotFound           = errors.New("user not found")
	ErrInvalidUserID          = errors.New("invalid user id")
	ErrEmailAlreadyExists     = errors.New("email already exists")
	ErrUsernameTaken          = errors.New("username already taken")
	ErrInvalidCredentials     = errors.New("invalid credentials")
	ErrUnauthorized           = errors.New("unauthorized access")
	ErrForbidden              = errors.New("forbidden action")
	ErrValidationFailed       = errors.New("validation failed")
	ErrInternalServer         = errors.New("internal server error")
	ErrInvalidToken           = errors.New("invalid token")
	ErrUserIsUnverified       = errors.New("user is unverified")
	ErrUserAlreadyAdmin       = errors.New("user is already an admin")
	ErrUserNotAdmin           = errors.New("user is not an admin")
	ErrCannotDemoteThemselves = errors.New("admin can not demote themself")
	ErrInvalidPostID          = errors.New("invalid post id")
	ErrNameCannotEmpty        = errors.New("name cannot be less that three alphabet")
	ErrInvalidUsername        = errors.New("invalid username")
	ErrInvalidInput           = errors.New("invalid input parameters")
	ErrContentBlocked         = errors.New("content blocked by safety filters")
	ErrPIILeak                = errors.New("potential PII detected")
	ErrIllegalContent         = errors.New("illegal content request")
	ErrFailedToParse          = errors.New("failed to parse content")
)
