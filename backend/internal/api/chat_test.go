package api

import (
	"bufio"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"tutorai/backend/config"
	"tutorai/backend/internal/retrieval"
)

// mockOllama returns "general" on the first call (intent classification) and
// a canned LLM response on every subsequent call.
// Responses are newline-delimited JSON (Ollama streaming format).
func mockOllama(t *testing.T) *httptest.Server {
	t.Helper()
	var callCount atomic.Int32
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := callCount.Add(1)
		w.Header().Set("Content-Type", "application/x-ndjson")
		content := "general"
		if n > 1 {
			content = "Deathtouch is a keyword ability that causes damage dealt by a source with deathtouch to be considered lethal."
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"message": map[string]string{"content": content},
			"done":    true,
		})
	}))
}

// mockDataService returns an empty results list for any request.
func mockDataService(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"results": []any{}})
	}))
}

func newTestHandler(t *testing.T, ollamaURL, dataURL string) *ChatHandler {
	t.Helper()
	cfg := config.Config{
		OllamaBaseURL:     ollamaURL,
		OllamaLLMModel:    "llama3.1",
		DataServiceURL:    dataURL,
		DataServiceAPIKey: "testkey",
	}
	lookup, err := retrieval.LoadColorLookup("../../../data/color_identity_lookup.json")
	if err != nil {
		t.Fatalf("LoadColorLookup: %v", err)
	}
	client := retrieval.NewClient(cfg, http.DefaultClient)
	return NewChatHandler(cfg, http.DefaultClient, lookup, client)
}

// parseSSEChunks reads all SSE events from an httptest.ResponseRecorder body
// and returns the concatenated chunk text and the done event payload.
func parseSSEChunks(t *testing.T, body *strings.Reader) (text string, done map[string]any, errEvent string) {
	t.Helper()
	scanner := bufio.NewScanner(body)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		var event map[string]any
		if err := json.Unmarshal([]byte(line[6:]), &event); err != nil {
			t.Fatalf("parse SSE event: %v", err)
		}
		if chunk, ok := event["chunk"].(string); ok {
			text += chunk
		}
		if _, ok := event["done"]; ok {
			done = event
		}
		if errMsg, ok := event["error"].(string); ok {
			errEvent = errMsg
		}
	}
	return
}

func TestChatHandler_ReturnsNonEmptyResponse(t *testing.T) {
	ollama := mockOllama(t)
	defer ollama.Close()
	ds := mockDataService(t)
	defer ds.Close()

	handler := newTestHandler(t, ollama.URL, ds.URL)

	body := strings.NewReader(`{"query": "how does deathtouch work?"}`)
	req := httptest.NewRequest(http.MethodPost, "/chat", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	if ct := w.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/event-stream") {
		t.Fatalf("expected text/event-stream, got %q", ct)
	}

	text, _, errEvent := parseSSEChunks(t, strings.NewReader(w.Body.String()))
	if errEvent != "" {
		t.Fatalf("unexpected SSE error: %s", errEvent)
	}
	if text == "" {
		t.Error("expected non-empty streamed text")
	}
}

func TestChatHandler_EmptyQueryReturns400(t *testing.T) {
	handler := newTestHandler(t, "http://localhost:1", "http://localhost:2")

	body := strings.NewReader(`{"query": ""}`)
	req := httptest.NewRequest(http.MethodPost, "/chat", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestChatHandler_OllamaUnreachableReturns503(t *testing.T) {
	ds := mockDataService(t)
	defer ds.Close()

	// Use a port nothing is listening on for Ollama.
	handler := newTestHandler(t, "http://localhost:19998", ds.URL)

	body := strings.NewReader(`{"query": "how does deathtouch work?"}`)
	req := httptest.NewRequest(http.MethodPost, "/chat", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// Intent classification fails → falls back to general, then LLM stream
	// fails before any chunk → handler returns plain JSON 503.
	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503, got %d: %s", w.Code, w.Body.String())
	}
	var errResp errorResponse
	if err := json.NewDecoder(w.Body).Decode(&errResp); err != nil {
		t.Fatalf("decode error response: %v", err)
	}
	if errResp.Error != "LLM unavailable" {
		t.Errorf("expected 'LLM unavailable', got %q", errResp.Error)
	}
}

func TestChatHandler_DataServiceUnreachableReturns503(t *testing.T) {
	ollama := mockOllama(t)
	defer ollama.Close()

	// Use a port nothing is listening on for the data service.
	handler := newTestHandler(t, ollama.URL, "http://localhost:19997")

	body := strings.NewReader(`{"query": "how does deathtouch work?"}`)
	req := httptest.NewRequest(http.MethodPost, "/chat", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503, got %d: %s", w.Code, w.Body.String())
	}
	var errResp errorResponse
	if err := json.NewDecoder(w.Body).Decode(&errResp); err != nil {
		t.Fatalf("decode error response: %v", err)
	}
	if errResp.Error != "data service unavailable" {
		t.Errorf("expected 'data service unavailable', got %q", errResp.Error)
	}
}

func TestExtractCardFilters_ColorIdentity(t *testing.T) {
	lookup, err := retrieval.LoadColorLookup("../../../data/color_identity_lookup.json")
	if err != nil {
		t.Fatalf("LoadColorLookup: %v", err)
	}

	cases := []struct {
		query  string
		colors []string
	}{
		{"build me a golgari commander deck", []string{"B", "G"}},
		{"I want a simic deck under $50", []string{"G", "U"}},
		{"build me a black green aristocrats deck", []string{"B", "G"}},
		{"build me a random deck", nil},
	}

	for _, tc := range cases {
		filters := extractCardFilters(tc.query, lookup)
		if tc.colors == nil {
			if filters.ColorIdentity != nil {
				t.Errorf("query %q: expected nil color identity, got %v", tc.query, filters.ColorIdentity)
			}
		} else {
			if len(filters.ColorIdentity) != len(tc.colors) {
				t.Errorf("query %q: color identity = %v, want %v", tc.query, filters.ColorIdentity, tc.colors)
			}
		}
	}
}

func TestExtractCardFilters_FormatAndPrice(t *testing.T) {
	lookup, _ := retrieval.LoadColorLookup("../../../data/color_identity_lookup.json")

	filters := extractCardFilters("build a golgari commander deck under $75", lookup)
	if filters.Format != "commander" {
		t.Errorf("expected format 'commander', got %q", filters.Format)
	}
	if filters.MaxPriceUSD != 75 {
		t.Errorf("expected price 75, got %f", filters.MaxPriceUSD)
	}
}
