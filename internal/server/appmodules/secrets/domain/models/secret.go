package models

import (
	"time"
)

type Data map[string]any

type Secret struct {
	ID          int
	UserID      int
	Type        string
	Description string
	Data        Data
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
