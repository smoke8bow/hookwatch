package server_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"hookwatch/internal/server"
	"hookwatch/internal/storage"
)

func newTestRouter(t *testing.T) http.Handler {
	t.Helper()
	store := storage.NewStore(50)
	return server.NewRouter(store)
}

func TestCaptureWebhook(t *testing.T) {
	router := newTestRouter(t)

	body := bytes.NewBufferString(`{"event":"push"}`)
	req := httptest.NewRequest(http.MethodPost, "/hooks/github", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp["id"] == "" {
		t.Error("expected non-empty id in response")
	}
}

func TestListRequests(t *testing.T) {
	router := newTestRouter(t)

	// Capture one request first
	router.ServeHTTP(httptest.NewRecorder(),
		httptest.NewRequest(http.MethodPost, "/hooks/test", bytes.NewBufferString(`{}`)))

	req := httptest.NewRequest(http.MethodGet, "/api/requests", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var requests []*storage.Request
	if err := json.NewDecoder(rec.Body).Decode(&requests); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(requests) != 1 {
		t.Errorf("expected 1 request, got %d", len(requests))
	}
}

func TestGetRequest_NotFound(t *testing.T) {
	router := newTestRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/api/requests/nonexistent", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rec.Code)
	}
}
