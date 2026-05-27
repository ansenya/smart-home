<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { push } from 'notivue'
import FloatingMenu from '@/components/FloatingMenu.vue'
import DeviceCard from '@/components/DeviceCard.vue'
import DeviceDetailModal from '@/components/modals/DeviceDetailModal.vue'
import { listDevices, setCapability, subscribeDeviceStream } from '@/api/devices'
import type { DeviceEvent, StreamSubscription } from '@/api/devices'
import type { Device } from '@/types/device'

const devices = ref<Device[]>([])
const loading = ref(true)
const error = ref<string | null>(null)
const pendingIds = ref<Set<string>>(new Set())
const search = ref('')
const isLive = ref(false)

const selectedDevice = ref<Device | null>(null)
const detailOpen = ref(false)

let streamSub: StreamSubscription | null = null

const groupedByRoom = computed<{ room: string; devices: Device[] }[]>(() => {
  const filtered = search.value.trim()
    ? devices.value.filter(d => {
        const q = search.value.trim().toLowerCase()
        return (
          d.name?.toLowerCase().includes(q) ||
          d.room?.toLowerCase().includes(q) ||
          d.type?.toLowerCase().includes(q)
        )
      })
    : devices.value

  const groups = new Map<string, Device[]>()
  for (const d of filtered) {
    const r = d.room?.trim() || 'No room'
    if (!groups.has(r)) groups.set(r, [])
    groups.get(r)!.push(d)
  }
  return Array.from(groups.entries())
    .sort((a, b) => {
      if (a[0] === 'No room') return 1
      if (b[0] === 'No room') return -1
      return a[0].localeCompare(b[0])
    })
    .map(([room, devices]) => ({ room, devices }))
})

const stats = computed(() => {
  const total = devices.value.length
  const online = devices.value.filter(d => {
    if (!d.last_seen) return false
    return Date.now() - new Date(d.last_seen).getTime() < 60_000
  }).length
  const on = devices.value.filter(d => {
    const cap = d.capabilities.find(c => c.type === 'devices.capabilities.on_off')
    return Boolean(cap?.state?.value)
  }).length
  return { total, online, on }
})

async function loadDevices() {
  loading.value = true
  error.value = null
  try {
    const r = await listDevices()
    devices.value = r.data.devices ?? []
  } catch (e: unknown) {
    const err = e as { response?: { status?: number } }
    if (err.response?.status === 401) {
      error.value = 'auth'
    } else {
      error.value = 'Failed to load devices'
    }
  } finally {
    loading.value = false
  }
}

async function handleToggle(device: Device, next: boolean) {
  if (pendingIds.value.has(device.id)) return
  pendingIds.value.add(device.id)

  // Optimistic update
  const cap = device.capabilities.find(c => c.type === 'devices.capabilities.on_off')
  const prev = cap?.state?.value
  if (cap) cap.state = { ...(cap.state ?? {}), value: next }

  try {
    await setCapability(device.id, 'devices.capabilities.on_off', { value: next })
  } catch {
    if (cap) cap.state = { ...(cap.state ?? {}), value: prev }
    push.error({ title: 'Failed', message: 'Command was not delivered' })
  } finally {
    pendingIds.value.delete(device.id)
  }
}

function openDetail(device: Device) {
  selectedDevice.value = device
  detailOpen.value = true
}

function handleDeviceUpdated(updated: Device) {
  const idx = devices.value.findIndex(d => d.id === updated.id)
  if (idx >= 0) devices.value[idx] = updated
  selectedDevice.value = updated
}

function handleDeviceDeleted(id: string) {
  devices.value = devices.value.filter(d => d.id !== id)
  detailOpen.value = false
  selectedDevice.value = null
}

function handleStreamEvent(ev: DeviceEvent) {
  const device = devices.value.find(d => d.id === ev.device_id)
  if (!device) return

  switch (ev.type) {
    case 'capability_state': {
      if (!ev.capability) break
      const cap = device.capabilities.find(c => c.type === ev.capability)
      if (!cap) break
      const payload = ev.payload as Record<string, unknown> | null
      cap.state = {
        ...(cap.state ?? {}),
        ...(payload && typeof payload === 'object' ? payload : { value: payload }),
      }
      // Mirror state into selectedDevice if it's the same one
      if (selectedDevice.value?.id === device.id) {
        const sCap = selectedDevice.value.capabilities.find(c => c.id === cap.id)
        if (sCap) sCap.state = { ...cap.state }
      }
      break
    }
    case 'property_state': {
      if (!ev.property) break
      const prop = device.properties.find(p => p.type === ev.property)
      if (!prop) break
      const payload = ev.payload as Record<string, unknown> | null
      prop.state = {
        ...(prop.state ?? {}),
        ...(payload && typeof payload === 'object' ? payload : { value: payload }),
      }
      if (selectedDevice.value?.id === device.id) {
        const sProp = selectedDevice.value.properties.find(p => p.id === prop.id)
        if (sProp) sProp.state = { ...prop.state }
      }
      break
    }
    case 'device_status':
    case 'device_state': {
      device.last_seen = new Date().toISOString()
      break
    }
  }
}

