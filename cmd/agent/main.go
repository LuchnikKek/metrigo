package main

import (
	"net/http"
	"time"

	"github.com/LuchnikKek/metrigo/internal/agent"
)

func main() {
	metricsAgent := agent.NewMetricsAgent(&http.Client{Timeout: 5 * time.Second}, "http://localhost:8080")

	stop := make(chan struct{})
	defer close(stop)

	metricsAgent.Start()

	<-stop
}
