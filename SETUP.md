# Setup

## Prerequisites

- [Go](https://go.dev/dl/) 1.22+
- Python 3.11+ (for ingestion scripts)
- [Ollama](https://ollama.ai) installed and running
- Node.js 18+ (for frontend)
- `rag-data-service` running locally (see its own SETUP.md)

## 1. Pull the Ollama LLM Model

```bash
ollama pull llama3.1
```

`llama3.1` is used by the Go backend for LLM responses.

> **Note:** The embedding model (`nomic-embed-text`) is owned by `rag-data-service`, not by TutorAI. Pull it as part of the data service setup — TutorAI itself never embeds anything. The Python ingestion scripts send raw text to the data service, which handles embedding at ingest and at query time.

## 2. Start the Data Service

TutorAI depends on `rag-data-service` for all retrieval. Make sure it's running before starting the Go backend:

```bash
# In the rag-data-service repo
uvicorn main:app --reload --port 8001
```

See the [rag-data-service SETUP.md](https://github.com/megajess/rag-data-service) for full setup instructions.

## 3. Clone and Install TutorAI

```bash
git clone https://github.com/megajess/tutorai
cd tutorai
go mod download
```

## 4. Configure Environment

```bash
cp .env.example .env
```

Edit `.env`:

```env
DATA_SERVICE_URL=http://localhost:8001
DATA_SERVICE_API_KEY=your-api-key-here
OLLAMA_BASE_URL=http://localhost:11434
OLLAMA_LLM_MODEL=llama3.1
```

The `DATA_SERVICE_API_KEY` must match a key scoped to the `tutorai` namespace in the data service's `API_KEYS` env var.

## 5. Ingest Data

Ingestion scripts push corpus data into the data service. Run these from the `/scripts` directory:

```bash
cd scripts
pip install -r requirements.txt
```

Make sure the data service is running before ingesting.

### Card Data (Scryfall)

```bash
python ingest_cards.py
```

This will:
- Fetch the latest bulk data URL from `https://api.scryfall.com/bulk-data`
- Download the Oracle Cards JSON (~270MB)
- Push card records to the data service (`POST /ingest/tutorai` with corpus `cards`)

**Expected time:** 15–30 minutes on first run.

### Comprehensive Rules

```bash
python ingest_rules.py
```

Downloads the plain text rules from Wizards, chunks by rule number, and pushes to the data service.

**Expected time:** 5–10 minutes.

### Slang Glossary

```bash
python ingest_slang.py
```

Pushes the curated slang and archetype glossary from `data/slang_glossary.json` to the data service.

**Expected time:** Under 1 minute.

## 6. Run the Go Backend

```bash
go run ./backend/cmd/server
```

API will be available at `http://localhost:8000`.

To verify it's running:

```bash
curl http://localhost:8000/health
# {"status": "ok"}
```

## 7. Run the Frontend

```bash
cd frontend
npm install
npm run dev
```

Frontend will be available at `http://localhost:5173`.

## Re-ingesting Data

Scryfall publishes new bulk data daily. To refresh the card corpus:

```bash
cd scripts
python ingest_cards.py --refresh
```

The `--refresh` flag sends a `delete_all` operation to the data service before re-ingesting, ensuring stale records are cleared.

## Running Tests

```bash
# Go backend tests
go test ./...

# Python ingestion script tests
cd scripts && pytest
```

## Troubleshooting

**Data service not reachable:** Make sure `rag-data-service` is running on port 8001 and `DATA_SERVICE_URL` in `.env` matches.

**401 from data service:** The `DATA_SERVICE_API_KEY` in TutorAI's `.env` doesn't match any key in the data service's `API_KEYS` env var.

**403 from data service:** The API key is valid but not scoped to the `tutorai` namespace. Check the `key:namespace` pairing in the data service's `API_KEYS`.

**Ollama not responding:** Make sure Ollama is running (`ollama serve` or the desktop app) and `llama3.1` is pulled.

**Scryfall download fails:** Scryfall rate-limits unauthenticated requests. Wait a moment and retry, or add a `User-Agent` header in `ingest_cards.py`.
