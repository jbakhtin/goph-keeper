package acesstoken

import (
	types2 "github.com/jbakhtin/goph-keeper/internal/server/appmodules/auth/domain/types"
	"github.com/jbakhtin/goph-keeper/internal/server/appmodules/auth/ports/secondary"
	"time"

	"github.com/go-faster/errors"
	"github.com/golang-jwt/jwt/v5"
)

type Config interface {
	GetAppKey() string
}

type AccessTokenAppService struct {
	cfg Config
	lgr ports.Logger
}

func NewAccessTokenAppService(cfg Config, lgr ports.Logger) (*AccessTokenAppService, error) {
	return &AccessTokenAppService{
		cfg: cfg,
		lgr: lgr,
	}, nil
}

func (a *AccessTokenAppService) Create(userID int, sessionID int, duration time.Duration) (*string, error) {
	expireAt := time.Now().Add(duration)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, types2.JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireAt),
		},
		Data: types2.UserData{
			UserID:    userID,
			SessionID: sessionID,
		},
	})

	rawAccessToken, err := token.SignedString([]byte(a.cfg.GetAppKey()))
	if err != nil {
		return nil, errors.Wrap(err, "sign auth jwt")
	}

	accessToken := rawAccessToken

	return &accessToken, nil
}
