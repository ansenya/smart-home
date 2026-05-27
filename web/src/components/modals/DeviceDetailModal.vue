<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { push } from 'notivue'
import type { Device, Capability } from '@/types/device'
import { CAPABILITY_LABELS, deviceTypeLabel } from '@/types/device'
import { updateDevice, deleteDevice, setCapability, getDevice } from '@/api/devices'

const props = defineProps<{
  modelValue: boolean
  device: Device | null
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'updated', device: Device): void
  (e: 'deleted', id: string): void
  (e: 'refresh'): void
}>()

const localDevice = ref<Device | null>(null)
const editName = ref('')
const editRoom = ref('')
const editDescription = ref('')
const renaming = ref(false)
const deleting = ref(false)
const capPending = ref<Set<string>>(new Set())
const confirmDelete = ref(false)

watch(() => props.device, (d) => {
  localDevice.value = d ? { ...d } : null
  if (d) {
    editName.value = d.name ?? ''
    editRoom.value = d.room ?? ''
    editDescription.value = d.description ?? ''
  }
  confirmDelete.value = false
}, { immediate: true })

watch(() => props.modelValue, (open) => {
  if (open && props.device) {
    refreshDevice()
  }
})

async function refreshDevice() {
  if (!props.device) return
  try {
    const r = await getDevice(props.device.id)
    localDevice.value = r.data
    editName.value = r.data.name ?? ''
    editRoom.value = r.data.room ?? ''
    editDescription.value = r.data.description ?? ''
  } catch {
    // silent — modal will keep cached
  }
}

const isOnline = computed(() => {
  if (!localDevice.value?.last_seen) return false
  return Date.now() - new Date(localDevice.value.last_seen).getTime() < 60_000
})

function close() {
  emit('update:modelValue', false)
}

async function saveMeta() {
  if (!localDevice.value) return
  const payload: Record<string, string> = {}
  if (editName.value !== localDevice.value.name) payload.name = editName.value
  if (editRoom.value !== localDevice.value.room) payload.room = editRoom.value
  if (editDescription.value !== localDevice.value.description) payload.description = editDescription.value
  if (Object.keys(payload).length === 0) return

  renaming.value = true
  try {
    const r = await updateDevice(localDevice.value.id, payload)
    localDevice.value = r.data
    emit('updated', r.data)
    push.success({ title: 'Saved', message: 'Changes applied' })
  } catch {
    push.error({ title: 'Failed', message: 'Could not save changes' })
  } finally {
    renaming.value = false
  }
}

async function handleDelete() {
  if (!localDevice.value) return
  if (!confirmDelete.value) {
    confirmDelete.value = true
    return
  }
  deleting.value = true
  try {
    await deleteDevice(localDevice.value.id)
    emit('deleted', localDevice.value.id)
    push.success({ title: 'Device deleted' })
  } catch {
    push.error({ title: 'Failed to delete' })
  } finally {
    deleting.value = false
    confirmDelete.value = false
  }
}

async function setCap(cap: Capability, value: unknown, instance?: string) {
  if (!localDevice.value) return
  if (capPending.value.has(cap.id)) return

  capPending.value.add(cap.id)
  const prev = cap.state?.value

  // Optimistic update
  cap.state = { ...(cap.state ?? {}), value, instance: instance ?? cap.state?.instance }

  try {
    await setCapability(localDevice.value.id, cap.type, { value, instance })
    emit('updated', localDevice.value)
  } catch {
    cap.state = { ...(cap.state ?? {}), value: prev }
    push.error({ title: 'Command was not delivered' })
  } finally {
    capPending.value.delete(cap.id)
  }
}

const onOffCap = computed(() => localDevice.value?.capabilities.find(c => c.type === 'devices.capabilities.on_off'))
const colorCap = computed(() => localDevice.value?.capabilities.find(c => c.type === 'devices.capabilities.color_setting'))
const rangeCap = computed(() => localDevice.value?.capabilities.find(c => c.type === 'devices.capabilities.range'))
const modeCaps = computed(() => localDevice.value?.capabilities.filter(c => c.type === 'devices.capabilities.mode') ?? [])
const toggleCaps = computed(() => localDevice.value?.capabilities.filter(c => c.type === 'devices.capabilities.toggle') ?? [])

