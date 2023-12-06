package models

import (
	"github.com/jbakhtin/goph-keeper/internal/server/core/domain/types"
)

type User struct {
	ID        *types.Id
	Email     string
	Password  string
	CreatedAt *types.TimeStamp
	UpdatedAt *types.TimeStamp
}
