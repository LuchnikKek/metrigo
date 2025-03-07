package agent

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
)

const (
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
)

type MetricsAgent struct {
	Metrics        Metrics
	PollInterval   time.Duration
	ReportInterval time.Duration
}

func NewMetricsAgent() *MetricsAgent {
	return &MetricsAgent{
		Metrics:        Metrics{},
		PollInterval:   pollInterval,
		ReportInterval: reportInterval,
	}
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

func (ag *MetricsAgent) Send() {
	log.Println("Sending started")
	for {
		process(ag.Metrics)
		time.Sleep(ag.ReportInterval)
	}
	// log.Println("Sending finished")
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

		mName := field.Name
		mType, inSendingList := MetricTypes[mName]
		if !inSendingList {
			continue
		}

		log.Printf("Sending metric: %s\n", mName)
		go send(mName, value.Interface(), string(mType))
	}
}
