package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// NewRouter builds and returns the application router.
func NewRouter(h *Handler) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	// Capture all incoming webhooks
	r.Post("/webhook/{path:.*}", h.handleCapture)

	// Inspection
	r.Get("/requests", h.handleList)
	r.Get("/requests/{id}", h.handleGet)

	// Filter, search, tag
	r.Get("/requests/filter", h.handleFilterRequests)
	r.Get("/requests/search", h.handleSearchRequests)
	r.Get("/requests/tag", h.handleTagRequests)

	// Replay
	r.Post("/requests/{id}/replay", h.handleReplay)

	// Export
	r.Get("/requests/export", handleExportRequests(h.store))

	// Stats
	r.Get("/stats", h.handleStats)

	return r
}
