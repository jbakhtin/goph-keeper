package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/jbakhtin/goph-keeper/internal/server/appmodules/key-value/domain/models"
)

type QuerySpecification interface {
	Query() string
	Value() []any
}

type KeyValueQuerySpecification interface {
	Limit(QuerySpecification, int) QuerySpecification
	Where(QuerySpecification) QuerySpecification
	Or(...QuerySpecification) QuerySpecification
	And(...QuerySpecification) QuerySpecification
	UserID(int) QuerySpecification
}

type KeyValueRepository interface {
	Create(context.Context, models.KeyValue) (*models.KeyValue, error)
	Get(context.Context, uuid.UUID) (*models.KeyValue, error)
	Update(context.Context, models.KeyValue) (*models.KeyValue, error)
	Delete(context.Context, models.KeyValue) (*models.KeyValue, error)
	Search(context.Context, QuerySpecification) ([]*models.KeyValue, error)
}
