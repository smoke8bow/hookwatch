package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/hookwatch/internal/storage"
)

func TestHandleReplay_NotFound(t *testing.T) {
	router := newTestRouter(t)

	req := httptest.NewRequest(http.MethodPost, "/requests/nonexistent/replay?target=http://example.com", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestHandleReplay_MissingTarget(t *testing.T) {
	router, store := newTestRouterWithStore(t)

	sr := &storage.Request{
		ID:     "abc123",
		Method: http.MethodPost,
		Path:   "/hook",
		Body:   []byte(`{"event":"test"}`),
		Headers: map[string][]string{
			"Content-Type": {"application/json"},
		},
	}
	store.Add(sr)

	req := httptest.NewRequest(http.MethodPost, "/requests/abc123/replay", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestHandleReplay_Success(t *testing.T) {
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer target.Close()

	router, store := newTestRouterWithStore(t)

	sr := &storage.Request{
		ID:     "xyz789",
		Method: http.MethodPost,
		Path:   "/hook",
		Body:   []byte(`{"event":"ping"}`),
		Headers: map[string][]string{
			"Content-Type": {"application/json"},
		},
	}
	store.Add(sr)

	req := httptest.NewRequest(http.MethodPost, "/requests/xyz789/replay?target="+target.URL, nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var resp replayResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.RequestID != "xyz789" {
		t.Errorf("expected request_id xyz789, got %s", resp.RequestID)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status_code 200, got %d", resp.StatusCode)
	}
	if resp.Error != "" {
		t.Errorf("unexpected error: %s", resp.Error)
	}
}
