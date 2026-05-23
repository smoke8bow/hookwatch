package server

import (
	"net/http"

	"hookwatch/internal/storage"
)

// NewRouter builds and returns the application's HTTP router.
func NewRouter(store *storage.Store) http.Handler {
	h := NewHandler(store)

	mux := http.NewServeMux()

	// Webhook capture endpoint — accepts any method under /hooks/
	mux.HandleFunc("/hooks/", h.CaptureWebhook)

	// Inspection API
	mux.HandleFunc("GET /api/requests", h.ListRequests)
	mux.HandleFunc("GET /api/requests/{id}", h.GetRequest)

	return mux
}
