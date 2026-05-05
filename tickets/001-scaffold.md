# Ticket 001 — Project Scaffold & Tooling Setup

## Status
`Todo`

## Type
`Chore`

## Summary
Set up the initial project structure, install dependencies, and configure dev tooling for both the Go backend and Python ingestion scripts. Creates the folder layout defined in `architecture.md` and ensures both the Go server and Vue frontend start without errors. No business logic — just a clean foundation for every subsequent ticket.

## Acceptance Criteria
- [ ] Folder structure matches `docs/architecture.md` exactly
- [ ] Go module initialised (`go mod init`) with Chi router added as a dependency
- [ ] Go backend starts with `go run ./backend/cmd/server` and returns `{"status": "ok"}` from `GET /health`
- [ ] `golangci-lint` runs without errors on the empty codebase
- [ ] `go test ./...` passes with one placeholder test
- [ ] Python `scripts/requirements.txt` includes: `httpx`, `python-dotenv`, `tqdm`, `black`, `ruff`, `pytest` (no `chromadb` — scripts only POST raw text to the data service, they never embed or write to Chroma directly)
- [ ] `black .` and `ruff check .` pass on empty scripts directory
- [ ] Frontend scaffolded with Vite + Vue 3 + TypeScript + Tailwind — runs with `npm run dev`
- [ ] `.env.example` lists all required environment variables (no values): `DATA_SERVICE_URL`, `DATA_SERVICE_API_KEY`, `OLLAMA_BASE_URL`, `OLLAMA_LLM_MODEL`
- [ ] `CLAUDE.md` is present at project root

## Implementation Notes
- Follow folder structure in `docs/architecture.md` exactly — use `backend/cmd/server/main.go` as the entry point
- Go backend: `github.com/go-chi/chi/v5` for routing, `github.com/joho/godotenv` for env loading
- Vue 3 scaffold: `npm create vite@latest frontend -- --template vue-ts`, then add Tailwind via postcss
- Python scripts are standalone — no shared module, just a flat `/scripts/` directory with `requirements.txt`

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
