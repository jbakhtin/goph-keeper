package usecases

import (
	"context"

	"github.com/jbakhtin/goph-keeper/internal/server/core/domain/models"
	"github.com/jbakhtin/goph-keeper/internal/server/core/domain/types"
)

type AuthUseCaseInterface interface {
	RegisterUser(ctx context.Context, email, rawPassword string) (*models.User, error)
	LoginUser(ctx context.Context, email string, password string, fingerprint types.FingerPrint) (*types.TokensPair, error)
	RefreshToken(ctx context.Context, refreshToken types.RefreshToken) (*types.TokensPair, error)
	Logout(ctx context.Context, logoutType types.LogoutType) ([]*models.Session, error)
}
