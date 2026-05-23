package server

import (
	"net/http"

	"github.com/user/hookwatch/internal/storage"
)

func NewRouter(store *storage.Store, stats *Stats) http.Handler {
	mux := http.NewServeMux()

	h := NewHandler(store, stats)

	// Capture incoming webhooks
	mux.HandleFunc("/webhook/", h.CaptureWebhook)

	// List all captured requests
	mux.HandleFunc("/requests", h.ListRequests)

	// Get a specific request by ID
	mux.HandleFunc("/requests/", h.GetRequest)

	// Filter requests
	mux.HandleFunc("/requests/filter", handleFilterRequests(store))

	// Export requests
	mux.HandleFunc("/requests/export", handleExportRequests(store))

	// Replay a request
	mux.HandleFunc("/replay/", handleReplay(store))

	// Stats
	mux.HandleFunc("/stats", handleStats(stats))

	return mux
}
