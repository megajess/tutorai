# Ticket 006 — Sealed Deck Building

## Status
`Todo`

## Type
`Feature`

## Summary
Add a sealed deck building flow to TutorAI. A user can list the cards they opened from booster packs and receive a deck suggestion built from that specific pool. This is distinct from the existing deck building flow, which searches the full card database — sealed constrains the answer to exactly the cards the user has.

## Background / Context
At pre-release events, players open ~6 booster packs and must construct a 40-card deck (plus basic lands) from only those cards to play in a small tournament. The system needs to:
1. Recognise this as a different intent from open-ended deck building
2. Extract the card names the user listed
3. Look up those exact cards in the data service
4. Suggest a deck from that specific pool

The data service name filter is implemented in `rag-data-service` ticket #008.

## Acceptance Criteria
- [ ] `IntentSealedBuilding` added to the `Intent` type in `backend/internal/retrieval/intent.go`. Intent classification system prompt updated to include `sealed_building` as a valid label with a clear description.
- [ ] `backend/internal/retrieval/names.go` — `ExtractCardNames(ctx, client, ollamaBaseURL, model, query) ([]string, error)` uses Ollama to extract a list of card names from the user's message. Returns a deduplicated, trimmed slice. Falls back to an empty slice (not an error) if no names are found.
- [ ] `Client.RetrieveCardsByNames(ctx, names []string, query string)` added to `backend/internal/retrieval/client.go`. Sends `{"corpus": "cards", "filters": {"names": [...]}, "query": query, "top_k": 25}` to the data service. Returns `[]Result`.
- [ ] `context/assemble.go` — `AssembleSealed(query string, results []Result) string` builds a prompt instructing the LLM to suggest a 40-card Commander/sealed deck from the provided card pool, including which cards to include and a brief rationale.
- [ ] The chat handler (ticket #004) routes `sealed_building` intent through the new path: extract names → retrieve by names → assemble sealed prompt → call Ollama.
- [ ] If `ExtractCardNames` returns an empty slice, the handler falls back to `IntentGeneral` and responds that it couldn't identify any card names in the message.
- [ ] Table-driven tests cover: `ExtractCardNames` parses a well-formatted card list (mock Ollama), `RetrieveCardsByNames` sends correct payload (mock data service), `AssembleSealed` includes card names and sealed deck instructions in prompt.

## Implementation Notes
- `ExtractCardNames` system prompt: instruct the model to return only a newline-separated list of card names, nothing else. Parse by splitting on newlines and trimming whitespace. This is a short, structured extraction task — stream: false, short timeout (10s).
- `RetrieveCardsByNames` sets `top_k: 25` since sealed pools are typically 45–60 cards and we want most of them back.
- `AssembleSealed` prompt should remind the LLM that basic lands are always available (the user can add as many as needed), that the deck must be exactly 40 cards, and that it should suggest a 2-colour strategy unless the pool strongly supports 3 colours.
- The user doesn't need to format their input specially — natural language like "I opened Viscera Seer, Llanowar Elves, ..." is fine. `ExtractCardNames` handles the parsing.
- No new UI needed — the existing chat interface handles this naturally.

## Example Flow
```
User: "I just opened these cards from my pre-release packs: Sheoldred the Apocalypse,
       Atraxa Praetors Voice, Llanowar Elves, [... more cards ...]
       What deck should I build?"

Intent → sealed_building
ExtractCardNames → ["Sheoldred the Apocalypse", "Atraxa Praetors Voice", "Llanowar Elves", ...]
RetrieveCardsByNames → card details for each
AssembleSealed → prompt with full pool + instructions
Ollama → "Based on your pool, I'd build a Green/Black aristocrats shell with Sheoldred as your cornerstone..."
```

## Relevant Areas
- `backend/internal/retrieval/intent.go`
- `backend/internal/retrieval/names.go` (new)
- `backend/internal/retrieval/client.go`
- `backend/internal/context/assemble.go`
- `backend/internal/api/chat.go` (ticket #004)

## Dependencies
- Requires: `rag-data-service` ticket #008 (name filter)
- Requires: tutorai #004 (chat endpoint — handler wiring)
- Requires: tutorai #005 (frontend — no UI changes needed but must be built)
- Blocks: nothing

## Out of Scope
- Do not store the user's card pool between sessions — each query is stateless
- Do not implement draft (pick-by-pick) suggestions — sealed only
- Do not validate that card names exist in the database — unknown names simply return no results
- Do not add a structured card list input UI — plain text is sufficient for v1
