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

// Usage holds token counts returned by Ollama.
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// Generate sends a prompt to Ollama and returns the response text and token usage.
// The prompt is sent as the user message; Ollama's own default behaviour
// handles the system context embedded in the prompt string.
func Generate(ctx context.Context, client *http.Client, baseURL, model, prompt string) (string, Usage, error) {
	reqBody, err := json.Marshal(map[string]any{
		"model": model,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
		"stream": false,
	})
	if err != nil {
		return "", Usage{}, fmt.Errorf("marshal llm request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+"/api/chat", bytes.NewReader(reqBody))
	if err != nil {
		return "", Usage{}, fmt.Errorf("create llm request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", Usage{}, &UnavailableError{Cause: err}
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", Usage{}, &UnavailableError{Cause: fmt.Errorf("status %d", resp.StatusCode)}
	}

	var result struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
		PromptEvalCount int `json:"prompt_eval_count"`
		EvalCount       int `json:"eval_count"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", Usage{}, fmt.Errorf("decode llm response: %w", err)
	}

	usage := Usage{
		PromptTokens:     result.PromptEvalCount,
		CompletionTokens: result.EvalCount,
		TotalTokens:      result.PromptEvalCount + result.EvalCount,
	}

	return result.Message.Content, usage, nil
}
