package api

import "github.com/narekps/go2/day2/storage"

type Config struct {
	BindAddr string
	LogLevel string
	Storage  *storage.Config
}

func NewConfig() *Config {
	return &Config{
		BindAddr: "localhost:8080",
		LogLevel: "debug",
		Storage:  storage.NewConfig(),
	}
}
