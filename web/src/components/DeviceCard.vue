<script setup lang="ts">
import { computed } from 'vue'
import type { Device } from '@/types/device'
import { deviceTypeLabel, isOn, getCapability } from '@/types/device'

const props = defineProps<{
  device: Device
  pending?: boolean
}>()

const emit = defineEmits<{
  (e: 'open', device: Device): void
  (e: 'toggle', device: Device, next: boolean): void
}>()

const on = computed(() => isOn(props.device))

const isOnline = computed(() => {
  if (!props.device.last_seen) return false
  const last = new Date(props.device.last_seen).getTime()
  return Date.now() - last < 60_000
})

const brightness = computed(() => {
  const cap = getCapability(props.device, 'devices.capabilities.range')
  if (!cap?.state?.value) return undefined
  const v = Number(cap.state.value)
  return Number.isFinite(v) ? v : undefined
})

const colorRGB = computed<string | undefined>(() => {
  const cap = getCapability(props.device, 'devices.capabilities.color_setting')
  if (!cap?.state?.value) return undefined
  const v = cap.state.value
  if (typeof v === 'number') {
    return '#' + v.toString(16).padStart(6, '0')
  }
  if (v && typeof v === 'object') {
    const o = v as Record<string, unknown>
    if ('h' in o && 's' in o && 'v' in o) {
      return `hsl(${Number(o.h)}, ${Number(o.s)}%, ${Math.round(Number(o.v) * 0.5)}%)`
    }
    if ('r' in o && 'g' in o && 'b' in o) {
      return `rgb(${Number(o.r)}, ${Number(o.g)}, ${Number(o.b)})`
    }
  }
  return undefined
})

const icon = computed(() => {
  switch (props.device.type) {
    case 'devices.types.light':
      return 'M9 18h6m-7 0a4 4 0 008 0c0-1.1.9-2 2-2v0a6 6 0 10-12 0v0c1.1 0 2 .9 2 2z'
    case 'devices.types.socket':
      return 'M9 17v-4m6 4v-4M5 9V5a2 2 0 012-2h10a2 2 0 012 2v4M3 9h18v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z'
    case 'devices.types.thermostat':
      return 'M14 14.76V3a2 2 0 10-4 0v11.76a4 4 0 104 0z'
    case 'devices.types.sensor':
      return 'M19 11H5m14 0a2 2 0 100-4H5a2 2 0 100 4m14 0v4a2 2 0 01-2 2H7a2 2 0 01-2-2v-4'
    default:
      return 'M21 16V8a2 2 0 00-1-1.73l-7-4a2 2 0 00-2 0l-7 4A2 2 0 003 8v8a2 2 0 001 1.73l7 4a2 2 0 002 0l7-4A2 2 0 0021 16z'
  }
})

function toggle(e: MouseEvent) {
  e.stopPropagation()
  if (on.value === undefined) return
  emit('toggle', props.device, !on.value)
}
</script>

<template>
  <div
    class="card"
    :class="{ 'card--on': on, 'card--pending': pending }"
    @click="emit('open', device)"
  >
    <div class="card-header">
      <div class="card-icon" :class="{ 'card-icon--on': on }">
        <svg width="20" height="20" fill="none" stroke="currentColor" stroke-width="1.8" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" :d="icon" />
        </svg>
      </div>
      <div class="status-dot" :class="{ 'status-dot--online': isOnline }" :title="isOnline ? 'Online' : 'Offline'" />
    </div>

    <div class="card-body">
      <div class="card-name" :title="device.name || 'Unnamed device'">
        {{ device.name || 'Unnamed device' }}
      </div>
      <div class="card-meta">
        <span class="meta-type">{{ deviceTypeLabel(device.type) }}</span>
        <span v-if="device.room" class="meta-dot">·</span>
        <span v-if="device.room" class="meta-room">{{ device.room }}</span>
      </div>
    </div>

    <div class="card-footer">
      <div v-if="brightness !== undefined && on" class="bright-indicator">
        <div class="bright-bar"><div class="bright-fill" :style="{ width: brightness + '%' }" /></div>
      </div>

      <div v-if="colorRGB && on" class="color-dot" :style="{ background: colorRGB }" />

      <button
        v-if="on !== undefined"
        class="toggle"
        :class="{ 'toggle--on': on }"
        @click="toggle"
        :disabled="pending"
      >
        <span class="toggle-knob" />
      </button>
    </div>
  </div>
</template>

<style scoped>
.card {
  background: #111;
  border: 1px solid #1e1e1e;
  border-radius: 14px;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  cursor: pointer;
  transition: border-color 0.15s, background 0.15s, transform 0.08s;
  min-height: 140px;
}
.card:hover { border-color: #2a2a2a; background: #131313; }
.card:active { transform: scale(0.99); }
.card--on { border-color: #2c2860; background: linear-gradient(180deg, #14122a 0%, #111 70%); }
.card--on:hover { border-color: #3a3470; }
.card--pending { opacity: 0.7; pointer-events: none; }

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.card-icon {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  background: #1a1a1a;
  border: 1px solid #222;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #737373;
  transition: background 0.15s, color 0.15s;
}
.card-icon--on { background: #2c2860; color: #c4b5fd; border-color: #3a3470; }

.status-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: #404040;
}
.status-dot--online { background: #22c55e; box-shadow: 0 0 8px rgba(34, 197, 94, 0.6); }

.card-body { display: flex; flex-direction: column; gap: 4px; flex: 1; }
.card-name {
  font-size: 14px;
  font-weight: 600;
  color: #e5e5e5;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.card-meta {
  font-size: 12px;
  color: #525252;
  display: flex;
  align-items: center;
  gap: 6px;
}
.meta-dot { color: #2a2a2a; }
.meta-room { color: #737373; }

.card-footer {
  display: flex;
  align-items: center;
  gap: 10px;
}

.bright-indicator { flex: 1; min-width: 0; }
.bright-bar {
  height: 4px;
  border-radius: 2px;
  background: #222;
  overflow: hidden;
}
.bright-fill {
  height: 100%;
  background: linear-gradient(90deg, #6366f1, #c084fc);
  transition: width 0.2s;
}

.color-dot {
  width: 14px;
  height: 14px;
  border-radius: 50%;
  border: 1px solid #2a2a2a;
  box-shadow: inset 0 0 4px rgba(0,0,0,0.4);
}

.toggle {
  width: 38px;
  height: 22px;
  border-radius: 11px;
  background: #2a2a2a;
  border: none;
  position: relative;
  cursor: pointer;
  margin-left: auto;
  padding: 0;
  transition: background 0.15s;
}
.toggle--on { background: #6366f1; }
.toggle:hover { filter: brightness(1.1); }
.toggle:disabled { cursor: not-allowed; opacity: 0.6; }

.toggle-knob {
  position: absolute;
  top: 2px;
  left: 2px;
  width: 18px;
  height: 18px;
  border-radius: 50%;
  background: #fff;
  transition: transform 0.15s ease;
  box-shadow: 0 1px 3px rgba(0,0,0,0.4);
}
.toggle--on .toggle-knob { transform: translateX(16px); }
</style>
