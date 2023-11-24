package apperror

import (
	"errors"
)

var (
	UserNotFound      = errors.New("user not found")
	NotAuthorized      = errors.New("user not authorized")
	InvalidPassword      = errors.New("invalid password")
	UserAlreadyExists      = errors.New("user already exists")
)