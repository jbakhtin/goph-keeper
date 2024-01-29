package auth

import (
	"github.com/jbakhtin/goph-keeper/internal/server/appmodules/auth/ports/primary"
	"github.com/jbakhtin/goph-keeper/internal/server/appmodules/auth/ports/secondary"
	"github.com/jbakhtin/goph-keeper/internal/server/appmodules/auth/services/acesstoken"
	"github.com/jbakhtin/goph-keeper/internal/server/appmodules/auth/services/password"
	"github.com/jbakhtin/goph-keeper/internal/server/appmodules/auth/services/usecase"
)

type Module struct {
	useCase primary_ports.UseCase
}

type Config interface {
	acesstoken.Config
	password.Config
	usecase.Config
}

func NewModule(
	cfg Config,
	logger secondary_ports.Logger,
	userRepository secondary_ports.UserRepository,
	sessionRepository secondary_ports.SessionRepository,
	) (*Module, error) {

	passwordAppService, err := password.NewPasswordAppService(cfg, logger)
	if err != nil {
		return nil, err
	}

	accessAppService, err := acesstoken.NewAccessTokenAppService(cfg, logger)
	if err != nil {
		return nil, err
	}

	useCase, err := usecase.NewAuthUseCase(cfg, logger, passwordAppService, accessAppService, sessionRepository, userRepository)

	return &Module{
		useCase: useCase,
	}, nil
}

func (m *Module) GetUseCase() primary_ports.UseCase {
	return m.useCase
}
