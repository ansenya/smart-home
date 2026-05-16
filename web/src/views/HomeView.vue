<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'

const authStore = useAuthStore()
const router = useRouter()
</script>

<template>
  <div class="home">
    <div class="hero">
      <div class="hero-icon">
        <svg width="32" height="32" fill="currentColor" viewBox="0 0 20 20">
          <path d="M10.707 2.293a1 1 0 00-1.414 0l-7 7a1 1 0 001.414 1.414L4 10.414V17a1 1 0 001 1h2a1 1 0 001-1v-2a1 1 0 011-1h2a1 1 0 011 1v2a1 1 0 001 1h2a1 1 0 001-1v-6.586l.293.293a1 1 0 001.414-1.414l-7-7z"/>
        </svg>
      </div>
      <h1 class="hero-title">Добро пожаловать в<br><span class="gradient-text">Hiphome</span></h1>
      <p class="hero-sub">Управляйте устройствами, создавайте автоматизации и общайтесь с AI-ассистентом умного дома.</p>
      <div class="hero-actions">
        <button v-if="!authStore.isAuthenticated" class="btn-primary" @click="authStore.login()">
          Войти
        </button>
        <button v-else class="btn-primary" @click="router.push('/chats')">
          Открыть чаты
        </button>
        <button v-if="authStore.isAuthenticated" class="btn-secondary" @click="router.push('/devices')">
          Мои устройства
        </button>
      </div>
    </div>

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
    title: 'AI-ассистент',
    desc: 'Управляйте домом голосом и текстом через умного ассистента.',
    icon: 'M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z',
    bg: 'linear-gradient(135deg, #6366f1, #8b5cf6)',
  },
  {
    title: 'Устройства',
    desc: 'Контролируйте все умные устройства из единого интерфейса.',
    icon: 'M9 3H5a2 2 0 00-2 2v4m6-6h10a2 2 0 012 2v4M9 3v18m0 0h10a2 2 0 002-2V9M9 21H5a2 2 0 01-2-2V9m0 0h18',
    bg: 'linear-gradient(135deg, #0ea5e9, #6366f1)',
  },
  {
    title: 'Автоматизация',
    desc: 'Создавайте сценарии, которые работают без вашего участия.',
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
  gap: 64px;
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
  .hero-title { font-size: 28px; }
}
</style>
