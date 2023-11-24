package models

import (
	"github.com/jbakhtin/goph-keeper/internal/server/application/types"
)

type User struct {
	ID        *types.Id
	Email     string
	Password  string
	CreatedAt *types.TimeStamp
	UpdatedAt *types.TimeStamp
}
