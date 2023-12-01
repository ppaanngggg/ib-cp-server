package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/caarlos0/env/v10"
)

var conf = &struct {
	Account struct {
		Username string `env:"ACCOUNT_USERNAME"`
		Password string `env:"ACCOUNT_PASSWORD"`
	}
	Server struct {
		Host     string        `env:"SERVER_HOST" envDefault:"0.0.0.0"`
		Port     string        `env:"SERVER_PORT" envDefault:"8000"`
		Timeout  time.Duration `env:"SERVER_TIMEOUT" envDefault:"60s"`
		Throttle int           `env:"SERVER_THROTTLE" envDefault:"100"`
	}
}{}

func init() {
	if err := env.Parse(conf); err != nil {
		logger.Fatal(context.Background(), "failed to parse config", err)
	}
	// log the config as pretty-printed JSON
	tmp, err := json.MarshalIndent(conf, "", "  ")
	if err != nil {
		logger.Fatal(context.Background(), "failed to marshal config", err)
	}
	logger.Info(context.Background(), "", "config", string(tmp))
}
