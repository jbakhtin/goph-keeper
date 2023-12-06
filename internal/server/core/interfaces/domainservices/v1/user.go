package domainservices

import (
	"context"

	"github.com/jbakhtin/goph-keeper/internal/server/core/domain/models"
)

// UserDomainServiceInterface контракт описывающий функциональность предметной области модели User
type UserDomainServiceInterface interface {
	CreateUser(ctx context.Context, email, password string) (*models.User, error)
	// ...
}
