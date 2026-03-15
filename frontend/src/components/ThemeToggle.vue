<template>
  <button
    @click="toggle()"
    class="btn-ghost btn-sm p-2 rounded-xl relative overflow-hidden"
    :title="theme === 'dark' ? 'Switch to light mode' : 'Switch to dark mode'"
    :aria-label="theme === 'dark' ? 'Switch to light mode' : 'Switch to dark mode'"
  >
    <Transition name="theme-icon" mode="out-in">
      <i v-if="theme === 'dark'" key="sun" class="pi pi-sun text-base" style="color: var(--amber)" />
      <i v-else key="moon" class="pi pi-moon text-base" style="color: var(--text-muted)" />
    </Transition>
  </button>
</template>

<script setup>
import { useTheme } from '@/composables/useTheme.js'
const { theme, toggle } = useTheme()
</script>

<style scoped>
.theme-icon-enter-active,
.theme-icon-leave-active {
  transition: opacity 0.15s ease, transform 0.15s ease;
}
.theme-icon-enter-from { opacity: 0; transform: rotate(-30deg) scale(0.7); }
.theme-icon-leave-to   { opacity: 0; transform: rotate(30deg) scale(0.7); }
</style>
