# Ticket 003 — Retrieval Layer & Data Service Client

## Status
`Todo`

## Type
`Feature`

## Summary
Build the retrieval layer in `backend/internal/retrieval/`. This includes intent classification, the local color identity lookup, and an HTTP client that calls the `rag-data-service` for card and rules retrieval. All written in Go. For local development, the data service is assumed to be running at the URL in `.env`.

## Background / Context
See `docs/architecture.md` for the full retrieval design and repo split. The Go backend never talks to Chroma or SQLite directly — all vector search and structured filtering is delegated to the data service over HTTP. This keeps the proprietary data layer fully decoupled from the open-source app code.

## Acceptance Criteria
- [ ] `backend/internal/retrieval/intent.go` — calls Ollama HTTP API to classify a query as `deck_building`, `card_lookup`, `rules_question`, or `general`. Returns a typed constant, not a raw string.
- [ ] `backend/internal/retrieval/lookup.go` — reads `data/color_identity_lookup.json` at startup and resolves guild/shard/wedge names to color identity slices. Returns `nil` if no match.
- [ ] `backend/internal/retrieval/client.go` — HTTP client wrapping the data service API. Methods: `RetrieveCards(ctx, filters, query)`, `RetrieveRules(ctx, query)`, `RetrieveSlang(ctx, query)`. Reads `DATA_SERVICE_URL` and `DATA_SERVICE_API_KEY` from config. Sends `X-API-Key` header on all requests.
- [ ] Client returns a typed error (not a raw `error` string) when the data service is unreachable
- [ ] `backend/internal/context/assemble.go` — assembles retrieved results into a prompt string ready to pass to Ollama
- [ ] Table-driven tests cover: intent returns valid constants, color lookup resolves known guild/shard/wedge names, lookup returns nil for unknown terms, client returns error on connection failure (use `httptest` to mock)

## Implementation Notes
- Ollama intent classification: `POST /api/chat` with a minimal system prompt and `stream: false`. Parse the response message content for the label.
- Color identity JSON is loaded once at startup into a `map[string][]string` — do not re-read the file on every request
- `RetrieveCards` filter struct should include: `ColorIdentity []string`, `Format string`, `MaxPriceUSD float64`, `Query string`
- Use `context.Context` on all client methods for timeout/cancellation support
- Keep the client thin — no retrieval logic, just HTTP marshalling/unmarshalling

## Relevant Areas
- `docs/architecture.md` (Repository Split, Data Flow sections)
- `backend/config/config.go`
- `data/color_identity_lookup.json`
- `.env.example`

## Dependencies
- Requires: #001
- Blocks: #004

## Out of Scope
- Do not talk to Chroma or SQLite directly
- Do not build the data service — separate private repo
- Do not implement conversation history
