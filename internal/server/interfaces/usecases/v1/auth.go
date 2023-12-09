package usecases

import (
	"context"

	models2 "github.com/jbakhtin/goph-keeper/internal/server/domain/models"
	types2 "github.com/jbakhtin/goph-keeper/internal/server/domain/types"
)

type AuthUseCaseInterface interface {
	RegisterUser(ctx context.Context, email, rawPassword string) (*models2.User, error)
	LoginUser(ctx context.Context, email string, password string, fingerprint types2.FingerPrint) (*types2.TokensPair, error)
	RefreshToken(ctx context.Context, refreshToken types2.RefreshToken) (*types2.TokensPair, error)
	Logout(ctx context.Context, logoutType types2.LogoutType) ([]*models2.Session, error)
}
