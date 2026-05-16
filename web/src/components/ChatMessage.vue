<script setup lang="ts">
import { computed, ref } from 'vue'
import type { Message } from '@/types/chat'

const props = defineProps<{ message: Message }>()

const expanded = ref(false)

// Parse and pretty-print tool result content
const toolResultFormatted = computed(() => {
  const raw = props.message.tool_result
  if (!raw) return ''
  // tool_result is stored as { result: "json string" }
  const inner = (raw as any).result
  if (typeof inner === 'string') {
    try {
      return JSON.stringify(JSON.parse(inner), null, 2)
    } catch {
      return inner
    }
  }
  return JSON.stringify(raw, null, 2)
})

const toolArgsFormatted = computed(() => {
  const args = props.message.tool_args
  if (!args || !Array.isArray(args)) return null
  // args is ToolCall[] — take the first one's arguments
  const first = args[0] as any
  if (!first?.Arguments) return null
  try {
    const parsed = JSON.parse(first.Arguments)
    const entries = Object.entries(parsed)
    return entries.length > 0 ? entries : null
  } catch {
    return null
  }
})

const toolIcon = computed(() => {
  const name = props.message.tool_name ?? ''
  if (name.includes('list')) return 'list'
  if (name.includes('control')) return 'control'
  if (name.includes('state') || name.includes('get')) return 'info'
  return 'tool'
})
</script>

<template>
  <!-- Tool result message -->
  <div v-if="message.role === 'tool'" class="tool-wrap">
    <button class="tool-card" :class="{ 'tool-card--open': expanded }" @click="expanded = !expanded">
      <div class="tool-card-left">
        <div class="tool-status-icon tool-status-icon--ok">
          <svg width="11" height="11" fill="none" stroke="currentColor" stroke-width="2.5" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7"/>
          </svg>
        </div>
        <div class="tool-card-meta">
          <span class="tool-card-name">{{ message.tool_name }}</span>
          <span class="tool-card-label">результат получен</span>
        </div>
      </div>
      <svg class="tool-chevron" :class="{ 'tool-chevron--open': expanded }"
        width="13" height="13" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" d="M19 9l-7 7-7-7"/>
      </svg>
    </button>
    <div v-if="expanded" class="tool-body">
      <pre class="tool-pre">{{ toolResultFormatted }}</pre>
    </div>
  </div>

  <!-- System message -->
  <div v-else-if="message.role === 'system'" class="system-row">
    <span>{{ message.content }}</span>
  </div>

  <!-- User message -->
  <div v-else-if="message.role === 'user'" class="user-row">
    <div class="user-bubble">
      <span class="message-text">{{ message.content }}</span>
    </div>
  </div>

  <!-- Assistant message -->
  <div v-else class="assistant-row">
    <div class="assistant-avatar">
      <svg width="14" height="14" fill="currentColor" viewBox="0 0 20 20">
        <path d="M10.707 2.293a1 1 0 00-1.414 0l-7 7a1 1 0 001.414 1.414L4 10.414V17a1 1 0 001 1h2a1 1 0 001-1v-2a1 1 0 011-1h2a1 1 0 011 1v2a1 1 0 001 1h2a1 1 0 001-1v-6.586l.293.293a1 1 0 001.414-1.414l-7-7z"/>
      </svg>
    </div>
    <div class="assistant-content">

      <!-- Tool call card (assistant is calling a tool) -->
      <div v-if="message.tool_args" class="tool-calling-card">
        <div class="tool-calling-header">
          <div class="tool-calling-icon">
            <!-- list -->
            <svg v-if="toolIcon === 'list'" width="12" height="12" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" d="M4 6h16M4 10h16M4 14h16M4 18h16"/>
            </svg>
            <!-- control -->
            <svg v-else-if="toolIcon === 'control'" width="12" height="12" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" d="M13 10V3L4 14h7v7l9-11h-7z"/>
            </svg>
            <!-- info -->
            <svg v-else-if="toolIcon === 'info'" width="12" height="12" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
            </svg>
            <!-- default gear -->
            <svg v-else width="12" height="12" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"/>
              <path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
            </svg>
          </div>
          <span class="tool-calling-name">{{ message.tool_name }}</span>
          <div class="tool-calling-spinner" />
        </div>
        <div v-if="toolArgsFormatted" class="tool-calling-args">
          <span v-for="[k, v] in toolArgsFormatted" :key="k" class="tool-arg-chip">
            <span class="tool-arg-key">{{ k }}</span>
            <span class="tool-arg-val">{{ String(v).length > 24 ? String(v).slice(0, 24) + '…' : v }}</span>
          </span>
        </div>
      </div>

      <!-- Text content -->
      <div v-if="message.content || (message.status === 'pending' && !message.tool_args)" class="assistant-text">
        {{ message.content }}<span v-if="message.status === 'pending' && !message.content" class="thinking">
          <span /><span /><span />
        </span><span v-else-if="message.status === 'pending'" class="cursor" />
      </div>
    </div>
  </div>
</template>

<style scoped>
/* ── User ───────────────────────────────────────────────────── */
.user-row {
  display: flex;
  justify-content: flex-end;
  padding: 3px 0;
}

