package auth

import (
	"context"
	"fmt"
	"github.com/go-faster/errors"
	"github.com/jbakhtin/goph-keeper/internal/server/application/types"
	"github.com/jbakhtin/goph-keeper/internal/server/domain/models"
	"github.com/jbakhtin/goph-keeper/internal/server/domain/services/auth"
	"github.com/jbakhtin/goph-keeper/internal/server/domain/services/session"
	pb "github.com/jbakhtin/goph-keeper/internal/server/infastructure/network/grpc/gen/auth/v1"
	"google.golang.org/grpc/peer"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

type Handler struct {
	pb.UnimplementedAuthServiceServer
	authService auth.Service
	sessionService session.Service
}

func (h *Handler) Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error)  {
	user, err := h.authService.LoginUser(ctx, request.Email, request.Password)
	if err != nil {
		return nil, err
	}

	_, ok := peer.FromContext(ctx)
	if !ok {
		fmt.Println("peer is not ok")
	}

	tokensPair, err := h.sessionService.Create(ctx, *user, types.FingerPrint{
		"test": "test",
	})
	if err != nil {
		return nil, err
	}

	return &pb.LoginResponse{
		AccessToken: tokensPair.AccessToken,
		RefreshToken: tokensPair.RefreshToken,
	}, nil
}

func (h *Handler) Register(ctx context.Context, request *pb.RegisterRequest) (*emptypb.Empty, error)  {
	fmt.Println("test:", request)
	user := models.User{
		Email: request.Email,
		Password: request.Password,
	}

	_, err := h.authService.RegisterUser(ctx, user)
	if err != nil {
		return nil, errors.Wrap(err, "register handlers")
	}

	return &emptypb.Empty{}, nil
}

func (h *Handler) RefreshAccessToken(ctx context.Context, request *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	fmt.Println(request.RefreshToken)
	tokensPair, err := h.sessionService.Update(ctx, request.RefreshToken)
	if err != nil {
		return nil, errors.Wrap(err, "refresh access token handlers")
	}

	return &pb.RefreshTokenResponse{
		AccessToken: tokensPair.AccessToken,
		RefreshToken: tokensPair.RefreshToken,
	}, nil
}

func (h *Handler) Logout(ctx context.Context, request *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	var logoutType models.LogoutType
	switch request.Type {
		case pb.LogoutType_TYPE_ALL:
			logoutType = models.LogoutType_ALL
		case pb.LogoutType_TYPE_UNSPECIFIED:
			logoutType = models.LogoutType_THIS
	default:
		return nil, errors.New("not allowed logout type")
	}

	sessions, err := h.sessionService.Close(ctx, logoutType)
	if err != nil {
		return nil, err
	}

	var pbSessions []*pb.Session

	for _, session := range sessions {
		pbSessions = append(pbSessions, &pb.Session{
			SessionId: uint64(*session.ID),
			UserId: uint64(session.UserId),
			CreatedAt: time.Time(*session.CreatedAt).String(),
			ClosedAt: time.Time(*session.ClosedAt).String(),
		})
	}

	return &pb.LogoutResponse{
		Type: request.Type,
		Sessions: pbSessions,
	}, nil
}