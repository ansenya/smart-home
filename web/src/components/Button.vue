<script setup lang="ts">
interface ButtonProps {
  variant?: 'primary' | 'secondary'
}

defineProps<ButtonProps>()
</script>

<template>
  <button
      type="button"
      class="text-white bg-blue-500 hover:bg-blue-600 border rounded-md text-sm px-4 py-2.5"
  >
    <slot>Default Button</slot>
  </button>
</template>