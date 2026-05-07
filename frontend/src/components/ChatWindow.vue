<script setup lang="ts">
import { ref, nextTick } from 'vue'
import { sendMessage, type TokenUsage } from '../api/chat'

const DEV = import.meta.env.DEV

interface Message {
  role: 'user' | 'assistant'
  text: string
  usage?: TokenUsage
}

const messages = ref<Message[]>([])
const input = ref('')
const loading = ref(false)
const error = ref<string | null>(null)
const messagesEl = ref<HTMLElement | null>(null)

async function scrollToBottom() {
  await nextTick()
  if (messagesEl.value) {
    messagesEl.value.scrollTop = messagesEl.value.scrollHeight
  }
}

function scrollSync() {
  if (messagesEl.value) {
    messagesEl.value.scrollTop = messagesEl.value.scrollHeight
  }
}

async function submit() {
  const query = input.value.trim()
  if (!query || loading.value) return

  error.value = null
  input.value = ''
  messages.value.push({ role: 'user', text: query })
  await scrollToBottom()

  // Add empty assistant message — text fills in as chunks arrive.
  loading.value = true
  messages.value.push({ role: 'assistant', text: '' })
  const assistantIdx = messages.value.length - 1
  let firstChunk = true

  try {
    const usage = await sendMessage(query, (chunk) => {
      if (firstChunk) {
        loading.value = false
        firstChunk = false
      }
      messages.value[assistantIdx].text += chunk
      scrollSync()
    })
    messages.value[assistantIdx].usage = usage
  } catch (err) {
    // Remove the empty assistant placeholder on error.
    messages.value.splice(assistantIdx, 1)
    error.value = err instanceof Error ? err.message : 'An unexpected error occurred.'
  } finally {
    loading.value = false
    await scrollToBottom()
  }
}

function handleKeydown(e: KeyboardEvent) {
  if (e.key === 'Enter' && !e.shiftKey) {
    e.preventDefault()
    submit()
  }
}
</script>

<template>
  <div class="chat-shell">
    <header class="chat-header">
      <span class="chat-header__logo">⚡</span>
      <h1 class="chat-header__title">TutorAI</h1>
      <p class="chat-header__sub">Magic: The Gathering assistant</p>
    </header>

    <div ref="messagesEl" class="chat-messages">
      <div v-if="messages.length === 0" class="chat-empty">
        Ask me anything about MTG — rules, card lookup, deck ideas.
      </div>

      <div
        v-for="(msg, i) in messages"
        :key="i"
        class="chat-message"
        :class="msg.role === 'user' ? 'chat-message--user' : 'chat-message--assistant'"
      >
        <div class="chat-message__body">
          <div v-if="msg.text" class="chat-bubble">{{ msg.text }}</div>
          <div
            v-if="DEV && msg.role === 'assistant' && msg.usage"
            class="chat-token-badge"
            :title="`Prompt: ${msg.usage.prompt_tokens} · Completion: ${msg.usage.completion_tokens}`"
          >
            {{ msg.usage.total_tokens }} tokens
          </div>
        </div>
      </div>

      <div v-if="loading" class="chat-message chat-message--assistant">
        <div class="chat-bubble chat-bubble--loading">
          <div class="mtg-spinner">
            <span class="mtg-dot mtg-dot--w" />
            <span class="mtg-dot mtg-dot--u" />
            <span class="mtg-dot mtg-dot--b" />
            <span class="mtg-dot mtg-dot--r" />
            <span class="mtg-dot mtg-dot--g" />
          </div>
        </div>
      </div>
    </div>

    <div v-if="error" class="chat-error">
      {{ error }}
    </div>

    <form class="chat-form" @submit.prevent="submit">
      <textarea
        v-model="input"
        class="chat-input"
        placeholder="Ask a question…"
        rows="1"
        :disabled="loading"
        @keydown="handleKeydown"
      />
      <button type="submit" class="chat-send" :disabled="loading || !input.trim()">
        Send
      </button>
    </form>
  </div>
</template>

<style scoped>
.chat-shell {
  display: flex;
  flex-direction: column;
  height: 100svh;
  max-width: 760px;
  margin: 0 auto;
  width: 100%;
}

/* ── Header ── */
.chat-header {
  padding: 20px 24px 16px;
  border-bottom: 1px solid var(--border);
  display: flex;
  align-items: baseline;
  gap: 10px;
  flex-shrink: 0;
}
.chat-header__logo {
  font-size: 20px;
}
.chat-header__title {
  font-size: 20px;
  font-weight: 600;
  color: var(--text-h);
  margin: 0;
  letter-spacing: -0.3px;
}
.chat-header__sub {
  font-size: 13px;
  color: var(--text);
  margin: 0;
}

