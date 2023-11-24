package main

import (
	"context"
	"github.com/go-faster/errors"
	"github.com/jbakhtin/goph-keeper/internal/server/application/config"
	authService "github.com/jbakhtin/goph-keeper/internal/server/domain/services/auth"
	sessionService "github.com/jbakhtin/goph-keeper/internal/server/domain/services/session"
	"github.com/jbakhtin/goph-keeper/internal/server/infastructure/database/postgres"
	sessionRepo "github.com/jbakhtin/goph-keeper/internal/server/infastructure/database/postgres/repositories/session"
	userRepo "github.com/jbakhtin/goph-keeper/internal/server/infastructure/database/postgres/repositories/user"
	grpcServer "github.com/jbakhtin/goph-keeper/internal/server/infastructure/network/grpc"
	authServer "github.com/jbakhtin/goph-keeper/internal/server/infastructure/network/grpc/handlers/v1/auth"
	"github.com/jbakhtin/goph-keeper/internal/server/infastructure/network/grpc/interceptors/auth"
	"github.com/jbakhtin/rtagent/pkg/closer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"os/signal"
	"syscall"
	"time"
)

var (
	cfg  *config.Config
	pgClient *postgres.Postgres
	server *grpcServer.Server
	clr  *closer.Closer
)

func accessibleRoles() map[string][]string {
	const authService = "/v1.auth.AuthService/"

	return map[string][]string{
		//authService + "Login": {},
		//authService + "Register": {},
		authService + "RefreshToken": {},
		authService + "Logout": {},
	}
}

func init() {
	var err error

	cfg, err = config.New().ParseEnv().Build()
	if err != nil {
		log.Fatal(err)
	}

	pgClient, err = postgres.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	userRepo, err := userRepo.New(*pgClient)
	if err != nil {
		log.Fatal(err)
	}

	sessionRepo, err := sessionRepo.New(*pgClient)
	if err != nil {
		log.Fatal(err)
	}

	authService, err := authService.New(cfg, userRepo)
	if err != nil {
		log.Fatal(err)
	}

	sessionService, err := sessionService.New(cfg, sessionRepo)
	if err != nil {
		log.Fatal(err)
	}

	authServer, err := authServer.New(*authService, *sessionService).Build()
	if err != nil {
		log.Fatal(err)
	}

	authInterceptor := auth.Interceptor{
		Cfg: cfg,
		AccessibleRoles: accessibleRoles(),
	}

	serverOptions := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(authInterceptor.Unary),
	}

	server, err = grpcServer.New(cfg, serverOptions...).
		WithAuthHandler(authServer).
		Build()
	if err != nil {
		log.Fatal(err)
	}

	reflection.Register(server)

	if clr, err = closer.New().WithFuncs(server.Shutdown).Build(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	osCtx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer cancel()

	if err := server.Start(osCtx); err != nil {
		log.Fatal(errors.Wrap(err, "start server"))
	}

	// Gracefully shut down
	<-osCtx.Done()
	withTimeout, cancelShutdownProc := context.WithTimeout(context.Background(), time.Second * 10)
	defer cancelShutdownProc()

	if err := clr.Close(withTimeout); err != nil {
		log.Fatal(errors.Wrap(err, "shutdown"))
	}
}
