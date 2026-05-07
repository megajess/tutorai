const API_BASE_URL = import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8000'

export interface TokenUsage {
  prompt_tokens: number
  completion_tokens: number
  total_tokens: number
}

interface ErrorResponse {
  error: string
}

/**
 * Streams a chat response from the backend, calling onChunk for each partial
 * text token as it arrives. Returns token usage from the final SSE event.
 */
export async function sendMessage(
  query: string,
  onChunk: (text: string) => void,
): Promise<TokenUsage | undefined> {
  let res: Response
  try {
    res = await fetch(`${API_BASE_URL}/chat`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ query }),
    })
  } catch {
    throw new Error('Could not reach the server. Is the backend running?')
  }

  // Pre-stream errors (e.g. 400, 503) are plain JSON.
  if (!res.ok) {
    const body = (await res.json().catch(() => ({}))) as Partial<ErrorResponse>
    throw new Error(body.error ?? `Server error: ${res.status}`)
  }

  // Parse SSE stream.
  const reader = res.body!.getReader()
  const decoder = new TextDecoder()
  let buffer = ''
  let usage: TokenUsage | undefined

  while (true) {
    const { done, value } = await reader.read()
    if (done) break

    buffer += decoder.decode(value, { stream: true })

    // SSE events are separated by double newlines.
    const events = buffer.split('\n\n')
    buffer = events.pop() ?? ''

    for (const event of events) {
      const line = event.trim()
      if (!line.startsWith('data: ')) continue

      const payload = JSON.parse(line.slice(6)) as Record<string, unknown>

      if (typeof payload.error === 'string') {
        throw new Error(payload.error)
      }
      if (typeof payload.chunk === 'string') {
        onChunk(payload.chunk)
      }
      if (payload.done === true) {
        usage = payload.usage as TokenUsage | undefined
      }
    }
  }

  return usage
}
