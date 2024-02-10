package accesstoken

import (
	"time"

	"github.com/jbakhtin/goph-keeper/internal/server/appmodules/auth/domain/types"
	secondaryports "github.com/jbakhtin/goph-keeper/internal/server/appmodules/auth/ports/secondary"

	"github.com/go-faster/errors"
	"github.com/golang-jwt/jwt/v5"
)

type Config interface {
	GetAppKey() string
}

type Service struct {
	cfg Config
	lgr secondaryports.Logger
}

func New(cfg Config, lgr secondaryports.Logger) (*Service, error) {
	return &Service{
		cfg: cfg,
		lgr: lgr,
	}, nil
}

func (a *Service) Create(userID int, sessionID int, duration time.Duration) (*string, error) {
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

	accessToken := rawAccessToken

	return &accessToken, nil
}
