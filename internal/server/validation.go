package server

import (
	"errors"

	"github.com/LuchnikKek/metrigo/internal/models"
)

func ValidateMetric(m models.Metric) error {
	if m.GetName() == "" {
		return errors.New("metric name cannot be empty")
	}
	if m.GetType() != models.Gauge && m.GetType() != models.Counter {
		return errors.New("invalid metric type")
	}
	return nil
}
