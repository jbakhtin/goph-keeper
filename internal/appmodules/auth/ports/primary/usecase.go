package ports

import (
	"context"

	"github.com/jbakhtin/goph-keeper/internal/appmodules/auth/domain/models"
	"github.com/jbakhtin/goph-keeper/internal/appmodules/auth/domain/types"
)

type LogOutType int

const (
	LogoutTypeThis LogOutType = iota
	LogoutTypeAll
)

// UseCase ToDo: разобраться:
// 1. возвращаемыми и принимаемыми параметрами метода и как это связано с хранимым состоянием
// 2. может ли это повлиять на поведение системы
type UseCase interface {
	RegisterUser(ctx context.Context, email, rawPassword string) (*models.User, error)
	LoginUser(ctx context.Context, email string, password string, fingerPrint models.FingerPrint) (*types.TokensPair, error)
	RefreshToken(ctx context.Context, refreshToken string) (*types.TokensPair, error)
	Logout(ctx context.Context, logoutType LogOutType) (sessions []*models.Session, err error)
}
