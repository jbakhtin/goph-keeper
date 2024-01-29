package secondary_ports

import (
	"context"
	"github.com/google/uuid"
	"github.com/jbakhtin/goph-keeper/internal/server/appmodules/key-value/domain/models"
)

type KeyValueRepository interface {
	SaveKeyValue(ctx context.Context, secret models.KeyValue) (*models.KeyValue, error)
	GetKeyValue(ctx context.Context, UUID uuid.UUID) (*models.KeyValue, error)
	UpdateKeyValue(ctx context.Context, secret models.KeyValue) (*models.KeyValue, error)
	DeleteKeyValue(ctx context.Context, secret models.KeyValue) (*models.KeyValue, error)
	SearchKeyValue(ctx context.Context) ([]models.KeyValue, error) // ToDo: need add specifications
}