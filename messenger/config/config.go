package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Port string `default:"7004" envconfig:"PORT"`

	Redis struct {
		Addr     string `default:"localhost:6379" envconfig:"REDIS_ADDR"`
		Password string `default:"" envconfig:"REDIS_PASSWORD"`
		DB       int    `default:"0" envconfig:"REDIS_DB"`
	}
}

func InitConfig() (Config, error) {
	var cfg Config
	return cfg, envconfig.Process("", &cfg)
}
