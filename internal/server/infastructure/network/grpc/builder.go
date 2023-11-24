package grpc

import (
	pb "github.com/jbakhtin/goph-keeper/internal/server/infastructure/network/grpc/gen/auth/v1"
	"github.com/jbakhtin/goph-keeper/internal/server/infastructure/network/grpc/handlers/v1/auth"
	"google.golang.org/grpc"
)

type Builder struct {
	server *Server
	err error
}

func New(config IConfig, options ...grpc.ServerOption) *Builder {
	return &Builder{
		server: &Server{
			*grpc.NewServer(options...),
			config,
		},
		err: nil,
	}
}

func (b *Builder) WithAuthHandler(handler *auth.Handler) *Builder {
	pb.RegisterAuthServiceServer(b.server, handler)
	return b
}

func (b *Builder) Build() (*Server, error) {
	if b.err != nil {
		return nil, b.err
	}

	return b.server, nil
}