function startStream() {
  if (streamSub) return
  streamSub = subscribeDeviceStream(
    (ev) => {
      isLive.value = true
      handleStreamEvent(ev)
    },
    () => {
      // Browsers auto-reconnect EventSource; just flag live=false until next event
      isLive.value = false
    },
  )
}

function stopStream() {
  if (streamSub) {
    streamSub.close()
    streamSub = null
  }
  isLive.value = false
}

onMounted(async () => {
  await loadDevices()
  if (!error.value) startStream()
})
onUnmounted(stopStream)
</script>

<template>
  <div class="devices-view">
    <header class="hero">
      <div class="hero-text">
        <h1>
          My devices
          <span
            v-if="!loading && !error"
            class="live-pill"
            :class="{ 'live-pill--active': isLive }"
            :title="isLive ? 'Connected to live stream' : 'Connecting...'"
          >
            <span class="live-dot" />
            {{ isLive ? 'live' : 'idle' }}
          </span>
        </h1>
        <p>Manage every connected device in one place.</p>
      </div>

      <div class="hero-stats" v-if="!loading && !error">
        <div class="stat">
          <div class="stat-value">{{ stats.total }}</div>
          <div class="stat-label">total</div>
        </div>
        <div class="stat stat--online">
          <div class="stat-value">{{ stats.online }}</div>
          <div class="stat-label">online</div>
        </div>
        <div class="stat stat--on">
          <div class="stat-value">{{ stats.on }}</div>
          <div class="stat-label">active</div>
        </div>
      </div>
    </header>

    <div v-if="devices.length > 0" class="search-row">
      <div class="search-input">
        <svg width="14" height="14" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-4.35-4.35M11 19a8 8 0 100-16 8 8 0 000 16z" />
        </svg>
        <input
          v-model="search"
          type="text"
          placeholder="Search by name, room or type"
        />
        <button v-if="search" class="clear-btn" @click="search = ''" aria-label="Clear">
          <svg width="12" height="12" fill="none" stroke="currentColor" stroke-width="2.5" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>
    </div>

    <main class="content">
      <div v-if="loading" class="loading-grid">
        <div v-for="i in 6" :key="i" class="skeleton" />
      </div>

      <div v-else-if="error === 'auth'" class="empty">
        <div class="empty-icon">
          <svg width="28" height="28" fill="none" stroke="currentColor" stroke-width="1.5" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
          </svg>
        </div>
        <div class="empty-title">Sign in required</div>
        <div class="empty-sub">Sign in to see your devices.</div>
      </div>

      <div v-else-if="error" class="empty">
        <div class="empty-icon empty-icon--error">
          <svg width="28" height="28" fill="none" stroke="currentColor" stroke-width="1.5" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v2m0 4h.01M5.07 19h13.86c1.54 0 2.5-1.67 1.73-3L13.73 4a2 2 0 00-3.46 0L3.34 16c-.77 1.33.19 3 1.73 3z" />
          </svg>
        </div>
        <div class="empty-title">Error</div>
        <div class="empty-sub">{{ error }}</div>
        <button class="btn-primary" @click="loadDevices">Retry</button>
      </div>

      <div v-else-if="devices.length === 0" class="empty">
        <div class="empty-icon">
          <svg width="28" height="28" fill="none" stroke="currentColor" stroke-width="1.5" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" d="M21 16V8a2 2 0 00-1-1.73l-7-4a2 2 0 00-2 0l-7 4A2 2 0 003 8v8a2 2 0 001 1.73l7 4a2 2 0 002 0l7-4A2 2 0 0021 16z" />
          </svg>
        </div>
        <div class="empty-title">No devices yet</div>
        <div class="empty-sub">Tap the round button at the bottom-right and choose &ldquo;Add device&rdquo; to add your first one.</div>
      </div>

      <div v-else-if="groupedByRoom.length === 0" class="empty">
        <div class="empty-title">Nothing found</div>
        <div class="empty-sub">Try a different search.</div>
      </div>

      <template v-else>
        <section v-for="group in groupedByRoom" :key="group.room" class="room-section">
          <h2 class="room-title">{{ group.room }}</h2>
          <div class="cards-grid">
            <DeviceCard
              v-for="device in group.devices"
              :key="device.id"
              :device="device"
              :pending="pendingIds.has(device.id)"
              @open="openDetail"
              @toggle="handleToggle"
            />
          </div>
        </section>
      </template>
    </main>

    <DeviceDetailModal
      v-model="detailOpen"
      :device="selectedDevice"
      @updated="handleDeviceUpdated"
      @deleted="handleDeviceDeleted"
      @refresh="loadDevices"
    />

    <FloatingMenu />
  </div>
