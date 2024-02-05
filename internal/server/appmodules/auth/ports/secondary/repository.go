package ports

import (
	"context"
	"github.com/jbakhtin/goph-keeper/internal/server/appmodules/auth/domain/models"
)

type QuerySpecification interface {
	Query() string
	Value() []any
}

type SessionQuerySpecification interface {
	Limit(QuerySpecification, int) QuerySpecification
	Where(QuerySpecification) QuerySpecification
	Or(...QuerySpecification) QuerySpecification
	And(...QuerySpecification) QuerySpecification
	UserID(int) QuerySpecification
	IsNotClosed() QuerySpecification
	FingerPrint(models.FingerPrint) QuerySpecification
	RefreshToken(string) QuerySpecification
}

type UserQuerySpecification interface {
	Limit(QuerySpecification, int) QuerySpecification
	Where(QuerySpecification) QuerySpecification
	Or(...QuerySpecification) QuerySpecification
	And(...QuerySpecification) QuerySpecification
	Email(string) QuerySpecification
}

type UserRepository interface {
	Get(ctx context.Context, id int) (*models.User, error)
	Create(ctx context.Context, user models.User) (*models.User, error)
	Update(ctx context.Context, user models.User) (*models.User, error)
	Delete(ctx context.Context, user models.User) (*models.User, error)
	Search(context.Context, QuerySpecification) ([]*models.User, error) // ToDo: need add specification
}

type SessionRepository interface {
	Get(ctx context.Context, ID int) (*models.Session, error)
	Create(ctx context.Context, session models.Session) (*models.Session, error)
	Update(ctx context.Context, session models.Session) (*models.Session, error)
	Delete(ctx context.Context, session models.Session) (*models.Session, error)
	Search(context.Context, QuerySpecification) ([]*models.Session, error)
}
