package grpc

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"net"
)

type IConfig interface {

}

type Server struct {
	grpc.Server
	cfg IConfig
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