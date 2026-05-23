package server

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"hookwatch/internal/storage"
)

// ReplayResult holds the outcome of a replayed webhook request.
type ReplayResult struct {
	RequestID  string        `json:"request_id"`
	TargetURL  string        `json:"target_url"`
	StatusCode int           `json:"status_code"`
	Duration   time.Duration `json:"duration_ms"`
	Error      string        `json:"error,omitempty"`
}

// replayRequest replays a captured webhook request to a given target URL.
// POST /requests/{id}/replay?target=<url>
func (h *Handler) replayRequest(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	targetURL := r.URL.Query().Get("target")

	if targetURL == "" {
		http.Error(w, `{"error":"target query parameter is required"}`, http.StatusBadRequest)
		return
	}

	req, ok := h.store.GetByID(id)
	if !ok {
		http.Error(w, `{"error":"request not found"}`, http.StatusNotFound)
		return
	}

	result, err := doReplay(req, targetURL)
	if err != nil {
		result.Error = err.Error()
		writeJSON(w, http.StatusBadGateway, result)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// doReplay performs the actual HTTP replay of a stored request.
func doReplay(req *storage.Request, targetURL string) (*ReplayResult, error) {
	result := &ReplayResult{
		RequestID: req.ID,
		TargetURL: targetURL,
	}

	body := bytes.NewReader(req.Body)

	outbound, err := http.NewRequest(req.Method, targetURL, body)
	if err != nil {
		return result, fmt.Errorf("building replay request: %w", err)
	}

	// Copy original headers, excluding hop-by-hop headers.
	for key, values := range req.Headers {
		for _, v := range values {
			outbound.Header.Add(key, v)
		}
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	start := time.Now()
	resp, err := client.Do(outbound)
	result.Duration = time.Since(start)

	if err != nil {
		return result, fmt.Errorf("replaying request: %w", err)
	}
	defer resp.Body.Close()

	// Drain the response body so the connection can be reused.
	_, _ = io.Copy(io.Discard, resp.Body)

	result.StatusCode = resp.StatusCode
	return result, nil
}
