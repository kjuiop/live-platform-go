package config

import (
	"github.com/kelseyhightower/envconfig"
)

type EnvConfig struct {
	Logger Logger
	Server Server
}
type Logger struct {
	Level       string `envconfig:"LCK_LOG_LEVEL" default:"debug"`
	Path        string `envconfig:"LCK_LOG_PATH" default:"./logs/access.log"`
	PrintStdOut bool   `envconfig:"LCK_LOG_STDOUT" default:"true"`
}

type Server struct {
	Mode           string `envconfig:"LCK_ENV" default:"dev"`
	Port           string `envconfig:"LCK_SERVER_PORT" default:"8090"`
	TrustedProxies string `envconfig:"LCK_TRUSTED_PROXIES" default:"127.0.0.1/32"`
}

func LoadEnvConfig() (*EnvConfig, error) {
	var config EnvConfig
	if err := envconfig.Process("lpg", &config); err != nil {
		return nil, err
	}
	return &config, nil
}
