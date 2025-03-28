package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/caarlos0/env"
)

type TimeInterval struct {
	Duration time.Duration
}

func (ti *TimeInterval) String() string {
	return fmt.Sprint(ti.Duration.Seconds())
}

func (ti *TimeInterval) Set(flagValue string) error {
	duration, err := strconv.ParseFloat(flagValue, 64)
	if err != nil {
		return err
	}
	ti.Duration = time.Duration(duration * float64(time.Second))
	return nil
}

type Config struct {
	Addr           string
	ReportInterval time.Duration
	PollInterval   time.Duration
	RequestTimeout time.Duration
}

func NewConfig() *Config {
	return &Config{}
}

func (cfg *Config) ParseFlags() {
	pollInterval := &TimeInterval{Duration: 2 * time.Second}
	reportInterval := &TimeInterval{Duration: 10 * time.Second}
	requestTimeout := &TimeInterval{Duration: 5 * time.Second}

	flag.StringVar(&cfg.Addr, "a", "localhost:8080", "адрес и порт сервера куда отправлять метрики")
	flag.Var(pollInterval, "p", "частота опроса метрик из пакета runtime")
	flag.Var(reportInterval, "r", "частота отправки метрик на сервер")
	flag.Var(requestTimeout, "t", "таймаут запроса на отправку метрики")

	flag.Parse()

	cfg.PollInterval = pollInterval.Duration
	cfg.ReportInterval = reportInterval.Duration
	cfg.RequestTimeout = requestTimeout.Duration
}

type EnvConfig struct {
	Addr                string  `env:"ADDRESS"`
	ReportInterval      float64 `env:"REPORT_INTERVAL"`
	PollInterval        float64 `env:"POLL_INTERVAL"`
	RequestTimeout      float64 `env:"REQUEST_TIMEOUT"`
}

func (cfg *Config) ParseEnvs() {
	var ec EnvConfig

	err := env.Parse(&ec)
	if err != nil {
		log.Fatal(err)
	}

	if ec.Addr != "" {
		cfg.Addr = ec.Addr
	}
	if ec.PollInterval != 0.0 {
		cfg.PollInterval = time.Duration(ec.PollInterval * float64(time.Second))
	}
	if ec.ReportInterval != 0.0 {
		cfg.ReportInterval = time.Duration(ec.ReportInterval * float64(time.Second))
	}
	if ec.RequestTimeout != 0.0 {
		cfg.RequestTimeout = time.Duration(ec.RequestTimeout * float64(time.Second))
	}
}
