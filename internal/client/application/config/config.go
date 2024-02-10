package config

type Config struct {
	appKey string `env:"APP_KEY"`
}

func (c *Config) SetAppKey(value string) {
	c.AppKey = value
}

func (c Config) GetAppKey() string {
	return c.AppKey
}
