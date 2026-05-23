package server

import (
	"strings"

	"github.com/user/hookwatch/internal/storage"
)

// TagOptions holds the tag name to filter by.
type TagOptions struct {
	Tag string
}

// ParseTagOptions extracts tag query param from a map.
func ParseTagOptions(params map[string]string) TagOptions {
	return TagOptions{
		Tag: strings.TrimSpace(params["tag"]),
	}
}

// ApplyTagFilter returns only requests that carry the given tag header value.
// Requests are tagged via the X-Hookwatch-Tag request header at capture time.
func ApplyTagFilter(requests []storage.Request, opts TagOptions) []storage.Request {
	if opts.Tag == "" {
		return requests
	}
	var out []storage.Request
	for _, r := range requests {
		for _, v := range r.Headers["X-Hookwatch-Tag"] {
			if strings.EqualFold(v, opts.Tag) {
				out = append(out, r)
				break
			}
		}
	}
	return out
}
