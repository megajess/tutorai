// Package retrieval handles intent classification and data service communication.
package retrieval

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Intent represents the classified type of a user query.
type Intent string

const (
	IntentDeckBuilding  Intent = "deck_building"
	IntentCardLookup    Intent = "card_lookup"
	IntentRulesQuestion Intent = "rules_question"
	IntentGeneral       Intent = "general"
)

const intentSystemPrompt = `You are a query classifier for a Magic: The Gathering assistant.
Classify the user's query into exactly one of these categories:
- deck_building: the user wants help constructing or improving a deck
- card_lookup: the user is asking about a specific card or searching for cards with certain properties
- rules_question: the user is asking about game rules, interactions, or mechanics
- general: anything else, including greetings or questions that don't fit the above

Respond with exactly one of these words and nothing else:
deck_building, card_lookup, rules_question, general`

// ClassifyIntent sends the query to Ollama and returns the classified Intent.
// Falls back to IntentGeneral if the response cannot be parsed.
func ClassifyIntent(ctx context.Context, client *http.Client, ollamaBaseURL, model, query string) (Intent, error) {
	reqBody, err := json.Marshal(map[string]any{
		"model": model,
		"messages": []map[string]string{
			{"role": "system", "content": intentSystemPrompt},
			{"role": "user", "content": query},
		},
		"stream": false,
	})
	if err != nil {
		return IntentGeneral, fmt.Errorf("marshal intent request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, ollamaBaseURL+"/api/chat", bytes.NewReader(reqBody))
	if err != nil {
		return IntentGeneral, fmt.Errorf("create intent request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return IntentGeneral, fmt.Errorf("ollama intent request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	var result struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return IntentGeneral, fmt.Errorf("decode intent response: %w", err)
	}

	return parseIntent(result.Message.Content), nil
}

// parseIntent extracts a known Intent from the model response.
// Returns IntentGeneral if the response does not match any known label.
func parseIntent(raw string) Intent {
	switch strings.TrimSpace(strings.ToLower(raw)) {
	case string(IntentDeckBuilding):
		return IntentDeckBuilding
	case string(IntentCardLookup):
		return IntentCardLookup
	case string(IntentRulesQuestion):
		return IntentRulesQuestion
	default:
		return IntentGeneral
	}
}
