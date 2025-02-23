package storage

import "github.com/LuchnikKek/metrigo/internal/models"

type Storage interface {
	Save(m models.Metric) error
	Get(name string) (models.Metric, error)
}
