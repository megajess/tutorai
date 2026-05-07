# Architectural Decisions

---

## 2026-04-27 — Use Llama 3.1 8B via Ollama instead of a hosted API

**Decision:** Run Llama 3.1 8B locally via Ollama for both development and initial deployment.

**Context:** Considered using the Claude or OpenAI APIs for the LLM layer. Both would give higher quality responses out of the box.

**Reasoning:** This project is explicitly a rehearsal for the metal chatbot, which will use a self-hosted open model. Using Ollama here means the architecture, prompt patterns, and integration code transfer directly. Local inference also means zero API cost during development and no rate limits while iterating on retrieval quality.

**Consequences:** Response quality is lower than a frontier model. Acceptable because retrieval does the heavy lifting — the LLM is mostly summarising and formatting grounded context, not generating from parametric knowledge. Can be swapped for a hosted model without touching the retrieval layer if quality becomes an issue.

**Status:** Accepted

---

## 2026-04-27 — Roll a custom retrieval pipeline instead of LangChain/LlamaIndex

**Decision:** Build retrieval logic directly using Chroma's Python client, SQLite, and the Ollama API. No orchestration framework.

**Context:** LangChain and LlamaIndex would accelerate initial setup and provide pre-built RAG chains.

**Reasoning:** Both frameworks add significant abstraction debt. Debugging becomes harder because failures can occur inside framework internals. More importantly, the goal is to build something portable to the metal chatbot — a custom pipeline means we understand and control every layer, making it much easier to replicate. The retrieval logic here isn't complex enough to justify a framework.

**Consequences:** More code to write upfront. Worth it for the learning value and portability.

**Status:** Accepted

---

## 2026-04-27 — Tiered retrieval with exact lookup table for color identity

**Decision:** Guild/shard/wedge names (golgari, sultai, etc.) are resolved via a JSON lookup table, not vector search.

**Context:** These terms could be embedded and retrieved semantically like other slang.

**Reasoning:** Color identity mappings are deterministic. "Golgari" always means {B, G} — there's no fuzziness. Using vector search for something with an exact answer adds latency and introduces the possibility of a wrong match. A lookup table is instant and never wrong.

**Consequences:** New color combination names need to be added to the JSON file manually. Acceptable given how rarely new ones are introduced.

**Status:** Accepted

---

## 2026-04-27 — SQLite for structured card data, Chroma for oracle text

**Decision:** Store structured card fields (legality, color identity, price, type) in SQLite and embed oracle text in Chroma as a separate collection.

**Context:** Could have put everything in Chroma using metadata filters, or everything in SQLite using FTS5.

**Reasoning:** Structured filtering (format legality, color identity, price ceiling) is naturally a SQL problem — fast, precise, and easy to compose. Semantic search over oracle text is naturally a vector problem. Using each tool for what it's good at gives better results than forcing one tool to do both jobs. Chroma metadata filters exist but are slower and less expressive than SQL for multi-field filtering.

**Consequences:** Two data stores to maintain and keep in sync during ingestion. Ingestion scripts handle this atomically.

**Status:** Accepted

---

## 2026-04-27 — No conversation history in v1

**Decision:** Each chat query is stateless. No message history is passed to the LLM.

**Context:** A conversational assistant would naturally maintain context across turns ("and add more ramp to that deck").

**Reasoning:** Stateful conversations significantly complicate the context window management and retrieval logic. For v1, the priority is getting the retrieval quality right. Conversation history can be added in v2 once the core RAG pipeline is solid.

**Consequences:** Users cannot refer back to previous turns. Each question must be self-contained.

**Status:** Accepted

---

---

## 2026-04-27 — Open-core release strategy: public app, private data service

**Decision:** Split the project into two repos. `tutorai` is open source and contains all application code. `rag-data-service` is private and contains the vector databases, corpora, and auth layer.

**Context:** The project serves dual purposes — a portfolio piece (needs to be publicly visible) and a potential monetization vehicle (curated data is the moat). Releasing everything open source gives away the data; releasing nothing undermines the portfolio value.

**Reasoning:** The application code (retrieval pipeline, intent classification, prompt assembly) is not a competitive advantage — RAG patterns are well understood. The curated and cleaned corpus (card data, rules chunking decisions, slang glossary) is the real work. Keeping the data private while open-sourcing the app gives users something to learn from and fork, while preserving the value of the data layer.

