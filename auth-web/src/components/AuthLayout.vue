<template>
  <div class="min-h-screen flex bg-white">

    <!-- Left: form panel -->
    <div class="w-full lg:w-[480px] xl:w-[540px] flex flex-col flex-shrink-0">
      <div class="px-8 sm:px-14 pt-8 pb-4">
        <div class="flex items-center gap-2.5">
          <div class="w-8 h-8 rounded-xl bg-indigo-600 flex items-center justify-center flex-shrink-0">
            <svg class="w-4 h-4 text-white" fill="currentColor" viewBox="0 0 20 20">
              <path d="M10.707 2.293a1 1 0 00-1.414 0l-7 7a1 1 0 001.414 1.414L4 10.414V17a1 1 0 001 1h2a1 1 0 001-1v-2a1 1 0 011-1h2a1 1 0 011 1v2a1 1 0 001 1h2a1 1 0 001-1v-6.586l.293.293a1 1 0 001.414-1.414l-7-7z"/>
            </svg>
          </div>
          <span class="font-semibold text-slate-900 tracking-tight text-base">Smart Home</span>
        </div>
      </div>

      <div class="flex-1 flex items-center px-8 sm:px-14 pb-16 pt-4">
        <div class="w-full max-w-sm mx-auto">
          <slot/>
        </div>
      </div>
    </div>

    <!-- Right: landing panel -->
    <div class="hidden lg:flex flex-1 relative overflow-hidden" style="background: linear-gradient(145deg, #0f172a 0%, #1e1b4b 55%, #0f172a 100%)">
      <!-- dot grid -->
      <div class="absolute inset-0" style="background-image: radial-gradient(circle, rgba(255,255,255,0.07) 1px, transparent 1px); background-size: 28px 28px;"></div>
      <!-- glow top-right -->
      <div class="absolute -top-48 -right-48 w-[520px] h-[520px] rounded-full bg-indigo-600 opacity-20" style="filter: blur(120px);"></div>
      <!-- glow bottom-left -->
      <div class="absolute -bottom-48 -left-24 w-[440px] h-[440px] rounded-full bg-violet-600 opacity-15" style="filter: blur(100px);"></div>

      <div class="relative z-10 flex flex-col justify-center px-14 xl:px-20 py-16 text-white w-full">

        <!-- badge -->
        <div class="inline-flex items-center gap-2 px-3 py-1.5 rounded-full text-sm mb-8 self-start"
             style="background: rgba(99,102,241,0.15); border: 1px solid rgba(99,102,241,0.3); color: #a5b4fc;">
          <span class="w-1.5 h-1.5 rounded-full bg-indigo-400" style="animation: pulse 2s infinite;"></span>
          Система умного дома
        </div>

        <h1 class="text-4xl xl:text-5xl font-bold leading-tight mb-5">
          Управляйте домом<br>
          <span style="background: linear-gradient(90deg, #818cf8, #c084fc); -webkit-background-clip: text; -webkit-text-fill-color: transparent; background-clip: text;">
            из одного места
          </span>
        </h1>

        <p class="text-base xl:text-lg leading-relaxed mb-14 max-w-md" style="color: #94a3b8;">
          Подключайте устройства, создавайте сценарии автоматизации и следите за состоянием вашего дома в реальном времени.
        </p>

        <!-- features -->
        <div class="space-y-3 max-w-md">
          <div v-for="f in features" :key="f.title"
               class="flex items-start gap-4 p-4 rounded-2xl"
               style="background: rgba(255,255,255,0.04); border: 1px solid rgba(255,255,255,0.07);">
            <div class="w-10 h-10 rounded-xl flex items-center justify-center flex-shrink-0" :style="{ background: f.bg }">
              <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" stroke-width="1.8" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" :d="f.icon"/>
              </svg>
            </div>
            <div>
              <div class="font-medium text-white mb-0.5 text-sm">{{ f.title }}</div>
              <div class="text-sm leading-relaxed" style="color: #64748b;">{{ f.desc }}</div>
            </div>
          </div>
        </div>

      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const features = [
  {
    title: 'Автоматизация',
    desc: 'Создавайте умные сценарии — дом реагирует на ваш распорядок дня.',
    icon: 'M13 10V3L4 14h7v7l9-11h-7z',
    bg: 'linear-gradient(135deg, #6366f1, #8b5cf6)',
  },
  {
    title: 'Управление устройствами',
    desc: 'Освещение, климат, безопасность — всё под контролем в одном интерфейсе.',
    icon: 'M9 3H5a2 2 0 00-2 2v4m6-6h10a2 2 0 012 2v4M9 3v18m0 0h10a2 2 0 002-2V9M9 21H5a2 2 0 01-2-2V9m0 0h18',
    bg: 'linear-gradient(135deg, #0ea5e9, #6366f1)',
  },
  {
    title: 'Мониторинг в реальном времени',
    desc: 'Следите за состоянием дома из любой точки мира.',
    icon: 'M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z',
    bg: 'linear-gradient(135deg, #10b981, #0ea5e9)',
  },
]
</script>
