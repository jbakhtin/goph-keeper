package repositories

import (
	"context"

	"github.com/jbakhtin/goph-keeper/internal/server/core/domain/models"
	"github.com/jbakhtin/goph-keeper/internal/server/core/domain/types"
)

type SessionRepositoryInterface interface {
	SaveSession(ctx context.Context, UserID types.ID, FingerPrint types.FingerPrint, ExpireAt types.TimeStamp) (*models.Session, error)
	UpdateSessionByID(ctx context.Context, id types.ID, session models.Session) (*models.Session, error)
	GetSessionByID(ctx context.Context, ID types.ID) (*models.Session, error)
	GetSessionsByUserID(ctx context.Context, UserID types.ID) ([]*models.Session, error)
	GetSessionByUserIDAndFingerPrint(ctx context.Context, userID types.ID, fingerPrint types.FingerPrint) (*models.Session, error)
	GetSessionByRefreshToken(ctx context.Context, refreshToken types.RefreshToken) (*models.Session, error)
	UpdateRefreshTokenByID(ctx context.Context, id types.ID) (*models.Session, error)
	CloseSessionByID(ctx context.Context, id types.ID) (*models.Session, error)
	CloseSessionsByUserID(ctx context.Context, userID types.ID) ([]*models.Session, error)
}
