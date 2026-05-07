const API_BASE_URL = import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8000'

interface ChatResponse {
  response: string
}

interface ErrorResponse {
  error: string
}

export async function sendMessage(query: string): Promise<string> {
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
  return body.response
}
