package auth

import (
	primary_ports "github.com/jbakhtin/goph-keeper/internal/appmodules/auth/ports/primary"
	secondary_ports "github.com/jbakhtin/goph-keeper/internal/appmodules/auth/ports/secondary"
	"github.com/jbakhtin/goph-keeper/internal/appmodules/auth/services/accesstoken"
	"github.com/jbakhtin/goph-keeper/internal/appmodules/auth/services/password"
	"github.com/jbakhtin/goph-keeper/internal/appmodules/auth/services/usecase"
)

type Module struct {
	useCase primary_ports.UseCase
}

type Config interface {
	accesstoken.Config
	password.Config
	usecase.Config
}

func NewModule(
	cfg Config,
	logger secondary_ports.Logger,
	userRepository secondary_ports.UserRepository,
	sessionRepository secondary_ports.SessionRepository,
	sessionQuerySpecification secondary_ports.SessionQuerySpecification,
	userQuerySpecification secondary_ports.UserQuerySpecification,
) (*Module, error) {
	passwordAppService, err := password.New(cfg, logger)
	if err != nil {
		return nil, err
	}

	accessAppService, err := accesstoken.New(cfg, logger)
	if err != nil {
		return nil, err
	}

	useCase, err := usecase.New(
		cfg,
		logger,
		passwordAppService,
		accessAppService,
		sessionRepository,
		sessionQuerySpecification,
		userQuerySpecification,
		userRepository,
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
