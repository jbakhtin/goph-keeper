package grpc

import (
	"context"
	"net"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type Config interface {
	GetGRPCServerAddress() string
	GetGRPCServerNetwork() string
}

type Server struct {
	cfg Config
	grpc.Server
}

type Option func([]grpc.ServerOption) []grpc.ServerOption

func WithUnaryInterceptor(unaryInterceptor grpc.UnaryServerInterceptor) Option {
	return func(options []grpc.ServerOption) []grpc.ServerOption {
		return append(options, grpc.ChainUnaryInterceptor(unaryInterceptor))
	}
}

func NewServer(cfg Config, options ...Option) (*Server, error) {
	var serverOptions []grpc.ServerOption
	for _, option := range options {
		serverOptions = option(serverOptions)
	}

	return &Server{
		cfg:    cfg,
		Server: *grpc.NewServer(serverOptions...),
	}, nil
}

func (s *Server) Start(ctx context.Context) (err error) {
	listen, err := net.Listen(s.cfg.GetGRPCServerNetwork(), s.cfg.GetGRPCServerAddress())
	if err != nil {
		return err
	}

	go func() { // ToDo: need to add context handle
		if err = s.Serve(listen); err != nil {
			err = errors.Wrap(err, "serve")
		}
	}()

	return err
}

func (s *Server) Shutdown(ctx context.Context) error {
	newCtx, finish := context.WithCancel(ctx)

	go func() {
		s.GracefulStop()
		finish()
	}()

	select {
	case <-newCtx.Done():
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