function hexToRgb(hex: string): { r: number; g: number; b: number } {
  const m = hex.replace('#', '')
  const n = parseInt(m.length === 3 ? m.split('').map(c => c + c).join('') : m, 16)
  return { r: (n >> 16) & 0xff, g: (n >> 8) & 0xff, b: n & 0xff }
}

function rgbToHsv(r: number, g: number, b: number): { h: number; s: number; v: number } {
  const rf = r / 255, gf = g / 255, bf = b / 255
  const max = Math.max(rf, gf, bf), min = Math.min(rf, gf, bf)
  const d = max - min
  let h = 0
  if (d !== 0) {
    if (max === rf) h = ((gf - bf) / d) % 6
    else if (max === gf) h = (bf - rf) / d + 2
    else h = (rf - gf) / d + 4
    h = Math.round(h * 60)
    if (h < 0) h += 360
  }
  const s = max === 0 ? 0 : Math.round((d / max) * 100)
  const v = Math.round(max * 100)
  return { h, s, v }
}

function hsvToRgb(h: number, s: number, v: number): { r: number; g: number; b: number } {
  const sf = s / 100, vf = v / 100
  const c = vf * sf
  const hh = h / 60
  const x = c * (1 - Math.abs((hh % 2) - 1))
  let r = 0, g = 0, b = 0
  if (hh < 1) { r = c; g = x }
  else if (hh < 2) { r = x; g = c }
  else if (hh < 3) { g = c; b = x }
  else if (hh < 4) { g = x; b = c }
  else if (hh < 5) { r = x; b = c }
  else { r = c; b = x }
  const m = vf - c
  return {
    r: Math.round((r + m) * 255),
    g: Math.round((g + m) * 255),
    b: Math.round((b + m) * 255),
  }
}

const colorValue = computed<string>(() => {
  const v = colorCap.value?.state?.value
  const hex = (n: number) => n.toString(16).padStart(2, '0')
  if (typeof v === 'number') return '#' + v.toString(16).padStart(6, '0')
  if (v && typeof v === 'object') {
    const o = v as Record<string, unknown>
    if ('h' in o && 's' in o && 'v' in o) {
      const rgb = hsvToRgb(Number(o.h), Number(o.s), Number(o.v))
      return '#' + hex(rgb.r) + hex(rgb.g) + hex(rgb.b)
    }
    if ('r' in o && 'g' in o && 'b' in o) {
      return '#' + hex(Number(o.r)) + hex(Number(o.g)) + hex(Number(o.b))
    }
  }
  return '#ffffff'
})

function setColor(hex: string) {
  if (!colorCap.value) return
  const rgb = hexToRgb(hex)
  const hsv = rgbToHsv(rgb.r, rgb.g, rgb.b)
  setCap(colorCap.value, hsv, 'hsv')
}

function rangeMin(): number {
  const p = rangeCap.value?.parameters as Record<string, unknown> | undefined
  if (p?.range && typeof p.range === 'object' && 'min' in (p.range as Record<string, unknown>)) {
    return Number((p.range as { min: number }).min)
  }
  return 0
}
function rangeMax(): number {
  const p = rangeCap.value?.parameters as Record<string, unknown> | undefined
  if (p?.range && typeof p.range === 'object' && 'max' in (p.range as Record<string, unknown>)) {
    return Number((p.range as { max: number }).max)
  }
  return 100
}
function rangeStep(): number {
  const p = rangeCap.value?.parameters as Record<string, unknown> | undefined
  if (p?.range && typeof p.range === 'object' && 'precision' in (p.range as Record<string, unknown>)) {
    return Number((p.range as { precision: number }).precision)
  }
  return 1
}
function rangeInstance(): string {
  const p = rangeCap.value?.parameters as Record<string, unknown> | undefined
  if (p?.instance) return String(p.instance)
  return rangeCap.value?.state?.instance ?? 'brightness'
}

