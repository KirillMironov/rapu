package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port                string `default:"7002" envconfig:"PORT"`
	UsersServiceAddress string `required:"true" envconfig:"USERS_SERVICE_ADDRESS"`
	PostsServiceAddress string `required:"true" envconfig:"POSTS_SERVICE_ADDRESS"`
}

func InitConfig() (Config, error) {
	var cfg Config
	return cfg, envconfig.Process("", &cfg)
}
