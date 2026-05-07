package context

import (
	"strings"
	"testing"

	"tutorai/backend/internal/retrieval"
)

func TestSanitizeQuery_TruncatesLongInput(t *testing.T) {
	long := strings.Repeat("a", maxQueryLen+100)
	got := SanitizeQuery(long)
	if len(got) > maxQueryLen {
		t.Errorf("SanitizeQuery did not truncate: len=%d, want <=%d", len(got), maxQueryLen)
	}
}

func TestSanitizeQuery_StripsControlCharacters(t *testing.T) {
	input := "hello\x00world\x01test"
	got := SanitizeQuery(input)
	if strings.ContainsAny(got, "\x00\x01") {
		t.Errorf("SanitizeQuery did not strip control chars: %q", got)
	}
	if !strings.Contains(got, "hello") || !strings.Contains(got, "world") {
		t.Errorf("SanitizeQuery stripped too much: %q", got)
	}
}

func TestAssemble_UserQueryAppearsLast(t *testing.T) {
	results := []retrieval.Result{
		{ID: "1", Text: "Deathtouch rules text", Score: 0.9},
	}
	prompt := Assemble("how does deathtouch work", results)

	contextIdx := strings.Index(prompt, "CONTEXT")
	queryIdx := strings.Index(prompt, "how does deathtouch work")
	if contextIdx < 0 || queryIdx < 0 {
		t.Fatalf("missing expected sections in prompt: %q", prompt)
	}
	if queryIdx < contextIdx {
		t.Error("user query appears before context — prompt injection risk")
	}
}

func TestAssemble_NoResultsIncludesFallbackMessage(t *testing.T) {
	prompt := Assemble("some query", nil)
	if !strings.Contains(prompt, "No relevant context") {
		t.Errorf("expected fallback message in prompt, got: %q", prompt)
	}
}

func TestAssemble_CardNameIncludedWhenPresent(t *testing.T) {
	results := []retrieval.Result{
		{ID: "abc", Name: "Viscera Seer", Text: "Sacrifice a creature: Scry 1.", Score: 0.95},
	}
	prompt := Assemble("what does viscera seer do", results)
	if !strings.Contains(prompt, "Viscera Seer") {
		t.Errorf("card name missing from prompt: %q", prompt)
	}
}
