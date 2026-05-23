package server

import (
	"net/http"
	"testing"
	"time"

	"github.com/user/hookwatch/internal/storage"
)

func makeTagRequest(tag string) storage.Request {
	headers := http.Header{}
	if tag != "" {
		headers.Set("X-Hookwatch-Tag", tag)
	}
	return storage.Request{
		ID:        "abc",
		Method:    "POST",
		Path:      "/hook",
		Headers:   headers,
		Body:      []byte(`{}`),
		Timestamp: time.Now(),
	}
}

func TestParseTagOptions(t *testing.T) {
	opts := ParseTagOptions(map[string]string{"tag": "  payments  "})
	if opts.Tag != "payments" {
		t.Fatalf("expected 'payments', got %q", opts.Tag)
	}
}

func TestApplyTagFilter_NoTag(t *testing.T) {
	reqs := []storage.Request{makeTagRequest("orders"), makeTagRequest("")}
	out := ApplyTagFilter(reqs, TagOptions{})
	if len(out) != 2 {
		t.Fatalf("expected 2 results, got %d", len(out))
	}
}

func TestApplyTagFilter_Match(t *testing.T) {
	reqs := []storage.Request{
		makeTagRequest("orders"),
		makeTagRequest("payments"),
		makeTagRequest("orders"),
	}
	out := ApplyTagFilter(reqs, TagOptions{Tag: "orders"})
	if len(out) != 2 {
		t.Fatalf("expected 2, got %d", len(out))
	}
}

func TestApplyTagFilter_CaseInsensitive(t *testing.T) {
	reqs := []storage.Request{makeTagRequest("Orders")}
	out := ApplyTagFilter(reqs, TagOptions{Tag: "orders"})
	if len(out) != 1 {
		t.Fatalf("expected 1, got %d", len(out))
	}
}

func TestApplyTagFilter_NoMatch(t *testing.T) {
	reqs := []storage.Request{makeTagRequest("payments")}
	out := ApplyTagFilter(reqs, TagOptions{Tag: "orders"})
	if len(out) != 0 {
		t.Fatalf("expected 0, got %d", len(out))
	}
}
