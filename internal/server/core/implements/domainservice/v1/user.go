package domainservice

import (
	"context"

	"github.com/jbakhtin/goph-keeper/internal/server/core/interfaces/domainservices/v1"
	"github.com/jbakhtin/goph-keeper/internal/server/core/interfaces/ports/output/repositories/v1"

	"github.com/jbakhtin/goph-keeper/internal/server/config"
	"github.com/jbakhtin/goph-keeper/internal/server/core/domain/models"
	"github.com/jbakhtin/goph-keeper/internal/server/core/domain/types"
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

func (u *userDomainService) GetUserByID(ctx context.Context, userID types.ID) (*models.User, error) {
	user, err := u.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userDomainService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := u.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return user, nil
}
