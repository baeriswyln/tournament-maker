<template>
  <div class="min-h-screen flex" :data-theme="theme">

    <!-- ── Sidebar (desktop) ── -->
    <aside
        class="hidden lg:flex flex-col w-64 fixed inset-y-0 left-0 z-30"
        style="background: var(--surface); border-right: 1px solid var(--border);"
    >
      <!-- Brand -->
      <div class="flex items-center gap-3 px-5 py-4" style="border-bottom: 1px solid var(--border);">
        <div class="w-9 h-9 rounded-xl flex items-center justify-center shrink-0"
             style="background: rgba(var(--brand-rgb), 0.12); border: 1px solid rgba(var(--brand-rgb), 0.25);">
          <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"
               style="color: var(--brand)">
            <path stroke-linecap="round" stroke-linejoin="round"
                  d="M9 12l2 2 4-4M7.835 4.697a3.42 3.42 0 001.946-.806 3.42 3.42 0 014.438 0 3.42 3.42 0 001.946.806 3.42 3.42 0 013.138 3.138 3.42 3.42 0 00.806 1.946 3.42 3.42 0 010 4.438 3.42 3.42 0 00-.806 1.946 3.42 3.42 0 01-3.138 3.138 3.42 3.42 0 00-1.946.806 3.42 3.42 0 01-4.438 0 3.42 3.42 0 00-1.946-.806 3.42 3.42 0 01-3.138-3.138 3.42 3.42 0 00-.806-1.946 3.42 3.42 0 010-4.438 3.42 3.42 0 00.806-1.946 3.42 3.42 0 013.138-3.138z"/>
          </svg>
        </div>
        <div class="leading-tight">
          <p class="font-display font-bold text-sm" style="color: var(--text-primary)">Tournament</p>
          <p class="font-display font-bold text-sm" style="color: var(--brand)">Maker</p>
        </div>
      </div>

      <!-- Nav -->
      <nav class="flex-1 p-3 space-y-0.5 overflow-y-auto">
        <RouterLink to="/" class="nav-item" :class="{ active: route.path === '/' }">
          <i class="pi pi-home text-base"/>
          {{ t('nav_home') }}
        </RouterLink>
        <RouterLink v-if="auth.isAdmin" to="/tournaments/new" class="nav-item"
                    :class="{ active: route.path === '/tournaments/new' }">
          <i class="pi pi-plus-circle text-base"/>
          {{ t('new_tournament') }}
        </RouterLink>
        <RouterLink to="/profile" class="nav-item" :class="{ active: route.path === '/profile' }">
          <i class="pi pi-user text-base"/>
          {{ t('nav_profile') }}
        </RouterLink>
      </nav>

      <!-- Footer: theme + lang + user -->
      <div class="p-3 space-y-2" style="border-top: 1px solid var(--border);">
        <!-- Controls row -->
        <div class="flex items-center gap-1 px-1">
          <ThemeToggle/>
          <LanguageSwitcher/>
        </div>
        <!-- User card -->
        <div class="flex items-center gap-2.5 px-3 py-2.5 rounded-xl" style="background: var(--bg-subtle);">
          <div class="w-8 h-8 rounded-lg flex items-center justify-center font-display font-bold text-sm shrink-0"
               style="background: rgba(var(--brand-rgb), 0.15); color: var(--brand)">
            {{ userInitial }}
          </div>
          <div class="flex-1 min-w-0">
            <p class="text-sm font-display font-semibold truncate" style="color: var(--text-primary)">
              {{ auth.user?.visible_name || 'Guest' }}</p>
            <p class="text-xs truncate" style="color: var(--text-muted)">
              {{ auth.user?.username ? `@${auth.user.username}` : t('participant_only') }}</p>
          </div>
          <button @click="handleLogout" class="btn-ghost btn-sm p-1.5 rounded-lg" :title="t('sign_out')">
            <i class="pi pi-sign-out text-sm" style="color: var(--text-muted)"/>
          </button>
        </div>
      </div>
    </aside>

    <!-- ── Mobile top bar ── -->
    <header
        class="lg:hidden fixed top-0 inset-x-0 z-30 flex items-center justify-between px-4 py-3"
        style="background: var(--surface); border-bottom: 1px solid var(--border); backdrop-filter: blur(12px);"
    >
      <div class="flex items-center gap-2">
        <div class="w-8 h-8 rounded-lg flex items-center justify-center"
             style="background: rgba(var(--brand-rgb), 0.12); border: 1px solid rgba(var(--brand-rgb), 0.2);">
          <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"
               style="color: var(--brand)">
            <path stroke-linecap="round" stroke-linejoin="round"
                  d="M9 12l2 2 4-4M7.835 4.697a3.42 3.42 0 001.946-.806 3.42 3.42 0 014.438 0 3.42 3.42 0 001.946.806 3.42 3.42 0 013.138 3.138 3.42 3.42 0 00.806 1.946 3.42 3.42 0 010 4.438 3.42 3.42 0 00-.806 1.946 3.42 3.42 0 01-3.138 3.138 3.42 3.42 0 00-1.946.806 3.42 3.42 0 01-4.438 0 3.42 3.42 0 00-1.946-.806 3.42 3.42 0 01-3.138-3.138 3.42 3.42 0 00-.806-1.946 3.42 3.42 0 010-4.438 3.42 3.42 0 00.806-1.946 3.42 3.42 0 013.138-3.138z"/>
          </svg>
        </div>
        <span class="font-display font-bold text-sm" style="color: var(--text-primary)">Tournament Maker</span>
      </div>
      <div class="flex items-center gap-1">
        <ThemeToggle/>
        <LanguageSwitcher/>
        <RouterLink to="/profile"
                    class="w-8 h-8 rounded-lg flex items-center justify-center font-display font-bold text-sm ml-1"
                    style="background: rgba(var(--brand-rgb), 0.15); color: var(--brand)">
          {{ userInitial }}
        </RouterLink>
      </div>
    </header>

    <!-- ── Main content ── -->
    <main class="flex-1 lg:ml-64 pt-[57px] lg:pt-0">
      <RouterView v-slot="{ Component, route: r }">
        <Transition name="page" mode="out-in">
          <component :is="Component" :key="r.path"/>
        </Transition>
      </RouterView>
    </main>

    <!-- ── Mobile bottom nav ── -->
    <nav
        class="lg:hidden fixed bottom-0 inset-x-0 z-30 flex"
        style="background: var(--surface); border-top: 1px solid var(--border); backdrop-filter: blur(12px);"
    >
      <RouterLink v-for="item in mobileNavItems.filter(e => e.visible)" :key="item.to" :to="item.to"
                  class="flex-1 flex flex-col items-center gap-0.5 py-3 text-xs font-display font-medium transition-colors"
                  :style="isActive(item) ? 'color: var(--brand)' : 'color: var(--text-muted)'"
      >
        <i :class="`pi pi-${item.icon} text-xl`"/>
        <span class="text-[10px]">{{ item.label }}</span>
      </RouterLink>
    </nav>
  </div>
</template>

<script setup>
import {computed} from 'vue'
import {RouterView, RouterLink, useRoute, useRouter} from 'vue-router'
import {useAuthStore} from '@/stores/auth'
import {useTheme} from '@/composables/useTheme.js'
import {useI18n} from '@/i18n/index.js'
import ThemeToggle from '@/components/ThemeToggle.vue'
import LanguageSwitcher from '@/components/LanguageSwitcher.vue'

const auth = useAuthStore()
const route = useRoute()
const router = useRouter()
const {theme} = useTheme()
const {t} = useI18n()

const userInitial = computed(() => (auth.user?.visible_name?.[0] || '?').toUpperCase())

const mobileNavItems = computed(() => [
  {to: '/', label: t('nav_home'), icon: 'home', visible: true},
  {to: '/tournaments/new', label: t('new'), icon: 'plus-circle', visible: auth.isAdmin},
  {to: '/profile', label: t('nav_profile'), icon: 'user', visible: true},
])

function isActive(item) {
  if (item.to === '/') return route.path === '/'
  return route.path.startsWith(item.to)
}

function handleLogout() {
  auth.logout()
  router.push('/login')
}
</script>
