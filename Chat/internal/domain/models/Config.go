package models

import "github.com/QutaqKicker/ChatParser/Common/dbHelper"

type Config struct {
	Env string            `yaml:"env" env-default:"local"`
	Db  dbHelper.DbConfig `yaml:"db"`
}
