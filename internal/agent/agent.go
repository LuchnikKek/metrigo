package agent

import (
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/LuchnikKek/metrigo/internal/models"
)

type MetricsAgent struct {
	client         *http.Client
	baseURL        string
	PollInterval   time.Duration
	ReportInterval time.Duration
	metricsMap     map[string]models.Metric
	SendBufferSize int
	mu sync.RWMutex
}

func NewMetricsAgent(baseURL string, pollInterval, reportInterval, requestTimeout time.Duration) *MetricsAgent {
	return &MetricsAgent{
		client:         &http.Client{Timeout: requestTimeout},
		baseURL:        baseURL,
		PollInterval:   pollInterval,
		ReportInterval: reportInterval,
		metricsMap: 	make(map[string]models.Metric),
		SendBufferSize: 5,
		mu:             sync.RWMutex{},
	}
}

func (ag *MetricsAgent) Start() {
	ag.InitMetrics()

	go func() {
		log.Println("Polling started")
		for {
			ag.Poll()
			time.Sleep(ag.PollInterval)
		}
		// log.Println("Polling finished")
	}()

	time.Sleep(ag.PollInterval) // wait for first poll

	go func() {
		log.Println("Sending started")
		for {
			ag.Process()
			time.Sleep(ag.ReportInterval)
		}
		// log.Println("Sending finished")
	}()
}

func (ag *MetricsAgent) InitMetrics() {
	ag.mu.Lock()
	defer ag.mu.Unlock()

	for mName, mType := range MetricTypes {
		switch mType {
		case models.Gauge:
			ag.metricsMap[mName] = models.NewGaugeMetric(mName, 0)
		case models.Counter:
			ag.metricsMap[mName] = models.NewCounterMetric(mName, 0)
		}
	}
}

func (ag *MetricsAgent) Poll() {
	memStats := runtime.MemStats{}
	runtime.ReadMemStats(&memStats)

	ag.mu.Lock()
	defer ag.mu.Unlock()

	ag.metricsMap["PollCount"].Update(1)
	ag.metricsMap["RandomValue"].Update(rand.Float64())
	ag.metricsMap["Alloc"].Update(float64(memStats.Alloc))
	ag.metricsMap["BuckHashSys"].Update(float64(memStats.BuckHashSys))
	ag.metricsMap["Frees"].Update(float64(memStats.Frees))
	ag.metricsMap["GCCPUFraction"].Update(float64(memStats.GCCPUFraction))
	ag.metricsMap["GCSys"].Update(float64(memStats.GCSys))
	ag.metricsMap["HeapAlloc"].Update(float64(memStats.HeapAlloc))
	ag.metricsMap["HeapIdle"].Update(float64(memStats.HeapIdle))
	ag.metricsMap["HeapInuse"].Update(float64(memStats.HeapInuse))
	ag.metricsMap["HeapObjects"].Update(float64(memStats.HeapObjects))
	ag.metricsMap["HeapReleased"].Update(float64(memStats.HeapReleased))
	ag.metricsMap["HeapSys"].Update(float64(memStats.HeapSys))
	ag.metricsMap["LastGC"].Update(float64(memStats.LastGC))
	ag.metricsMap["Lookups"].Update(float64(memStats.Lookups))
	ag.metricsMap["MCacheInuse"].Update(float64(memStats.MCacheInuse))
	ag.metricsMap["MCacheSys"].Update(float64(memStats.MCacheSys))
	ag.metricsMap["MSpanInuse"].Update(float64(memStats.MSpanInuse))
	ag.metricsMap["MSpanSys"].Update(float64(memStats.MSpanSys))
	ag.metricsMap["Mallocs"].Update(float64(memStats.Mallocs))
	ag.metricsMap["NextGC"].Update(float64(memStats.NextGC))
	ag.metricsMap["NumForcedGC"].Update(float64(memStats.NumForcedGC))
	ag.metricsMap["NumGC"].Update(float64(memStats.NumGC))
	ag.metricsMap["OtherSys"].Update(float64(memStats.OtherSys))
	ag.metricsMap["PauseTotalNs"].Update(float64(memStats.PauseTotalNs))
	ag.metricsMap["StackInuse"].Update(float64(memStats.StackInuse))
	ag.metricsMap["StackSys"].Update(float64(memStats.StackSys))
	ag.metricsMap["Sys"].Update(float64(memStats.Sys))
	ag.metricsMap["TotalAlloc"].Update(float64(memStats.TotalAlloc))
}

func (ag *MetricsAgent) Process() {
	ch := make(chan struct{}, ag.SendBufferSize)

	ag.mu.RLock()
	defer ag.mu.RUnlock()

	for name, value := range ag.metricsMap {
		ch <- struct{}{}
		go func(n string, v models.Metric) {
			defer func() { <-ch }()
			if err := ag.SendMetric(v); err != nil {
				log.Println(err)
			}
		}(name, value)
	}
}

func (ag *MetricsAgent) SendMetric(m models.Metric) error {
	suffix := fmt.Sprintf("/update/%s/%s/%v", m.GetType(), m.GetName(), m.GetValue())
	req, err := http.NewRequest(http.MethodPost, ag.baseURL+suffix, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := ag.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %s", resp.Status)
	}

	log.Printf("Response \"%s\" - Status: %s\r\n", suffix, resp.Status)
	return nil
}
