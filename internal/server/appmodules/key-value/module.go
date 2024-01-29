package key_value

import (
	"github.com/jbakhtin/goph-keeper/internal/server/appmodules/key-value/ports/primary"
	"github.com/jbakhtin/goph-keeper/internal/server/appmodules/key-value/ports/secondary"
	"github.com/jbakhtin/goph-keeper/internal/server/appmodules/key-value/services"
)

type Config interface {}

type Module struct {
	useCase primary_ports.UseCase
}

func NewModule(cfg Config, logger secondary_ports.Logger, repository secondary_ports.KeyValueRepository) (*Module, error) {
	useCase, err := services.NewKeyValueUseCase(logger, repository)
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
