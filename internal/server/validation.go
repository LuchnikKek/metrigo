package server

import (
	"errors"
	"strconv"

	"github.com/LuchnikKek/metrigo/internal/models"
)

// Ошибки при валидации
var (
	ErrMissingMetricName  = errors.New("metric name is required")
	ErrMissingMetricType  = errors.New("metric type is required")
	ErrMissingMetricValue = errors.New("metric value is required")
	ErrInvalidMetricType  = errors.New("invalid metric type")
	ErrInvalidMetricValue = errors.New("invalid metric value")
)

// DTO-шка
type MetricRequest struct {
	Type  string
	Name  string
	Value string
}

// ValidateMetricVars проверяет, что в URL есть нужные параметры
func ValidateMetricVars(vars map[string]string) (MetricRequest, error) {
	metricType, okType := vars["type"]
	metricName, okName := vars["name"]
	metricValue, okValue := vars["value"]

	if !okName || metricName == "" {
		return MetricRequest{}, ErrMissingMetricName
	}
	if !okType || metricType == "" {
		return MetricRequest{}, ErrMissingMetricType
	}
	if !okValue || metricValue == "" {
		return MetricRequest{}, ErrMissingMetricValue
	}

	return MetricRequest{
		Type:  metricType,
		Name:  metricName,
		Value: metricValue,
	}, nil
}

// ParseMetric создаёт экземпляр GaugeMetric или CounterMetric
// на основе URL-параметров (type, name, value)
func ParseMetric(req MetricRequest) (models.Metric, error) {
	switch req.Type {
	case string(models.Gauge):
		val, err := strconv.ParseFloat(req.Value, 64)
		if err != nil {
			return nil, ErrInvalidMetricValue
		}
		return &models.GaugeMetric{
			Name:  req.Name,
			Value: val,
		}, nil

	case string(models.Counter):
		val, err := strconv.ParseInt(req.Value, 10, 64)
		if err != nil {
			return nil, ErrInvalidMetricValue
		}
		return &models.CounterMetric{
			Name:  req.Name,
			Value: val,
		}, nil

	default:
		return nil, ErrInvalidMetricType
	}
}
