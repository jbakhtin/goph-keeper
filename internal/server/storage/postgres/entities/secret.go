package entities

import (
	"time"

	"github.com/google/uuid"
)

type Secret struct {
	ID        uuid.UUID
	UserID    int
	Type      string
	MetaData  string
	Data      map[any]any
	CreatedAt time.Time
	UpdatedAt time.Time
}
