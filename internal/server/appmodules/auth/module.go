package auth

import (
	primaryports "github.com/jbakhtin/goph-keeper/internal/server/appmodules/auth/ports/primary"
	ports "github.com/jbakhtin/goph-keeper/internal/server/appmodules/auth/ports/secondary"
	"github.com/jbakhtin/goph-keeper/internal/server/appmodules/auth/services/acesstoken"
	"github.com/jbakhtin/goph-keeper/internal/server/appmodules/auth/services/password"
	"github.com/jbakhtin/goph-keeper/internal/server/appmodules/auth/services/usecase"
)

type Module struct {
	useCase primaryports.UseCase
}

type Config interface {
	acesstoken.Config
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
	passwordAppService, err := password.NewPasswordAppService(cfg, logger)
	if err != nil {
		return nil, err
	}

	accessAppService, err := acesstoken.NewAccessTokenAppService(cfg, logger)
	if err != nil {
		return nil, err
	}

	useCase, err := usecase.NewAuthUseCase(
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
