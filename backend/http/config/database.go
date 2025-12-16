package config

import (
	"os"
)

type DBConfig struct {
	Path string
}

func LoadDBConfig() *DBConfig {
	return &DBConfig{
		Path: os.Getenv("DB_PATH"),
	}
}
