package common

import (
	"github.com/caarlos0/env"
)

type Environment struct {
	ServerAddress string `env:"SERVER_URL,required"`
	User          string `env:"USER,required"`
}

func NewEnvironment() (*Environment, error) {
	cfg := new(Environment)
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
