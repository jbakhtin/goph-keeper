package main

import (
	"context"
	"github.com/jbakhtin/goph-keeper/internal/server/implements/adapters/output/logger/v1/zap"
	"github.com/jbakhtin/goph-keeper/internal/server/interfaces/ports/input/config/v1"
	"github.com/jbakhtin/goph-keeper/internal/server/interfaces/ports/output/logger/v1"
	"os/signal"
	"syscall"

	"github.com/jbakhtin/goph-keeper/internal/server/implements/adapters/output/repositories/postgres/v1"
	"github.com/jbakhtin/goph-keeper/internal/server/implements/adapters/output/repositories/postgres/v1/repositories"

	"github.com/jbakhtin/goph-keeper/internal/server/implements/adapters/input/config/v1/drivers"
	"github.com/jbakhtin/goph-keeper/internal/server/implements/adapters/input/grpc/v1"
	"github.com/jbakhtin/goph-keeper/internal/server/implements/adapters/input/grpc/v1/handlers"
	"github.com/jbakhtin/goph-keeper/internal/server/implements/adapters/input/grpc/v1/interceptors"
	"github.com/jbakhtin/goph-keeper/internal/server/implements/appservices/v1"
	"github.com/jbakhtin/goph-keeper/internal/server/implements/domainservice/v1"
	"github.com/jbakhtin/goph-keeper/internal/server/implements/usecase/v1"

	"github.com/jbakhtin/rtagent/pkg/closer"
	"google.golang.org/grpc/reflection"
)

var (
	server   *grpc.Server
	clr      *closer.Closer
	lgr logger.Interface
	cfg config.Interface
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

	cfg, err = drivers.NewFormENV()
	if err != nil {
		panic(err)
	}

	lgr, err = zap.NewLogger(cfg)
	if err != nil {
		panic(err)
	}

	lgr.Info("app configuration file", cfg)

	pgClient, err := postgres.New(cfg) // ToDo: need to move inside the constructor
	if err != nil {
		panic(err)
	}

	userRepo, err := repositories.NewUserRepository(lgr, *pgClient)
	if err != nil {
		panic(err)
	}

	sessionRepo, err := repositories.NewSessionRepository(lgr, *pgClient)
	if err != nil {
		panic(err)
	}

	userDomainService, err := domainservice.NewUserDomainService(cfg, lgr, userRepo)
	if err != nil {
		panic(err)
	}

	sessionDomainService, err := domainservice.NewSessionDomainService(cfg, lgr,  sessionRepo)
	if err != nil {
		panic(err)
	}

	accessTokenAppService, err := appservices.NewAccessTokenAppService(cfg, lgr)
	if err != nil {
		panic(err)
	}

	passwordAppService, err := appservices.NewPasswordAppService(cfg, lgr)
	if err != nil {
		panic(err)
	}

	authUseCases, err := usecase.NewAuthUseCase(
		cfg,
		lgr,
		userDomainService,
		sessionDomainService,
		passwordAppService,
		accessTokenAppService,
		sessionRepo,
		userRepo,
	)
	if err != nil {
		panic(err)
	}

	authHandler, err := handlers.NewAuthHandler(lgr, authUseCases)
	if err != nil {
		panic(err)
	}

	authInterceptor, err := interceptors.NewAuthInterceptor(cfg, lgr, accessibleRoles())
	if err != nil {
		panic(err)
	}

	server, err = grpc.NewServer(cfg, lgr, authHandler, grpc.WithUnaryInterceptor(authInterceptor.Unary))
	if err != nil {
		panic(err)
	}

	reflection.Register(server)

	if clr, err = closer.New().WithFuncs(server.Shutdown).Build(); err != nil {
		panic(err)
	}
}

func main() {
	osCtx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer cancel()

	if err := server.Start(osCtx); err != nil {
		lgr.Error(err.Error())
	}

	lgr.Info("server started", &server)

	// Gracefully shut down
	<-osCtx.Done()
	withTimeout, cancelShutdownProc := context.WithTimeout(context.Background(), cfg.GetShutdownTimeout())
	defer cancelShutdownProc()

	if err := clr.Close(withTimeout); err != nil {
		lgr.Fatal(err.Error(), nil)
	}
}
