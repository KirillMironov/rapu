package config

import (
	"github.com/kelseyhightower/envconfig"
	"time"
)

type Config struct {
	Port string `default:"7001" envconfig:"PORT"`

	Postgres struct {
		ConnectionString string `required:"true" envconfig:"POSTGRES_CONNECTION_STRING"`
	}

	Security struct {
		TokenTTL time.Duration `required:"true" envconfig:"TOKEN_TTL"`
		JWTKey   string        `required:"true" envconfig:"JWT_KEY"`
	}
}

func InitConfig() (Config, error) {
	var cfg Config
	return cfg, envconfig.Process("", &cfg)
}
