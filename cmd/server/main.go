package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/jbakhtin/goph-keeper/internal/server/appmodules/secrets"
	keyvalue2 "github.com/jbakhtin/goph-keeper/internal/server/storage/postgres/specifications/keyvalue"

	auth2 "github.com/jbakhtin/goph-keeper/gen/go/v1/auth"
	secrets2 "github.com/jbakhtin/goph-keeper/gen/go/v1/secrets"
	"github.com/jbakhtin/goph-keeper/internal/server/appmodules/auth"
	"github.com/jbakhtin/goph-keeper/internal/server/config"
	"github.com/jbakhtin/goph-keeper/internal/server/config/drivers"
	"github.com/jbakhtin/goph-keeper/internal/server/logger/zap"
	"github.com/jbakhtin/goph-keeper/internal/server/storage/postgres"
	"github.com/jbakhtin/goph-keeper/internal/server/storage/postgres/repositories"
	"github.com/jbakhtin/goph-keeper/internal/server/storage/postgres/specifications/session"
	"github.com/jbakhtin/goph-keeper/internal/server/storage/postgres/specifications/user"
	"github.com/jbakhtin/goph-keeper/internal/server/webserver/grpc"
	"github.com/jbakhtin/goph-keeper/internal/server/webserver/grpc/handlers"
	"github.com/jbakhtin/goph-keeper/internal/server/webserver/grpc/interceptors"

	"github.com/jbakhtin/rtagent/pkg/closer"
)

var (
	server *grpc.Server
	clr    *closer.Closer
	lgr    *zap.Logger
	cfg    *config.Config
)

// accessibleRoles возвращает список GRPC обработчиков которые должны быть проверены на аутентификацию пользователя
// NOTE список нужно обновлять при добавлении новых обработчиков, если требуется
func accessibleRoles() map[string][]string {
	const authService = "/v1.auth.AuthService/"
	const keyValueService = "/v1.secrets.SecretsService/"

	return map[string][]string{
		authService + "RefreshToken": {},
		authService + "Logout":       {},
		keyValueService + "Create":   {},
	}
}

func init() {
	var err error

	if cfg, err = drivers.NewConfigFormENV(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if lgr, err = zap.NewLogger(cfg); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	sqlClient, err := postgres.NewSQLClient(cfg)
	if err != nil {
		lgr.Fatal(err.Error())
	}

	// Init repositories
	// NOTE there are we can implement switch repositories according to env

	userRepository, err := repositories.NewUserRepository(lgr, sqlClient)
	if err != nil {
		lgr.Fatal(err.Error())
	}

	sessionRepository, err := repositories.NewSessionRepository(lgr, sqlClient)
	if err != nil {
		lgr.Fatal(err.Error())
	}

	keyValueRepository, err := repositories.NewKeyValueRepository(lgr, sqlClient)
	if err != nil {
		lgr.Fatal(err.Error())
	}

	// Init query specifications

	sessionQuerySpecification, err := session.NewSessionQuerySpecification()
	if err != nil {
		lgr.Fatal(err.Error())
	}

	userQuerySpecification, err := user.NewUserQuerySpecification()
	if err != nil {
		lgr.Fatal(err.Error())
	}

	keyValueQuerySpecification, err := keyvalue2.NewKeyValueQuerySpecification() // ToDo: need refactoring
	if err != nil {
		lgr.Fatal(err.Error())
	}

	// Init app modules
	// NOTE The app modules are built on a hexagonal architecture

	authModule, err := auth.NewModule(cfg, lgr, userRepository, sessionRepository, sessionQuerySpecification, userQuerySpecification)
	if err != nil {
		lgr.Fatal(err.Error())
	}

	keyValueModule, err := secrets.NewModule(cfg, lgr, keyValueRepository, keyValueQuerySpecification)
	if err != nil {
		lgr.Fatal(err.Error())
	}

	// GRPC Handlers

	authHandler, err := handlers.NewAuthHandler(lgr, authModule.GetUseCase())
	if err != nil {
		lgr.Fatal(err.Error())
	}

	keyValueHandler, err := handlers.NewKeyValueHandler(lgr, keyValueModule.GetUseCase())
	if err != nil {
		lgr.Fatal(err.Error())
	}

	// GRPC Interceptors
	authInterceptor, err := interceptors.NewAuthInterceptor(cfg, accessibleRoles())
	if err != nil {
		lgr.Fatal(err.Error())
	}

	// Init GRPC Server and connect handlers to it
	if server, err = grpc.NewServer(cfg, grpc.WithUnaryInterceptor(authInterceptor.Unary)); err != nil {
		lgr.Fatal(err.Error())
	}

	auth2.RegisterAuthServiceServer(server, authHandler) // ToDo: need to move to another place
	secrets2.RegisterSecretsServiceServer(server, keyValueHandler)

	// Init closer
	if clr, err = closer.New().WithFuncs(server.Shutdown).Build(); err != nil {
		lgr.Fatal(err.Error())
	}
}

func main() {
	osCtx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer cancel()

	if err := server.Start(osCtx); err != nil {
		lgr.Fatal(err.Error())
	}

	lgr.Info("server started", &server) // ToDo: display logs with Debug level

	// Gracefully shut down
	<-osCtx.Done()
	withTimeout, cancelShutdownProc := context.WithTimeout(context.Background(), cfg.GetShutdownTimeout())
	defer cancelShutdownProc()

	if err := clr.Close(withTimeout); err != nil {
		lgr.Fatal(err.Error(), nil)
	}
}
