package services

import (
	"context"
	"github.com/jbakhtin/goph-keeper/internal/server/appmodules/key-value/domain/models"
	"github.com/jbakhtin/goph-keeper/internal/server/appmodules/key-value/ports/secondary"
)

type KeyValueUseCase struct {
	logger     secondary_ports.Logger
	repository secondary_ports.KeyValueRepository
}

func NewKeyValueUseCase(logger secondary_ports.Logger, repository secondary_ports.KeyValueRepository) (*KeyValueUseCase, error) {
	return &KeyValueUseCase{
		logger: logger,
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