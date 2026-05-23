package server

import (
	"encoding/json"
	"net/http"
)

// handleStats serves a JSON snapshot of the current request statistics.
func (h *Handler) handleStats(w http.ResponseWriter, r *http.Request) {
	snap := h.stats.Snapshot()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(snap); err != nil {
		http.Error(w, "failed to encode stats", http.StatusInternalServerError)
	}
}
