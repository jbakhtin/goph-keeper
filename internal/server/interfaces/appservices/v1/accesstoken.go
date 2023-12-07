package appservices

import (
	"time"

	types2 "github.com/jbakhtin/goph-keeper/internal/server/domain/types"
)

type AccessTokenAppServiceInterface interface {
	Create(userID types2.ID, sessionID types2.ID, duration time.Duration) (*types2.AccessToken, error)
}
