package keyvalue

import (
	primaryports "github.com/jbakhtin/goph-keeper/internal/appmodules/key-value/ports/primary"
	ports "github.com/jbakhtin/goph-keeper/internal/appmodules/key-value/ports/secondary"
	"github.com/jbakhtin/goph-keeper/internal/appmodules/key-value/services"
)

type Config interface{}

type Module struct {
	useCase primaryports.UseCase
}

func NewModule(cfg Config, logger ports.Logger, repository ports.KeyValueRepository) (*Module, error) {
	useCase, err := services.NewKeyValueUseCase(logger, repository)
	if err != nil {
		return nil, err
	}

	return &Module{
		useCase: useCase,
	}, nil
}

func (m *Module) GetUseCase() primaryports.UseCase {
	return m.useCase
}
