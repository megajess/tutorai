const API_BASE_URL = import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8000'

export interface TokenUsage {
  prompt_tokens: number
  completion_tokens: number
  total_tokens: number
}

interface ChatResponse {
  response: string
  usage?: TokenUsage
}

interface ErrorResponse {
  error: string
}

export interface SendResult {
  text: string
  usage?: TokenUsage
}

export async function sendMessage(query: string): Promise<SendResult> {
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

  if (!res.ok) {
    const body = (await res.json().catch(() => ({}))) as Partial<ErrorResponse>
    throw new Error(body.error ?? `Server error: ${res.status}`)
  }

  const body = (await res.json()) as ChatResponse
  return { text: body.response, usage: body.usage }
}
