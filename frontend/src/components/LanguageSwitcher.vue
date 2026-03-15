<template>
  <div class="relative" ref="root">
    <button
      @click="toggleLocale"
      class="btn-ghost btn-sm px-2.5 py-2 rounded-xl flex items-center gap-1.5 font-mono text-xs font-semibold uppercase tracking-wider"
      :style="{ color: 'var(--text-muted)' }"
      aria-label="Change language"
    >
      <span>{{ currentFlag }}</span>
      <span class="hidden sm:inline">{{ locale.toUpperCase() }}</span>
      <i class="pi pi-chevron-down text-xs transition-transform duration-200" :class="{ 'rotate-180': open }" />
    </button>

    <Transition name="dropdown">
      <div v-if="open"
        class="absolute right-0 z-50 rounded-xl overflow-hidden shadow-xl"
        :class="openAbove ? 'bottom-full mb-1.5' : 'top-full mt-1.5'"
        style="background: var(--surface); border: 1px solid var(--border-strong); min-width: 140px;"
      >
        <button
          v-for="loc in availableLocales" :key="loc.code"
          @click="selectLocale(loc.code)"
          class="w-full flex items-center gap-2.5 px-4 py-2.5 text-sm text-left transition-colors"
          :style="locale === loc.code
            ? 'color: var(--brand); background: rgba(var(--brand-rgb), 0.07)'
            : 'color: var(--text-secondary)'"
          :class="{ 'hover:bg-[var(--bg-subtle)]': locale !== loc.code }"
        >
          <span class="text-base">{{ loc.flag }}</span>
          <span class="font-body font-medium">{{ loc.label }}</span>
          <i v-if="locale === loc.code" class="pi pi-check ml-auto text-xs" style="color: var(--brand)" />
        </button>
      </div>
    </Transition>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { onClickOutside } from '@vueuse/core'
import { useI18n } from '@/i18n/index.js'

const { locale, setLocale, availableLocales } = useI18n()
const open = ref(false)
const openAbove = ref(false)
const root = ref(null)

const currentFlag = ref(availableLocales.find(l => l.code === locale.value)?.flag || '🌐')

function toggleLocale() {
  // Measure before opening so we can pick the right direction
  const rect = root.value?.getBoundingClientRect()
  if (rect) {
    // Dropdown height: 3 items × ~44px + 2px border = ~134px
    const dropdownHeight = availableLocales.length * 44 + 4
    const spaceBelow = window.innerHeight - rect.bottom
    openAbove.value = spaceBelow < dropdownHeight && rect.top > dropdownHeight
  }
  open.value = !open.value
}

function selectLocale(code) {
  setLocale(code)
  currentFlag.value = availableLocales.find(l => l.code === code)?.flag || '🌐'
  open.value = false
}

// Close on outside click — graceful fallback if @vueuse/core not available
try {
  onClickOutside(root, () => { open.value = false })
} catch {
  // fallback: close on document click
  document.addEventListener('click', (e) => {
    if (root.value && !root.value.contains(e.target)) open.value = false
  })
}
</script>

<style scoped>
.dropdown-enter-active { animation: dropIn 0.18s cubic-bezier(0.16, 1, 0.3, 1); }
.dropdown-leave-active { animation: dropIn 0.12s reverse ease-in; }
@keyframes dropIn {
  from { opacity: 0; transform: translateY(-6px) scale(0.97); }
  to   { opacity: 1; transform: translateY(0) scale(1); }
}
</style>
