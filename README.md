# TutorAI

A RAG-powered Magic: The Gathering assistant for deck building, rules lookup, and card search.

> Part of the [`megajess/tutorai-workspace`](https://github.com/megajess/tutorai-workspace) — see there for full system architecture and context.

## Features

- **Deck Building** — Get card recommendations by color identity, strategy, format, and budget
- **Rules Lookup** — Ask natural language questions grounded in the official Comprehensive Rules
- **Card Search** — Look up cards by name, effect, type, or mechanic
- **MTG Domain Knowledge** — Understands guild names, archetypes, slang, and community terminology

## How It Works

TutorAI uses a tiered retrieval pipeline: your query is classified by intent, routed to the right retrieval strategy, and grounded context is injected into a local LLM (Llama 3.1 via Ollama) before generating a response.

Retrieval is handled by a separate [rag-data-service](https://github.com/megajess/rag-data-service) — a private backend that hosts the vector databases. This repo contains all the application code; the data layer is kept separate.

## Tech Stack

| Layer | Tool |
|---|---|
| LLM | Llama 3.1 8B (via Ollama) |
| Embeddings | `nomic-embed-text` via Ollama (owned by rag-data-service) |
| Data Backend | rag-data-service (REST API) |
| App Backend | Go 1.22+ with Chi router |
| Ingestion Scripts | Python 3.11+ |
| Frontend | Vue 3 (TypeScript) |

## Self-Hosting

You can run TutorAI against your own data service instance. Set `DATA_SERVICE_URL` and `DATA_SERVICE_API_KEY` in your `.env`, run the ingestion scripts to build your own corpus, and you're fully independent of the hosted service.

See [SETUP.md](./SETUP.md) for full instructions.

## Project Structure

```
tutorai/                     (this repo — public)
├── CLAUDE.md
├── README.md
├── SETUP.md
├── .env.example
├── go.mod
├── go.sum
├── docs/
│   ├── project-overview.md
│   ├── tech-stack.md
│   ├── architecture.md      # Includes data service design
│   ├── decisions.md
│   └── data.md
├── tickets/
├── backend/                 # Go app backend
│   ├── cmd/server/          # Entry point (main.go)
│   ├── internal/
│   │   ├── api/             # POST /chat handler
│   │   ├── retrieval/       # Intent classification, color lookup, data service client
│   │   ├── llm/             # Ollama HTTP client
│   │   └── context/         # Prompt / context assembly
│   └── config/              # Env var loading
├── scripts/                 # Python ingestion scripts (push raw text to data service)
├── data/
│   ├── color_identity_lookup.json
│   └── slang_glossary.json
└── frontend/                # Vue 3 + TypeScript + Tailwind
```
