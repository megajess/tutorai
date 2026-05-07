package retrieval

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"tutorai/backend/config"
)

// ---------------------------------------------------------------------------
// Intent classification
// ---------------------------------------------------------------------------

func TestParseIntent_KnownLabels(t *testing.T) {
	cases := []struct {
		input string
		want  Intent
	}{
		{"deck_building", IntentDeckBuilding},
		{"DECK_BUILDING", IntentDeckBuilding},
		{"card_lookup", IntentCardLookup},
		{"rules_question", IntentRulesQuestion},
		{"general", IntentGeneral},
		{"something_unexpected", IntentGeneral},
		{"", IntentGeneral},
	}
	for _, tc := range cases {
		got := parseIntent(tc.input)
		if got != tc.want {
			t.Errorf("parseIntent(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestClassifyIntent_ReturnsValidConstant(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"message": map[string]string{"content": "deck_building"},
		})
	}))
	defer srv.Close()

	intent, err := ClassifyIntent(context.Background(), srv.Client(), srv.URL, "llama3.1", "build me a golgari deck")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if intent != IntentDeckBuilding {
		t.Errorf("got intent %q, want %q", intent, IntentDeckBuilding)
	}
}

// ---------------------------------------------------------------------------
// Color identity lookup
// ---------------------------------------------------------------------------

func writeTempLookup(t *testing.T, data map[string][]string) string {
	t.Helper()
	b, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "lookup.json")
	if err := os.WriteFile(path, b, 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestColorLookup_ResolvesKnownNames(t *testing.T) {
	cases := []struct {
		term   string
		colors []string
	}{
		{"golgari", []string{"B", "G"}},
		{"GOLGARI", []string{"B", "G"}},
		{"bg", []string{"B", "G"}},
		{"esper", []string{"W", "U", "B"}},
		{"sultai", []string{"B", "G", "U"}},
	}

	lookup, err := LoadColorLookup("../../../data/color_identity_lookup.json")
	if err != nil {
		t.Fatalf("LoadColorLookup: %v", err)
	}

	for _, tc := range cases {
		got := lookup.Resolve(tc.term)
		if got == nil {
			t.Errorf("Resolve(%q) returned nil, want %v", tc.term, tc.colors)
			continue
		}
		if len(got) != len(tc.colors) {
			t.Errorf("Resolve(%q) = %v, want %v", tc.term, got, tc.colors)
		}
	}
}

func TestColorLookup_ReturnsNilForUnknown(t *testing.T) {
	path := writeTempLookup(t, map[string][]string{
		"golgari": {"B", "G"},
	})
	lookup, err := LoadColorLookup(path)
	if err != nil {
		t.Fatal(err)
	}
	if got := lookup.Resolve("izzet"); got != nil {
		t.Errorf("Resolve(unknown) = %v, want nil", got)
	}
	if got := lookup.Resolve(""); got != nil {
		t.Errorf("Resolve(\"\") = %v, want nil", got)
	}
}

// ---------------------------------------------------------------------------
// Data service client
// ---------------------------------------------------------------------------

func TestClient_ReturnsServiceErrorOnConnectionFailure(t *testing.T) {
	cfg := config.Config{
		DataServiceURL:    "http://localhost:19999", // nothing listening here
		DataServiceAPIKey: "testkey",
	}
	client := NewClient(cfg, &http.Client{})

	_, err := client.RetrieveRules(context.Background(), "how does deathtouch work")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var svcErr *ServiceError
	if !isServiceError(err, &svcErr) {
		t.Errorf("expected *ServiceError, got %T: %v", err, err)
	}
	if svcErr.StatusCode != 0 {
		t.Errorf("expected StatusCode 0 for connection failure, got %d", svcErr.StatusCode)
	}
}

func TestClient_SendsAPIKeyHeader(t *testing.T) {
	var gotKey string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotKey = r.Header.Get("X-API-Key")
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"results": []any{}})
	}))
	defer srv.Close()

	cfg := config.Config{DataServiceURL: srv.URL, DataServiceAPIKey: "secret-key"}
	client := NewClient(cfg, srv.Client())

	_, err := client.RetrieveSlang(context.Background(), "mana rock")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotKey != "secret-key" {
		t.Errorf("X-API-Key = %q, want %q", gotKey, "secret-key")
	}
}

// isServiceError checks if err is (or wraps) a *ServiceError and sets target.
func isServiceError(err error, target **ServiceError) bool {
	if se, ok := err.(*ServiceError); ok {
		*target = se
		return true
	}
	return false
}
