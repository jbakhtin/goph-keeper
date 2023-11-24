package auth

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jbakhtin/goph-keeper/internal/server/application/apperror"
	"github.com/jbakhtin/goph-keeper/internal/server/application/types"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"strings"
)

// ToDo: need to refactoring:
// builder

type IConfig interface {
	GetAppKey() string
}

type Interceptor struct {
	Cfg IConfig
	AccessibleRoles map[string][]string
}

func (i *Interceptor) Unary(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {

	ctx, err := i.authorize(ctx, info.FullMethod)
	if err != nil {
		return nil, errors.Wrap(err, "unary")
	}

	return handler(ctx, req)
}

func (i *Interceptor) authorize(ctx context.Context, method string) (context.Context, error) {
	if _, ok := i.AccessibleRoles[method]; !ok {
		return ctx, nil
	}

	rawAccessToken := metadata.ValueFromIncomingContext(ctx, "authorization")
	if len(rawAccessToken) == 0 {
		return nil, apperror.NotAuthorized
	}

	bearerToken := strings.Split(rawAccessToken[0], " ")
	if bearerToken[0] != "Bearer" {
		return nil, apperror.NotAuthorized
	}

	token, err := jwt.ParseWithClaims(bearerToken[1], &types.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(i.Cfg.GetAppKey()), nil
	})

	if err != nil || !token.Valid {
		return nil, apperror.NotAuthorized
	}

	customClaims, ok := token.Claims.(*types.JWTClaims)
	if !ok {
		return nil, apperror.NotAuthorized
	}

	ctx = context.WithValue(ctx, types.ContextKeyUserID, customClaims.Data.UserId)
	ctx = context.WithValue(ctx, types.ContextKeySessionID, customClaims.Data.SessionID)

	return ctx, nil
}