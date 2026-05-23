package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleFilterRequests_NoResults(t *testing.T) {
	router := newTestRouter()
	req := httptest.NewRequest(http.MethodGet, "/requests/filter?method=DELETE", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	count := int(resp["count"].(float64))
	if count != 0 {
		t.Errorf("expected 0 results, got %d", count)
	}
}

func TestHandleFilterRequests_MatchMethod(t *testing.T) {
	router := newTestRouter()

	// First capture a webhook
	body := `{"event":"test"}`
	captureReq := httptest.NewRequest(http.MethodPost, "/hooks/test", stringReader(body))
	captureReq.Header.Set("Content-Type", "application/json")
	captureW := httptest.NewRecorder()
	router.ServeHTTP(captureW, captureReq)
	if captureW.Code != http.StatusOK {
		t.Fatalf("capture failed: %d", captureW.Code)
	}

	// Now filter by POST
	filterReq := httptest.NewRequest(http.MethodGet, "/requests/filter?method=POST", nil)
	filterW := httptest.NewRecorder()
	router.ServeHTTP(filterW, filterReq)

	if filterW.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", filterW.Code)
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(filterW.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	count := int(resp["count"].(float64))
	if count != 1 {
		t.Errorf("expected 1 POST result, got %d", count)
	}
}

func stringReader(s string) *strings.Reader {
	return strings.NewReader(s)
}
