<template>
  <div
      :class="{'input-error': error || formError}"
      class="input-container">
    <input
        :value="modelValue"
        @input="$emit('update:modelValue', $event.target.value)"
        :type="type"
        :placeholder="placeholder"
        :required="required"
    />
    <img :src="src" alt="input icon"/>
  </div>
  <span v-if="error" class="error-text">{{ error }}</span>
</template>

<style scoped>

.input-container {
  display: flex;
  flex-direction: row;
  margin: 6px auto;
  padding: 16px;
  background-color: #1e1e1e;
  border: 1px solid #505050;
  border-radius: 10px;
  color: #fff;
  font-size: 14px;
  outline: none;
  transition: border 0.2s;
  width: 100%;
  max-width: 280px;
}

.input-container input {
  background: rgba(255, 255, 255, 0);
  border: none;
  width: 100%;
  height: 100%;
  margin-right: 10px;
  outline: none;
  color: #fff;
  padding: 0;
}

.input-container img {
  max-height: 25px;
}

.input-error {
  border: 1px solid #ff4d4f !important;
}

.error-text {
  color: #ff4d4f;
  font-size: 12px;
  max-width: 300px;
  margin-top: 4px;
  width: 100%;
  text-align: center;
}

</style>

<script setup>
defineProps({
  modelValue: {
    type: String,
    required: true
  },
  placeholder: {
    type: String,
    default: ''
  },
  type: {
    type: String,
    default: 'text'
  },
  src: {
    type: String,
    default: null
  },
  error: {
    type: String,
    default: null
  },
  formError: {
    type: String,
    default: null
  },
  required: {
    type: Boolean,
    default: false
  }
})

defineEmits(['update:modelValue'])

</script>