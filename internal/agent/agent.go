package agent

import (
	"fmt"
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
	baseURL        string
	Metrics        Metrics
	PollInterval   time.Duration
	ReportInterval time.Duration
}

func NewMetricsAgent(client *http.Client, baseURL string) *MetricsAgent {
	return &MetricsAgent{
		client:         client,
		baseURL:        baseURL,
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
	url := fmt.Sprintf("%s/update/%s/%s/%v", ag.baseURL, mType, mName, mValue)
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := ag.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	log.Printf("Metric: %s, value: %v, type: %s. Response: %s\n", mName, mValue, mType, resp.Status)
	return nil
}
