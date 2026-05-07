# Ticket 005 — Chat Frontend

## Status
`Done`

## Type
`Feature`

## Summary
Build the Vue 3 chat UI that lets users type queries and see responses. This is intentionally minimal — a working chat interface, not a polished product. The priority is a clean, functional UI that makes the API usable without Swagger.

## Acceptance Criteria
- [x] User can type a query into an input field and submit it
- [x] Submitted query appears in the message history as a user message
- [x] API response appears in the message history as an assistant message
- [x] A loading indicator is shown while the API request is in flight
- [x] Error state is shown if the API returns an error or is unreachable
- [x] Message history scrolls correctly as new messages are added
- [x] Input field is cleared after submission
- [x] The UI is usable on a standard desktop browser width

## Implementation Notes
- Single component is fine: `ChatWindow.vue` handles state, input, and message display
- `src/api/chat.ts` should export a typed `sendMessage(query: string): Promise<string>` function — keeps API logic out of the component
- Use Pinia for chat message state if it grows complex, but a simple `ref([])` in the component is acceptable for v1
- Tailwind for styling — keep it clean but don't over-design
- No markdown rendering needed — display responses as plain text for now
- The API base URL should come from a Vite env variable (`VITE_API_BASE_URL`) defaulting to `http://localhost:8000`

## Relevant Areas
- `frontend/src/`
- `backend/internal/api/chat.go` (for request/response shape reference)
- `docs/tech-stack.md`

## Dependencies
- Requires: #004
- Blocks: nothing — this is the last ticket for v1

## Out of Scope
- Do not implement card image display
- Do not implement markdown rendering in responses
- Do not add user accounts or saved conversations
- Do not build a mobile layout
- Do not add syntax highlighting or special formatting for card names
