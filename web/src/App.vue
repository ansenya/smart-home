<script setup lang="ts">

import Header from "./components/Header.vue";
import PairingButton from "./components/PairingButton.vue";
</script>

<template>
  <Header></Header>
  <RouterView/>
  <PairingButton/>
</template>
