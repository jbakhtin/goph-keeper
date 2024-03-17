package auth

import (
	primaryports "github.com/jbakhtin/goph-keeper/internal/appmodules/auth/ports/primary"
	ports "github.com/jbakhtin/goph-keeper/internal/appmodules/auth/ports/secondary"
	"github.com/jbakhtin/goph-keeper/internal/appmodules/auth/services/accesstoken"
	"github.com/jbakhtin/goph-keeper/internal/appmodules/auth/services/password"
	"github.com/jbakhtin/goph-keeper/internal/appmodules/auth/services/usecase"
)

type Module struct {
	useCase primaryports.UseCase
}

type Config interface {
	accesstoken.Config
	password.Config
	usecase.Config
}

func NewModule(
	cfg Config,
	logger ports.Logger,
	userRepository ports.UserRepository,
	sessionRepository ports.SessionRepository,
	sessionQuerySpecification ports.SessionQuerySpecification,
	userQuerySpecification ports.UserQuerySpecification,
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

func (m *Module) GetUseCase() primaryports.UseCase {
	return m.useCase
}
