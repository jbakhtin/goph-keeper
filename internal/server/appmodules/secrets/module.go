package secrets

import (
	primary_ports "github.com/jbakhtin/goph-keeper/internal/server/appmodules/secrets/ports/primary"
	secondary_ports "github.com/jbakhtin/goph-keeper/internal/server/appmodules/secrets/ports/secondary"
	"github.com/jbakhtin/goph-keeper/internal/server/appmodules/secrets/services"
)

type Config interface{}

type Module struct {
	useCase primary_ports.UseCase
}

func NewModule(
	cfg Config,
	logger secondary_ports.Logger,
	repository secondary_ports.SecretRepository,
	repositoryQuerySpecification secondary_ports.SecretQuerySpecification,
) (*Module, error) {
	useCase, err := services.NewKeyValueUseCase(
		logger,
		repository,
		repositoryQuerySpecification,
	)
	if err != nil {
		return nil, err
	}

	return &Module{
		useCase: useCase,
	}, nil
}

func (m *Module) GetUseCase() primary_ports.UseCase {
	return m.useCase
}
