package agent

import (
	"runtime"

	"github.com/LuchnikKek/metrigo/internal/models"
)

type Metrics struct {
	runtime.MemStats
	PollCount   int
	RandomValue float64
}

var MetricTypes = map[string]models.MetricType{
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
