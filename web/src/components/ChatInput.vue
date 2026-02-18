<script setup lang="ts">
import { ref, nextTick } from 'vue'
import ArrowUpIcon from '@/assets/icons/ArrowUpIcon.vue'
import BaseButton from '@/components/BaseButton.vue'

const textarea = ref<HTMLTextAreaElement | null>(null)

const autoResize = () => {
  if (textarea.value) {
    textarea.value.style.height = 'auto'
    textarea.value.style.height = textarea.value.scrollHeight + 'px'
  }
}

const handleKeydown = (event: KeyboardEvent) => {
  if ((event.ctrlKey || event.metaKey) && event.key === 'Enter') {
    event.preventDefault()
  }
}
</script>

<template>
  <div class="input-container">
    <div class="input-wrapper">
      <textarea
        ref="textarea"
        placeholder="Message..."
        rows="1"
        @input="autoResize"
        @keydown="handleKeydown"
      ></textarea>
      <button class="send-btn">
        <ArrowUpIcon class="icon" />
      </button>
    </div>
    <p class="disclaimer">AI can make mistakes. Check important info.</p>
  </div>
</template>

<style scoped>
.input-container {
  padding: 24px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  gap: 12px;
}

.input-wrapper {
  width: 100%;
  max-width: 768px;
  position: relative;
  display: flex;
  align-items: flex-end;
  background: #2f2f2f;
  border-radius: 8px;
  border: 1px solid #404040;
  transition: border-color 0.2s;
  padding: 6px;
}

.input-wrapper:focus-within {
  border-color: #666;
}

.input-wrapper textarea {
  flex: 1;
  background: transparent;
  border: none;
  color: #fff;
  padding: 12px 48px 12px 16px;
  font-size: 14px;
  resize: none;
  max-height: 480px;
  line-height: 1.5;
  font-family: inherit;
  overflow-y: auto;
}

.input-wrapper textarea:focus {
  outline: none;
}

.input-wrapper textarea::placeholder {
  color: #888;
}

.input-wrapper textarea::-webkit-scrollbar {
  width: 6px;
}

.input-wrapper textarea::-webkit-scrollbar-track {
  background: transparent;
}

.input-wrapper textarea::-webkit-scrollbar-thumb {
  background: #555;
  border-radius: 3px;
}

.input-wrapper textarea::-webkit-scrollbar-thumb:hover {
  background: #666;
}

.send-btn {
  position: absolute;
  right: 11px;
  bottom: 11px;
  background: #fff;
  border: none;
  border-radius: 6px;
  width: 36px;
  height: 36px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: opacity 0.2s;
  color: #000;
  font-size: 14px;
}

.send-btn:hover {
  opacity: 0.8;
}

.send-btn:disabled {
  opacity: 0.3;
  cursor: not-allowed;
}

.icon {
  height: 25px;
  width: 25px;
}

.disclaimer {
  font-size: 12px;
  color: #888;
  text-align: center;
  margin: 0;
}
</style>