</template>

<style scoped>
.devices-view {
  min-height: 100%;
  padding: 48px 24px 96px;
  display: flex;
  flex-direction: column;
  gap: 32px;
  max-width: 1100px;
  margin: 0 auto;
  width: 100%;
  background: #0f0f0f;
}

.hero {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  gap: 24px;
  flex-wrap: wrap;
}

.hero-text { display: flex; flex-direction: column; gap: 8px; }
.hero h1 {
  font-size: 28px;
  font-weight: 700;
  color: #e5e5e5;
  letter-spacing: -0.02em;
  margin: 0;
  display: flex;
  align-items: center;
  gap: 10px;
}

.live-pill {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  font-size: 10px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  padding: 3px 8px;
  border-radius: 6px;
  background: #1a1a1a;
  border: 1px solid #222;
  color: #525252;
  vertical-align: middle;
}
.live-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #525252;
}
.live-pill--active {
  background: #052e16;
  border-color: #0f4d2a;
  color: #4ade80;
}
.live-pill--active .live-dot {
  background: #4ade80;
  box-shadow: 0 0 6px rgba(74, 222, 128, 0.6);
  animation: pulse 1.6s ease-in-out infinite;
}
@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.55; }
}
.hero p {
  font-size: 14px;
  color: #555;
  margin: 0;
}

.hero-stats { display: flex; gap: 8px; }
.stat {
  padding: 10px 14px;
  background: #111;
  border: 1px solid #1e1e1e;
  border-radius: 10px;
  display: flex;
  flex-direction: column;
  align-items: center;
  min-width: 64px;
}
.stat-value {
  font-size: 18px;
  font-weight: 700;
  color: #e5e5e5;
}
.stat-label {
  font-size: 11px;
  color: #555;
  text-transform: lowercase;
  margin-top: 2px;
}
.stat--online .stat-value { color: #4ade80; }
.stat--on .stat-value { color: #c4b5fd; }

.search-row {
  display: flex;
  gap: 12px;
}
.search-input {
  flex: 1;
  position: relative;
  display: flex;
  align-items: center;
  background: #111;
  border: 1px solid #1e1e1e;
  border-radius: 10px;
  padding: 0 12px;
  color: #555;
  transition: border-color 0.15s;
}
.search-input:focus-within { border-color: #2a2a2a; }
.search-input input {
  flex: 1;
  background: transparent;
  border: none;
  outline: none;
  padding: 10px 8px;
  color: #e5e5e5;
  font-size: 13px;
  font-family: inherit;
}
.search-input input::placeholder { color: #404040; }
.clear-btn {
  background: transparent;
  border: none;
  color: #555;
  cursor: pointer;
  width: 20px;
  height: 20px;
  border-radius: 5px;
  display: flex;
  align-items: center;
  justify-content: center;
}
.clear-btn:hover { background: #1a1a1a; color: #a3a3a3; }

.content {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.room-section {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.room-title {
  font-size: 12px;
  font-weight: 600;
  color: #525252;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  margin: 0 0 4px 4px;
}

.cards-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
  gap: 12px;
}

.loading-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
  gap: 12px;
}
.skeleton {
  height: 140px;
  border-radius: 14px;
  background: linear-gradient(90deg, #111 0%, #181818 50%, #111 100%);
  background-size: 200% 100%;
  animation: skeleton-shimmer 1.4s linear infinite;
}
@keyframes skeleton-shimmer {
  0% { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}

.empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  background: #111;
  border: 1px solid #1e1e1e;
  border-radius: 14px;
  text-align: center;
  gap: 8px;
}
.empty-icon {
  width: 56px;
  height: 56px;
  border-radius: 14px;
  background: #1a1a1a;
  border: 1px solid #222;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #525252;
  margin-bottom: 8px;
}
.empty-icon--error { background: #1f0f0f; border-color: #401818; color: #f87171; }
.empty-title { font-size: 15px; font-weight: 600; color: #e5e5e5; }
.empty-sub {
  font-size: 13px;
  color: #555;
  max-width: 380px;
  line-height: 1.55;
}

.btn-primary {
  padding: 8px 18px;
  border-radius: 8px;
  background: #6366f1;
  color: #fff;
  font-size: 13px;
  font-weight: 500;
  border: none;
  cursor: pointer;
  margin-top: 12px;
  font-family: inherit;
}
.btn-primary:hover { background: #4f46e5; }

@media (max-width: 640px) {
  .devices-view { padding: 24px 16px 96px; }
  .hero { align-items: flex-start; }
  .hero h1 { font-size: 22px; }
  .hero-stats { width: 100%; justify-content: space-between; }
  .stat { flex: 1; }
  .cards-grid, .loading-grid {
    grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
  }
}
</style>
