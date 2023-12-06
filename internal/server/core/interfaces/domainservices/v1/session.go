package domainservices

import (
	"context"

	"github.com/jbakhtin/goph-keeper/internal/server/core/domain/models"
	"github.com/jbakhtin/goph-keeper/internal/server/core/domain/types"
)

// SessionDomainServiceInterface контракт описывающий функциональность предметной области модели Session
type SessionDomainServiceInterface interface {
	CreateSession(ctx context.Context, UserID types.ID, fingerPrint types.FingerPrint, expireAt types.TimeStamp) (*models.Session, error)
	CloseSession(ctx context.Context, session models.Session) (*models.Session, error)
	UpdateRefreshToken(ctx context.Context, session models.Session) (*models.Session, error)
	// ...
}
