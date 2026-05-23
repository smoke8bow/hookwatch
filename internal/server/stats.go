package server

import (
	"sync"
	"time"
)

// Stats tracks aggregate metrics about captured webhook requests.
type Stats struct {
	mu           sync.RWMutex
	TotalCount   int            `json:"total_count"`
	MethodCounts map[string]int `json:"method_counts"`
	PathCounts   map[string]int `json:"path_counts"`
	LastReceived *time.Time     `json:"last_received,omitempty"`
	ReplayCount  int            `json:"replay_count"`
}

// NewStats initialises an empty Stats tracker.
func NewStats() *Stats {
	return &Stats{
		MethodCounts: make(map[string]int),
		PathCounts:   make(map[string]int),
	}
}

// RecordRequest updates counters for a newly captured request.
func (s *Stats) RecordRequest(method, path string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.TotalCount++
	s.MethodCounts[method]++
	s.PathCounts[path]++
	now := time.Now().UTC()
	s.LastReceived = &now
}

// RecordReplay increments the replay counter.
func (s *Stats) RecordReplay() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ReplayCount++
}

// Snapshot returns a copy of the current stats safe for serialisation.
func (s *Stats) Snapshot() Stats {
	s.mu.RLock()
	defer s.mu.RUnlock()

	methodCopy := make(map[string]int, len(s.MethodCounts))
	for k, v := range s.MethodCounts {
		methodCopy[k] = v
	}
	pathCopy := make(map[string]int, len(s.PathCounts))
	for k, v := range s.PathCounts {
		pathCopy[k] = v
	}

	snap := Stats{
		TotalCount:   s.TotalCount,
		MethodCounts: methodCopy,
		PathCounts:   pathCopy,
		ReplayCount:  s.ReplayCount,
	}
	if s.LastReceived != nil {
		t := *s.LastReceived
		snap.LastReceived = &t
	}
	return snap
}
