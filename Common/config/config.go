package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

func MustLoad[T any]() *T {
	return MustLoadPath[T]("config.yaml")
}

func MustLoadPath[T any](configPath string) *T {
	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg T

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}
