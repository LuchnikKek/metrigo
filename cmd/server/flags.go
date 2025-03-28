package main

import (
	"flag"
	"os"
)

type Config struct {
	Addr string
}

func NewConfig() *Config {
	return &Config{}
}

func (cfg *Config) ParseFlags() {
	flag.StringVar(&cfg.Addr, "a", "localhost:8080", "адрес и порт, на котором будет запущен сервер")
	flag.Parse()
}

func (cfg *Config) ParseEnvs() {
	if envRunAddr, isSet := os.LookupEnv("ADDRESS"); isSet {
		cfg.Addr = envRunAddr
	}
}
