package drivers

import (
	"github.com/caarlos0/env/v6"
	"github.com/jbakhtin/goph-keeper/internal/server/config"
)

func NewConfigFormENV() (*config.Config, error) {
	var cfg config.Config
	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
