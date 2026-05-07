# Ticket 001 — Project Scaffold & Tooling Setup

## Status
`Done`

## Type
`Chore`

## Summary
Set up the initial project structure, install dependencies, and configure dev tooling for the Go backend and Vue frontend. Creates the folder layout defined in `architecture.md` and ensures both the Go server and Vue frontend start without errors. No business logic — just a clean foundation for every subsequent ticket.

> **Note:** The Python ingestion scripts (`ingest_cards.py`, `ingest_rules.py`, `ingest_slang.py`) and their dependencies live in the private `rag-data-service` repo (ticket 007), not here.

## Acceptance Criteria
- [x] Folder structure matches `docs/architecture.md` exactly
- [x] Go module initialised (`go mod init`) with Chi router added as a dependency
- [x] Go backend starts with `go run ./backend/cmd/server` and returns `{"status": "ok"}` from `GET /health`
- [x] `golangci-lint` runs without errors on the empty codebase
- [x] `go test ./...` passes with one placeholder test
- [x] Frontend scaffolded with Vite + Vue 3 + TypeScript + Tailwind — runs with `npm run dev`
- [x] `.env.example` lists all required environment variables (no values): `DATA_SERVICE_URL`, `DATA_SERVICE_API_KEY`, `OLLAMA_BASE_URL`, `OLLAMA_LLM_MODEL`
- [x] `CLAUDE.md` is present at project root

## Implementation Notes
- Follow folder structure in `docs/architecture.md` exactly — use `backend/cmd/server/main.go` as the entry point
- Go backend: `github.com/go-chi/chi/v5` for routing, `github.com/joho/godotenv` for env loading
- Vue 3 scaffold: `npm create vite@latest frontend -- --template vue-ts`, then add Tailwind via postcss

## Relevant Areas
- `docs/architecture.md`
- `docs/tech-stack.md`
- `CLAUDE.md`

## Dependencies
- Requires: nothing — this is first
- Blocks: all other tickets

## Out of Scope
- No retrieval logic
- No Ollama integration
- No data service calls
- No chat endpoint
