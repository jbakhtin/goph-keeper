package drivers

import (
	"flag"
	"github.com/jbakhtin/goph-keeper/internal/server/config"
	"time"
)

const (
	configSessionExpireIntervalFlag     = "session_expire_interval"
	configAccessTokenExpireIntervalFlag = "access_token_expire_interval"
	configGRPCServerAddressFlag         = "grpc_server_address"
	configAppKeyFlag                    = "app_key"
	configDatabaseDSNFlag               = "database_dsn"
	configDatabaseDriverFlag            = "database_driver"
	configShutdownTimeoutFlag           = "shutdown_timeout"
)

var configDefaultValues map[string]any = map[string]any{
	configSessionExpireIntervalFlag:     time.Minute * 60 * 24 * 30,
	configAccessTokenExpireIntervalFlag: time.Minute * 30,
	configGRPCServerAddressFlag:         "127.0.0.1:8080",
	configAppKeyFlag:                    "",
	configDatabaseDSNFlag:               "",
	configDatabaseDriverFlag:            "pgx",
	configShutdownTimeoutFlag:           time.Second * 10,
}

var configDefaultUsage map[string]string = map[string]string{
	configSessionExpireIntervalFlag:     "Период истечения сессии клинта",
	configAccessTokenExpireIntervalFlag: "Период истечения JWT токена",
	configGRPCServerAddressFlag:         "Адрес gRPC сервера",
	configAppKeyFlag:                    "Ключ приложения",
	configDatabaseDSNFlag:               "DSN строка подключения к базе данных",
	configDatabaseDriverFlag:            "Драйвер базы данных",
	configShutdownTimeoutFlag:           "Время на остановку приложения (graceful shutdown)",
}

func NewConfigFormFlags() (*config.Config, error) {
	sessionExpireInterval := flag.Duration(
		configSessionExpireIntervalFlag,
		configDefaultValues[configSessionExpireIntervalFlag].(time.Duration),
		configDefaultUsage[configSessionExpireIntervalFlag])

	accessTokenExpireInterval := flag.Duration(
		configAccessTokenExpireIntervalFlag,
		configDefaultValues[configAccessTokenExpireIntervalFlag].(time.Duration),
		configDefaultUsage[configAccessTokenExpireIntervalFlag])

	grpcServerAddress := flag.String(
		configGRPCServerAddressFlag,
		configDefaultValues[configGRPCServerAddressFlag].(string),
		configDefaultUsage[configGRPCServerAddressFlag])

	appKey := flag.String(
		configAppKeyFlag,
		configDefaultValues[configAppKeyFlag].(string),
		configDefaultUsage[configAppKeyFlag])

	databaseDSN := flag.String(
		configDatabaseDSNFlag,
		configDefaultValues[configDatabaseDSNFlag].(string),
		configDefaultUsage[configDatabaseDSNFlag])

	databaseDriver := flag.String(
		configDatabaseDriverFlag,
		configDefaultValues[configDatabaseDriverFlag].(string),
		configDefaultUsage[configDatabaseDriverFlag])

	shutdownTimeout := flag.Duration(
		configShutdownTimeoutFlag,
		configDefaultValues[configShutdownTimeoutFlag].(time.Duration),
		configDefaultUsage[configShutdownTimeoutFlag])

	flag.Parse()

	var cfg config.Config

	cfg.SetSessionExpire(*sessionExpireInterval)
	cfg.SetAccessTokenExpire(*accessTokenExpireInterval)
	cfg.SetAppKey(*appKey)
	cfg.SetDataBaseDSN(*databaseDSN)
	cfg.SetDataBaseDriver(*databaseDriver)
	cfg.SetGRPCServerAddress(*grpcServerAddress)
	cfg.SetShutdownTimeout(*shutdownTimeout)

	return &cfg, nil
}
