package server

import (
	"strings"

	"github.com/user/hookwatch/internal/storage"
)

// SearchOptions holds criteria for full-text search across captured requests.
type SearchOptions struct {
	BodyContains   string
	HeaderContains string
	PathContains   string
}

// ParseSearchOptions extracts search parameters from query string values.
func ParseSearchOptions(query map[string][]string) SearchOptions {
	get := func(key string) string {
		if vals, ok := query[key]; ok && len(vals) > 0 {
			return strings.TrimSpace(vals[0])
		}
		return ""
	}
	return SearchOptions{
		BodyContains:   get("body"),
		HeaderContains: get("header_value"),
		PathContains:   get("path"),
	}
}

// ApplySearch filters requests using full-text search across body, headers, and path.
func ApplySearch(requests []*storage.Request, opts SearchOptions) []*storage.Request {
	if opts.BodyContains == "" && opts.HeaderContains == "" && opts.PathContains == "" {
		return requests
	}

	var result []*storage.Request
	for _, r := range requests {
		if opts.PathContains != "" && !strings.Contains(r.Path, opts.PathContains) {
			continue
		}
		if opts.BodyContains != "" && !strings.Contains(string(r.Body), opts.BodyContains) {
			continue
		}
		if opts.HeaderContains != "" {
			found := false
			for _, vals := range r.Headers {
				for _, v := range vals {
					if strings.Contains(v, opts.HeaderContains) {
						found = true
						break
					}
				}
				if found {
					break
				}
			}
			if !found {
				continue
			}
		}
		result = append(result, r)
	}
	return result
}
