package models

import (
	"github.com/google/uuid"
	"time"
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
