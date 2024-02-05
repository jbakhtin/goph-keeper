package apperror

import (
	"errors"
)

// NOTE Данное расположение ошибок будет пересмотрено так как на данный момент некоторые
// более глубокие (../appmodules) компоненты приложения являются зависимыми от данного пакета.
// ToDo: изучить хорошие практики обработки ошибокб
// ToDo: разобраться где лучше и удобнее располагать пользовательские ошибки
var (
	ErrUserNotFound      = errors.New("user not found")
	ErrNotAuthorized     = errors.New("user not authorized")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrUserAlreadyExists = errors.New("user already exists")
)
