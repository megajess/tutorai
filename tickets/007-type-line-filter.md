# Ticket 007 — Creature Type Filter in extractCardFilters

## Status
`Todo`

## Type
`Feature`

## Summary
Extract creature type keywords from user queries and pass them as a `type_line` filter to the data service. Without this, queries like "show me sliver cards" find cascade/commander synergy cards instead of actual Sliver creatures because the embedding signal from "cascade" and "commander" outweighs "sliver" in the vector search.

## Acceptance Criteria
- [ ] `extractCardFilters()` detects creature type keywords in the query and sets a `TypeLine` field on `CardFilters`
- [ ] `CardFilters` struct has a `TypeLine string` field
- [ ] `RetrieveCards()` includes `type_line` in the request body filters when set
- [ ] At minimum, common subtype keywords are detected: sliver, dragon, elf, goblin, zombie, human, merfolk, angel, demon, vampire, warrior, wizard, soldier

## Implementation Notes
- Detection is case-insensitive word match, same pattern as format detection
- Only one creature type per query (first match wins)
- Requires rag-data-service #009 to be deployed first

## Dependencies
- Requires: rag-data-service #009
- Requires: #004

## Relevant Areas
- `backend/internal/api/chat.go` (`extractCardFilters`)
- `backend/internal/retrieval/client.go` (`CardFilters`, `RetrieveCards`)