**Consequences:** The TutorAI backend communicates with the data service over HTTP rather than calling Chroma/SQLite directly. Adds a network hop in the local dev setup, but this is acceptable and actually makes the production architecture cleaner. The data service is designed to be reusable across future RAG projects.

**Status:** Accepted

---

## 2026-04-27 — Shared private data service for all RAG projects

**Decision:** The `rag-data-service` will serve as the data backend for TutorAI and all future RAG projects, with each app isolated in its own namespace.

**Context:** Future RAG projects (e.g. a metal chatbot) will have the same architectural need — a private vector store with a retrieval API. Building a separate data service per project would mean duplicating infrastructure.

**Reasoning:** A single multi-tenant data service means one deployment to maintain, one place to update embedding models or retrieval logic, and one auth system. Each app gets its own namespace so data is isolated. Adding a new project is just a new namespace config and a new API key scope — the service itself doesn't change.

**Consequences:** The data service needs to be designed with namespacing from the start. App-specific retrieval logic (like the MTG card SQLite filtering) needs to be handled carefully — either as app-specific endpoints or as generic filter parameters the service accepts. The latter is cleaner for reuse.

**Status:** Accepted

---

---

## 2026-05-06 — Ingestion scripts and slang glossary moved to private repo

**Decision:** The Python ingestion scripts (`ingest_cards.py`, `ingest_rules.py`, `ingest_slang.py`) and `slang_glossary.json` live in `rag-data-service`, not `tutorai`.

**Context:** Originally planned to include ingestion scripts in the public `tutorai` repo so users could self-host with their own corpus.

**Reasoning:** Publishing the ingestion scripts alongside the private data repo defeats the purpose of keeping the data private. The scripts are a precise recipe for recreating the curated corpus from freely accessible sources. Anyone with the scripts can reproduce the dataset without ever accessing the private repo. The ingestion pipeline is part of the proprietary layer, not the application layer.

**Consequences:** Self-hosters cannot use the hosted ingestion pipeline — they must build their own. The `tutorai` public repo contains no data artifacts or ingestion logic.

**Status:** Accepted

---

## 2026-05-07 — Tailwind CSS v4 with @tailwindcss/vite instead of PostCSS CLI config

**Decision:** Use Tailwind CSS v4 with the `@tailwindcss/vite` plugin. Add `@import "tailwindcss"` to the main CSS entry point. No `tailwind.config.js` or `postcss.config.js` created.

**Context:** The ticket specified "add Tailwind via postcss." `npm install tailwindcss` installs v4, which dropped the `tailwindcss init` CLI and the JS config file in favour of a CSS-first configuration approach.

**Reasoning:** Tailwind v4 is the current version and is the correct target going forward. The v3 PostCSS approach (`npx tailwindcss init -p`) no longer exists in v4. The Vite plugin (`@tailwindcss/vite`) is the recommended v4 integration for Vite projects — it replaces the PostCSS plugin and gives faster HMR.

**Consequences:** No `tailwind.config.js` needed for basic use. Theme customisation (if needed later) is done in CSS with `@theme` directives, not in a JS config file.

**Status:** Accepted

---

## 2026-05-05 — Go for app backend, Python for ingestion scripts

**Decision:** The TutorAI app backend is written in Go. The ingestion scripts stay in Python. The private `rag-data-service` is also Python.

**Context:** Originally planned as all-Python. Reconsidered because Go is already familiar from previous backend projects, and the app backend is fundamentally an HTTP orchestration layer — a problem Go handles idiomatically.

**Reasoning:** The app backend (receive query → classify intent → call data service → call Ollama → return response) is pure HTTP orchestration with no ML tooling dependencies. Go's `net/http` + Chi is a natural fit, and familiarity means faster iteration on what actually matters: retrieval quality and prompt design. Python's ecosystem advantage (Chroma clients, Ollama SDK, data wrangling libraries) is real but only relevant in the ingestion pipeline, which is an offline process. Keeping ingestion in Python captures that advantage where it counts. The split is clean — Go owns the runtime request path, Python owns the offline data pipeline. They share no code.

**Consequences:** The repo contains two languages. Go modules in `/backend/`, Python scripts in `/scripts/` with their own `requirements.txt`. This is a minor operational overhead but a natural separation of concerns.

**Status:** Accepted

---
