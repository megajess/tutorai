// Package llm contains the Ollama HTTP client for LLM inference.
package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// UnavailableError is returned when Ollama cannot be reached.
type UnavailableError struct {
	Cause error
}

func (e *UnavailableError) Error() string {
	return fmt.Sprintf("ollama unavailable: %v", e.Cause)
}

func (e *UnavailableError) Unwrap() error {
	return e.Cause
}

// Generate sends a prompt to Ollama and returns the response text.
// The prompt is sent as the user message; Ollama's own default behaviour
// handles the system context embedded in the prompt string.
func Generate(ctx context.Context, client *http.Client, baseURL, model, prompt string) (string, error) {
	reqBody, err := json.Marshal(map[string]any{
		"model": model,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
		"stream": false,
	})
	if err != nil {
		return "", fmt.Errorf("marshal llm request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+"/api/chat", bytes.NewReader(reqBody))
	if err != nil {
		return "", fmt.Errorf("create llm request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", &UnavailableError{Cause: err}
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", &UnavailableError{Cause: fmt.Errorf("status %d", resp.StatusCode)}
	}

	var result struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode llm response: %w", err)
	}

	return result.Message.Content, nil
}
