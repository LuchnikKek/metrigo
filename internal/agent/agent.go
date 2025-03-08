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
	client         *http.Client
	Metrics        Metrics
	PollInterval   time.Duration
	ReportInterval time.Duration
}

func NewMetricsAgent(client *http.Client) *MetricsAgent {
	return &MetricsAgent{
		client:         client,
		Metrics:        Metrics{},
		PollInterval:   pollInterval,
		ReportInterval: reportInterval,
	}
}

func (ag *MetricsAgent) Start() {
	go func() {
		log.Println("Polling started")
		for {
			ag.Poll(&ag.Metrics)
			time.Sleep(ag.PollInterval)
		}
		// log.Println("Polling finished")
	}()

	time.Sleep(pollInterval) // wait for first poll

	go func() {
		log.Println("Sending started")
		for {
			ag.Process(ag.Metrics)
			time.Sleep(ag.ReportInterval)
		}
		// log.Println("Sending finished")
	}()
}

func (ag *MetricsAgent) Poll(m *Metrics) {
	m.PollCount += 1
	m.RandomValue = rand.Float64()
	runtime.ReadMemStats(&m.MemStats)
}

func (ag *MetricsAgent) Process(m interface{}) {
	v := reflect.ValueOf(m)
	t := reflect.TypeOf(m)

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		if value.Kind() == reflect.Struct {
			ag.Process(value.Interface())
			continue
		}

		mName := field.Name
		mType, inSendingList := MetricTypes[mName]
		if !inSendingList {
			continue
		}

		log.Printf("Sending metric: %s\n", mName)
		go ag.Send(mName, value.Interface(), string(mType))
	}
}

func (ag *MetricsAgent) Send(mName string, mValue any, mType string) error {
	url := fmt.Sprintf("http://localhost:8080/update/%s/%s/%v", mType, mName, mValue)

	body := bytes.NewBufferString(fmt.Sprintf("%v", mValue))

	resp, err := ag.client.Post(url, "text/plain", body)
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
