package storage

import (
	"fmt"
	"testing"
	"time"
)

func makeRequest(id string) *Request {
	return &Request{
		ID:         id,
		ReceivedAt: time.Now(),
		Method:     "POST",
		Path:       "/webhook",
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       []byte(`{"event":"test"}`),
		SourceIP:   "127.0.0.1",
	}
}

func TestStore_AddAndLen(t *testing.T) {
	s := NewStore(10)
	s.Add(makeRequest("req-1"))
	s.Add(makeRequest("req-2"))
	if got := s.Len(); got != 2 {
		t.Errorf("expected Len()=2, got %d", got)
	}
}

func TestStore_GetAll_NewestFirst(t *testing.T) {
	s := NewStore(10)
	s.Add(makeRequest("req-1"))
	s.Add(makeRequest("req-2"))
	s.Add(makeRequest("req-3"))
	all := s.GetAll()
	if len(all) != 3 {
		t.Fatalf("expected 3 requests, got %d", len(all))
	}
	if all[0].ID != "req-3" {
		t.Errorf("expected newest first, got %s", all[0].ID)
	}
}

func TestStore_GetByID(t *testing.T) {
	s := NewStore(10)
	s.Add(makeRequest("abc-123"))
	r := s.GetByID("abc-123")
	if r == nil {
		t.Fatal("expected to find request, got nil")
	}
	if r.ID != "abc-123" {
		t.Errorf("expected ID abc-123, got %s", r.ID)
	}
	if s.GetByID("nonexistent") != nil {
		t.Error("expected nil for missing ID")
	}
}

func TestStore_Eviction(t *testing.T) {
	s := NewStore(3)
	for i := 1; i <= 5; i++ {
		s.Add(makeRequest(fmt.Sprintf("req-%d", i)))
	}
	if s.Len() != 3 {
		t.Errorf("expected Len()=3 after eviction, got %d", s.Len())
	}
	if s.GetByID("req-1") != nil || s.GetByID("req-2") != nil {
		t.Error("oldest requests should have been evicted")
	}
	if s.GetByID("req-5") == nil {
		t.Error("newest request should still be present")
	}
}

func TestStore_Clear(t *testing.T) {
	s := NewStore(10)
	s.Add(makeRequest("req-1"))
	s.Clear()
	if s.Len() != 0 {
		t.Errorf("expected empty store after Clear, got %d", s.Len())
	}
}
