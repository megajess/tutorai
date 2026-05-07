<script setup lang="ts">
import { ref, nextTick } from 'vue'
import { sendMessage } from '../api/chat'

interface Message {
  role: 'user' | 'assistant'
  text: string
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

async function submit() {
  const query = input.value.trim()
  if (!query || loading.value) return

  error.value = null
  input.value = ''
  messages.value.push({ role: 'user', text: query })
  await scrollToBottom()

  loading.value = true
  try {
    const response = await sendMessage(query)
    messages.value.push({ role: 'assistant', text: response })
  } catch (err) {
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
        <div class="chat-bubble">{{ msg.text }}</div>
      </div>

      <div v-if="loading" class="chat-message chat-message--assistant">
        <div class="chat-bubble chat-bubble--loading">
          <span class="dot" />
          <span class="dot" />
          <span class="dot" />
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

/* ── Bubbles ── */
.chat-bubble {
  max-width: 72%;
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

/* ── Loading dots ── */
.chat-bubble--loading {
  display: flex;
  align-items: center;
  gap: 5px;
  padding: 12px 16px;
}
.dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: var(--text);
  animation: blink 1.2s infinite ease-in-out;
}
.dot:nth-child(2) { animation-delay: 0.2s; }
.dot:nth-child(3) { animation-delay: 0.4s; }
@keyframes blink {
  0%, 80%, 100% { opacity: 0.2; }
  40% { opacity: 1; }
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
