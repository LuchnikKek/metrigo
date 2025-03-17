package main

import (
	"flag"
	"log"
)

var Options struct {
	Addr string
}

func ParseFlags() {
	flag.StringVar(&Options.Addr, "a", "localhost:8080", "адрес и порт, на котором будет запущен сервер")

	flag.Parse()

	log.Printf("Config parsed: %+v\r\n", Options)
}
