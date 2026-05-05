# Tech Stack

## Languages & Runtimes

### Go (App Backend)
- **Version:** Go 1.22+
- **Why:** Already familiar with Go from previous projects, which means faster iteration on the parts that matter — retrieval quality and prompt design. The app backend is fundamentally an HTTP orchestration layer (receive query → call data service → call Ollama → return response), which is squarely in Go's wheelhouse. Performance is a bonus; productivity is the real reason.

### Python (Ingestion Scripts)
- **Version:** Python 3.11+
- **Why:** Data wrangling — chunking 250KB of plain text rules, processing 270MB of Scryfall JSON, iterating over embeddings — is more ergonomic in Python. The Ollama and Chroma client libraries are also first-class in Python with better documentation and more community examples. Ingestion is an offline, one-time-ish process, so language familiarity matters less than ecosystem fit.

The split is clean: Go owns the request path (runtime), Python owns the data pipeline (offline). They never need to share code.

## Frontend
- **Framework:** Vue 3 (Composition API) with TypeScript
- **Styling:** Tailwind CSS
- **State Management:** Pinia
- **Why:** Vue 3 is already being learned for other projects, so this reinforces that investment.

## App Backend (Go — `tutorai` public repo)
- **Framework:** `net/http` + [Chi](https://github.com/go-chi/chi) router
- **HTTP Client:** `net/http` (standard library) for data service and Ollama calls
- **Config:** `godotenv` for `.env` loading
- **Testing:** `testing` package + `httptest`
- **Why Chi:** Lightweight, idiomatic Go, no magic. Chi is the router pattern you'd recognize from previous Go backend work.

## RAG Data Service (Python — `rag-data-service` private repo)
- **Framework:** FastAPI
- **Why:** Python-native, async, generates OpenAPI docs automatically. The data service is the right place for Python since it sits closest to the Chroma and SQLite tooling.

## Ingestion Scripts (Python — `tutorai` public repo, `/scripts/`)
- **Why Python:** Scryfall JSON processing and rules text chunking are more ergonomic in Python. Scripts are offline tools, not part of the request path.
- **What they do:** Fetch source data, chunk/normalise it, POST raw text + metadata to the data service. They never embed or talk to Chroma/SQLite directly — the data service owns those concerns.
- **Key libraries:** `httpx` (POST to data service), `python-dotenv`, `tqdm` (progress on long ingests)

## LLM
- **Model:** Llama 3.1 8B (local via Ollama)
- **Integration:** Direct HTTP calls to `http://localhost:11434/api/chat` from the Go backend
- **Why:** No Go-specific Ollama SDK needed — Ollama exposes a simple REST API. Go's `net/http` handles it cleanly.
- **Production path:** Swap the Ollama base URL for a hosted endpoint without changing any other code.

## Embeddings
- **Model:** `nomic-embed-text` (via Ollama)
- **Owned by:** `rag-data-service` — embedding happens inside the data service at both ingest time (for documents) and retrieval time (for queries). Centralising it there guarantees the same model is used end-to-end, which is a hard requirement for meaningful similarity scores.
- **Not used by:** TutorAI's Go backend or the Python ingestion scripts. Both send raw text only.

## Vector Store
- **Tool:** Chroma (development) / pgvector (production path)
- **Lives in:** `rag-data-service` (private, Python)
- **Accessed by:** Data service only — TutorAI Go backend never touches Chroma directly

## Structured Card Storage
- **Tool:** SQLite with FTS5
- **Lives in:** `rag-data-service` (private, Python)
- **Accessed by:** Data service only

## Inter-Service Communication
- **Protocol:** HTTP REST (JSON)
- **Auth:** `X-API-Key` header
- **Why:** Simple, easy to test with curl, no extra dependencies in Go. JWT with scoped claims is the v2 upgrade path.

## Data Sources
| Source | Purpose | Notes |
|--------|---------|-------|
| Scryfall Bulk Data API | Card data (oracle text, legalities, prices) | Free, updated daily, ~270MB JSON |
| WotC Comprehensive Rules | Rules text for RAG | Plain text, chunked by rule number |
| Hand-curated slang glossary | MTG community terminology | JSON file in `/data/` |

## Development Tools

### Go
- **Module management:** Go modules (`go.mod`)
- **Formatting:** `gofmt` (built-in)
- **Linting:** `golangci-lint`
- **Testing:** `go test ./...`

### Python (scripts only)
- **Package management:** pip + `requirements.txt`
- **Formatting:** Black
- **Linting:** Ruff

### Frontend
- **Tooling:** Vite + npm

## What We're Deliberately NOT Using
- **LangChain / LlamaIndex** — abstraction debt, harder to port to the metal chatbot. Custom pipeline in Go means we understand every layer.
- **A Go Chroma SDK** — TutorAI never talks to Chroma directly; the data service owns that. No SDK needed.
- **OpenAI / Anthropic API** — local Ollama keeps costs zero during dev.
- **gRPC** — HTTP/JSON is sufficient for two local services talking to each other.
- **Docker (for now)** — adds friction during early development. Worth adding before deployment.
- **JWT for v1 auth** — API key is sufficient for a single-client setup.
