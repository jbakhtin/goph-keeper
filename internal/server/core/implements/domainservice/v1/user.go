package domainservice

import (
	"context"

	"github.com/jbakhtin/goph-keeper/internal/server/core/interfaces/domainservices/v1"
	"github.com/jbakhtin/goph-keeper/internal/server/core/interfaces/ports/output/repositories/v1"

	"github.com/jbakhtin/goph-keeper/internal/server/config"
	"github.com/jbakhtin/goph-keeper/internal/server/core/domain/models"
)

var _ domainservices.UserDomainServiceInterface = &userDomainService{}

// UserDomainService - implementation of UserDomainServiceInterface
type userDomainService struct {
	repo repositories.UserRepositoryInterface
	cfg  config.Config
}

func NewUserDomainService(cfg config.Config, repo repositories.UserRepositoryInterface) (*userDomainService, error) {
	return &userDomainService{
		cfg:  cfg,
		repo: repo,
	}, nil
}

func (u *userDomainService) CreateUser(ctx context.Context, email, password string) (*models.User, error) {
	user, err := u.repo.SaveUser(ctx, email, password)
	if err != nil {
		return nil, err
	}

	return user, nil
}
