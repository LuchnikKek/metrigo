package storage

import (
	"log"
	"sync"

	"github.com/LuchnikKek/metrigo/internal/models"
)

type InMemoryStorage struct {
	mu      sync.RWMutex
	metrics map[string]models.Metric
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		metrics: make(map[string]models.Metric),
	}
}

func (s *InMemoryStorage) Save(m models.Metric) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	name := m.GetName()

	stored, exists := s.metrics[name]
	if !exists {
		s.metrics[name] = m
		log.Printf("Metric created: %v, val=%#v\n", name, m)
		return nil
	}
	if err := stored.Update(m); err == nil {
		s.metrics[name] = stored
		log.Printf("Metric updated: %v, val=%#v\n", name, stored)
		return nil
	} else {
		log.Printf("Error updating metric: %v, val=%#v\n", name, stored)
		return err
	}
}

func (s *InMemoryStorage) Get(name string) (models.Metric, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m, exists := s.metrics[name]
	if !exists {
		return nil, ErrMetricNotFound
	}
	return m, nil
}

func (s *InMemoryStorage) GetAll() []models.Metric {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return mapToValues(s.metrics)
}

func mapToValues[M ~map[K]V, K comparable, V any](m M) []V {
	r := make([]V, 0, len(m))
	for _, v := range m {
		r = append(r, v)
	}
	return r
}
