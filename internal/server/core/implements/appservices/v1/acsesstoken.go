package appservices

import (
	"time"

	"github.com/jbakhtin/goph-keeper/internal/server/core/interfaces/appservices/v1"

	"github.com/go-faster/errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jbakhtin/goph-keeper/internal/server/config"
	"github.com/jbakhtin/goph-keeper/internal/server/core/domain/types"
)

var _ appservices.AccessTokenAppServiceInterface = &AccessTokenAppService{}

type AccessTokenAppService struct {
	cfg config.Config
}

func NewAccessTokenAppService(cfg config.Config) (*AccessTokenAppService, error) {
	return &AccessTokenAppService{
		cfg: cfg,
	}, nil
}

func (a *AccessTokenAppService) Create(userID types.ID, sessionID types.ID, duration time.Duration) (*types.AccessToken, error) {
	expireAt := time.Now().Add(duration)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, types.JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireAt),
		},
		Data: types.UserData{
			UserID:    userID,
			SessionID: sessionID,
		},
	})

	rawAccessToken, err := token.SignedString([]byte(a.cfg.GetAppKey()))
	if err != nil {
		return nil, errors.Wrap(err, "sign auth jwt")
	}

	accessToken := types.AccessToken(rawAccessToken)

	return &accessToken, nil
}
