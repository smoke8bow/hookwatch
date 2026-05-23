package storage

import (
	"sync"
	"time"
)

// Request represents a captured incoming webhook request.
type Request struct {
	ID        string            `json:"id"`
	ReceivedAt time.Time        `json:"received_at"`
	Method    string            `json:"method"`
	Path      string            `json:"path"`
	Headers   map[string]string `json:"headers"`
	Body      []byte            `json:"body"`
	SourceIP  string            `json:"source_ip"`
}

// Store holds captured requests in memory.
type Store struct {
	mu       sync.RWMutex
	requests []*Request
	maxSize  int
}

// NewStore creates a new in-memory request store with a given max capacity.
func NewStore(maxSize int) *Store {
	if maxSize <= 0 {
		maxSize = 500
	}
	return &Store{
		requests: make([]*Request, 0, maxSize),
		maxSize:  maxSize,
	}
}

// Add appends a new request to the store, evicting the oldest if at capacity.
func (s *Store) Add(r *Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.requests) >= s.maxSize {
		s.requests = s.requests[1:]
	}
	s.requests = append(s.requests, r)
}

// GetAll returns a copy of all stored requests, newest first.
func (s *Store) GetAll() []*Request {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]*Request, len(s.requests))
	for i, r := range s.requests {
		copy := *r
		result[len(s.requests)-1-i] = &copy
	}
	return result
}

// GetByID returns a single request by its ID, or nil if not found.
func (s *Store) GetByID(id string) *Request {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, r := range s.requests {
		if r.ID == id {
			copy := *r
			return &copy
		}
	}
	return nil
}

// Clear removes all stored requests.
func (s *Store) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.requests = s.requests[:0]
}

// Len returns the number of stored requests.
func (s *Store) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.requests)
}
