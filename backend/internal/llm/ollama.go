// Package llm contains the Ollama HTTP client for LLM inference.
package llm

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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

// ollamaChunk is a single line from Ollama's streaming response.
type ollamaChunk struct {
	Message struct {
		Content string `json:"content"`
	} `json:"message"`
	Done            bool `json:"done"`
	PromptEvalCount int  `json:"prompt_eval_count"`
	EvalCount       int  `json:"eval_count"`
}

func ollamaRequest(model, prompt string) ([]byte, error) {
	return json.Marshal(map[string]any{
		"model": model,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
		"stream": true,
	})
}

// Stream sends a prompt to Ollama and calls onChunk for each partial response.
// Returns token usage from the final chunk.
func Stream(ctx context.Context, client *http.Client, baseURL, model, prompt string, onChunk func(string)) (Usage, error) {
	reqBody, err := ollamaRequest(model, prompt)
	if err != nil {
		return Usage{}, fmt.Errorf("marshal llm request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+"/api/chat", bytes.NewReader(reqBody))
	if err != nil {
		return Usage{}, fmt.Errorf("create llm request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return Usage{}, &UnavailableError{Cause: err}
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return Usage{}, &UnavailableError{Cause: fmt.Errorf("status %d", resp.StatusCode)}
	}

	var usage Usage
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		var chunk ollamaChunk
		if err := json.Unmarshal(scanner.Bytes(), &chunk); err != nil {
			continue
		}
		if chunk.Message.Content != "" {
			onChunk(chunk.Message.Content)
		}
		if chunk.Done {
			usage = Usage{
				PromptTokens:     chunk.PromptEvalCount,
				CompletionTokens: chunk.EvalCount,
				TotalTokens:      chunk.PromptEvalCount + chunk.EvalCount,
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return usage, fmt.Errorf("read llm stream: %w", err)
	}

	return usage, nil
}

// Generate is a convenience wrapper around Stream that returns the full
// response text at once. Used in tests and intent classification fallbacks.
func Generate(ctx context.Context, client *http.Client, baseURL, model, prompt string) (string, Usage, error) {
	var sb strings.Builder
	usage, err := Stream(ctx, client, baseURL, model, prompt, func(chunk string) {
		sb.WriteString(chunk)
	})
	return sb.String(), usage, err
}
