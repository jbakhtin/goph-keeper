package models

import (
	"github.com/jbakhtin/goph-keeper/internal/server/appmodules/auth/domain/types"
	"time"
)

type LogoutType int32

const (
	LogoutTypeThis LogoutType = 0
	LogoutTypeAll  LogoutType = 1
)

var (
	LogoutTypeName = map[LogoutType]string{
		LogoutTypeThis: "this",
		LogoutTypeAll:  "all",
	}

	LogoutTypeValue = map[string]LogoutType{
		LogoutTypeName[0]: LogoutTypeThis,
		LogoutTypeName[1]: LogoutTypeAll,
	}
)

type Session struct {
	ID           *types.ID
	UserID       types.ID
	RefreshToken types.RefreshToken
	FingerPrint  types.FingerPrint
	ExpireAt     types.TimeStamp
	CreatedAt    *types.TimeStamp
	ClosedAt     *types.TimeStamp
	UpdatedAt    *types.TimeStamp
}

func (s *Session) IsExpired() bool {
	return time.Now().After(time.Time(s.ExpireAt))
}

func (s *Session) IsClosed() bool {
	return s.ClosedAt != nil
}
