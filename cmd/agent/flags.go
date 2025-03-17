package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"time"
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

func ParseFlags() {
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

	log.Printf("Config parsed: %+v\r\n", Options)
}
