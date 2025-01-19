package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type Config struct {
	Db DbConfig `yaml:"db"`
}

type DbConfig struct {
	Host     string `yaml:"host" env-required:"true"`
	Port     int    `yaml:"port" env-required:"true"`
	User     string `yaml:"user" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
	DbName   string `yaml:"dbName" env-required:"true"`
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
