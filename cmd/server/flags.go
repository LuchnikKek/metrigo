package main

import (
	"flag"
	"log"
	"os"
)

var Options struct {
	Addr string
}

func InitOptions() {
	parseFlags()
	parseEnvs()

	log.Printf("Config parsed: %+v\r\n", Options)
}

func parseFlags() {
	flag.StringVar(&Options.Addr, "a", "localhost:8080", "адрес и порт, на котором будет запущен сервер")
	flag.Parse()
}

func parseEnvs() {
	if envRunAddr, isSet := os.LookupEnv("ADDRESS"); isSet {
		Options.Addr = envRunAddr
	}
}
