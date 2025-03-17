package models

import (
	"errors"
	"log"
	"strconv"
)

type MetricType string

// Enum с возможными значениями
const (
	Gauge   MetricType = "gauge"
	Counter MetricType = "counter"
)

var (
	ErrInvalidMetricType  = errors.New("invalid metric type")
	ErrInvalidMetricValue = errors.New("invalid metric value")
)

func IsValidMetricType(mType string) bool {
	m := MetricType(mType)
	return m == Gauge || m == Counter
}

// Интерфейс метрики
type Metric interface {
	GetName() string
	GetType() MetricType
	Update(value any) error
	GetValue() string
}

func NewMetric(mType, name, value string) (Metric, error) {
	switch MetricType(mType) {
	case Gauge:
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, err
		}
		return NewGaugeMetric(name, v), nil

	case Counter:
		v, err := strconv.Atoi(value)
		if err != nil {
			return nil, err
		}
		return NewCounterMetric(name, v), nil

	default:
		return nil, ErrInvalidMetricType
	}

}

type GaugeMetric struct {
	Name  string     `json:"name"`
	Value float64    `json:"value"`
	Type  MetricType `json:"type"`
}

func NewGaugeMetric(name string, value float64) *GaugeMetric {
	return &GaugeMetric{Name: name, Value: value, Type: Gauge}
}

func (m *GaugeMetric) GetName() string     { return m.Name }
func (m *GaugeMetric) GetType() MetricType { return m.Type }
func (m *GaugeMetric) GetValue() string    { return strconv.FormatFloat(m.Value, 'f', -1, 64) }

func (m *GaugeMetric) Update(value any) error {
	log.Printf("%T and %#v", value, value)
	v, ok := value.(*GaugeMetric)
	if !ok {
		return ErrInvalidMetricValue
	}
	m.Value = v.Value
	return nil
}

type CounterMetric struct {
	Name  string
	Value int
	Type  MetricType
}

func NewCounterMetric(name string, value int) *CounterMetric {
	return &CounterMetric{Name: name, Value: value, Type: Counter}
}

func (m *CounterMetric) GetName() string     { return m.Name }
func (m *CounterMetric) GetType() MetricType { return m.Type }
func (m *CounterMetric) GetValue() string    { return strconv.Itoa(m.Value) }

func (m *CounterMetric) Update(value any) error {
	v, ok := value.(*CounterMetric)
	if !ok {
		return ErrInvalidMetricValue
	}
	m.Value += v.Value
	return nil
}
