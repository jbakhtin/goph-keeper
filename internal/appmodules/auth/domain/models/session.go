package models

import (
	"crypto/md5"
	"encoding/hex"
	"time"
)

type FingerPrint map[string]any

type Session struct {
	ID           *int
	UserID       int
	RefreshToken string
	FingerPrint  FingerPrint
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

func (s *Session) Close() {
	now := time.Now()
	s.ClosedAt = &now
}

func (s *Session) UpdateRefreshToken() {
	hash := md5.Sum([]byte(time.Now().String()))
	s.RefreshToken = hex.EncodeToString(hash[:])
}
