package main

import (
	"github.com/caarlos0/env/v6"
)

// Config confuguration struct.
type Config struct {
	Logger struct {
		Level string `env:"LOG_LEVEL" envDefault:"info"`
	}
	DB struct {
		Username  string `env:"DB_USERNAME,required"`
		Password  string `env:"DB_PASSWORD,required"`
		Host      string `env:"DB_HOST" envDefault:"localhost"`
		Port      int    `env:"DB_PORT" envDefault:"5432"`
		SSLEnable bool   `env:"DB_SSL_ENABLE" envDefault:"false"`
	}
}

// NewConfig returns configuration.
func NewConfig() (*Config, error) {
	cfg := Config{}
	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
