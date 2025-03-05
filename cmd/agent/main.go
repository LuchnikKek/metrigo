package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math/rand/v2"
	"net/http"
	"reflect"
	"runtime"
	"time"

	"github.com/LuchnikKek/metrigo/internal/models"
)

const (
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
)

type Metrics struct {
	runtime.MemStats
	PollCount   int
	RandomValue float64
}

type MetricsAgent struct {
	Metrics        Metrics
	PollInterval   time.Duration
	ReportInterval time.Duration
}

func (ag *MetricsAgent) Poll() {
	log.Println("Polling started")
	for {
		poll(&ag.Metrics)
		time.Sleep(ag.PollInterval)
	}
	// log.Println("Polling finished")
}

func poll(m *Metrics) {
	m.PollCount += 1
	m.RandomValue = rand.Float64()
	runtime.ReadMemStats(&m.MemStats)
}

var metricTypes = map[string]models.MetricType{
	"Alloc":         models.Gauge,
	"BuckHashSys":   models.Gauge,
	"Frees":         models.Gauge,
	"GCCPUFraction": models.Gauge,
	"GCSys":         models.Gauge,
	"HeapAlloc":     models.Gauge,
	"HeapIdle":      models.Gauge,
	"HeapInuse":     models.Gauge,
	"HeapObjects":   models.Gauge,
	"HeapReleased":  models.Gauge,
	"HeapSys":       models.Gauge,
	"LastGC":        models.Gauge,
	"Lookups":       models.Gauge,
	"MCacheInuse":   models.Gauge,
	"MCacheSys":     models.Gauge,
	"MSpanInuse":    models.Gauge,
	"MSpanSys":      models.Gauge,
	"Mallocs":       models.Gauge,
	"NextGC":        models.Gauge,
	"NumForcedGC":   models.Gauge,
	"NumGC":         models.Gauge,
	"OtherSys":      models.Gauge,
	"PauseTotalNs":  models.Gauge,
	"StackInuse":    models.Gauge,
	"StackSys":      models.Gauge,
	"Sys":           models.Gauge,
	"TotalAlloc":    models.Gauge,
	"PollCount":     models.Counter,
	"RandomValue":   models.Gauge,
}

func send(mName string, mValue any, mType string) error {
	url := fmt.Sprintf("http://localhost:8080/update/%s/%s/%v", mType, mName, mValue)

	body := bytes.NewBufferString(fmt.Sprintf("%v", mValue))

	resp, err := http.Post(url, "text/plain", body)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	log.Printf("Metric: %s, value: %v, type: %s. Response: %s\n", mName, mValue, mType, resp.Status)
	return nil
}

func process(m interface{}) {
	v := reflect.ValueOf(m)
	t := reflect.TypeOf(m)

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		if value.Kind() == reflect.Struct {
			process(value.Interface())
			continue
		}

		metricName := field.Name
		metricType, ok := metricTypes[metricName]
		if !ok {
			continue
		}

		log.Printf("Sending metric: %s\n", metricName)
		go send(metricName, value.Interface(), string(metricType))
	}
}

func (ag *MetricsAgent) Send() {
	log.Println("Sending started")
	for {
		process(ag.Metrics)
		time.Sleep(ag.ReportInterval)
	}
	// log.Println("Sending finished")
}

func main() {
	agent := MetricsAgent{
		Metrics:        Metrics{},
		PollInterval:   pollInterval,
		ReportInterval: reportInterval,
	}

	go agent.Poll()

	time.Sleep(2 * time.Second)
	agent.Send()

}
