package storage

import (
	"errors"

	"github.com/LuchnikKek/metrigo/internal/models"
)

type Storage interface {
	SaveMetric(m models.Metric) error
	GetMetricByName(name string) (models.Metric, error)
	GetMetrics() []models.Metric
}

var (
	ErrMetricNotFound = errors.New("metric not found")
)
