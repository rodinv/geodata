package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config is a config
type Config struct {
	Http Http `yaml:"http"`
	Db   Db   `yaml:"db"`
}

type Db struct {
	Path string `yaml:"path" env-default:"data/geodata.dat"`
}

type Http struct {
	Port    string        `yaml:"port" env-default:"8081"`
	Timeout time.Duration `yaml:"timeout" env-default:"1m"`
}

// Get reads config from environment. Once.
func Get(path string) (*Config, error) {
	var cfg Config
	err := cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
