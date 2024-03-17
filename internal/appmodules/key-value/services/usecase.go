package services

import (
	"context"

	"github.com/jbakhtin/goph-keeper/internal/appmodules/key-value/domain/models"
	ports "github.com/jbakhtin/goph-keeper/internal/appmodules/key-value/ports/secondary"
)

type KeyValueUseCase struct {
	logger     ports.Logger
	repository ports.KeyValueRepository
}

func NewKeyValueUseCase(logger ports.Logger, repository ports.KeyValueRepository) (*KeyValueUseCase, error) {
	return &KeyValueUseCase{
		logger:     logger,
		repository: repository,
	}, nil
}

func (uc *KeyValueUseCase) Create(ctx context.Context, model models.KeyValue) error {
	_, err := uc.repository.SaveKeyValue(ctx, model)
	if err != nil {
		return err
	}

	return nil
}