/* ── Message list ── */
.chat-messages {
  flex: 1;
  overflow-y: auto;
  padding: 24px 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  scroll-behavior: smooth;
}

.chat-empty {
  text-align: center;
  color: var(--text);
  font-size: 15px;
  margin: auto;
  max-width: 340px;
  line-height: 1.6;
}

/* ── Message rows ── */
.chat-message {
  display: flex;
}
.chat-message--user {
  justify-content: flex-end;
}
.chat-message--assistant {
  justify-content: flex-start;
}
.chat-message__body {
  display: flex;
  flex-direction: column;
  gap: 4px;
  max-width: 72%;
}
.chat-message--user .chat-message__body {
  align-items: flex-end;
}
.chat-message--assistant .chat-message__body {
  align-items: flex-start;
}

/* ── Token badge (dev only) ── */
.chat-token-badge {
  font-size: 11px;
  color: var(--text);
  opacity: 0.6;
  font-family: var(--mono);
  padding: 1px 6px;
  border-radius: 4px;
  background: var(--code-bg);
  cursor: default;
}

/* ── Bubbles ── */
.chat-bubble {
  padding: 10px 14px;
  border-radius: 16px;
  font-size: 15px;
  line-height: 1.55;
  white-space: pre-wrap;
  word-break: break-word;
}
.chat-message--user .chat-bubble {
  background: var(--accent);
  color: #fff;
  border-bottom-right-radius: 4px;
}
.chat-message--assistant .chat-bubble {
  background: var(--code-bg);
  color: var(--text-h);
  border-bottom-left-radius: 4px;
}

/* ── MTG spinner ── */
.chat-bubble--loading {
  padding: 10px 14px;
  display: flex;
  align-items: center;
  justify-content: center;
}
.mtg-spinner {
  /* Change these two values to resize the whole spinner */
  --dot: 4px;
  --orbit: 6px;

  width: calc(var(--orbit) * 2 + var(--dot));
  height: calc(var(--orbit) * 2 + var(--dot));
  position: relative;
  animation: mtg-rotate 1.4s linear infinite;
}
.mtg-dot {
  width: var(--dot);
  height: var(--dot);
  border-radius: 50%;
  position: absolute;
  top: 50%;
  left: 50%;
  margin: calc(var(--dot) / -2) 0 0 calc(var(--dot) / -2);
  box-shadow: 0 0 0 1.5px rgba(255, 255, 255, 0.35);
}
.mtg-dot--w { background: #f0c040; transform: rotate(0deg)   translateY(calc(-1 * var(--orbit))); }
.mtg-dot--u { background: #3b82f6; transform: rotate(72deg)  translateY(calc(-1 * var(--orbit))); }
.mtg-dot--b { background: #1a1a1a; transform: rotate(144deg) translateY(calc(-1 * var(--orbit))); }
.mtg-dot--r { background: #ef4444; transform: rotate(216deg) translateY(calc(-1 * var(--orbit))); }
.mtg-dot--g { background: #16a34a; transform: rotate(288deg) translateY(calc(-1 * var(--orbit))); }
@keyframes mtg-rotate {
  to { transform: rotate(360deg); }
}

/* ── Error banner ── */
.chat-error {
  margin: 0 16px 8px;
  padding: 10px 14px;
  border-radius: 8px;
  font-size: 14px;
  background: rgba(239, 68, 68, 0.12);
  color: #ef4444;
  border: 1px solid rgba(239, 68, 68, 0.3);
  flex-shrink: 0;
}

/* ── Input form ── */
.chat-form {
  display: flex;
  align-items: flex-end;
  gap: 8px;
  padding: 12px 16px 16px;
  border-top: 1px solid var(--border);
  flex-shrink: 0;
}
.chat-input {
  flex: 1;
  resize: none;
  font-family: var(--sans);
  font-size: 15px;
  line-height: 1.5;
  padding: 10px 14px;
  border-radius: 12px;
  border: 1px solid var(--border);
  background: var(--bg);
  color: var(--text-h);
  outline: none;
  transition: border-color 0.15s;
  field-sizing: content;
  max-height: 160px;
  overflow-y: auto;
}
.chat-input:focus {
  border-color: var(--accent);
}
.chat-input:disabled {
  opacity: 0.6;
}
.chat-send {
  font-family: var(--sans);
  font-size: 14px;
  font-weight: 500;
  padding: 10px 18px;
  border-radius: 12px;
  border: none;
  background: var(--accent);
  color: #fff;
  cursor: pointer;
  transition: opacity 0.15s;
  white-space: nowrap;
  flex-shrink: 0;
}
.chat-send:hover:not(:disabled) {
  opacity: 0.88;
}
.chat-send:disabled {
  opacity: 0.45;
  cursor: not-allowed;
}
</style>
