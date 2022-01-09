package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Port            string `default:"7003" envconfig:"PORT"`
	MaxPostsPerPage int64  `default:"10" envconfig:"MAX_POSTS_PER_PAGE"`

	Mongo struct {
		ConnectionString string `required:"true" envconfig:"MONGO_CONNECTION_STRING"`
		DBName           string `required:"true" envconfig:"MONGO_DB_NAME"`
		Collection       string `required:"true" envconfig:"MONGO_COLLECTION"`
	}
}

func InitConfig() (Config, error) {
	var cfg Config
	return cfg, envconfig.Process("", &cfg)
}
