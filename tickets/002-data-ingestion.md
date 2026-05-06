# Ticket 002 — Data Ingestion

## Status
`Moved`

## Type
`Feature`

## Summary
~~Moved to `rag-data-service` ticket 007.~~ The ingestion scripts and slang glossary are proprietary assets that live in the private repo. See `rag-data-service/tickets/007-data-ingestion.md`.

## Background / Context
See `docs/data.md` for the full breakdown of each data source, field mappings, and chunking strategy. The ingestion layer is the foundation of the RAG system — if the data is poorly structured here, retrieval quality suffers downstream. The data service must be running before any script is invoked.

## Acceptance Criteria
- [ ] `scripts/ingest_cards.py` fetches the Scryfall Oracle Cards bulk data file and POSTs records to `POST {DATA_SERVICE_URL}/ingest/tutorai` with `corpus: "cards"`. Each record's `metadata` includes `name`, `color_identity` (array), `legalities` (object — e.g. `{"commander": "legal", ...}`), and `price_usd` (nullable float).
- [ ] `scripts/ingest_rules.py` downloads the WotC Comprehensive Rules plain text file, chunks by top-level rule number (keeping sub-rules together), and POSTs to the data service with `corpus: "rules"`. Each record's `metadata` includes the rule range (e.g. `"rule_range": "702.2–702.2b"`) and `section` name.
- [ ] `scripts/ingest_slang.py` reads `data/slang_glossary.json` and POSTs to the data service with `corpus: "slang"`.
- [ ] `data/color_identity_lookup.json` exists and covers all 10 guilds, 5 shards, and 5 wedges.
- [ ] `data/slang_glossary.json` exists with the schema defined in `docs/data.md` (terms, archetypes, shorthand).
- [ ] All three scripts support a `--refresh` flag that first sends `{"operation": "delete_all", "corpus": "<name>"}` to the data service, then upserts.
- [ ] Each script reads `DATA_SERVICE_URL` and `DATA_SERVICE_API_KEY` from `.env` and sends `X-API-Key` on every request.
- [ ] Records are POSTed in batches (suggest 500/batch) with progress shown via `tqdm`.
- [ ] Running all three scripts from scratch against a running data service completes without errors.
- [ ] A pytest test stubs the data service with `httpx.MockTransport` and confirms each script POSTs the expected payload shape.

## Implementation Notes
- Scryfall bulk data URL: fetch from `https://api.scryfall.com/bulk-data` first, then download the `oracle_cards` download_uri — don't hardcode the download URL.
- Set a `User-Agent` header on Scryfall requests to avoid throttling.
- Rules chunking: group by top-level rule number (e.g. all of 702.x together). A rule section + all its lettered sub-rules should be one chunk.
- Each ingest record sent to the data service has the shape `{ "id": str, "text": str, "metadata": {...} }`. The data service embeds `text` and writes structured fields from `metadata` to SQLite for the `cards` corpus.
- Use a long HTTP timeout on the POST (the data service embeds synchronously and large card batches can take minutes).
- No direct Chroma or SQLite calls from the scripts — they're a network client only.

## Relevant Areas
- `docs/data.md`
- `docs/tech-stack.md`
- `data/slang_glossary.json` (create this file as part of the ticket)
- `data/color_identity_lookup.json` (create this file as part of the ticket)

## Dependencies
- Requires: #001
- Blocks: #003

## Out of Scope
- Do not build any retrieval logic here — just ingestion
- Do not add any API endpoints
- Do not implement automatic refresh scheduling — manual `--refresh` flag only
