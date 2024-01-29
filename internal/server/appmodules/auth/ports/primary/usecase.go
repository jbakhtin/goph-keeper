package primary_ports

import (
	"context"
	models2 "github.com/jbakhtin/goph-keeper/internal/server/appmodules/auth/domain/models"
	types2 "github.com/jbakhtin/goph-keeper/internal/server/appmodules/auth/domain/types"
)

// UseCase ToDo: разобраться:
// 1. возвращаемыми и принимаемыми параметрами метода и как это связано с хранимым состоянимем
// 2. может ли это повлиять на поведение системы
type UseCase interface {
	RegisterUser(ctx context.Context, email, rawPassword string) (*models2.User, error)
	LoginUser(ctx context.Context, email string, password string, fingerPrint types2.FingerPrint) (*types2.TokensPair, error)
	RefreshToken(ctx context.Context, refreshToken types2.RefreshToken) (*types2.TokensPair, error)
	Logout(ctx context.Context, logoutType types2.LogoutType) (sessions []*models2.Session, err error)
}

