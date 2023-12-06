package models

import (
	"time"

	"github.com/jbakhtin/goph-keeper/internal/server/core/domain/types"
)

type LogoutType int32

const (
	LogoutType_THIS LogoutType = 0
	LogoutType_ALL  LogoutType = 1
)

var (
	LogoutType_name = map[LogoutType]string{
		LogoutType_THIS: "this",
		LogoutType_ALL:  "all",
	}

	LogoutType_value = map[string]LogoutType{
		LogoutType_name[0]: LogoutType_THIS,
		LogoutType_name[1]: LogoutType_ALL,
	}
)

type Session struct {
	ID           *types.Id
	UserId       types.Id
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
