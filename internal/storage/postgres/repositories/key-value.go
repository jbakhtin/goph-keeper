package repositories

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jbakhtin/goph-keeper/internal/server/appmodules/key-value/domain/models"
	ports "github.com/jbakhtin/goph-keeper/internal/server/appmodules/key-value/ports/secondary"
	"github.com/jbakhtin/goph-keeper/internal/server/logger/zap"
	"github.com/jbakhtin/goph-keeper/internal/server/storage/postgres/entities"
	"github.com/jbakhtin/goph-keeper/internal/server/storage/postgres/query"
)

var _ ports.KeyValueRepository = &KeyValueRepository{}

type KeyValueRepository struct {
	*sql.DB
	lgr *zap.Logger
}

func NewKeyValueRepository(lgr *zap.Logger, client *sql.DB) (*KeyValueRepository, error) {
	return &KeyValueRepository{
		DB:  client,
		lgr: lgr,
	}, nil
}

// ModelToEntity фабрика сущности базы данных
func ModelToEntity(model models.KeyValue) (entities.Secret, error) {
	return entities.Secret{
		ID:     model.ID,
		UserID: model.UserID,
		Type:   "kv",
		Data: map[any]any{
			"key": "test",
		},
		MetaData:  model.Metadata,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}, nil
}

// EntityToModel фабрика модели ядра приложения (модуля)
// NOTE фабрика модели ядра должна располагаться снаружи, что бы ядро не зависела от внешних компонентов
// можно использовать внутренние фабрики (Но они не должны зависеть от внешних пакетов) и спецификации предписывающие правила создания модели
func EntityToModel(entity entities.Secret) (models.KeyValue, error) {
	return models.KeyValue{
		ID:        entity.ID,
		UserID:    entity.UserID,
		Key:       "key",
		Value:     "test",
		Metadata:  entity.MetaData,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}, nil
}

func (r *KeyValueRepository) SaveKeyValue(ctx context.Context, model models.KeyValue) (*models.KeyValue, error) {
	row, err := ModelToEntity(model)
	if err != nil {
		return nil, err
	}

	err = r.QueryRowContext(ctx, query.CreateSecret, row.ID, row.UserID, row.Type, row.MetaData, row.Data, row.CreatedAt, row.UpdatedAt).
		Scan(&row.ID,
			&row.UserID,
			&row.Type,
			&row.MetaData,
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

func (r *KeyValueRepository) GetKeyValue(ctx context.Context, UUID uuid.UUID) (*models.KeyValue, error) {
	return nil, nil
}

func (r *KeyValueRepository) UpdateKeyValue(ctx context.Context, secret models.KeyValue) (*models.KeyValue, error) {
	return nil, nil
}

func (r *KeyValueRepository) DeleteKeyValue(ctx context.Context, secret models.KeyValue) (*models.KeyValue, error) {
	return nil, nil
}

func (r *KeyValueRepository) SearchKeyValue(ctx context.Context) ([]models.KeyValue, error) {
	return []models.KeyValue{}, nil
}
