<script setup lang="ts">
interface Props {
  variant?: 'primary' | 'secondary' | 'danger' | 'ghost'
  size?: 'sm' | 'md' | 'lg' | 'icon'
  type?: 'button' | 'submit' | 'reset'
  disabled?: boolean
  loading?: boolean
  highlight?: boolean
}

withDefaults(defineProps<Props>(), {
  variant: 'secondary',
  size: 'md',
  type: 'button',
  disabled: false,
  loading: false,
  highlight: false,
})
</script>

<template>
  <button
    :type="type"
    :disabled="disabled || loading"
    class="btn"
    :class="[`btn--${variant}`, `btn--${size}`, { 'btn--highlight': highlight }]"
  >
    <span v-if="loading" class="btn__loader" />
    <slot />
  </button>
</template>

<style scoped>
.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  font-weight: 500;
  cursor: pointer;
  transition: background 0.15s, color 0.15s, border-color 0.15s;
  border: 1px solid transparent;
  border-radius: 8px;
  font-family: inherit;
  white-space: nowrap;
}

.btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

/* Variants */
.btn--primary {
  background: #6366f1;
  color: #fff;
  border-color: #6366f1;
}
.btn--primary:hover:not(:disabled) {
  background: #4f46e5;
  border-color: #4f46e5;
}

.btn--secondary {
  background: #1e1e1e;
  color: #a3a3a3;
  border-color: #2a2a2a;
}
.btn--secondary:hover:not(:disabled) {
  background: #252525;
  color: #e5e5e5;
  border-color: #333;
}

.btn--danger {
  background: transparent;
  color: #f87171;
  border-color: #3f1515;
}
.btn--danger:hover:not(:disabled) {
  background: #1f0f0f;
  border-color: #7f1d1d;
}

.btn--ghost {
  background: transparent;
  color: #737373;
  border-color: transparent;
}
.btn--ghost:hover:not(:disabled) {
  background: #1a1a1a;
  color: #e5e5e5;
}

/* Sizes */
.btn--sm  { padding: 5px 10px; font-size: 12px; }
.btn--md  { padding: 7px 14px; font-size: 14px; }
.btn--lg  { padding: 10px 20px; font-size: 15px; }
.btn--icon { padding: 8px; border-radius: 50%; }

/* Loader */
.btn__loader {
  width: 14px;
  height: 14px;
  border: 2px solid currentColor;
  border-top-color: transparent;
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
  flex-shrink: 0;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* Highlight (shake + red for login prompt) */
.btn--highlight {
  background: #7f1d1d;
  border-color: #991b1b;
  color: #fca5a5;
  animation: shake 0.5s ease-in-out;
}

@keyframes shake {
  0%, 100% { transform: translateX(0); }
  20%, 60% { transform: translateX(-4px); }
  40%, 80% { transform: translateX(4px); }
}
</style>
