package entities

import (
	"github.com/google/uuid"
	"time"
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
