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

var Options struct {
	Addr           string
	ReportInterval time.Duration
	PollInterval   time.Duration
	RequestTimeout time.Duration
}

func InitOptions() {
	parseFlags()
	parseEnvs()

	log.Printf("Config parsed: %+v\r\n", Options)
}

func parseFlags() {
	pollInterval := &TimeInterval{Duration: 2 * time.Second}
	reportInterval := &TimeInterval{Duration: 10 * time.Second}
	requestTimeout := &TimeInterval{Duration: 5 * time.Second}

	flag.StringVar(&Options.Addr, "a", "localhost:8080", "адрес и порт сервера куда отправлять метрики")
	flag.Var(pollInterval, "p", "частота опроса метрик из пакета runtime")
	flag.Var(reportInterval, "r", "частота отправки метрик на сервер")
	flag.Var(requestTimeout, "t", "таймаут запроса на отправку метрики")

	flag.Parse()

	Options.PollInterval = pollInterval.Duration
	Options.ReportInterval = reportInterval.Duration
	Options.RequestTimeout = requestTimeout.Duration
}

type EnvConfig struct {
	Addr           string  `env:"ADDRESS"`
	ReportInterval float64 `env:"REPORT_INTERVAL"`
	PollInterval   float64 `env:"POLL_INTERVAL"`
}

func parseEnvs() {
	var cfg EnvConfig

	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	if cfg.Addr != "" {
		Options.Addr = cfg.Addr
	}
	if cfg.PollInterval != 0.0 {
		Options.PollInterval = time.Duration(cfg.PollInterval * float64(time.Second))
	}
	if cfg.ReportInterval != 0.0 {
		Options.ReportInterval = time.Duration(cfg.ReportInterval * float64(time.Second))
	}
}
