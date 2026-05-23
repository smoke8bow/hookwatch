package server

import (
	"net/http"
	"strings"

	"github.com/user/hookwatch/internal/storage"
)

// FilterOptions holds criteria for filtering stored requests.
type FilterOptions struct {
	Method      string
	HeaderKey   string
	HeaderValue string
	PathPrefix  string
}

// ParseFilterOptions reads query parameters from a request and builds FilterOptions.
func ParseFilterOptions(r *http.Request) FilterOptions {
	q := r.URL.Query()
	return FilterOptions{
		Method:      strings.ToUpper(q.Get("method")),
		HeaderKey:   q.Get("header_key"),
		HeaderValue: q.Get("header_value"),
		PathPrefix:  q.Get("path_prefix"),
	}
}

// ApplyFilter returns only those requests that match all non-empty filter criteria.
func ApplyFilter(requests []*storage.Request, opts FilterOptions) []*storage.Request {
	var out []*storage.Request
	for _, req := range requests {
		if opts.Method != "" && req.Method != opts.Method {
			continue
		}
		if opts.PathPrefix != "" && !strings.HasPrefix(req.Path, opts.PathPrefix) {
			continue
		}
		if opts.HeaderKey != "" {
			vals, ok := req.Headers[opts.HeaderKey]
			if !ok {
				continue
			}
			if opts.HeaderValue != "" {
				found := false
				for _, v := range vals {
					if strings.EqualFold(v, opts.HeaderValue) {
						found = true
						break
					}
				}
				if !found {
					continue
				}
			}
		}
		out = append(out, req)
	}
	return out
}
