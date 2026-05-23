package server

import (
	"encoding/json"
	"net/http"
)

// handleFilterRequests handles GET /requests/filter and returns filtered requests.
func (h *Handler) handleFilterRequests(w http.ResponseWriter, r *http.Request) {
	opts := ParseFilterOptions(r)
	all := h.store.GetAll()
	filtered := ApplyFilter(all, opts)

	type response struct {
		Count    int         `json:"count"`
		Requests interface{} `json:"requests"`
	}

	if filtered == nil {
		filtered = []*storageRequest{}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response{
		Count:    len(filtered),
		Requests: filtered,
	})
}
