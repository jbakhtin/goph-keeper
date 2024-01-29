package models

import (
	"github.com/jbakhtin/goph-keeper/internal/server/appmodules/auth/domain/types"
)

type User struct {
	ID        *types.ID
	Email     string
	Password  string
	CreatedAt *types.TimeStamp
	UpdatedAt *types.TimeStamp
}
