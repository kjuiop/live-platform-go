package config

import (
	"github.com/kelseyhightower/envconfig"
)

type EnvConfig struct {
	Logger Logger
	Server Server
	Policy Policy
	Redis  Redis
}
type Logger struct {
	Level       string `envconfig:"LOG_LEVEL" default:"debug"`
	Path        string `envconfig:"LOG_PATH" default:"./logs/access.log"`
	PrintStdOut bool   `envconfig:"LOG_STDOUT" default:"true"`
}

type Server struct {
	Mode           string `envconfig:"ENV" default:"dev"`
	Port           string `envconfig:"SERVER_PORT" default:"8080"`
	TrustedProxies string `envconfig:"TRUSTED_PROXIES" default:"127.0.0.1/32"`
}

type Redis struct {
	Mode     string `envconfig:"LCK_REDIS_MODE" default:"single"`
	Addr     string `envconfig:"LCK_REDIS_ADDR" default:":6379"`
	Password string `envconfig:"LCK_REDIS_PASSWORD"`
	Masters  string `envconfig:"LCK_REDIS_MASTERS"`
	PoolSize int    `envconfig:"LCK_REDIS_POOL_SIZE" default:"100"`
}

type Policy struct {
	ContextTimeout int `envconfig:"POLICY_CONTEXT_TIMEOUT" default:"60"`
}

func LoadEnvConfig() (*EnvConfig, error) {
	var config EnvConfig
	if err := envconfig.Process("lpg", &config); err != nil {
		return nil, err
	}
	return &config, nil
}
