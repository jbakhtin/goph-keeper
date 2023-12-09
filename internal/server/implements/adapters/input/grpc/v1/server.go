package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/jbakhtin/goph-keeper/internal/server/interfaces/ports/output/logger/v1"

	"github.com/jbakhtin/goph-keeper/internal/server/interfaces/ports/input/config/v1"
	"github.com/jbakhtin/goph-keeper/internal/server/interfaces/ports/input/grpc/v1/auth"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type Server struct {
	cfg config.Interface
	lgr logger.Interface
	grpc.Server
}

type Option func([]grpc.ServerOption) []grpc.ServerOption

func WithUnaryInterceptor(unaryInterceptor grpc.UnaryServerInterceptor) Option {
	return func(options []grpc.ServerOption) []grpc.ServerOption {
		return append(options, grpc.ChainUnaryInterceptor(unaryInterceptor))
	}
}

func NewServer(cfg config.Interface, lgr logger.Interface, authHandler auth.AuthServiceServer, setters ...Option) (*Server, error) {
	var serverOptions []grpc.ServerOption
	for _, setter := range setters {
		serverOptions = setter(serverOptions)
	}

	server := &Server{
		cfg:    cfg,
		lgr:    lgr,
		Server: *grpc.NewServer(serverOptions...),
	}

	auth.RegisterAuthServiceServer(server, authHandler)
	// ...

	return server, nil
}

func (s *Server) Start(ctx context.Context) (err error) {
	listen, err := net.Listen("tcp", ":3200")
	if err != nil {
		return err
	}

	go func() {
		if err = s.Serve(listen); err != nil {
			err = errors.Wrap(err, "serve")
		}
	}()

	return err
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.GracefulStop()
	fmt.Println("grpc server is stopped")
	return nil
}
