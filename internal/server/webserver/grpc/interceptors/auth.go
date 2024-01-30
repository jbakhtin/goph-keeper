package interceptors

import (
	"context"
	"github.com/jbakhtin/goph-keeper/internal/server/appmodules/auth/domain/types"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jbakhtin/goph-keeper/internal/server/apperror"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

// ToDo: need move to infrastructure
var (
	ContextKeyUserID    = contextKey("user_id")
	ContextKeySessionID = contextKey("session_id")
)

type Config interface {
	GetAppKey() string
}

type AuthInterceptor struct {
	cfg             Config
	accessibleRoles map[string][]string
}

func NewAuthInterceptor(cfg Config, accessibleRoles map[string][]string) (*AuthInterceptor, error) {
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

// ToDo: разобраться как оформить код в мидлвеере
func (i *AuthInterceptor) authorize(ctx context.Context, method string) (context.Context, error) {
	if _, ok := i.accessibleRoles[method]; !ok {
		return ctx, nil
	}

	rawAccessToken := metadata.ValueFromIncomingContext(ctx, "authorization")
	if len(rawAccessToken) == 0 {
		return nil, apperror.ErrNotAuthorized
	}

	bearerToken := strings.Split(rawAccessToken[0], " ")
	if bearerToken[0] != "Bearer" {
		return nil, apperror.ErrNotAuthorized
	}

	token, err := jwt.ParseWithClaims(bearerToken[1], &types.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(i.cfg.GetAppKey()), nil
	})

	if err != nil || !token.Valid {
		return nil, apperror.ErrNotAuthorized
	}

	customClaims, ok := token.Claims.(*types.JWTClaims)
	if !ok {
		return nil, apperror.ErrNotAuthorized
	}

	ctx = context.WithValue(ctx, ContextKeyUserID, customClaims.Data.UserID)
	ctx = context.WithValue(ctx, ContextKeySessionID, customClaims.Data.SessionID)

	return ctx, nil
}
