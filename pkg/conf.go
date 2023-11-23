package pkg

import (
	"context"
	"encoding/json"

	"github.com/caarlos0/env/v10"
)

var Conf = &struct {
	Account struct {
		Username string `env:"ACCOUNT_USERNAME"`
		Password string `env:"ACCOUNT_PASSWORD"`
	}
}{}

func init() {
	if err := env.Parse(Conf); err != nil {
		Logger.Fatal(context.Background(), "failed to parse config", err)
	}
	// log the config as pretty-printed JSON
	tmp, err := json.MarshalIndent(Conf, "", "  ")
	if err != nil {
		Logger.Fatal(context.Background(), "failed to marshal config", err)
	}
	Logger.Info(context.Background(), "", "config", string(tmp))
}
