package appservices

import (
	"github.com/jbakhtin/goph-keeper/internal/server/interfaces/ports/output/logger/v1"
	"time"

	types2 "github.com/jbakhtin/goph-keeper/internal/server/domain/types"
	"github.com/jbakhtin/goph-keeper/internal/server/interfaces/appservices/v1"
	"github.com/jbakhtin/goph-keeper/internal/server/interfaces/ports/input/config/v1"

	"github.com/go-faster/errors"
	"github.com/golang-jwt/jwt/v5"
)

var _ appservices.AccessTokenAppServiceInterface = &AccessTokenAppService{}

type AccessTokenAppService struct {
	cfg config.Interface
	lgr logger.Interface
}

func NewAccessTokenAppService(cfg config.Interface, lgr logger.Interface) (*AccessTokenAppService, error) {
	return &AccessTokenAppService{
		cfg: cfg,
		lgr: lgr,
	}, nil
}

func (a *AccessTokenAppService) Create(userID types2.ID, sessionID types2.ID, duration time.Duration) (*types2.AccessToken, error) {
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

	accessToken := types2.AccessToken(rawAccessToken)

	return &accessToken, nil
}
