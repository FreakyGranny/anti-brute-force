package main

import (
	"fmt"

	"github.com/caarlos0/env/v6"
)

// Config confuguration struct.
type Config struct {
	GRPC struct {
		Host string `env:"GRPC_HOST" envDefault:"localhost"`
		Port int    `env:"GRPC_PORT" envDefault:"50051"`
	}
	Logger struct {
		Level string `env:"LOG_LEVEL" envDefault:"info"`
	}
	Limits struct {
		User     int `env:"LIMIT_USER" envDefault:"10"`
		Password int `env:"LIMIT_PASSWORD" envDefault:"100"`
		IP       int `env:"LIMIT_IP" envDefault:"1000"`
	}
	DB struct {
		Username  string `env:"DB_USERNAME,required"`
		Password  string `env:"DB_PASSWORD,required"`
		Host      string `env:"DB_HOST" envDefault:"localhost"`
		Port      int    `env:"DB_PORT" envDefault:"5432"`
		Name      string `env:"DB_NAME,required"`
		SSLEnable bool   `env:"DB_SSL_ENABLE" envDefault:"false"`
	}
	Redis struct {
		Password string `env:"REDIS_PASSWORD,required"`
		Host     string `env:"REDIS_HOST" envDefault:"localhost"`
		Port     int    `env:"REDIS_PORT" envDefault:"6379"`
	}
	IPListRefresh string `env:"IP_LIST_REFRESH" envDefault:"15s"`
}

// NewConfig returns configuration.
func NewConfig() (*Config, error) {
	cfg := Config{}
	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}
	fmt.Println(cfg.IPListRefresh)

	return &cfg, nil
}
