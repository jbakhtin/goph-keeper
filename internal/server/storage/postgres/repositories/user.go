package repositories

import (
	"context"
	"database/sql"
	"github.com/jbakhtin/goph-keeper/internal/server/appmodules/auth/domain/models"
	"github.com/jbakhtin/goph-keeper/internal/server/appmodules/auth/domain/types"
	"github.com/jbakhtin/goph-keeper/internal/server/appmodules/auth/ports/secondary"
	"github.com/jbakhtin/goph-keeper/internal/server/logger/zap"
	"github.com/jbakhtin/goph-keeper/internal/server/storage/postgres/query"
)

var _ secondary_ports.UserRepository = &UserRepository{}

type UserRepository struct {
	*sql.DB
	lgr *zap.Logger
}

func NewUserRepository(lgr *zap.Logger, client *sql.DB) (*UserRepository, error) { // ToDo: need to remove mock client
	return &UserRepository{
		DB: client,
		lgr:      lgr,
	}, nil
}

func (u *UserRepository) SaveUser(ctx context.Context, email, password string) (*models.User, error) {
	var stored models.User
	err := u.QueryRowContext(ctx, query.CreateUser, email, password).
		Scan(&stored.ID,
			&stored.Email,
			&stored.Password,
			&stored.UpdatedAt,
			&stored.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &stored, nil
}

func (u *UserRepository) GetUserByID(ctx context.Context, id types.ID) (*models.User, error) {
	var user models.User
	err := u.QueryRowContext(ctx, query.GetUserByID, id).
		Scan(&user.ID,
			&user.Email,
			&user.Password,
			&user.CreatedAt,
			&user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserRepository) GetUserByEmail(ctx context.Context, login string) (*models.User, error) {
	var user models.User
	err := u.QueryRowContext(ctx, query.GetUserByEmail, login).
		Scan(&user.ID,
			&user.Email,
			&user.Password,
			&user.CreatedAt,
			&user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserRepository) UpdateUser(ctx context.Context, user models.User) (*models.User, error) {
	return nil, nil
}
