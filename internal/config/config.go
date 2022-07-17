package config

import (
	"github.com/jinzhu/configor"
)

type Config struct {
	Database *Database
	REST     *REST
	Logger   *Logger
}

type Database struct {
	Filename string `required:"true"`
	Capacity int    `required:"true"`
}

type REST struct {
	Port int `required:"true"`
}

type Logger struct {
	Level string `required:"true"`
}

func NewConfig(filename string) (*Config, error) {
	config := &Config{}
	configor.Load(config, filename)
	return config, nil
}
