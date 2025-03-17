package main

import (
	"net/http"

	"github.com/LuchnikKek/metrigo/internal/agent"
)

func main() {
	ParseFlags()
	metricsAgent := agent.NewMetricsAgent(
		&http.Client{Timeout: Options.RequestTimeout},
		"http://"+Options.Addr,
		Options.PollInterval,
		Options.ReportInterval,
	)

	stop := make(chan struct{})
	defer close(stop)

	metricsAgent.Start()

	<-stop
}
