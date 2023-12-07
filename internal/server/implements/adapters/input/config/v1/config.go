package config

import "time"

type Config struct {
	AppKey          string        `env:"APP_KEY" default:""`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" default:"10s"`
	GRPCServer      struct {
		Address string `env:"GRPC_SERVER_ADDRESS" default:"127.0.0.1:8080"`
	}
	DB struct {
		DSN    string `env:"DATABASE_DSN" default:""`
		Driver string `env:"DATABASE_DRIVER" default:""`
	}
	Session struct {
		Expire time.Duration `env:"SESSION_EXPIRE" default:"720h"`
	}
	AccessToken struct {
		Expire time.Duration `env:"ACCESS_TOKEN_EXPIRE" default:"10m"`
	}
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
