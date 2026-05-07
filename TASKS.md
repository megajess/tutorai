# TASKS — Work Log

## Ticket 005 — Chat Frontend

### Status: In Review

### Plan
- [x] Flip ticket to In Progress
- [x] `frontend/src/api/chat.ts` — typed `sendMessage(query)` with fetch, network + HTTP error handling
- [x] `frontend/src/components/ChatWindow.vue` — message history, loading dots, error banner, textarea input
- [x] `frontend/src/App.vue` — replace HelloWorld scaffold with ChatWindow
- [x] `frontend/src/style.css` — strip scaffold-specific rules, keep CSS design tokens + base styles
- [x] `frontend/.env.example` — `VITE_API_BASE_URL=http://localhost:8000`
- [x] `npm run build` passes clean
- [x] Check off acceptance criteria + flip ticket to In Review

### Log
- Ticket set to In Progress
- `sendMessage()` catches network failures separately from HTTP errors — gives the user a readable message if the backend is down
- Auto-scroll uses `nextTick` + `scrollTop = scrollHeight` — fires after user message and after assistant reply / loading state resolves
- Textarea uses `field-sizing: content` for auto-grow up to 160px; Enter submits, Shift+Enter is a newline
- Loading indicator: three bouncing dots (CSS `blink` animation, staggered delays)
- Error banner displayed below message list; cleared on next successful submit
- HelloWorld.vue scaffold left in place (unused, not imported) — can be deleted later
- Ticket set to In Review

---

## Ticket 004 — Chat API Endpoint

### Status: In Review

### Plan
- [x] Install `github.com/go-chi/cors`
- [x] `backend/config/config.go` — add `HTTPTimeout time.Duration` field (120s default)
- [x] `backend/internal/llm/ollama.go` — `Generate()`, `UnavailableError` typed error
- [x] `backend/internal/api/chat.go` — `ChatHandler`, `extractCardFilters()`, `writeJSON()`
- [x] `backend/cmd/server/main.go` — wire shared `http.Client`, `retrieval.Client`, `ChatHandler`, CORS, `POST /chat`
- [x] `backend/internal/api/chat_test.go` — integration tests with mock Ollama + data service
- [x] `go test ./...` passes, `golangci-lint` clean
- [x] Check off acceptance criteria + flip ticket to `In Review`

### Log
- Ticket set to `In Progress`
- CORS: used `github.com/go-chi/cors` (approved by user) — allows `http://localhost:5173`, `POST` + `OPTIONS`, `Content-Type` header
- Empty results flow through to LLM (LLM-in-the-loop approach) — LLM explains no matches found rather than a canned error response
- Filter extraction is best-effort: color identity via `ColorLookup`, format via string match against known names, price via regex `$X`
- Intent classification failure (Ollama unreachable during classification) falls back to `general` silently — intent errors never surface to the user
- `HTTPTimeout`: 120s — generous for local LLM inference; single shared `http.Client` for all outbound calls
- Decision logged for LLM-in-the-loop and filter extraction approach
- Ticket set to `In Review`

---

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
