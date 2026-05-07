# Setup

## Prerequisites

- [Go](https://go.dev/dl/) 1.22+
- [Ollama](https://ollama.ai) installed and running
- Node.js 18+ (for frontend)
- `rag-data-service` running locally (see its own SETUP.md)

## 1. Pull the Ollama LLM Model

```bash
ollama pull llama3.1
```

`llama3.1` is used by the Go backend for LLM responses.

> **Note:** The embedding model (`nomic-embed-text`) is owned by `rag-data-service`, not by TutorAI. Pull it as part of the data service setup — TutorAI itself never embeds anything.

## 2. Start the Data Service

TutorAI depends on `rag-data-service` for all retrieval. Make sure it's running before starting the Go backend:

```bash
# In the rag-data-service repo
uvicorn main:app --reload --port 8001
```

See the [rag-data-service SETUP.md](https://github.com/megajess/rag-data-service) for full setup instructions, including how to run the ingestion scripts that populate the corpus.

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

## 5. Run the Go Backend

```bash
go run ./backend/cmd/server
```

API will be available at `http://localhost:8080`.

To verify it's running:

```bash
curl http://localhost:8080/health
# {"status": "ok"}
```

## 6. Run the Frontend

```bash
cd frontend
npm install
npm run dev
```

Frontend will be available at `http://localhost:5173`.

## Running Tests

```bash
# Go backend tests
go test ./...
```

## Troubleshooting

**Data service not reachable:** Make sure `rag-data-service` is running on port 8001 and `DATA_SERVICE_URL` in `.env` matches.

**401 from data service:** The `DATA_SERVICE_API_KEY` in TutorAI's `.env` doesn't match any key in the data service's `API_KEYS` env var.

**403 from data service:** The API key is valid but not scoped to the `tutorai` namespace. Check the `key:namespace` pairing in the data service's `API_KEYS`.

**Ollama not responding:** Make sure Ollama is running (`ollama serve` or the desktop app) and `llama3.1` is pulled.