function modeOptions(cap: Capability): string[] {
  const p = cap.parameters as Record<string, unknown> | undefined
  if (Array.isArray(p?.modes)) {
    return (p.modes as Array<{ value?: string }>).map(m => m.value ?? '').filter(Boolean)
  }
  return []
}
function modeInstance(cap: Capability): string {
  const p = cap.parameters as Record<string, unknown> | undefined
  return (p?.instance as string) ?? cap.state?.instance ?? 'mode'
}
function toggleInstance(cap: Capability): string {
  const p = cap.parameters as Record<string, unknown> | undefined
  return (p?.instance as string) ?? cap.state?.instance ?? 'toggle'
}

function lastSeenLabel(): string {
  if (!localDevice.value?.last_seen) return '—'
  const d = new Date(localDevice.value.last_seen)
  if (isNaN(d.getTime())) return '—'
  return d.toLocaleString()
}
</script>

<template>
  <Teleport to="body">
    <transition name="modal-fade">
      <div v-if="modelValue && localDevice" class="backdrop" @click.self="close">
        <div class="modal">
          <header class="modal-header">
            <div class="header-left">
              <div class="header-icon" :class="{ 'header-icon--on': onOffCap?.state?.value }">
                <svg width="20" height="20" fill="none" stroke="currentColor" stroke-width="1.8" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M21 16V8a2 2 0 00-1-1.73l-7-4a2 2 0 00-2 0l-7 4A2 2 0 003 8v8a2 2 0 001 1.73l7 4a2 2 0 002 0l7-4A2 2 0 0021 16z" />
                </svg>
              </div>
              <div class="header-text">
                <div class="header-name">{{ localDevice.name || 'Unnamed device' }}</div>
                <div class="header-meta">
                  <span :class="['status-pill', isOnline ? 'status-pill--on' : 'status-pill--off']">
                    {{ isOnline ? 'online' : 'offline' }}
                  </span>
                  <span>{{ deviceTypeLabel(localDevice.type) }}</span>
                  <span v-if="localDevice.room">· {{ localDevice.room }}</span>
                </div>
              </div>
            </div>
            <button class="close-btn" @click="close" aria-label="Close">
              <svg width="16" height="16" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </header>

          <div class="modal-body">
            <!-- Capabilities -->
            <section v-if="onOffCap || colorCap || rangeCap || modeCaps.length || toggleCaps.length" class="section">
              <h3 class="section-title">Controls</h3>

              <!-- On/Off -->
              <div v-if="onOffCap" class="control-row">
                <div class="control-label">{{ CAPABILITY_LABELS['devices.capabilities.on_off'] }}</div>
                <button
                  class="toggle"
                  :class="{ 'toggle--on': onOffCap.state?.value }"
                  :disabled="capPending.has(onOffCap.id)"
                  @click="setCap(onOffCap, !onOffCap.state?.value, 'on')"
                >
                  <span class="toggle-knob" />
                </button>
              </div>

              <!-- Brightness / Range -->
              <div v-if="rangeCap" class="control-row control-row--column">
                <div class="control-header">
                  <div class="control-label">{{ CAPABILITY_LABELS['devices.capabilities.range'] }}</div>
                  <div class="control-value">{{ Number(rangeCap.state?.value ?? rangeMin()) }}</div>
                </div>
                <input
                  type="range"
                  :min="rangeMin()"
                  :max="rangeMax()"
                  :step="rangeStep()"
                  :value="Number(rangeCap.state?.value ?? rangeMin())"
                  :disabled="capPending.has(rangeCap.id)"
                  @change="(e) => setCap(rangeCap!, Number((e.target as HTMLInputElement).value), rangeInstance())"
                  class="slider"
                />
              </div>

              <!-- Color -->
              <div v-if="colorCap" class="control-row control-row--column">
                <div class="control-header">
                  <div class="control-label">{{ CAPABILITY_LABELS['devices.capabilities.color_setting'] }}</div>
                  <div class="control-value">{{ colorValue.toUpperCase() }}</div>
                </div>
                <div class="color-row">
                  <input
                    type="color"
                    :value="colorValue"
                    :disabled="capPending.has(colorCap.id)"
                    @change="(e) => setColor((e.target as HTMLInputElement).value)"
                    class="color-picker"
                  />
                  <div class="color-presets">
                    <button
                      v-for="preset in ['#ffe4b5', '#ffffff', '#6366f1', '#f87171', '#4ade80', '#fbbf24']"
                      :key="preset"
                      class="color-preset"
                      :style="{ background: preset }"
                      :title="preset"
                      @click="setColor(preset)"
                    />
                  </div>
                </div>
              </div>

              <!-- Mode -->
              <div v-for="cap in modeCaps" :key="cap.id" class="control-row control-row--column">
                <div class="control-header">
                  <div class="control-label">{{ modeInstance(cap) }}</div>
                </div>
                <div class="mode-buttons">
                  <button
                    v-for="opt in modeOptions(cap)"
                    :key="opt"
                    class="mode-btn"
                    :class="{ 'mode-btn--active': cap.state?.value === opt }"
                    :disabled="capPending.has(cap.id)"
                    @click="setCap(cap, opt, modeInstance(cap))"
                  >
                    {{ opt }}
                  </button>
                </div>
              </div>

              <!-- Toggle -->
              <div v-for="cap in toggleCaps" :key="cap.id" class="control-row">
                <div class="control-label">{{ toggleInstance(cap) }}</div>
                <button
                  class="toggle"
                  :class="{ 'toggle--on': cap.state?.value }"
                  :disabled="capPending.has(cap.id)"
                  @click="setCap(cap, !cap.state?.value, toggleInstance(cap))"
                >
                  <span class="toggle-knob" />
                </button>
              </div>
            </section>

            <section v-else class="section section--empty">
              <div class="empty-cap">This device has no controllable capabilities.</div>
            </section>

            <!-- Properties -->
            <section v-if="localDevice.properties?.length" class="section">
              <h3 class="section-title">Sensors</h3>
              <div class="properties">
                <div v-for="prop in localDevice.properties" :key="prop.id" class="property">
                  <div class="prop-label">{{ prop.state?.instance ?? prop.type.split('.').pop() }}</div>
                  <div class="prop-value">{{ prop.state?.value ?? '—' }}</div>
                </div>
              </div>
            </section>

            <!-- Settings -->
            <section class="section">
              <h3 class="section-title">Settings</h3>
              <div class="field">
                <label>Name</label>
                <input v-model="editName" type="text" placeholder="e.g. Bedroom lamp" />
              </div>
              <div class="field">
                <label>Room</label>
                <input v-model="editRoom" type="text" placeholder="e.g. Bedroom" />
              </div>
              <div class="field">
                <label>Description</label>
                <input v-model="editDescription" type="text" placeholder="—" />
              </div>
              <button class="btn-primary btn-block" @click="saveMeta" :disabled="renaming">
                {{ renaming ? 'Saving…' : 'Save' }}
              </button>
            </section>

            <!-- Info -->
            <section class="section">
              <h3 class="section-title">Info</h3>
              <div class="info-row">
                <span class="info-key">Device ID</span>
                <span class="info-val">{{ localDevice.id }}</span>
              </div>
              <div class="info-row">
                <span class="info-key">UID</span>
                <span class="info-val">{{ localDevice.device_uid }}</span>
              </div>
              <div class="info-row">
                <span class="info-key">Last seen</span>
                <span class="info-val">{{ lastSeenLabel() }}</span>
              </div>
            </section>

            <!-- Danger -->
            <section class="section section--danger">
              <button
                class="btn-danger"
                :class="{ 'btn-danger--confirm': confirmDelete }"
                @click="handleDelete"
                :disabled="deleting"
              >
                {{ deleting
                   ? 'Deleting…'
                   : confirmDelete
                     ? 'Confirm delete?'
                     : 'Delete device' }}
              </button>
            </section>
          </div>
        </div>
      </div>
    </transition>
  </Teleport>
