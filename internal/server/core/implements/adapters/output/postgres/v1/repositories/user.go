package repositories

import (
	"context"

	"github.com/jbakhtin/goph-keeper/internal/server/core/domain/types"
	"github.com/jbakhtin/goph-keeper/internal/server/core/implements/adapters/output/postgres/v1"
	"github.com/jbakhtin/goph-keeper/internal/server/core/interfaces/ports/output/repositories/v1"

	"github.com/jbakhtin/goph-keeper/internal/server/core/domain/models"
	"github.com/jbakhtin/goph-keeper/internal/server/core/implements/adapters/output/postgres/v1/query"
)

var _ repositories.UserRepositoryInterface = &UserRepository{}

type UserRepository struct {
	*postgres.Postgres
}

func NewUserRepository(client postgres.Postgres) (*UserRepository, error) { // ToDo: need to remove postgres client
	return &UserRepository{
		&client,
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

func (u *UserRepository) GetUserByID(ctx context.Context, id types.Id) (*models.User, error) {
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
