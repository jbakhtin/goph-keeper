package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/jbakhtin/goph-keeper/internal/server/config"
	"github.com/jbakhtin/goph-keeper/internal/server/core/implements/adapters/input/grpc/v1"
	"github.com/jbakhtin/goph-keeper/internal/server/core/implements/adapters/input/grpc/v1/handlers"
	"github.com/jbakhtin/goph-keeper/internal/server/core/implements/adapters/input/grpc/v1/interceptors"
	"github.com/jbakhtin/goph-keeper/internal/server/core/implements/adapters/output/postgres/v1"
	"github.com/jbakhtin/goph-keeper/internal/server/core/implements/adapters/output/postgres/v1/repositories"
	"github.com/jbakhtin/goph-keeper/internal/server/core/implements/appservices/v1"
	"github.com/jbakhtin/goph-keeper/internal/server/core/implements/domainservice/v1"
	"github.com/jbakhtin/goph-keeper/internal/server/core/implements/usecase/v1"

	"github.com/go-faster/errors"

	"github.com/jbakhtin/rtagent/pkg/closer"
	"google.golang.org/grpc/reflection"
)

var (
	cfg      *config.Config
	pgClient *postgres.Postgres
	server   *grpc.Server
	clr      *closer.Closer
)

func accessibleRoles() map[string][]string {
	const authService = "/v1.auth.AuthService/"

	return map[string][]string{
		//authService + "Login": {},
		//authService + "Register": {},
		authService + "RefreshToken": {},
		authService + "Logout":       {},
	}
}

func init() {
	var err error

	cfg, err = config.New().ParseEnv().Build()
	if err != nil {
		log.Fatal(err)
	}

	pgClient, err = postgres.New(cfg) // ToDo: need to move inside the constructor
	if err != nil {
		log.Fatal(err)
	}

	userRepo, err := repositories.NewUserRepository(*pgClient)
	if err != nil {
		log.Fatal(err)
	}

	sessionRepo, err := repositories.NewSessionRepository(*pgClient)
	if err != nil {
		log.Fatal(err)
	}

	userDomainService, err := domainservice.NewUserDomainService(*cfg, userRepo)
	if err != nil {
		log.Fatal(err)
	}

	sessionDomainService, err := domainservice.NewSessionDomainService(*cfg, sessionRepo)
	if err != nil {
		log.Fatal(err)
	}

	accessTokenAppService, err := appservices.NewAccessTokenAppService(*cfg)
	if err != nil {
		log.Fatal(err)
	}

	passwordAppService, err := appservices.NewPasswordAppService(*cfg)
	if err != nil {
		log.Fatal(err)
	}

	authUseCases, err := usecase.NewAuthUseCase(*cfg,
		userDomainService,
		sessionDomainService,
		passwordAppService,
		accessTokenAppService,
		sessionRepo,
	)
	if err != nil {
		log.Fatal(err)
	}

	authHandler, err := handlers.NewAuthHandler(authUseCases)
	if err != nil {
		log.Fatal(err)
	}

	authInterceptor, err := interceptors.NewAuthInterceptor(*cfg, accessibleRoles())
	if err != nil {
		log.Fatal(err)
	}

	server, err = grpc.NewServer(*cfg, authHandler, grpc.WithUnaryInterceptor(authInterceptor.Unary))
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
	withTimeout, cancelShutdownProc := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelShutdownProc()

	if err := clr.Close(withTimeout); err != nil {
		log.Fatal(errors.Wrap(err, "shutdown"))
	}
}
