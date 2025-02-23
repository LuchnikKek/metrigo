package storage

import (
	"errors"
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
		log.Printf("Metric created: %v\n", m)
		return nil
	}
	if err := stored.Update(m.GetValue()); err == nil {
		s.metrics[name] = stored
		log.Printf("Metric updated: %v\n", stored)
		return nil
	} else {
		log.Printf("Error updating metric: %v\n", m)
		return err
	}
}

func (s *InMemoryStorage) Get(name string) (models.Metric, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m, exists := s.metrics[name]
	if !exists {
		return nil, errors.New("metric not found")
	}
	return m, nil
}
