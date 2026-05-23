package server

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"hookwatch/internal/storage"
)

// Handler holds dependencies for HTTP handlers.
type Handler struct {
	store *storage.Store
}

// NewHandler creates a new Handler with the given store.
func NewHandler(store *storage.Store) *Handler {
	return &Handler{store: store}
}

// CaptureWebhook receives an inbound webhook and stores it.
func (h *Handler) CaptureWebhook(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20)) // 1 MB limit
	if err != nil {
		http.Error(w, "failed to read body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	headers := make(map[string]string)
	for key := range r.Header {
		headers[key] = r.Header.Get(key)
	}

	req := &storage.Request{
		ID:        uuid.NewString(),
		Method:    r.Method,
		Path:      r.URL.Path,
		Headers:   headers,
		Body:      body,
		ReceivedAt: time.Now().UTC(),
	}

	h.store.Add(req)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"id": req.ID})
}

// ListRequests returns all captured requests as JSON.
func (h *Handler) ListRequests(w http.ResponseWriter, r *http.Request) {
	requests := h.store.GetAll()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(requests)
}

// GetRequest returns a single captured request by ID.
func (h *Handler) GetRequest(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	req, ok := h.store.GetByID(id)
	if !ok {
		http.Error(w, "request not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(req)
}
