<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'
import { listDevices } from '@/api/devices'
import { getChats } from '@/api/chat'
import type { Device } from '@/types/device'

const authStore = useAuthStore()
const router = useRouter()

const devices = ref<Device[]>([])
const chatsCount = ref<number | null>(null)
const loading = ref(false)

async function loadOverview() {
  if (!authStore.isAuthenticated) return
  loading.value = true
  try {
    const [devicesR, chatsR] = await Promise.allSettled([listDevices(), getChats(50)])
    if (devicesR.status === 'fulfilled') {
      devices.value = devicesR.value.data.devices ?? []
    }
    if (chatsR.status === 'fulfilled') {
      chatsCount.value = chatsR.value.data.chats?.length ?? 0
    }
  } finally {
    loading.value = false
  }
}

onMounted(loadOverview)
watch(() => authStore.isAuthenticated, (auth) => { if (auth) loadOverview() })

const stats = (() => ({
  total: () => devices.value.length,
  online: () => devices.value.filter(d => {
    if (!d.last_seen) return false
    return Date.now() - new Date(d.last_seen).getTime() < 60_000
  }).length,
  on: () => devices.value.filter(d => {
    const cap = d.capabilities?.find(c => c.type === 'devices.capabilities.on_off')
    return Boolean(cap?.state?.value)
  }).length,
}))()
</script>

<template>
  <div class="home">
    <div class="hero">
      <div class="hero-icon">
        <svg width="32" height="32" fill="currentColor" viewBox="0 0 20 20">
          <path d="M10.707 2.293a1 1 0 00-1.414 0l-7 7a1 1 0 001.414 1.414L4 10.414V17a1 1 0 001 1h2a1 1 0 001-1v-2a1 1 0 011-1h2a1 1 0 011 1v2a1 1 0 001 1h2a1 1 0 001-1v-6.586l.293.293a1 1 0 001.414-1.414l-7-7z"/>
        </svg>
      </div>
      <h1 class="hero-title">
        <template v-if="authStore.isAuthenticated">
          Welcome back,<br>
          <span class="gradient-text">{{ authStore.user?.name || 'friend' }}</span>
        </template>
        <template v-else>
          Welcome to<br><span class="gradient-text">Hiphome</span>
        </template>
      </h1>
      <p class="hero-sub">
        Control your devices, build automations, and chat with your smart home assistant.
      </p>
      <div class="hero-actions">
        <button v-if="!authStore.isAuthenticated" class="btn-primary" @click="authStore.login()">Sign in</button>
        <template v-else>
          <button class="btn-primary" @click="router.push('/chats')">Open chat</button>
          <button class="btn-secondary" @click="router.push('/devices')">My devices</button>
        </template>
      </div>
    </div>

    <!-- Overview cards for authenticated users -->
    <section v-if="authStore.isAuthenticated && !loading" class="overview">
      <button class="overview-card" @click="router.push('/devices')">
        <div class="overview-header">
          <div class="overview-icon" style="background: linear-gradient(135deg, #0ea5e9, #6366f1)">
            <svg width="16" height="16" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" d="M21 16V8a2 2 0 00-1-1.73l-7-4a2 2 0 00-2 0l-7 4A2 2 0 003 8v8a2 2 0 001 1.73l7 4a2 2 0 002 0l7-4A2 2 0 0021 16z" />
            </svg>
          </div>
          <div class="overview-label">Devices</div>
        </div>
        <div class="overview-stats">
          <div class="stat-big">
            <div class="stat-num">{{ stats.total() }}</div>
            <div class="stat-cap">total</div>
          </div>
          <div class="stat-pair">
            <div class="stat-line">
              <span class="stat-dot stat-dot--online" />
              <span class="stat-val">{{ stats.online() }} online</span>
            </div>
            <div class="stat-line">
              <span class="stat-dot stat-dot--on" />
              <span class="stat-val">{{ stats.on() }} active</span>
            </div>
          </div>
        </div>
      </button>

      <button class="overview-card" @click="router.push('/chats')">
        <div class="overview-header">
          <div class="overview-icon" style="background: linear-gradient(135deg, #6366f1, #8b5cf6)">
            <svg width="16" height="16" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
            </svg>
          </div>
          <div class="overview-label">AI Chats</div>
        </div>
        <div class="overview-stats">
          <div class="stat-big">
            <div class="stat-num">{{ chatsCount ?? '—' }}</div>
            <div class="stat-cap">sessions</div>
          </div>
          <div class="stat-pair">
            <div class="stat-hint">Talk to the assistant and control your home through conversation.</div>
          </div>
        </div>
      </button>
    </section>

    <div class="features">
      <div class="feature-card" v-for="f in features" :key="f.title">
        <div class="feature-icon" :style="{ background: f.bg }">
          <svg width="18" height="18" fill="none" stroke="currentColor" stroke-width="1.8" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" :d="f.icon"/>
          </svg>
        </div>
        <div class="feature-title">{{ f.title }}</div>
        <div class="feature-desc">{{ f.desc }}</div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
