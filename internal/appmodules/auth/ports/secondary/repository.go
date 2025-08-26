package ports

import (
	"context"

	"github.com/jbakhtin/goph-keeper/pkg/queryspec"

	"github.com/jbakhtin/goph-keeper/internal/appmodules/auth/domain/models"
)

type SessionQuerySpecification interface {
	Limit(queryspec.QuerySpecification, int) queryspec.QuerySpecification
	Where(queryspec.QuerySpecification) queryspec.QuerySpecification
	Or(...queryspec.QuerySpecification) queryspec.QuerySpecification
	And(...queryspec.QuerySpecification) queryspec.QuerySpecification
	UserID(int) queryspec.QuerySpecification
	IsNotClosed() queryspec.QuerySpecification
	FingerPrint(models.FingerPrint) queryspec.QuerySpecification
	RefreshToken(string) queryspec.QuerySpecification
}

type UserQuerySpecification interface {
	Limit(queryspec.QuerySpecification, int) queryspec.QuerySpecification
	Where(queryspec.QuerySpecification) queryspec.QuerySpecification
	Or(...queryspec.QuerySpecification) queryspec.QuerySpecification
	And(...queryspec.QuerySpecification) queryspec.QuerySpecification
	Email(string) queryspec.QuerySpecification
}

type UserRepository interface {
	Get(ctx context.Context, id int) (*models.User, error)
	Create(ctx context.Context, user models.User) (*models.User, error)
	Update(ctx context.Context, user models.User) (*models.User, error)
	Delete(ctx context.Context, user models.User) (*models.User, error)
	Search(context.Context, queryspec.QuerySpecification) ([]*models.User, error) // ToDo: need add specification
}

type SessionRepository interface {
	Get(ctx context.Context, ID int) (*models.Session, error)
	Create(ctx context.Context, session models.Session) (*models.Session, error)
	Update(ctx context.Context, session models.Session) (*models.Session, error)
	Delete(ctx context.Context, session models.Session) (*models.Session, error)
	Search(context.Context, queryspec.QuerySpecification) ([]*models.Session, error)
}
