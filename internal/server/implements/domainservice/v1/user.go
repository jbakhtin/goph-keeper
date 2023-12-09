package domainservice

import (
	"context"

	"github.com/jbakhtin/goph-keeper/internal/server/interfaces/ports/output/logger/v1"

	"github.com/jbakhtin/goph-keeper/internal/server/domain/models"
	"github.com/jbakhtin/goph-keeper/internal/server/interfaces/domainservices/v1"
	"github.com/jbakhtin/goph-keeper/internal/server/interfaces/ports/input/config/v1"
	"github.com/jbakhtin/goph-keeper/internal/server/interfaces/ports/output/repositories/v1"
)

var _ domainservices.UserDomainServiceInterface = &userDomainService{}

// UserDomainService - implementation of UserDomainServiceInterface
type userDomainService struct {
	repo repositories.UserRepositoryInterface
	cfg  config.Interface
	lgr  logger.Interface
}

func NewUserDomainService(cfg config.Interface, lgr logger.Interface, repo repositories.UserRepositoryInterface) (*userDomainService, error) {
	return &userDomainService{
		cfg:  cfg,
		repo: repo,
		lgr:  lgr,
	}, nil
}

func (u *userDomainService) CreateUser(ctx context.Context, email, password string) (*models.User, error) {
	user, err := u.repo.SaveUser(ctx, email, password)
	if err != nil {
		return nil, err
	}

	return user, nil
}
