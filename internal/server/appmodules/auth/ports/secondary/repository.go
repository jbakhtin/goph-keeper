package secondary_ports

import (
	"context"
	"github.com/jbakhtin/goph-keeper/internal/server/appmodules/auth/domain/models"
	"time"
)

type UserRepository interface {
	GetUser(ctx context.Context, id int) (*models.User, error)
	SaveUser(ctx context.Context, user models.User) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error) // ToDo: need remove and realize with SearchUser()
	UpdateUser(ctx context.Context, user models.User) (*models.User, error)
	DeleteUser(ctx context.Context, user models.User) (*models.User, error)
	SearchUser(ctx context.Context) ([]models.User, error) // ToDo: need add specification
}

type SessionRepository interface {
	SaveSession(ctx context.Context, UserID int, FingerPrint models.FingerPrint, ExpireAt time.Time) (*models.Session, error)
	UpdateSessionByID(ctx context.Context, id int, session models.Session) (*models.Session, error)
	GetSessionByID(ctx context.Context, ID int) (*models.Session, error)
	GetSessionsByUserID(ctx context.Context, UserID int) ([]*models.Session, error)
	GetSessionByUserIDAndFingerPrint(ctx context.Context, userID int, fingerPrint models.FingerPrint) (*models.Session, error)
	GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*models.Session, error)
	UpdateRefreshTokenByID(ctx context.Context, id int) (*models.Session, error)
	CloseSessionByID(ctx context.Context, id int) (*models.Session, error)
	CloseSessionsByUserID(ctx context.Context, userID int) ([]*models.Session, error)
}
