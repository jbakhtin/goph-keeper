package models

import (
	"time"
)

type FingerPrint map[any]any

type Session struct {
	ID           *int
	UserID       int
	RefreshToken string
	FingerPrint FingerPrint
	ExpireAt     time.Time
	CreatedAt    *time.Time
	ClosedAt     *time.Time
	UpdatedAt    *time.Time
}

func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpireAt)
}

func (s *Session) IsClosed() bool {
	return s.ClosedAt != nil
}
