package services

import (
	"context"

	"github.com/jbakhtin/goph-keeper/internal/server/appmodules/secrets/domain/models"
	secondary_ports "github.com/jbakhtin/goph-keeper/internal/server/appmodules/secrets/ports/secondary"
)

type KeyValueUseCase struct {
	logger                       secondary_ports.Logger
	repository                   secondary_ports.SecretRepository
	repositoryQuerySpecification secondary_ports.SecretQuerySpecification
}

func NewKeyValueUseCase(
	logger secondary_ports.Logger,
	repository secondary_ports.SecretRepository,
	repositoryQuerySpecification secondary_ports.SecretQuerySpecification,
) (*KeyValueUseCase, error) {
	return &KeyValueUseCase{
		logger:                       logger,
		repository:                   repository,
		repositoryQuerySpecification: repositoryQuerySpecification,
	}, nil
}

func (uc *KeyValueUseCase) Create(ctx context.Context, model models.Secret) (*models.Secret, error) {
	newModel, err := uc.repository.Create(ctx, model)
	if err != nil {
		return nil, err
	}

	return newModel, nil
}
