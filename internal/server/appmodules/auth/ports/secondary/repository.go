package secondary_ports

import (
	"context"
	"github.com/jbakhtin/goph-keeper/internal/server/appmodules/auth/domain/models"
	"github.com/jbakhtin/goph-keeper/internal/server/appmodules/auth/domain/types"
)

type UserRepository interface {
	GetUserByID(ctx context.Context, id types.ID) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	SaveUser(ctx context.Context, email, password string) (*models.User, error)
	UpdateUser(ctx context.Context, user models.User) (*models.User, error)
}

type SessionRepository interface {
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
