package ports

import (
	"context"

	"github.com/jbakhtin/goph-keeper/internal/appmodules/key-value/domain/models"
)

type UseCase interface {
	Create(ctx context.Context, model models.KeyValue) error
}
