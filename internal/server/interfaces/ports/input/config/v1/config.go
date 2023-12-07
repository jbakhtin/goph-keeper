package config

import "time"

type Interface interface {
	GetterInterface
	SetterInterface
}

type GetterInterface interface {
	GetAppKey() string
	GetDataBaseDSN() string
	GetDataBaseDriver() string
	GetSessionExpire() time.Duration
	GetAccessTokenExpire() time.Duration
	GetShutdownTimeout() time.Duration
	GetGRPCServerAddress() string
}

type SetterInterface interface {
	SetAppKey(string)
	SetDataBaseDSN(string)
	SetDataBaseDriver(string)
	SetSessionExpire(time.Duration)
	SetAccessTokenExpire(time.Duration)
	SetGRPCServerAddress(string)
	SetShutdownTimeout(time.Duration)
}
