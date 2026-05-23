package server

import (
	"net/http"
	"testing"

	"github.com/user/hookwatch/internal/storage"
)

func makeSearchRequest(method, path, body string, headers map[string][]string) *storage.Request {
	return &storage.Request{
		ID:      "test-id",
		Method:  method,
		Path:    path,
		Headers: headers,
		Body:    []byte(body),
	}
}

func TestParseSearchOptions(t *testing.T) {
	query := map[string][]string{
		"body":         {"hello"},
		"header_value": {"application/json"},
		"path":         {"/webhook"},
	}
	opts := ParseSearchOptions(query)
	if opts.BodyContains != "hello" {
		t.Errorf("expected body 'hello', got %q", opts.BodyContains)
	}
	if opts.HeaderContains != "application/json" {
		t.Errorf("expected header_value 'application/json', got %q", opts.HeaderContains)
	}
	if opts.PathContains != "/webhook" {
		t.Errorf("expected path '/webhook', got %q", opts.PathContains)
	}
}

func TestApplySearch_Body(t *testing.T) {
	reqs := []*storage.Request{
		makeSearchRequest("POST", "/a", `{"event":"push"}`, nil),
		makeSearchRequest("POST", "/b", `{"event":"pull_request"}`, nil),
	}
	result := ApplySearch(reqs, SearchOptions{BodyContains: "push"})
	if len(result) != 1 || result[0].Path != "/a" {
		t.Errorf("expected 1 result with path /a, got %v", result)
	}
}

func TestApplySearch_Path(t *testing.T) {
	reqs := []*storage.Request{
		makeSearchRequest("GET", "/webhook/github", "", nil),
		makeSearchRequest("GET", "/webhook/stripe", "", nil),
		makeSearchRequest("GET", "/other", "", nil),
	}
	result := ApplySearch(reqs, SearchOptions{PathContains: "/webhook"})
	if len(result) != 2 {
		t.Errorf("expected 2 results, got %d", len(result))
	}
}

func TestApplySearch_HeaderValue(t *testing.T) {
	reqs := []*storage.Request{
		makeSearchRequest("POST", "/a", "", http.Header{"Content-Type": {"application/json"}}),
		makeSearchRequest("POST", "/b", "", http.Header{"Content-Type": {"text/plain"}}),
	}
	result := ApplySearch(reqs, SearchOptions{HeaderContains: "application/json"})
	if len(result) != 1 || result[0].Path != "/a" {
		t.Errorf("expected 1 result, got %v", result)
	}
}

func TestApplySearch_NoOptions(t *testing.T) {
	reqs := []*storage.Request{
		makeSearchRequest("GET", "/a", "", nil),
		makeSearchRequest("GET", "/b", "", nil),
	}
	result := ApplySearch(reqs, SearchOptions{})
	if len(result) != 2 {
		t.Errorf("expected all 2 results, got %d", len(result))
	}
}
