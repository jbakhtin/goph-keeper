package ports

import (
	"context"

	"github.com/jbakhtin/goph-keeper/pkg/queryspec"

	"github.com/google/uuid"
	"github.com/jbakhtin/goph-keeper/internal/server/appmodules/secrets/domain/models"
)

type SecretQuerySpecification interface {
	Limit(queryspec.QuerySpecification, int) queryspec.QuerySpecification
	Where(queryspec.QuerySpecification) queryspec.QuerySpecification
	Or(...queryspec.QuerySpecification) queryspec.QuerySpecification
	And(...queryspec.QuerySpecification) queryspec.QuerySpecification
	UserID(int) queryspec.QuerySpecification
}

type SecretRepository interface {
	Create(context.Context, models.Secret) (*models.Secret, error)
	Get(context.Context, uuid.UUID) (*models.Secret, error)
	Update(context.Context, models.Secret) (*models.Secret, error)
	Delete(context.Context, models.Secret) (*models.Secret, error)
	Search(context.Context, queryspec.QuerySpecification) ([]*models.Secret, error)
}
