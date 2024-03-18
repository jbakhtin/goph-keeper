package repositories

import (
	"context"
	"database/sql"

	"github.com/jbakhtin/goph-keeper/pkg/queryspec"

	"github.com/google/uuid"
	"github.com/jbakhtin/goph-keeper/internal/appmodules/secrets/domain/models"
	secondaryports "github.com/jbakhtin/goph-keeper/internal/appmodules/secrets/ports/secondary"
	"github.com/jbakhtin/goph-keeper/internal/logger/zap"
	"github.com/jbakhtin/goph-keeper/internal/storage/postgres/entities"
	"github.com/jbakhtin/goph-keeper/internal/storage/postgres/query"
)

var _ secondaryports.SecretRepository = &SecretRepository{}

type SecretRepository struct {
	*sql.DB
	lgr *zap.Logger
}

func NewKeyValueRepository(lgr *zap.Logger, client *sql.DB) (*SecretRepository, error) {
	return &SecretRepository{
		DB:  client,
		lgr: lgr,
	}, nil
}

// ModelToEntity фабрика сущности базы данных
func ModelToEntity(model models.Secret) (entities.Secret, error) {
	return entities.Secret{
		ID:     model.ID,
		UserID: model.UserID,
		Data: map[string]any{
			"key": "test",
		},
		Description: model.Description,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}, nil
}

// EntityToModel фабрика модели ядра приложения (модуля)
// NOTE фабрика модели ядра должна располагаться снаружи, что бы ядро не зависела от внешних компонентов
// можно использовать внутренние фабрики (Но они не должны зависеть от внешних пакетов) и спецификации предписывающие правила создания модели
func EntityToModel(entity entities.Secret) (models.Secret, error) {
	return models.Secret{
		ID:          entity.ID,
		UserID:      entity.UserID,
		Data:        models.Data(entity.Data),
		Description: entity.Description,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}, nil
}

func (r *SecretRepository) Create(ctx context.Context, model models.Secret) (*models.Secret, error) {
	row, err := ModelToEntity(model)
	if err != nil {
		return nil, err
	}

	err = r.QueryRowContext(ctx, query.CreateSecret, row.UserID, row.Description, row.Data).
		Scan(&row.ID,
			&row.UserID,
			&row.Description,
			&row.Data,
			&row.CreatedAt,
			&row.UpdatedAt) // ToDo: need figure it out how to scan row directly to structure
	if err != nil {
		return nil, err
	}

	model, err = EntityToModel(row)
	if err != nil {
		return nil, err
	}

	return &model, nil
}

func (r *SecretRepository) Get(ctx context.Context, UUID uuid.UUID) (*models.Secret, error) {
	return nil, nil
}

func (r *SecretRepository) Update(ctx context.Context, secret models.Secret) (*models.Secret, error) {
	return nil, nil
}

func (r *SecretRepository) Delete(ctx context.Context, secret models.Secret) (*models.Secret, error) {
	return nil, nil
}

func (r *SecretRepository) Search(ctx context.Context, query queryspec.QuerySpecification) ([]*models.Secret, error) {
	return []*models.Secret{}, nil
}
