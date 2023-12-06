package appservices

import (
	"time"

	"github.com/jbakhtin/goph-keeper/internal/server/core/domain/types"
)

type AccessTokenAppServiceInterface interface {
	Create(userID types.ID, sessionID types.ID, duration time.Duration) (*types.AccessToken, error)
}
