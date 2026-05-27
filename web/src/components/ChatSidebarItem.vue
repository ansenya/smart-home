<script setup lang="ts">
import { ref, nextTick } from 'vue'
import type { Chat } from '@/types/chat'

const props = defineProps<{ chat: Chat; active?: boolean }>()
const emit = defineEmits<{
  click: []
  delete: []
  rename: [title: string]
}>()

const editing = ref(false)
const editTitle = ref('')
const inputRef = ref<HTMLInputElement | null>(null)

const formatDate = (dateString: string) => {
  const date = new Date(dateString)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))
  if (days === 0) return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
  if (days === 1) return 'Yesterday'
  if (days < 7) return date.toLocaleDateString([], { weekday: 'short' })
  return date.toLocaleDateString([], { month: 'short', day: 'numeric' })
}

const startEdit = async (e: Event) => {
  e.stopPropagation()
  editTitle.value = props.chat.title || ''
  editing.value = true
  await nextTick()
  inputRef.value?.focus()
  inputRef.value?.select()
}

const commitEdit = () => {
  const title = editTitle.value.trim()
  if (title && title !== props.chat.title) {
    emit('rename', title)
  }
  editing.value = false
}

const cancelEdit = () => {
  editing.value = false
}

const onKeydown = (e: KeyboardEvent) => {
  if (e.key === 'Enter') commitEdit()
  if (e.key === 'Escape') cancelEdit()
}
</script>

<template>
  <div class="item" :class="{ 'item--active': active }" @click="emit('click')">
    <div class="item-icon">
      <svg width="13" height="13" fill="none" stroke="currentColor" stroke-width="1.8" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"/>
      </svg>
    </div>

    <div class="item-body">
      <input
        v-if="editing"
        ref="inputRef"
        v-model="editTitle"
        class="title-input"
        @blur="commitEdit"
        @keydown="onKeydown"
        @click.stop
      />
      <span v-else class="item-title">{{ chat.title || 'New chat' }}</span>
      <span class="item-date">{{ formatDate(chat.created_at) }}</span>
    </div>

    <div v-if="!editing" class="item-actions">
      <button class="action-btn" title="Rename" @click.stop="startEdit">
        <svg width="12" height="12" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"/>
        </svg>
      </button>
      <button class="action-btn action-btn--danger" title="Delete" @click.stop="emit('delete')">
        <svg width="12" height="12" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
        </svg>
      </button>
    </div>
  </div>
</template>

<style scoped>
.item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 10px;
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.12s;
  min-width: 0;
  position: relative;
}

.item:hover { background: #1a1a1a; }
.item--active { background: #212121; }

.item-icon {
  width: 24px;
  height: 24px;
  border-radius: 6px;
  background: #1e1e1e;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #555;
  flex-shrink: 0;
  transition: background 0.12s, color 0.12s;
}
.item--active .item-icon { background: #2a2a2a; color: #888; }
.item:hover .item-icon { background: #252525; color: #777; }

.item-body {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 1px;
}

.item-title {
  font-size: 13px;
  font-weight: 400;
  color: #a3a3a3;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  transition: color 0.12s;
}
.item--active .item-title,
.item:hover .item-title { color: #e5e5e5; }

.item-date {
  font-size: 11px;
  color: #404040;
}
.item--active .item-date { color: #525252; }

.title-input {
  font-size: 13px;
  font-weight: 400;
  color: #e5e5e5;
  background: #2a2a2a;
  border: 1px solid #404040;
  border-radius: 4px;
  padding: 1px 6px;
  width: 100%;
  outline: none;
  font-family: inherit;
}

/* Action buttons — visible on hover */
.item-actions {
  display: none;
  align-items: center;
  gap: 2px;
  flex-shrink: 0;
}
.item:hover .item-actions,
.item--active .item-actions { display: flex; }

.action-btn {
  width: 22px;
  height: 22px;
  border-radius: 5px;
  border: none;
  background: transparent;
  color: #555;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background 0.12s, color 0.12s;
}
.action-btn:hover { background: #2a2a2a; color: #a3a3a3; }
.action-btn--danger:hover { background: #2d1111; color: #f87171; }
</style>
