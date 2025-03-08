package models

import (
	"errors"
)

// Интерфейс метрики
type Metric interface {
	GetName() string
	GetType() MetricType
	Update(value any) error
	GetValue() any
}

type MetricType string

// Enum с возможными значениями
const (
	Gauge   MetricType = "gauge"
	Counter MetricType = "counter"
)

type GaugeMetric struct {
	Name  string
	Value float64
}

func (m *GaugeMetric) GetName() string     { return m.Name }
func (m *GaugeMetric) GetType() MetricType { return Gauge }
func (m *GaugeMetric) GetValue() any       { return m.Value }

func (m *GaugeMetric) Update(value any) error {
	v, ok := value.(float64)
	if !ok {
		return errors.New("invalid value type for gauge, expected float64")
	}
	m.Value = v
	return nil
}

type CounterMetric struct {
	Name  string
	Value int64
}

func (m *CounterMetric) GetName() string     { return m.Name }
func (m *CounterMetric) GetType() MetricType { return Counter }
func (m *CounterMetric) GetValue() any       { return m.Value }

func (m *CounterMetric) Update(value any) error {
	v, ok := value.(int64)
	if !ok {
		return errors.New("invalid value type for counter, expected int64")
	}
	m.Value += v
	return nil
}
