package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/pkg/errors"
)

type Builder struct {
	config *Config
	err    error
}

func New() *Builder {
	return &Builder{
		&Config{},
		nil,
	}
}

func (b Builder) Build() (*Config, error) {
	if b.err != nil {
		return nil, errors.Wrap(b.err, "build config")
	}

	return b.config, nil
}

func (b *Builder) ParseEnv() *Builder {
	err := env.Parse(b.config)
	if err != nil {
		b.err = err
	}

	return b
}
