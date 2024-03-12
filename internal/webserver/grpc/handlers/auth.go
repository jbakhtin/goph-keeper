package handlers

import (
	"context"

	"github.com/bufbuild/protovalidate-go"
	"github.com/go-faster/errors"
	"github.com/jbakhtin/goph-keeper/gen/go/v1/auth"
	"github.com/jbakhtin/goph-keeper/internal/server/appmodules/auth/domain/models"
	primary_ports "github.com/jbakhtin/goph-keeper/internal/server/appmodules/auth/ports/primary"
	"github.com/jbakhtin/goph-keeper/internal/server/logger/zap"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ auth.AuthServiceServer = &AuthHandler{}

type AuthHandler struct {
	lgr *zap.Logger
	auth.UnimplementedAuthServiceServer
	authUseCase primary_ports.UseCase
	validator   *protovalidate.Validator
}

func NewAuthHandler(lgr *zap.Logger, authUseCase primary_ports.UseCase) (*AuthHandler, error) {
	validator, err := protovalidate.New()
	if err != nil {
		return nil, err
	}

	return &AuthHandler{
		lgr:         lgr,
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

	tokensPair, err := h.authUseCase.LoginUser(ctx, request.Email, request.Password, models.FingerPrint{ // ToDo: need move to fabric
		"addr": fingerprint.Addr.String(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "login user")
	}

	return &auth.LoginResponse{
		AccessToken:  tokensPair.AccessToken,
		RefreshToken: tokensPair.RefreshToken,
	}, nil
}

func (h *AuthHandler) Register(ctx context.Context, request *auth.RegisterRequest) (*emptypb.Empty, error) {
	if err := h.validator.Validate(request); err != nil {
		return nil, errors.Wrap(err, "request validation")
	}

	_, err := h.authUseCase.RegisterUser(ctx, request.Email, request.Password)
	if err != nil {
		return nil, errors.Wrap(err, "register handler")
	}

	return &emptypb.Empty{}, nil
}

func (h *AuthHandler) RefreshAccessToken(ctx context.Context, request *emptypb.Empty) (*auth.RefreshTokenResponse, error) {
	if err := h.validator.Validate(request); err != nil {
		return nil, errors.Wrap(err, "request validation")
	}

	values := metadata.ValueFromIncomingContext(ctx, "refresh-token") // ToDo need move to another place

	tokensPair, err := h.authUseCase.RefreshToken(ctx, values[0])
	if err != nil {
		return nil, errors.Wrap(err, "refresh access token")
	}

	return &auth.RefreshTokenResponse{
		AccessToken:  tokensPair.AccessToken,
		RefreshToken: tokensPair.RefreshToken,
	}, nil
}

func (h *AuthHandler) Logout(ctx context.Context, request *auth.LogoutRequest) (*auth.LogoutResponse, error) {
	if err := h.validator.Validate(request); err != nil {
		return nil, errors.Wrap(err, "request validation")
	}

	var logoutType = primary_ports.LogOutType(request.Type)
	sessions, err := h.authUseCase.Logout(ctx, logoutType)
	if err != nil {
		return nil, errors.Wrap(err, "close sessions")
	}

	var pbSessions []*auth.Session
	for _, session := range sessions {
		pbSessions = append(pbSessions, &auth.Session{
			Id:        uint64(*session.ID),
			UserId:    uint64(session.UserID),
			CreatedAt: session.CreatedAt.String(),
			ClosedAt:  session.ClosedAt.String(),
		})
	}

	return &auth.LogoutResponse{
		Type:     request.Type,
		Sessions: pbSessions,
	}, nil
}
