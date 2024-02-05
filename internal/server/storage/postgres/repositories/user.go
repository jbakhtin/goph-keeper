package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jbakhtin/goph-keeper/internal/server/appmodules/auth/domain/models"
	"github.com/jbakhtin/goph-keeper/internal/server/appmodules/auth/ports/secondary"
	"github.com/jbakhtin/goph-keeper/internal/server/logger/zap"
	"github.com/jbakhtin/goph-keeper/internal/server/storage/postgres/query"
)

var _ ports.UserRepository = &UserRepository{}

type UserRepository struct {
	*sql.DB
	lgr *zap.Logger
}

func NewUserRepository(lgr *zap.Logger, client *sql.DB) (*UserRepository, error) {
	return &UserRepository{
		DB:  client,
		lgr: lgr,
	}, nil
}

func (u *UserRepository) Create(ctx context.Context, user models.User) (*models.User, error) {
	var stored models.User
	err := u.QueryRowContext(ctx, query.CreateUser, user.Email, user.Password).
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

func (u *UserRepository) Get(ctx context.Context, id int) (*models.User, error) {
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

func (u *UserRepository) Update(ctx context.Context, user models.User) (*models.User, error) {
	return nil, nil
}

func (u *UserRepository) Delete(ctx context.Context, user models.User) (*models.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u *UserRepository) Search(ctx context.Context, specification ports.QuerySpecification) ([]*models.User, error) {
	rows, err := u.QueryContext(ctx, fmt.Sprintf("%s %s", query.SearchUserTemp, specification.Query()))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*models.User, 0)
	for rows.Next() {
		var user models.User
		err = rows.Scan(&user.ID,
			&user.Email,
			&user.Password,
			&user.CreatedAt,
			&user.UpdatedAt)
		if err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
