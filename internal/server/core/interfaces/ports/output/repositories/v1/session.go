package repositories

import (
	"context"

	"github.com/jbakhtin/goph-keeper/internal/server/core/domain/models"
	"github.com/jbakhtin/goph-keeper/internal/server/core/domain/types"
)

type SessionRepositoryInterface interface {
	SaveSession(ctx context.Context, UserId types.Id, FingerPrint types.FingerPrint, ExpireAt types.TimeStamp) (*models.Session, error)
	UpdateSessionByID(ctx context.Context, id types.Id, session models.Session) (*models.Session, error)
	GetSessionByID(ctx context.Context, ID types.Id) (*models.Session, error)
	GetSessionsByUserID(ctx context.Context, UserID types.Id) ([]*models.Session, error)
	GetSessionByUserIDAndFingerPrint(ctx context.Context, userId types.Id, fingerPrint types.FingerPrint) (*models.Session, error)
	GetSessionByRefreshToken(ctx context.Context, refreshToken types.RefreshToken) (*models.Session, error)
	UpdateRefreshTokenById(ctx context.Context, id types.Id) (*models.Session, error)
	CloseSessionByID(ctx context.Context, id types.Id) (*models.Session, error)
	CloseSessionsByUserId(ctx context.Context, userId types.Id) ([]*models.Session, error)
}
