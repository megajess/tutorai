# Ticket 002 — Data Ingestion

## Status
`Done`

## Type
`Feature`

## Summary
The ingestion scripts and slang glossary moved to `rag-data-service` ticket 007 (private repo). The one tutorai-side deliverable — `data/color_identity_lookup.json` — has been created and populated as part of the ticket 001 scaffold.

## What Moved to `rag-data-service`
- `scripts/ingest_cards.py` — rag-data-service ticket 007 ✓
- `scripts/ingest_rules.py` — rag-data-service ticket 007 ✓
- `scripts/ingest_slang.py` — rag-data-service ticket 007 ✓
- `data/slang_glossary.json` — rag-data-service ticket 007 ✓

## What Stayed in `tutorai`
- [x] `data/color_identity_lookup.json` — flat map of guild/shard/wedge names and aliases to color identity arrays. Covers all 10 guilds, 5 shards, and 5 wedges. Used by `backend/internal/retrieval/lookup.go` (ticket 003).

## Dependencies
- Requires: #001
- Blocks: #003
