package repositories

import (
	"context"
	"github.com/jbakhtin/goph-keeper/internal/server/domain/models"
)

type UserRepository interface {
	GetUserByID(ctx context.Context, id int) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	SaveUser(ctx context.Context, user models.User) (*models.User, error)
	UpdateUser(ctx context.Context, user models.User) (*models.User, error)
}