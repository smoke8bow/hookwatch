package server

import (
	"encoding/json"
	"net/http"
)

// handleSearchRequests handles GET /search and returns requests matching
// full-text search criteria supplied via query parameters.
func (h *Handler) handleSearchRequests(w http.ResponseWriter, r *http.Request) {
	opts := ParseSearchOptions(r.URL.Query())

	all := h.store.GetAll()
	matched := ApplySearch(all, opts)

	w.Header().Set("Content-Type", "application/json")

	if len(matched) == 0 {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("[]"))
		return
	}

	if err := json.NewEncoder(w).Encode(matched); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
