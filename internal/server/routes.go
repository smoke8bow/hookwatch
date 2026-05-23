package server

import "net/http"

// NewRouter wires all HTTP routes to their respective handlers.
func NewRouter(h *Handler) http.Handler {
	mux := http.NewServeMux()

	// Capture incoming webhooks on any path under /hook/
	mux.HandleFunc("/hook/", h.handleCaptureWebhook)

	// Inspect captured requests
	mux.HandleFunc("/requests", h.handleListRequests)
	mux.HandleFunc("/requests/", h.handleGetRequest)

	// Filter requests by method, path prefix, or header
	mux.HandleFunc("/filter", h.handleFilterRequests)

	// Full-text search across body, headers, and path
	mux.HandleFunc("/search", h.handleSearchRequests)

	// Replay a captured request to a target URL
	mux.HandleFunc("/replay/", h.handleReplay)

	// Export captured requests as JSON or NDJSON
	mux.HandleFunc("/export", handleExportRequests(h.store))

	// Aggregated stats
	mux.HandleFunc("/stats", h.handleStats)

	return mux
}
