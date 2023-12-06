package repositories

import (
	"context"

	"github.com/jbakhtin/goph-keeper/internal/server/core/domain/types"

	"github.com/jbakhtin/goph-keeper/internal/server/core/domain/models"
)

type UserRepositoryInterface interface {
	GetUserByID(ctx context.Context, id types.ID) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	SaveUser(ctx context.Context, email, password string) (*models.User, error)
	UpdateUser(ctx context.Context, user models.User) (*models.User, error)
}
