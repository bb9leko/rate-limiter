package store

import (
	"context"
	"sync"
	"time"
)

type MemoryStore struct {
	data map[string]*entry
	mu   sync.RWMutex
}

type entry struct {
	count  int
	expiry time.Time
}

func NewMemoryStore() *MemoryStore {
	store := &MemoryStore{
		data: make(map[string]*entry),
	}
	go store.cleanup() // Goroutine para limpar entradas expiradas
	return store
}

func (s *MemoryStore) Increment(ctx context.Context, key string, window time.Duration) (int, time.Duration, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()

	// Verifica se a entrada existe e não expirou
	if entry, exists := s.data[key]; exists && now.Before(entry.expiry) {
		entry.count++
		ttl := entry.expiry.Sub(now)
		return entry.count, ttl, nil
	}

	// Cria nova entrada
	s.data[key] = &entry{
		count:  1,
		expiry: now.Add(window),
	}
	return 1, window, nil
}

func (s *MemoryStore) Reset(ctx context.Context, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
	return nil
}

func (s *MemoryStore) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = make(map[string]*entry)
	return nil
}

func (s *MemoryStore) Ping(ctx context.Context) error {
	return nil // Memory store sempre está "conectado"
}

func (s *MemoryStore) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		now := time.Now()
		for key, entry := range s.data {
			if now.After(entry.expiry) {
				delete(s.data, key)
			}
		}
		s.mu.Unlock()
	}
}
