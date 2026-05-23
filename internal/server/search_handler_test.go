package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/hookwatch/internal/storage"
)

func seedSearchStore(t *testing.T) *storage.Store {
	t.Helper()
	st := storage.NewStore(100)
	st.Add(&storage.Request{
		ID:      "id-1",
		Method:  "POST",
		Path:    "/webhook/github",
		Headers: http.Header{"Content-Type": {"application/json"}},
		Body:    []byte(`{"event":"push"}`),
	})
	st.Add(&storage.Request{
		ID:      "id-2",
		Method:  "POST",
		Path:    "/webhook/stripe",
		Headers: http.Header{"Content-Type": {"text/plain"}},
		Body:    []byte(`payment.succeeded`),
	})
	return st
}

func TestHandleSearch_NoMatch(t *testing.T) {
	router := newTestRouter(t)
	req := httptest.NewRequest(http.MethodGet, "/search?body=nonexistent", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if rec.Body.String() != "[]" {
		t.Errorf("expected empty array, got %s", rec.Body.String())
	}
}

func TestHandleSearch_MatchBody(t *testing.T) {
	st := seedSearchStore(t)
	h := NewHandler(st)
	router := NewRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/search?body=push", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var results []map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&results); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
}

func TestHandleSearch_MatchPath(t *testing.T) {
	st := seedSearchStore(t)
	h := NewHandler(st)
	router := NewRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/search?path=/webhook", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var results []map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&results); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
}
