# TASKS — Work Log

## Ticket 003 — Retrieval Layer & Data Service Client

### Status: In Review

### Plan
- [x] `backend/config/config.go` — `Config` struct + `Load()` from env vars
- [x] `backend/internal/retrieval/intent.go` — `Intent` type, constants, `ClassifyIntent()`, `parseIntent()`
- [x] `backend/internal/retrieval/lookup.go` — `ColorLookup`, `LoadColorLookup()`, `Resolve()`
- [x] `backend/internal/retrieval/client.go` — `Client`, `ServiceError`, `RetrieveCards/Rules/Slang()`
- [x] `backend/internal/context/assemble.go` — `Assemble()`, `SanitizeQuery()`
- [x] `backend/cmd/server/main.go` — wire up Config + ColorLookup at startup
- [x] `backend/internal/retrieval/retrieval_test.go` — table-driven tests for intent, lookup, client
- [x] `backend/internal/context/assemble_test.go` — tests for sanitization and prompt structure
- [x] `go test ./...` passes, `golangci-lint` clean
- [x] Check off acceptance criteria + flip ticket to `In Review`

### Log
- Ticket set to `In Progress`
- Intent classification: zero-shot prompt, defaults to `general` on unrecognised response
- Prompt injection defence: `SanitizeQuery` strips C0 control chars, caps at 500 chars; user input always placed after system prompt and retrieved context, clearly delimited
- `Config` struct passed to `NewClient` — no direct `os.Getenv` in retrieval packages
- `ServiceError` typed error (StatusCode=0 for connection failures, non-zero for HTTP errors)
- golangci-lint: fixed `errcheck` on `resp.Body.Close()` and `json.Encode` in tests; fixed staticcheck `fmt.Fprintf` vs `WriteString(fmt.Sprintf(...))`
- Ticket set to `In Review`

---

## Ticket 001 — Project Scaffold & Tooling Setup

### Status: In Review

### Plan
- [x] Install `golangci-lint` via Homebrew
- [x] Create full folder structure matching `docs/architecture.md`
- [x] `go mod init tutorai` — add Chi and godotenv dependencies
- [x] `backend/cmd/server/main.go` — Chi router, `/health` endpoint, godotenv load
- [x] `backend/cmd/server/main_test.go` — placeholder test for health handler
- [x] Stub files for all internal packages (config, api, retrieval, llm, context)
- [x] `data/color_identity_lookup.json` — empty placeholder (populated in ticket 003)
- [x] `.env.example` — all required env vars, no values
- [x] Scaffold Vue 3 + TypeScript frontend with Vite (`create-vite --template vue-ts`)
- [x] Install Tailwind CSS v4 + `@tailwindcss/vite` plugin; add `@import "tailwindcss"` to `style.css`
- [x] Verify `go build ./...`, `go test ./...`, `golangci-lint run ./...` all pass
- [x] Verify frontend builds with `npm run build`
- [x] Verify `go run ./backend/cmd/server` returns `{"status":"ok"}` from `GET /health`
- [x] Check off acceptance criteria + flip ticket to `In Review`

### Log
- Ticket set to `In Progress`
- Module name: `tutorai` (simple, not github.com prefixed — project is local-first for now)
- Tailwind v4 installed (latest `npm create vite` + `npm install -D tailwindcss`); uses `@tailwindcss/vite` plugin instead of `tailwind.config.js` + postcss CLI init (which is a v3 pattern). Decision logged.
- `data/color_identity_lookup.json` created as an empty `{}` placeholder; populated in ticket 003 when `lookup.go` is implemented
- `frontend/src/api/` directory created with `.gitkeep` — `chat.ts` added in ticket 005
- Ticket set to `In Review`

---
