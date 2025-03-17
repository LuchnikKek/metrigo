package storage

import (
	"errors"

	"github.com/LuchnikKek/metrigo/internal/models"
)

type Storage interface {
	Save(m models.Metric) error
	Get(name string) (models.Metric, error)
	GetAll() []models.Metric
}

var (
	ErrMetricNotFound = errors.New("metric not found")
)
