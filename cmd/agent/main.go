package main

import (
	"time"

	"github.com/LuchnikKek/metrigo/internal/agent"
)

func main() {
	metricsAgent := agent.NewMetricsAgent()

	go metricsAgent.Poll()

	time.Sleep(time.Second * 3)

	metricsAgent.Send()
}
