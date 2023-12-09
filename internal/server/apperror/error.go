package apperror

import (
	"errors"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrNotAuthorized     = errors.New("user not authorized")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrUserAlreadyExists = errors.New("user already exists")
)
