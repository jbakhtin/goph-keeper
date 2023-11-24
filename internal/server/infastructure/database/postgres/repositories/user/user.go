package user

import (
	"context"
	"github.com/jbakhtin/goph-keeper/internal/server/domain/models"
	"github.com/jbakhtin/goph-keeper/internal/server/domain/repositories"
	"github.com/jbakhtin/goph-keeper/internal/server/infastructure/database/postgres"
	"github.com/jbakhtin/goph-keeper/internal/server/infastructure/database/postgres/query"
)

var _ repositories.UserRepository = (*Repository)(nil)

type Repository struct {
	*postgres.Postgres
}

func New(client postgres.Postgres) (*Repository, error) {
	return &Repository{
		&client,
	}, nil
}

func (ur *Repository) SaveUser(ctx context.Context, user models.User) (*models.User, error) {
	var stored models.User
	err := ur.QueryRowContext(ctx, query.CreateUser, &user.Email, &user.Password).
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

func (ur *Repository) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	var user models.User
	err := ur.QueryRowContext(ctx, query.GetUserByID, id).
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

func (ur *Repository) GetUserByEmail(ctx context.Context, login string) (*models.User, error) {
	var user models.User
	err := ur.QueryRowContext(ctx, query.GetUserByEmail, login).
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

func (ur *Repository) UpdateUser(ctx context.Context, user models.User) (*models.User, error) {
	return nil, nil
}
