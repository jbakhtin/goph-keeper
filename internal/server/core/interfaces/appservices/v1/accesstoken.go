package appservices

import (
	"time"

	"github.com/jbakhtin/goph-keeper/internal/server/core/domain/types"
)

type AccessTokenAppServiceInterface interface {
	Create(userId types.Id, sessionId types.Id, duration time.Duration) (*types.AccessToken, error)
}
