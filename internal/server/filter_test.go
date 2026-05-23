package server

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/user/hookwatch/internal/storage"
)

func makeStorageRequest(id, method, path string, headers map[string][]string) *storage.Request {
	return &storage.Request{
		ID:        id,
		Method:    method,
		Path:      path,
		Headers:   headers,
		Body:      []byte(`{}`),
		ReceivedAt: time.Now(),
	}
}

func TestParseFilterOptions(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?method=post&path_prefix=/hooks&header_key=X-Event&header_value=push", nil)
	opts := ParseFilterOptions(req)
	if opts.Method != "POST" {
		t.Errorf("expected POST, got %s", opts.Method)
	}
	if opts.PathPrefix != "/hooks" {
		t.Errorf("expected /hooks, got %s", opts.PathPrefix)
	}
	if opts.HeaderKey != "X-Event" {
		t.Errorf("expected X-Event, got %s", opts.HeaderKey)
	}
	if opts.HeaderValue != "push" {
		t.Errorf("expected push, got %s", opts.HeaderValue)
	}
}

func TestApplyFilter_Method(t *testing.T) {
	requests := []*storage.Request{
		makeStorageRequest("1", "POST", "/hook", nil),
		makeStorageRequest("2", "GET", "/hook", nil),
	}
	result := ApplyFilter(requests, FilterOptions{Method: "POST"})
	if len(result) != 1 || result[0].ID != "1" {
		t.Errorf("expected 1 POST result, got %d", len(result))
	}
}

func TestApplyFilter_PathPrefix(t *testing.T) {
	requests := []*storage.Request{
		makeStorageRequest("1", "POST", "/hooks/github", nil),
		makeStorageRequest("2", "POST", "/api/other", nil),
	}
	result := ApplyFilter(requests, FilterOptions{PathPrefix: "/hooks"})
	if len(result) != 1 || result[0].ID != "1" {
		t.Errorf("expected 1 result with /hooks prefix, got %d", len(result))
	}
}

func TestApplyFilter_Header(t *testing.T) {
	requests := []*storage.Request{
		makeStorageRequest("1", "POST", "/hook", map[string][]string{"X-Event": {"push"}}),
		makeStorageRequest("2", "POST", "/hook", map[string][]string{"X-Event": {"pull_request"}}),
		makeStorageRequest("3", "POST", "/hook", nil),
	}
	result := ApplyFilter(requests, FilterOptions{HeaderKey: "X-Event", HeaderValue: "push"})
	if len(result) != 1 || result[0].ID != "1" {
		t.Errorf("expected 1 result matching header, got %d", len(result))
	}
}

func TestApplyFilter_Empty(t *testing.T) {
	_ = url.Values{} // ensure import used
	requests := []*storage.Request{
		makeStorageRequest("1", "POST", "/hook", nil),
		makeStorageRequest("2", "GET", "/other", nil),
	}
	result := ApplyFilter(requests, FilterOptions{})
	if len(result) != 2 {
		t.Errorf("expected all 2 results with empty filter, got %d", len(result))
	}
}
