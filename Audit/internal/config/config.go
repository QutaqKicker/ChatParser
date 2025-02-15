package config

import (
	"github.com/QutaqKicker/ChatParser/Common/dbHelper"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type Config struct {
	Db dbHelper.DbConfig `yaml:"db"`
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
