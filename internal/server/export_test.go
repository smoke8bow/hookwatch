package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/user/hookwatch/internal/storage"
)

func seedStore(t *testing.T, store *storage.Store, n int) {
	t.Helper()
	for i := 0; i < n; i++ {
		store.Add(storage.Request{
			ID:        fmt.Sprintf("id-%d", i),
			Timestamp: time.Now(),
			Method:    "POST",
			Path:      "/webhook",
			Headers:   map[string][]string{"Content-Type": {"application/json"}},
			Body:      `{"event":"test"}`,
		})
	}
}

func TestHandleExport_JSON(t *testing.T) {
	router := newTestRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/requests/export?format=json", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected application/json, got %s", ct)
	}

	var records []exportRecord
	if err := json.NewDecoder(rec.Body).Decode(&records); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
}

func TestHandleExport_NDJSON(t *testing.T) {
	router := newTestRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/requests/export?format=ndjson", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	lines := strings.Split(strings.TrimSpace(rec.Body.String()), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		var r exportRecord
		if err := json.Unmarshal([]byte(line), &r); err != nil {
			t.Errorf("invalid ndjson line: %v", err)
		}
	}
}

func TestHandleExport_InvalidFormat(t *testing.T) {
	router := newTestRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/requests/export?format=csv", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}
