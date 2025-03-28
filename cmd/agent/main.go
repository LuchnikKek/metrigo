package main

import (
	"log"

	"github.com/LuchnikKek/metrigo/internal/agent"
)

func main() {
	cfg := NewConfig()
	cfg.ParseFlags()
	cfg.ParseEnvs()
	log.Printf("Config parsed: %+v\r\n", cfg)

	metricsAgent := agent.NewMetricsAgent(
		"http://"+cfg.Addr,
		cfg.PollInterval,
		cfg.ReportInterval,
		cfg.RequestTimeout,
	)

	stop := make(chan struct{})
	defer close(stop)

	metricsAgent.Start()

	<-stop
}
