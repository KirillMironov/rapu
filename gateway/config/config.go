package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port          string `default:"7002" envconfig:"PORT"`
	ServerAddress string `required:"true" envconfig:"SERVER_ADDRESS"`
}

func InitConfig() (Config, error) {
	var cfg Config
	return cfg, envconfig.Process("", &cfg)
}
