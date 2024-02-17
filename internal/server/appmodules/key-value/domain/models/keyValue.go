package models

import (
	"time"

	"github.com/google/uuid"
)

type KeyValue struct {
	ID        uuid.UUID
	UserID    int
	Key       string
	Value     string
	Metadata  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
