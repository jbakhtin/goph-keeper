package handlers

import (
	"context"
	"fmt"
	"time"

	types2 "github.com/jbakhtin/goph-keeper/internal/server/domain/types"
	auth2 "github.com/jbakhtin/goph-keeper/internal/server/interfaces/ports/input/grpc/v1/auth"
	"github.com/jbakhtin/goph-keeper/internal/server/interfaces/usecases/v1"

	"github.com/bufbuild/protovalidate-go"
	"github.com/go-faster/errors"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ auth2.AuthServiceServer = &AuthHandler{}

type AuthHandler struct {
	auth2.UnimplementedAuthServiceServer
	authUseCase usecases.AuthUseCaseInterface
	validator   *protovalidate.Validator
}

func NewAuthHandler(authUseCase usecases.AuthUseCaseInterface) (*AuthHandler, error) {
	validator, err := protovalidate.New()
	if err != nil {
		return nil, err
	}

	return &AuthHandler{
		authUseCase: authUseCase,
		validator:   validator,
	}, nil
}

func (h *AuthHandler) Login(ctx context.Context, request *auth2.LoginRequest) (*auth2.LoginResponse, error) {
	if err := h.validator.Validate(request); err != nil {
		return nil, errors.Wrap(err, "request validation")
	}

	// ToDo: need move finger print declaration mere near to db
	fingerprint, ok := peer.FromContext(ctx)
	if !ok {
		return nil, errors.New("get fingerprint from context")
	}

	tokensPair, err := h.authUseCase.LoginUser(ctx, request.Email, request.Password, types2.FingerPrint{
		"addr": fingerprint.Addr.String(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "login user")
	}

	return &auth2.LoginResponse{
		AccessToken:  string(tokensPair.AccessToken),
		RefreshToken: string(tokensPair.RefreshToken),
	}, nil
}

func (h *AuthHandler) Register(ctx context.Context, request *auth2.RegisterRequest) (*emptypb.Empty, error) {
	if err := h.validator.Validate(request); err != nil {
		return nil, errors.Wrap(err, "request validation")
	}

	_, err := h.authUseCase.RegisterUser(ctx, request.Email, request.Password)
	if err != nil {
		return nil, errors.Wrap(err, "register handlers")
	}

	return &emptypb.Empty{}, nil
}

func (h *AuthHandler) RefreshAccessToken(ctx context.Context, request *emptypb.Empty) (*auth2.RefreshTokenResponse, error) {
	if err := h.validator.Validate(request); err != nil {
		return nil, errors.Wrap(err, "request validation")
	}

	values := metadata.ValueFromIncomingContext(ctx, "refresh-token")

	tokensPair, err := h.authUseCase.RefreshToken(ctx, types2.RefreshToken(values[0]))
	if err != nil {
		return nil, errors.Wrap(err, "refresh access token")
	}

	return &auth2.RefreshTokenResponse{
		AccessToken:  string(tokensPair.AccessToken),
		RefreshToken: string(tokensPair.RefreshToken),
	}, nil
}

func (h *AuthHandler) Logout(ctx context.Context, request *auth2.LogoutRequest) (*auth2.LogoutResponse, error) {
	if err := h.validator.Validate(request); err != nil {
		return nil, errors.Wrap(err, "request validation")
	}

	var logoutType types2.LogoutType
	switch request.Type {
	case auth2.LogoutType_TYPE_ALL:
		fmt.Println("LogoutType_TYPE_ALL")
		logoutType = types2.LogoutTypeAll
	case auth2.LogoutType_TYPE_UNSPECIFIED:
		logoutType = types2.LogoutTypeThis
	default:
		return nil, errors.New("logout type is not allowed")
	}

	sessions, err := h.authUseCase.Logout(ctx, logoutType)
	if err != nil {
		return nil, errors.Wrap(err, "close sessions")
	}

	var pbSessions []*auth2.Session
	for _, session := range sessions {
		pbSessions = append(pbSessions, &auth2.Session{
			ID:        uint64(*session.ID),
			UserID:    uint64(session.UserID),
			CreatedAt: time.Time(*session.CreatedAt).String(),
			ClosedAt:  time.Time(*session.ClosedAt).String(),
		})
	}

	fmt.Println("test")

	return &auth2.LogoutResponse{
		Type:     request.Type,
		Sessions: pbSessions,
	}, nil
}