.user-bubble {
  max-width: 70%;
  background: #2f2f2f;
  border-radius: 18px 18px 4px 18px;
  padding: 10px 16px;
}

.message-text {
  font-size: 14px;
  line-height: 1.65;
  color: #e5e5e5;
  white-space: pre-wrap;
  word-break: break-word;
}

/* ── Assistant ──────────────────────────────────────────────── */
.assistant-row {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  padding: 8px 0;
  max-width: 85%;
}

.assistant-avatar {
  width: 28px;
  height: 28px;
  border-radius: 8px;
  background: #6366f1;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  flex-shrink: 0;
  margin-top: 2px;
}

.assistant-content {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.assistant-text {
  font-size: 14px;
  line-height: 1.75;
  color: #d4d4d4;
  white-space: pre-wrap;
  word-break: break-word;
}

.cursor {
  display: inline-block;
  width: 2px;
  height: 14px;
  background: #6366f1;
  border-radius: 1px;
  margin-left: 2px;
  vertical-align: middle;
  animation: blink 1s step-end infinite;
}
@keyframes blink { 0%, 100% { opacity: 1; } 50% { opacity: 0; } }

.thinking {
  display: inline-flex;
  gap: 4px;
  margin-left: 4px;
  vertical-align: middle;
}
.thinking span {
  width: 5px; height: 5px;
  border-radius: 50%;
  background: #555;
  animation: dot 1.2s ease-in-out infinite;
}
.thinking span:nth-child(2) { animation-delay: 0.2s; }
.thinking span:nth-child(3) { animation-delay: 0.4s; }
@keyframes dot {
  0%, 80%, 100% { transform: scale(0.6); opacity: 0.4; }
  40% { transform: scale(1); opacity: 1; }
}

/* ── Tool calling card (inside assistant bubble) ────────────── */
.tool-calling-card {
  display: inline-flex;
  flex-direction: column;
  gap: 6px;
  padding: 8px 12px;
  background: #161616;
  border: 1px solid #252525;
  border-radius: 10px;
  max-width: 100%;
}

.tool-calling-header {
  display: flex;
  align-items: center;
  gap: 8px;
}

.tool-calling-icon {
  width: 22px;
  height: 22px;
  border-radius: 6px;
  background: #1e1e2e;
  border: 1px solid #2d2d44;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #818cf8;
  flex-shrink: 0;
}

.tool-calling-name {
  font-size: 12px;
  font-weight: 500;
  color: #a5b4fc;
  font-family: ui-monospace, monospace;
  flex: 1;
}

.tool-calling-spinner {
  width: 10px;
  height: 10px;
  border: 1.5px solid #2d2d44;
  border-top-color: #6366f1;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  flex-shrink: 0;
}
@keyframes spin { to { transform: rotate(360deg); } }

.tool-calling-args {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}

.tool-arg-chip {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 2px 7px;
  border-radius: 4px;
  background: #1c1c2e;
  border: 1px solid #2a2a3e;
  font-size: 11px;
  font-family: ui-monospace, monospace;
}

.tool-arg-key { color: #555; }
.tool-arg-val { color: #818cf8; }

/* ── Tool result card (standalone tool role message) ─────────── */
.tool-wrap {
  margin-left: 38px; /* align with assistant text, past the avatar */
  max-width: calc(85% - 38px);
  display: flex;
  flex-direction: column;
  gap: 1px;
  padding: 2px 0;
}

.tool-card {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 7px 10px;
  background: #141414;
  border: 1px solid #222;
  border-radius: 8px;
  cursor: pointer;
  text-align: left;
  width: 100%;
  transition: background 0.12s, border-color 0.12s;
  font-family: inherit;
  gap: 10px;
}
.tool-card:hover { background: #181818; border-color: #2a2a2a; }
.tool-card--open { border-radius: 8px 8px 0 0; border-bottom-color: transparent; }

.tool-card-left {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
}

.tool-status-icon {
  width: 18px;
  height: 18px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}
.tool-status-icon--ok {
  background: #052e16;
  color: #4ade80;
}

.tool-card-meta {
  display: flex;
  flex-direction: column;
  gap: 1px;
  min-width: 0;
}

.tool-card-name {
  font-size: 12px;
  font-weight: 500;
  color: #a3a3a3;
  font-family: ui-monospace, monospace;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.tool-card-label {
  font-size: 10px;
  color: #3f3f3f;
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

.tool-chevron {
  color: #404040;
  flex-shrink: 0;
  transition: transform 0.2s;
}
.tool-chevron--open { transform: rotate(180deg); }

.tool-body {
  background: #0f0f0f;
  border: 1px solid #222;
  border-top: none;
  border-radius: 0 0 8px 8px;
  padding: 10px 12px;
  overflow-x: auto;
}

.tool-pre {
  font-size: 11px;
  color: #5a5a5a;
  font-family: ui-monospace, monospace;
  white-space: pre;
  margin: 0;
  line-height: 1.65;
}

/* ── System ─────────────────────────────────────────────────── */
.system-row {
  display: flex;
  justify-content: center;
  padding: 8px 0;
}
.system-row span {
  font-size: 12px;
  color: #3f3f3f;
  background: #1a1a1a;
  padding: 4px 12px;
  border-radius: 20px;
}
</style>
