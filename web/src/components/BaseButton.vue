<script setup lang="ts">
interface Props {
  variant?: 'primary' | 'secondary' | 'danger' | 'ghost'
  size?: 'sm' | 'md' | 'lg' | 'icon'
  type?: 'button' | 'submit' | 'reset'
  disabled?: boolean
  loading?: boolean
}

withDefaults(defineProps<Props>(), {
  variant: 'secondary',
  size: 'md',
  type: 'button',
  disabled: false,
  loading: false,
})
</script>

<template>
  <button
    :type="type"
    :disabled="disabled || loading"
    class="btn"
    :class="[`btn--${variant}`, `btn--${size}`]"
  >
    <span v-if="loading" class="btn__loader">⟳</span>
    <slot />
  </button>
</template>

<style scoped>
.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  border: 1px solid transparent;
  border-radius: 6px;
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn--primary {
  background: #3b82f6;
  color: white;
}
.btn--primary:hover:not(:disabled) {
  background: #2563eb;
}

.btn--secondary {
  background: #f3f4f6;
  color: #1f2937;
}
.btn--secondary:hover:not(:disabled) {
  background: #e5e7eb;
}

.btn--danger {
  background: #fee2e2;
  color: #dc2626;
}
.btn--danger:hover:not(:disabled) {
  background: #fecaca;
}

.btn--ghost {
  background: transparent;
  color: #1f2937;
}
.btn--ghost:hover:not(:disabled) {
  background: #f3f4f6;
}

.btn--sm {
  padding: 4px 12px;
  font-size: 13px;
}

.btn--md {
  padding: 8px 16px;
  font-size: 14px;
}

.btn--lg {
  padding: 12px 24px;
  font-size: 16px;
}

.btn--icon {
  padding: 8px;
  border-radius: 50%;
}

.btn__loader {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}
</style>
