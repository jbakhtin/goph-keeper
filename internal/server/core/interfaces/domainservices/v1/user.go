package domainservices

import (
	"context"

	"github.com/jbakhtin/goph-keeper/internal/server/core/domain/models"
	"github.com/jbakhtin/goph-keeper/internal/server/core/domain/types"
)

// UserDomainServiceInterface контракт описывающий функциональность предметной области модели User
type UserDomainServiceInterface interface {
	CreateUser(ctx context.Context, email, password string) (*models.User, error)
	GetUserByID(ctx context.Context, userId types.Id) (*models.User, error)
	GetUserByEmail(ctx context.Context, string string) (*models.User, error)
	// ...
}
