package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/hookwatch/internal/storage"
)

func seedTagStore(tag string) *storage.Store {
	s := storage.NewStore(10)
	headers := http.Header{}
	if tag != "" {
		headers.Set("X-Hookwatch-Tag", tag)
	}
	s.Add(storage.Request{
		ID:      "t1",
		Method:  "POST",
		Path:    "/hook",
		Headers: headers,
		Body:    []byte(`{"event":"ping"}`),
	})
	return s
}

func TestHandleTagRequests_NoMatch(t *testing.T) {
	s := seedTagStore("payments")
	router := newTestRouter(s)

	req := httptest.NewRequest(http.MethodGet, "/requests/tag?tag=orders", nil)
	rw := httptest.NewRecorder()
	router.ServeHTTP(rw, req)

	if rw.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rw.Code)
	}
	var results []storage.Request
	if err := json.NewDecoder(rw.Body).Decode(&results); err != nil {
		t.Fatal(err)
	}
	if len(results) != 0 {
		t.Fatalf("expected 0, got %d", len(results))
	}
}

func TestHandleTagRequests_Match(t *testing.T) {
	s := seedTagStore("orders")
	router := newTestRouter(s)

	req := httptest.NewRequest(http.MethodGet, "/requests/tag?tag=orders", nil)
	rw := httptest.NewRecorder()
	router.ServeHTTP(rw, req)

	if rw.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rw.Code)
	}
	var results []storage.Request
	if err := json.NewDecoder(rw.Body).Decode(&results); err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1, got %d", len(results))
	}
}
