package handlers

import (
	"context"
	"fmt"
	"github.com/jbakhtin/goph-keeper/internal/server/interfaces/ports/output/logger/v1"
	"time"

	"github.com/jbakhtin/goph-keeper/internal/server/domain/types"
	"github.com/jbakhtin/goph-keeper/internal/server/interfaces/ports/input/grpc/v1/auth"
	"github.com/jbakhtin/goph-keeper/internal/server/interfaces/usecases/v1"

	"github.com/bufbuild/protovalidate-go"
	"github.com/go-faster/errors"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ auth.AuthServiceServer = &AuthHandler{}

type AuthHandler struct {
	lgr logger.Interface
	auth.UnimplementedAuthServiceServer
	authUseCase usecases.AuthUseCaseInterface
	validator   *protovalidate.Validator
}

func NewAuthHandler(lgr logger.Interface, authUseCase usecases.AuthUseCaseInterface) (*AuthHandler, error) {
	validator, err := protovalidate.New()
	if err != nil {
		return nil, err
	}

	return &AuthHandler{
		lgr: lgr,
		authUseCase: authUseCase,
		validator:   validator,
	}, nil
}

func (h *AuthHandler) Login(ctx context.Context, request *auth.LoginRequest) (*auth.LoginResponse, error) {
	if err := h.validator.Validate(request); err != nil {
		return nil, errors.Wrap(err, "request validation")
	}

	// ToDo: need move finger print declaration mere near to db
	fingerprint, ok := peer.FromContext(ctx)
	if !ok {
		return nil, errors.New("get fingerprint from context")
	}

	tokensPair, err := h.authUseCase.LoginUser(ctx, request.Email, request.Password, types.FingerPrint{
		"addr": fingerprint.Addr.String(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "login user")
	}

	return &auth.LoginResponse{
		AccessToken:  string(tokensPair.AccessToken),
		RefreshToken: string(tokensPair.RefreshToken),
	}, nil
}

func (h *AuthHandler) Register(ctx context.Context, request *auth.RegisterRequest) (*emptypb.Empty, error) {
	if err := h.validator.Validate(request); err != nil {
		return nil, errors.Wrap(err, "request validation")
	}

	_, err := h.authUseCase.RegisterUser(ctx, request.Email, request.Password)
	if err != nil {
		return nil, errors.Wrap(err, "register handlers")
	}

	return &emptypb.Empty{}, nil
}

func (h *AuthHandler) RefreshAccessToken(ctx context.Context, request *emptypb.Empty) (*auth.RefreshTokenResponse, error) {
	if err := h.validator.Validate(request); err != nil {
		return nil, errors.Wrap(err, "request validation")
	}

	values := metadata.ValueFromIncomingContext(ctx, "refresh-token")

	tokensPair, err := h.authUseCase.RefreshToken(ctx, types.RefreshToken(values[0]))
	if err != nil {
		return nil, errors.Wrap(err, "refresh access token")
	}

	return &auth.RefreshTokenResponse{
		AccessToken:  string(tokensPair.AccessToken),
		RefreshToken: string(tokensPair.RefreshToken),
	}, nil
}

func (h *AuthHandler) Logout(ctx context.Context, request *auth.LogoutRequest) (*auth.LogoutResponse, error) {
	if err := h.validator.Validate(request); err != nil {
		return nil, errors.Wrap(err, "request validation")
	}

	var logoutType types.LogoutType
	switch request.Type {
	case auth.LogoutType_TYPE_ALL:
		logoutType = types.LogoutTypeAll
	case auth.LogoutType_TYPE_UNSPECIFIED:
		logoutType = types.LogoutTypeThis
	default:
		return nil, errors.New("logout type is not allowed")
	}

	sessions, err := h.authUseCase.Logout(ctx, logoutType)
	if err != nil {
		return nil, errors.Wrap(err, "close sessions")
	}

	var pbSessions []*auth.Session
	for _, session := range sessions {
		pbSessions = append(pbSessions, &auth.Session{
			ID:        uint64(*session.ID),
			UserID:    uint64(session.UserID),
			CreatedAt: time.Time(*session.CreatedAt).String(),
			ClosedAt:  time.Time(*session.ClosedAt).String(),
		})
	}

	fmt.Println("test")

	return &auth.LogoutResponse{
		Type:     request.Type,
		Sessions: pbSessions,
	}, nil
}
