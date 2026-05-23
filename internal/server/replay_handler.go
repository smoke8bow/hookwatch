package server

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type replayResponse struct {
	RequestID  string `json:"request_id"`
	StatusCode int    `json:"status_code"`
	ForwardURL string `json:"forward_url"`
	Error      string `json:"error,omitempty"`
}

func (h *Handler) handleReplay(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "missing request id", http.StatusBadRequest)
		return
	}

	req, ok := h.store.GetByID(id)
	if !ok {
		http.Error(w, "request not found", http.StatusNotFound)
		return
	}

	targetURL := r.URL.Query().Get("target")
	if targetURL == "" {
		http.Error(w, "missing target query parameter", http.StatusBadRequest)
		return
	}

	statusCode, err := doReplay(req, targetURL)
	resp := replayResponse{
		RequestID:  id,
		StatusCode: statusCode,
		ForwardURL: targetURL,
	}
	if err != nil {
		resp.Error = err.Error()
		w.WriteHeader(http.StatusBadGateway)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
