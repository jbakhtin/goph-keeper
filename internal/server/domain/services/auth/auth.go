package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"github.com/jbakhtin/goph-keeper/internal/server/application/apperror"
	"github.com/jbakhtin/goph-keeper/internal/server/domain/models"
	"github.com/jbakhtin/goph-keeper/internal/server/domain/repositories"
	"github.com/pkg/errors"
)

type TokensPair struct {
	AccessToken string
	RefreshToken string
}

type IConfig interface {
	GetAppKey() string
}

type Service struct {
	cfg    IConfig
	repo   repositories.UserRepository
}

func New(cfg IConfig, repo repositories.UserRepository) (*Service, error) {
	return &Service{
		cfg: cfg,
		repo:   repo,
	}, nil
}

func (us *Service) RegisterUser(ctx context.Context, newUser models.User) (*models.User, error) {
	fmt.Println("test:", newUser)
	user, err := us.repo.GetUserByEmail(ctx, newUser.Email)
	if err != nil && err != sql.ErrNoRows {
		return nil, errors.Wrap(err, "get user by email")
	}

	if user != nil {
		return nil, apperror.UserAlreadyExists
	}

	h := hmac.New(sha256.New, []byte(us.cfg.GetAppKey()))
	h.Write([]byte(fmt.Sprintf("%s:%s", newUser.Email, newUser.Password)))
	dst := h.Sum(nil)

	newUser.Password = fmt.Sprintf("%x", dst)

	user, err = us.repo.SaveUser(ctx, newUser)
	if err != nil {
		return nil, errors.Wrap(err, "save new user")
	}

	return user, nil
}

func (us *Service) LoginUser(ctx context.Context, email string, password string) (*models.User, error) {
	user, err := us.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, errors.Wrap(err, "get user by email")
	}

	h := hmac.New(sha256.New, []byte(us.cfg.GetAppKey()))
	h.Write([]byte(fmt.Sprintf("%s:%s", email, password)))
	hashedPassword := h.Sum(nil)

	fmt.Println(user.Password, fmt.Sprintf("%x", hashedPassword))

	if user.Password != fmt.Sprintf("%x", hashedPassword) {
		return nil, apperror.InvalidPassword
	}

	return user, nil
}
