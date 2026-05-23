package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(h *Handler) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Webhook capture endpoint
	r.Post("/hooks/{path:.*}", h.handleCapture)

	// Request inspection endpoints
	r.Get("/requests", h.handleListRequests)
	r.Get("/requests/filter", h.handleFilterRequests)
	r.Get("/requests/{id}", h.handleGetRequest)

	// Replay endpoint
	r.Post("/requests/{id}/replay", h.handleReplay)

	return r
}