const features = [
  {
    title: 'AI Assistant',
    desc: 'Run your home with voice and chat through an AI that understands your devices.',
    icon: 'M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z',
    bg: 'linear-gradient(135deg, #6366f1, #8b5cf6)',
  },
  {
    title: 'Devices',
    desc: 'Control every connected device from one place, with live status updates.',
    icon: 'M9 3H5a2 2 0 00-2 2v4m6-6h10a2 2 0 012 2v4M9 3v18m0 0h10a2 2 0 002-2V9M9 21H5a2 2 0 01-2-2V9m0 0h18',
    bg: 'linear-gradient(135deg, #0ea5e9, #6366f1)',
  },
  {
    title: 'Automations',
    desc: 'Build routines that run on their own, so your home reacts before you ask.',
    icon: 'M13 10V3L4 14h7v7l9-11h-7z',
    bg: 'linear-gradient(135deg, #10b981, #0ea5e9)',
  },
]
</script>

<style scoped>
.home {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 100%;
  padding: 48px 24px;
  gap: 48px;
  overflow-y: auto;
  background: #0f0f0f;
}

.hero {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  gap: 16px;
  max-width: 540px;
}

.hero-icon {
  width: 64px;
  height: 64px;
  border-radius: 18px;
  background: #1a1a1a;
  border: 1px solid #2a2a2a;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #6366f1;
  margin-bottom: 8px;
}

.hero-title {
  font-size: 36px;
  font-weight: 700;
  color: #e5e5e5;
  line-height: 1.2;
  letter-spacing: -0.03em;
}

.gradient-text {
  background: linear-gradient(90deg, #818cf8, #c084fc);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.hero-sub {
  font-size: 15px;
  color: #555;
  line-height: 1.65;
  max-width: 420px;
}

.hero-actions {
  display: flex;
  gap: 10px;
  margin-top: 8px;
  flex-wrap: wrap;
  justify-content: center;
}

.btn-primary, .btn-secondary {
  padding: 10px 22px;
  border-radius: 10px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  font-family: inherit;
  transition: background 0.15s, transform 0.1s;
}

.btn-primary {
  background: #6366f1;
  color: #fff;
  border: 1px solid #6366f1;
}
.btn-primary:hover { background: #4f46e5; }

.btn-secondary {
  background: transparent;
  color: #737373;
  border: 1px solid #2a2a2a;
}
.btn-secondary:hover { background: #1a1a1a; color: #a3a3a3; }

/* Overview */
.overview {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
  width: 100%;
  max-width: 720px;
}
.overview-card {
  background: #111;
  border: 1px solid #1e1e1e;
  border-radius: 14px;
  padding: 18px 20px;
  display: flex;
  flex-direction: column;
  gap: 14px;
  cursor: pointer;
  transition: border-color 0.15s, background 0.15s;
  text-align: left;
  font-family: inherit;
  color: inherit;
}
.overview-card:hover { border-color: #2a2a2a; background: #131313; }

.overview-header { display: flex; align-items: center; gap: 10px; }
.overview-icon {
  width: 32px;
  height: 32px;
  border-radius: 9px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}
.overview-label {
  font-size: 13px;
  font-weight: 600;
  color: #e5e5e5;
}

.overview-stats {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 12px;
}
.stat-big { display: flex; align-items: baseline; gap: 6px; }
.stat-num {
  font-size: 32px;
  font-weight: 700;
  color: #e5e5e5;
  letter-spacing: -0.03em;
  line-height: 1;
}
.stat-cap { font-size: 12px; color: #525252; }

.stat-pair {
  display: flex;
  flex-direction: column;
  gap: 4px;
  align-items: flex-end;
}
.stat-line {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: #737373;
}
.stat-val { color: #a3a3a3; }
.stat-dot { width: 7px; height: 7px; border-radius: 50%; }
.stat-dot--online { background: #4ade80; box-shadow: 0 0 6px rgba(74, 222, 128, 0.5); }
.stat-dot--on { background: #c4b5fd; box-shadow: 0 0 6px rgba(196, 181, 253, 0.5); }
.stat-hint {
  font-size: 11px;
  color: #525252;
  max-width: 140px;
  text-align: right;
  line-height: 1.4;
}

/* Features grid */
.features {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
  max-width: 720px;
  width: 100%;
}

.feature-card {
  background: #111;
  border: 1px solid #1e1e1e;
  border-radius: 14px;
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 10px;
  transition: border-color 0.15s;
}

.feature-card:hover {
  border-color: #2a2a2a;
}

.feature-icon {
  width: 38px;
  height: 38px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.feature-title {
  font-size: 14px;
  font-weight: 600;
  color: #e5e5e5;
}

.feature-desc {
  font-size: 13px;
  color: #555;
  line-height: 1.55;
}

@media (max-width: 640px) {
  .features { grid-template-columns: 1fr; }
  .overview { grid-template-columns: 1fr; }
  .hero-title { font-size: 28px; }
}
</style>
