package interceptors

import (
	"context"
	"fmt"
	"strings"

	"github.com/jbakhtin/goph-keeper/internal/server/config"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jbakhtin/goph-keeper/internal/server/apperror"
	"github.com/jbakhtin/goph-keeper/internal/server/core/domain/types"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type IConfig interface {
	GetAppKey() string
}

type AuthInterceptor struct {
	cfg             IConfig
	accessibleRoles map[string][]string
}

func NewAuthInterceptor(cfg config.Config, accessibleRoles map[string][]string) (*AuthInterceptor, error) {
	return &AuthInterceptor{
		cfg:             cfg,
		accessibleRoles: accessibleRoles,
	}, nil
}

func (i *AuthInterceptor) Unary(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	ctx, err := i.authorize(ctx, info.FullMethod)
	if err != nil {
		return nil, errors.Wrap(err, "unary")
	}

	return handler(ctx, req)
}

func (i *AuthInterceptor) authorize(ctx context.Context, method string) (context.Context, error) {
	if _, ok := i.accessibleRoles[method]; !ok {
		return ctx, nil
	}

	rawAccessToken := metadata.ValueFromIncomingContext(ctx, "authorization")
	if len(rawAccessToken) == 0 {
		return nil, apperror.ErrNotAuthorized
	}

	fmt.Println(rawAccessToken)

	bearerToken := strings.Split(rawAccessToken[0], " ")
	if bearerToken[0] != "Bearer" {
		return nil, apperror.ErrNotAuthorized
	}

	fmt.Println(bearerToken[1])

	token, err := jwt.ParseWithClaims(bearerToken[1], &types.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(i.cfg.GetAppKey()), nil
	})

	fmt.Println(token)
	fmt.Println(token.Valid)

	if err != nil || !token.Valid {
		return nil, apperror.ErrNotAuthorized
	}

	customClaims, ok := token.Claims.(*types.JWTClaims)
	if !ok {
		return nil, apperror.ErrNotAuthorized
	}

	fmt.Println(customClaims)

	ctx = context.WithValue(ctx, types.ContextKeyUserID, customClaims.Data.UserID)
	ctx = context.WithValue(ctx, types.ContextKeySessionID, customClaims.Data.SessionID)

	return ctx, nil
}
