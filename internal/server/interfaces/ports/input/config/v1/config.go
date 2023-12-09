package config

import "time"

type Interface interface {
	GetterInterface
	SetterInterface
}

type GetterInterface interface {
	GetAppKey() string
	GetAppEnv() string
	GetShutdownTimeout() time.Duration

	GetDataBaseDSN() string
	GetDataBaseDriver() string
	GetSessionExpire() time.Duration
	GetAccessTokenExpire() time.Duration
	GetGRPCServerAddress() string

	GetLoggerFileDirectory() string
	GetLoggerFileMaxSize() int
	GetLoggerFileMaxBackups() int
	GetLoggerFileMaxAge() int
	GetLoggerFileCompress() bool
}

type SetterInterface interface {
	SetAppKey(string)
	SetAppEnv(string)
	SetShutdownTimeout(time.Duration)

	SetDataBaseDSN(string)
	SetDataBaseDriver(string)
	SetSessionExpire(time.Duration)
	SetAccessTokenExpire(time.Duration)
	SetGRPCServerAddress(string)

	SetLoggerFileDirectory(string)
	SetLoggerFileMaxSize(int)
	SetLoggerFileMaxBackups(int)
	SetLoggerFileMaxAge(int)
	SetLoggerFileCompress(bool)
}