</template>

<style scoped>
.backdrop {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.7);
  display: flex;
  align-items: flex-start;
  justify-content: center;
  padding: 48px 16px;
  z-index: 2000;
  backdrop-filter: blur(4px);
  overflow-y: auto;
}

.modal {
  background: #111;
  border: 1px solid #1e1e1e;
  border-radius: 16px;
  width: 100%;
  max-width: 420px;
  box-shadow: 0 24px 64px rgba(0,0,0,0.6);
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 16px 12px;
  border-bottom: 1px solid #1a1a1a;
}
.header-left {
  display: flex;
  align-items: center;
  gap: 10px;
  min-width: 0;
}
.header-icon {
  width: 38px;
  height: 38px;
  border-radius: 10px;
  background: #1a1a1a;
  border: 1px solid #222;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #737373;
  flex-shrink: 0;
}
.header-icon--on { background: #2c2860; color: #c4b5fd; border-color: #3a3470; }
.header-text { display: flex; flex-direction: column; gap: 2px; min-width: 0; }
.header-name {
  font-size: 14px;
  font-weight: 600;
  color: #e5e5e5;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.header-meta {
  display: flex;
  gap: 6px;
  align-items: center;
  font-size: 12px;
  color: #525252;
  flex-wrap: wrap;
}
.status-pill {
  font-size: 10px;
  padding: 2px 6px;
  border-radius: 5px;
  text-transform: uppercase;
  font-weight: 600;
  letter-spacing: 0.05em;
}
.status-pill--on { background: #052e16; color: #4ade80; }
.status-pill--off { background: #1f1f1f; color: #525252; }

.close-btn {
  width: 28px;
  height: 28px;
  border-radius: 7px;
  border: none;
  background: transparent;
  color: #555;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}
.close-btn:hover { background: #1a1a1a; color: #a3a3a3; }

.modal-body {
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 16px;
  max-height: calc(100vh - 180px);
  overflow-y: auto;
}

.section {
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding-bottom: 14px;
  border-bottom: 1px solid #1a1a1a;
}
.section:last-child { border-bottom: none; padding-bottom: 0; }
.section--empty .empty-cap {
  font-size: 12px;
  color: #525252;
  text-align: center;
  padding: 16px;
  background: #181818;
  border-radius: 8px;
}
.section--danger { border-bottom: none; padding-top: 4px; }

.section-title {
  font-size: 11px;
  font-weight: 600;
  color: #525252;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  margin: 0;
}

.control-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 8px;
}
.control-row--column {
  flex-direction: column;
  align-items: stretch;
}
.control-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}
.control-label {
  font-size: 13px;
  color: #a3a3a3;
}
.control-value {
  font-size: 12px;
  color: #737373;
  font-variant-numeric: tabular-nums;
}

.toggle {
  width: 38px;
  height: 22px;
  border-radius: 11px;
  background: #2a2a2a;
  border: none;
  position: relative;
  cursor: pointer;
  padding: 0;
  transition: background 0.15s;
  flex-shrink: 0;
}
.toggle--on { background: #6366f1; }
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

.slider {
  width: 100%;
  height: 6px;
  appearance: none;
  background: #1a1a1a;
  border: 1px solid #222;
  border-radius: 3px;
  outline: none;
  padding: 0;
}
.slider::-webkit-slider-thumb {
  appearance: none;
  width: 16px;
  height: 16px;
  border-radius: 50%;
  background: #c4b5fd;
  cursor: pointer;
  border: 2px solid #1a1a1a;
}
.slider::-moz-range-thumb {
  width: 16px;
  height: 16px;
  border-radius: 50%;
  background: #c4b5fd;
  cursor: pointer;
  border: 2px solid #1a1a1a;
}

.color-row { display: flex; gap: 12px; align-items: center; }
.color-picker {
  width: 44px;
  height: 36px;
  padding: 0;
  border: 1px solid #222;
  background: #1a1a1a;
  border-radius: 8px;
  cursor: pointer;
}
.color-presets { display: flex; gap: 6px; flex-wrap: wrap; }
.color-preset {
  width: 22px;
  height: 22px;
  border-radius: 50%;
  border: 1px solid #222;
  cursor: pointer;
  padding: 0;
  transition: transform 0.1s;
}
.color-preset:hover { transform: scale(1.15); }

.mode-buttons { display: flex; flex-wrap: wrap; gap: 6px; }
.mode-btn {
  background: #1a1a1a;
  border: 1px solid #222;
  color: #a3a3a3;
  padding: 6px 10px;
  border-radius: 7px;
  cursor: pointer;
  font-size: 12px;
  font-family: inherit;
  transition: background 0.12s, color 0.12s, border-color 0.12s;
}
.mode-btn:hover:not(:disabled) { background: #222; color: #e5e5e5; }
.mode-btn--active {
  background: #2c2860;
  color: #c4b5fd;
  border-color: #3a3470;
}
.mode-btn:disabled { opacity: 0.5; cursor: not-allowed; }

.properties {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 8px;
}
.property {
  display: flex;
  flex-direction: column;
  gap: 2px;
  padding: 10px 12px;
  background: #1a1a1a;
  border: 1px solid #222;
  border-radius: 8px;
}
.prop-label { font-size: 11px; color: #525252; }
.prop-value {
  font-size: 14px;
  color: #e5e5e5;
  font-variant-numeric: tabular-nums;
  font-weight: 600;
}

.field { display: flex; flex-direction: column; gap: 4px; }
.field label {
  font-size: 11px;
  color: #525252;
}
.field input {
  background: #1a1a1a;
  border: 1px solid #222;
  border-radius: 7px;
  padding: 8px 10px;
  color: #e5e5e5;
  font-size: 13px;
  font-family: inherit;
  outline: none;
  transition: border-color 0.15s;
}
.field input:focus { border-color: #3a3a3a; }

.info-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
  font-size: 12px;
}
.info-key { color: #525252; }
.info-val {
  color: #a3a3a3;
  font-variant-numeric: tabular-nums;
  font-family: ui-monospace, monospace;
  font-size: 11px;
  text-overflow: ellipsis;
  overflow: hidden;
  max-width: 60%;
  white-space: nowrap;
}

.btn-primary, .btn-danger {
  padding: 9px 18px;
  border-radius: 8px;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  font-family: inherit;
  border: none;
  transition: background 0.15s;
}
.btn-block { width: 100%; }
.btn-primary {
  background: #6366f1;
  color: #fff;
}
.btn-primary:hover:not(:disabled) { background: #4f46e5; }
.btn-primary:disabled { opacity: 0.6; cursor: not-allowed; }

.btn-danger {
  background: transparent;
  color: #f87171;
  border: 1px solid #401818;
  width: 100%;
}
.btn-danger:hover:not(:disabled) { background: #1f0f0f; }
.btn-danger--confirm {
  background: #7f1d1d;
  color: #fff;
  border-color: #991b1b;
}
.btn-danger--confirm:hover:not(:disabled) { background: #991b1b; }
.btn-danger:disabled { opacity: 0.6; cursor: not-allowed; }

.modal-fade-enter-active, .modal-fade-leave-active { transition: opacity 0.2s; }
.modal-fade-enter-active .modal, .modal-fade-leave-active .modal { transition: transform 0.2s, opacity 0.2s; }
.modal-fade-enter-from, .modal-fade-leave-to { opacity: 0; }
.modal-fade-enter-from .modal, .modal-fade-leave-to .modal { transform: scale(0.96) translateY(8px); }

@media (max-width: 480px) {
  .backdrop { padding: 16px 8px; }
  .modal { max-width: 100%; }
}
</style>
