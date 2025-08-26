package config

import (
	"time"

	"github.com/jbakhtin/goph-keeper/internal/appmodules/auth/services/accesstoken"
	"github.com/jbakhtin/goph-keeper/internal/appmodules/auth/services/password"
	"github.com/jbakhtin/goph-keeper/internal/logger/zap"
)

var _ accesstoken.Config = Config{}
var _ password.Config = Config{}
var _ zap.Config = Config{}

type Config struct {
	AppEnv          string        `env:"APP_ENV" envDefault:"development"`
	AppKey          string        `env:"APP_KEY" envDefault:""`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"10s"`
	GRPCServer      struct {
		Address string `env:"GRPC_SERVER_ADDRESS" envDefault:":3200"`
		Network string `env:"GRPC_SERVER_NETWORK" envDefault:"tcp"`
	}
	DB struct {
		DSN    string `env:"DATABASE_DSN" envDefault:""`
		Driver string `env:"DATABASE_DRIVER" envDefault:""`
	}
	Session struct {
		Expire time.Duration `env:"SESSION_EXPIRE" envDefault:"720h"`
	}
	AccessToken struct {
		Expire time.Duration `env:"ACCESS_TOKEN_EXPIRE" envDefault:"10m"`
	}
	Logger struct { // ToDo: продумать структуру конфига для логгера
		File struct {
			// Directory is the
			Directory string `env:"LOG_DIRECTORY" envDefault:"storage/logs/"`
			// MaxSize is the maximum size in megabytes of the log secrets before it gets rotated.
			MaxSize    int  `env:"LOGGER_FILE_MAX_SIZE" envDefault:"1"`
			MaxBackups int  `env:"LOGGER_FILE_MAX_BACKUPS" envDefault:"1"`
			MaxAge     int  `env:"LOGGER_FILE_MAX_AGE" envDefault:"1"`
			Compress   bool `env:"LOGGER_FILE_COMPRESS" envDefault:"true"`
		}
	}
}

func (c Config) GetLoggerFileDirectory() string {
	return c.Logger.File.Directory
}

func (c *Config) SetLoggerFileDirectory(s string) {
	c.Logger.File.Directory = s
}

func (c *Config) SetAppEnv(s string) {
	c.AppEnv = s
}

func (c Config) GetLoggerFileMaxSize() int {
	return c.Logger.File.MaxSize
}

func (c Config) GetLoggerFileMaxBackups() int {
	return c.Logger.File.MaxBackups
}

func (c Config) GetLoggerFileMaxAge() int {
	return c.Logger.File.MaxAge
}

func (c Config) GetLoggerFileCompress() bool {
	return c.Logger.File.Compress
}

func (c *Config) SetLoggerFileMaxSize(i int) {
	c.Logger.File.MaxSize = i
}

func (c *Config) SetLoggerFileMaxBackups(i int) {
	c.Logger.File.MaxBackups = i
}

func (c *Config) SetLoggerFileMaxAge(i int) {
	c.Logger.File.MaxAge = i
}

func (c *Config) SetLoggerFileCompress(b bool) {
	c.Logger.File.Compress = b
}

func (c Config) GetAppEnv() string {
	return c.AppEnv
}

func (c *Config) SetAppKey(s string) {
	c.AppKey = s
}

func (c *Config) SetDataBaseDSN(s string) {
	c.DB.DSN = s
}

func (c *Config) SetDataBaseDriver(s string) {
	c.DB.Driver = s
}

func (c *Config) SetSessionExpire(duration time.Duration) {
	c.Session.Expire = duration
}

func (c *Config) SetAccessTokenExpire(duration time.Duration) {
	c.AccessToken.Expire = duration
}

func (c *Config) SetGRPCServerAddress(s string) {
	c.GRPCServer.Address = s
}

func (c *Config) SetGRPCServerNetwork(s string) {
	c.GRPCServer.Network = s
}

func (c *Config) SetShutdownTimeout(duration time.Duration) {
	c.ShutdownTimeout = duration
}

func (c Config) GetAppKey() string {
	return c.AppKey
}

func (c Config) GetDataBaseDSN() string {
	return c.DB.DSN
}

func (c Config) GetDataBaseDriver() string {
	return c.DB.Driver
}

func (c Config) GetSessionExpire() time.Duration {
	return c.Session.Expire
}

func (c Config) GetAccessTokenExpire() time.Duration {
	return c.AccessToken.Expire
}

func (c Config) GetShutdownTimeout() time.Duration {
	return c.ShutdownTimeout
}

func (c Config) GetGRPCServerAddress() string {
	return c.GRPCServer.Address
}

func (c Config) GetGRPCServerNetwork() string {
	return c.GRPCServer.Network
}
