import { ref } from 'vue'
import { defineStore } from 'pinia'

export const useFloatingButtonDevicesButtonStore = defineStore(
  'floatingButtonDevicesButtonStore',
  () => {
    const counter = ref(0)
    const isFan = ref(false)

    async function incr() {
      counter.value++
      if (counter.value >= 10) {
        isFan.value = true
      }
      setTimeout(async () => {
        counter.value--
      }, 10000)
    }

    return {isFan, incr}
  }
)
