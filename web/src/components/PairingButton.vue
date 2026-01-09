<script setup lang="ts">

</script>

<template>
  <button class="pair-btn" @click="$router.push('/devices/add')">+</button>
</template>

<style scoped>
.pair-btn {
  position: fixed;
  right: 24px;
  bottom: 24px;
  width: 56px;
  height: 56px;
  border-radius: 50%;
  font-size: 28px;
  background: #2563eb;
  color: white;
  border: none;

  cursor: pointer;
}

.pair-btn:hover {
  background: #4a76d7;
}
</style>
