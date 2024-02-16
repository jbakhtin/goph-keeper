package ports

import (
	"context"

	"github.com/jbakhtin/goph-keeper/internal/server/appmodules/secrets/domain/models"
)

type UseCase interface {
	Create(ctx context.Context, model models.Secret) (*models.Secret, error)
}
