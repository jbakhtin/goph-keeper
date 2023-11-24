package repositories

import (
	"context"
	"github.com/jbakhtin/goph-keeper/internal/server/application/types"
	"github.com/jbakhtin/goph-keeper/internal/server/domain/models"
)

type SessionRepository interface {
	SaveSession(ctx context.Context, session models.Session) (*models.Session, error)
	UpdateSessionByID(ctx context.Context, id types.Id, session models.Session) (*models.Session, error)
	GetSessionByUserIDAndFingerPrint(ctx context.Context, userId types.Id, fingerPrint types.FingerPrint) (*models.Session, error)
	GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*models.Session, error)
	UpdateRefreshTokenById(ctx context.Context,  id types.Id) (*models.Session, error)
	CloseSessionByID(ctx context.Context, id types.Id) (*models.Session, error)
	CloseSessionsByUserId(ctx context.Context, userId types.Id) ([]*models.Session, error)
}