package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStats_RecordAndSnapshot(t *testing.T) {
	s := NewStats()

	s.RecordRequest("POST", "/hooks/github")
	s.RecordRequest("POST", "/hooks/github")
	s.RecordRequest("GET", "/hooks/stripe")
	s.RecordReplay()

	snap := s.Snapshot()

	if snap.TotalCount != 3 {
		t.Fatalf("expected TotalCount 3, got %d", snap.TotalCount)
	}
	if snap.MethodCounts["POST"] != 2 {
		t.Errorf("expected POST count 2, got %d", snap.MethodCounts["POST"])
	}
	if snap.MethodCounts["GET"] != 1 {
		t.Errorf("expected GET count 1, got %d", snap.MethodCounts["GET"])
	}
	if snap.PathCounts["/hooks/github"] != 2 {
		t.Errorf("expected /hooks/github count 2, got %d", snap.PathCounts["/hooks/github"])
	}
	if snap.ReplayCount != 1 {
		t.Errorf("expected ReplayCount 1, got %d", snap.ReplayCount)
	}
	if snap.LastReceived == nil {
		t.Error("expected LastReceived to be set")
	}
}

func TestHandleStats(t *testing.T) {
	router := newTestRouter(t)

	// Seed a request so stats are non-zero.
	body := `{"event":"push"}`
	req := httptest.NewRequest(http.MethodPost, "/hooks/test", stringReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Fetch stats.
	req2 := httptest.NewRequest(http.MethodGet, "/stats", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w2.Code)
	}

	var stats Stats
	if err := json.NewDecoder(w2.Body).Decode(&stats); err != nil {
		t.Fatalf("failed to decode stats response: %v", err)
	}
	if stats.TotalCount < 1 {
		t.Errorf("expected at least 1 total request, got %d", stats.TotalCount)
	}
}
