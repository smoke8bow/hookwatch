package server

import (
	"encoding/json"
	"net/http"

	"github.com/user/hookwatch/internal/storage"
)

// handleTagRequests handles GET /requests/tag?tag=<name>
// and returns all captured requests carrying that tag header.
func (h *Handler) handleTagRequests(w http.ResponseWriter, r *http.Request) {
	params := map[string]string{
		"tag": r.URL.Query().Get("tag"),
	}
	opts := ParseTagOptions(params)

	all := h.store.GetAll()
	matched := ApplyTagFilter(all, opts)

	if matched == nil {
		matched = []storage.Request{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(matched) //nolint:errcheck
}
