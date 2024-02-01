package secondary_ports

import (
	"context"
	"github.com/jbakhtin/goph-keeper/internal/server/appmodules/auth/domain/models"
	"time"
)

type Specification interface {
	Query() string
	Value() []any
}

type SessionSpecifications interface {
	Limit(Specification, int) Specification
	Where(Specification) Specification
	Or(...Specification) Specification
	And(...Specification) Specification
	UserID(int) Specification
	IsNotClosed() Specification
	FingerPrint(models.FingerPrint) Specification
	RefreshToken(string) Specification
}

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
	UpdateSession(ctx context.Context, session models.Session) (*models.Session, error)
	GetSession(ctx context.Context, ID int) (*models.Session, error)
	Search(ctx context.Context, specs Specification) ([]*models.Session, error)
	GetSessionByUserIDAndFingerPrint(ctx context.Context, userID int, fingerPrint models.FingerPrint) (*models.Session, error)
	GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*models.Session, error)
}
