package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env       string        `yaml:"env" env-default:"local"`
	ExportDir string        `yaml:"exportDir"`
	Timeout   time.Duration `yaml:"timeout"`
}

func MustLoad() *Config {
	path := "config.yaml"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config does not exists: " + path)
	}

	var config Config

	if err := cleanenv.ReadConfig(path, &config); err != nil {
		panic("failed to read config " + err.Error())
	}
	return &config
}
