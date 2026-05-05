# TutorAI — Claude Instructions

## What This Project Is
A RAG-powered Magic: The Gathering assistant that helps players with deck building, rules lookups, and card search. The app backend is written in Go; ingestion scripts are in Python. Retrieval is delegated to a separate private `rag-data-service`.

## How to Navigate This Repo
- `/docs/` — project documentation. Read these before doing anything.
  - `project-overview.md` — goals, scope, users
  - `tech-stack.md` — chosen technologies and why
  - `architecture.md` — system structure, components, folder conventions
  - `decisions.md` — settled architectural decisions. Do not re-litigate these.
- `/tickets/` — units of work. Each file is one ticket.
  - Tickets are numbered (`001`, `002`, ...) and should be worked in order unless dependencies say otherwise.
  - A ticket is "done" when all acceptance criteria are met and the code is in a passing state.

## How to Work a Ticket
When asked to work on a ticket (e.g. "work on ticket 003"):
1. Read the ticket file fully before writing any code.
2. Read any files listed in the ticket's "Relevant Areas" section.
3. Implement the work described, following the conventions in `tech-stack.md` and `architecture.md`.
4. Confirm each acceptance criterion is met before finishing.
5. Do not start the next ticket unless explicitly asked.

## Coding Conventions

### Go (app backend — `/backend/`)
- **Version:** Go 1.22+
- **Formatting:** `gofmt` — always run before committing
- **Linting:** `golangci-lint`
- **Naming:** idiomatic Go — `camelCase` for unexported, `PascalCase` for exported
- **Error handling:** always handle errors explicitly, no blank identifier discards on errors
- **Tests:** `go test ./...` — write table-driven tests for handlers and retrieval logic
- **Commits:** conventional commits — `feat:`, `fix:`, `chore:`, `docs:`

### Python (ingestion scripts — `/scripts/`)
- **Version:** Python 3.11+
- **Formatting:** Black (`black .`)
- **Linting:** Ruff (`ruff check .`)
- **Naming:** `snake_case` for variables and functions, `PascalCase` for classes
- **Type hints:** required on all function signatures
- **Tests:** pytest for any non-trivial logic

### Frontend (`/frontend/`)
- **Language:** TypeScript (strict mode)
- **Formatting:** Prettier
- **Naming:** `camelCase` for variables/functions, `PascalCase` for components

## Environment
- Run Go backend: `go run ./backend/cmd/server`
- Run ingestion scripts: `cd scripts && pip install -r requirements.txt && python ingest_cards.py`
- Run frontend: `cd frontend && npm install && npm run dev`
- Run Go tests: `go test ./...`
- Run Python tests: `pytest`
- Key environment variables (see `.env.example`):
  - `DATA_SERVICE_URL`
  - `DATA_SERVICE_API_KEY`
  - `OLLAMA_BASE_URL`
  - `OLLAMA_LLM_MODEL`

## Off-Limits
- Do not talk to Chroma or SQLite directly from Go — all data retrieval goes through the data service client.
- Do not swap out the LLM or embedding model without being asked.
- Do not install new Go modules or Python packages without confirming first.
- Do not implement authentication on the TutorAI backend — out of scope for v1.
- Do not add streaming responses — plain request/response only for now.
