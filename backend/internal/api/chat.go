// Package api contains HTTP handlers for the TutorAI backend.
package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	appcontext "tutorai/backend/internal/context"
	"tutorai/backend/internal/llm"
	"tutorai/backend/internal/retrieval"

	"tutorai/backend/config"
)

var priceRe = regexp.MustCompile(`\$(\d+(?:\.\d+)?)`)

var knownFormats = []string{
	"commander", "modern", "standard", "legacy", "vintage",
	"pioneer", "pauper", "historic", "brawl",
}

// ChatHandler handles POST /chat requests.
type ChatHandler struct {
	cfg         config.Config
	httpClient  *http.Client
	colorLookup *retrieval.ColorLookup
	dataClient  *retrieval.Client
}

// NewChatHandler creates a ChatHandler with the given dependencies.
func NewChatHandler(
	cfg config.Config,
	httpClient *http.Client,
	colorLookup *retrieval.ColorLookup,
	dataClient *retrieval.Client,
) *ChatHandler {
	return &ChatHandler{
		cfg:         cfg,
		httpClient:  httpClient,
		colorLookup: colorLookup,
		dataClient:  dataClient,
	}
}

type chatRequest struct {
	Query string `json:"query"`
}

type chatResponse struct {
	Response string    `json:"response"`
	Usage    *llm.Usage `json:"usage,omitempty"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func (h *ChatHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req chatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid request body"})
		return
	}

	query := appcontext.SanitizeQuery(req.Query)
	if query == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "query is required"})
		return
	}

	ctx := r.Context()

	// Classify intent — falls back to general on any error.
	intent, _ := retrieval.ClassifyIntent(ctx, h.httpClient, h.cfg.OllamaBaseURL, h.cfg.OllamaLLMModel, query)

	// Retrieve relevant context based on intent.
	var results []retrieval.Result
	var retrieveErr error

	switch intent {
	case retrieval.IntentDeckBuilding:
		filters := extractCardFilters(query, h.colorLookup)
		results, retrieveErr = h.dataClient.RetrieveCards(ctx, filters)
	case retrieval.IntentCardLookup:
		results, retrieveErr = h.dataClient.RetrieveCards(ctx, retrieval.CardFilters{Query: query})
	case retrieval.IntentRulesQuestion:
		results, retrieveErr = h.dataClient.RetrieveRules(ctx, query)
	default: // general
		results, retrieveErr = h.dataClient.RetrieveSlang(ctx, query)
	}

	if retrieveErr != nil {
		var svcErr *retrieval.ServiceError
		if errors.As(retrieveErr, &svcErr) {
			writeJSON(w, http.StatusServiceUnavailable, errorResponse{Error: "data service unavailable"})
			return
		}
	}

	// Assemble prompt and call the LLM.
	prompt := appcontext.Assemble(query, results)
	response, usage, err := llm.Generate(ctx, h.httpClient, h.cfg.OllamaBaseURL, h.cfg.OllamaLLMModel, prompt)
	if err != nil {
		var unavailable *llm.UnavailableError
		if errors.As(err, &unavailable) {
			writeJSON(w, http.StatusServiceUnavailable, errorResponse{Error: "LLM unavailable"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "internal error"})
		return
	}

	writeJSON(w, http.StatusOK, chatResponse{Response: response, Usage: &usage})
}

// extractCardFilters scans the query for color identity, format, and price
// hints. All fields are best-effort — missing hints are left as zero values
// and the data service will omit the corresponding filter.
func extractCardFilters(query string, lookup *retrieval.ColorLookup) retrieval.CardFilters {
	lower := strings.ToLower(query)
	words := strings.Fields(lower)
	filters := retrieval.CardFilters{Query: query}

	// Color identity: try each word and each adjacent word pair.
	for i, word := range words {
		if colors := lookup.Resolve(word); colors != nil {
			filters.ColorIdentity = colors
			break
		}
		if i < len(words)-1 {
			if colors := lookup.Resolve(word + " " + words[i+1]); colors != nil {
				filters.ColorIdentity = colors
				break
			}
		}
	}

	// Format: first known format name found in the query.
	for _, f := range knownFormats {
		if strings.Contains(lower, f) {
			filters.Format = f
			break
		}
	}

	// Price: first "$X" or "$X.XX" pattern found.
	if m := priceRe.FindStringSubmatch(lower); m != nil {
		if price, err := strconv.ParseFloat(m[1], 64); err == nil {
			filters.MaxPriceUSD = price
		}
	}

	return filters
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
