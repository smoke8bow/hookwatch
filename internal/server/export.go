package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/hookwatch/internal/storage"
)

type exportRecord struct {
	ID        string              `json:"id"`
	Timestamp time.Time           `json:"timestamp"`
	Method    string              `json:"method"`
	Path      string              `json:"path"`
	Headers   map[string][]string `json:"headers"`
	Body      string              `json:"body"`
}

func toExportRecord(r storage.Request) exportRecord {
	return exportRecord{
		ID:        r.ID,
		Timestamp: r.Timestamp,
		Method:    r.Method,
		Path:      r.Path,
		Headers:   r.Headers,
		Body:      r.Body,
	}
}

func handleExportRequests(store *storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requests := store.GetAll()

		format := r.URL.Query().Get("format")
		if format == "" {
			format = "json"
		}

		switch format {
		case "json":
			records := make([]exportRecord, len(requests))
			for i, req := range requests {
				records[i] = toExportRecord(req)
			}
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="hookwatch-export-%s.json"`, time.Now().Format("20060102-150405")))
			json.NewEncoder(w).Encode(records)

		case "ndjson":
			w.Header().Set("Content-Type", "application/x-ndjson")
			w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="hookwatch-export-%s.ndjson"`, time.Now().Format("20060102-150405")))
			enc := json.NewEncoder(w)
			for _, req := range requests {
				enc.Encode(toExportRecord(req))
			}

		default:
			http.Error(w, fmt.Sprintf("unsupported format: %q (use json or ndjson)", format), http.StatusBadRequest)
		}
	}
}
