# Project Overview

## Summary
TutorAI is a RAG-powered assistant for Magic: The Gathering players. It answers natural language questions about deck building, card rules, and card lookup by retrieving grounded context from Scryfall card data, the WotC Comprehensive Rules, and a curated community terminology glossary — then generating answers with a local LLM. It serves as a rehearsal project for a larger domain-specific AI chatbot, intentionally using the same stack and architecture.

## Goals
- [ ] Users can ask deck building questions in natural language and receive relevant card suggestions filtered by format, color identity, and budget
- [ ] Users can ask rules questions and receive answers grounded in the Comprehensive Rules — not hallucinated
- [ ] Users can look up cards by name, effect, or mechanic
- [ ] The system understands MTG community terminology (guild names, archetype slang, shorthand) without hallucinating
- [ ] The architecture is clean enough to port directly to the metal chatbot project

## Non-Goals
- Not building a full deck manager or deck tracker
- Not integrating with EDHREC, Moxfield, or other external deck sites
- Not implementing user accounts or saved decks in v1
- Not optimising for production scale — local Ollama is fine for now
- Not supporting real-time card price alerts or inventory tracking

## Target Users
MTG players — primarily Commander/EDH players — who want a conversational assistant for deck building help and rules clarification. Users are hobbyists, not developers. They speak in MTG community language ("build me a golgari aristocrats deck", "how does deathtouch work with trample").

## Core User Flows

1. **Deck Building** — User describes a deck strategy, color identity, format, and optional budget. The chatbot returns relevant card suggestions with brief explanations.
2. **Rules Question** — User asks how a mechanic or card interaction works. The chatbot retrieves the relevant rule sections and explains them in plain language.
3. **Card Lookup** — User asks about a specific card or asks for cards that do a particular thing. The chatbot retrieves matching cards and summarizes their relevant text.

## Out of Scope for v1
- User accounts and authentication
- Saving or exporting decklists
- Card image display
- Price tracking or alerts
- Integration with external deck sites
- Streaming chat responses
- Mobile app
