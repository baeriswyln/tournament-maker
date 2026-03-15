<template>
  <div class="p-4 lg:p-8 pb-24 lg:pb-8 max-w-6xl mx-auto">
    <!-- Header -->
    <div class="flex items-start justify-between mb-8 animate-slide-up">
      <div>
        <h1 class="font-display text-2xl lg:text-3xl font-bold" style="color: var(--text-primary)">
          {{ t('dashboard_greeting', { name: firstName }) }}
        </h1>
        <p class="text-sm mt-1" style="color: var(--text-muted)">{{ t('dashboard_subtitle') }}</p>
      </div>
      <RouterLink to="/tournaments/join" class="btn-primary hidden lg:inline-flex">
        <i class="pi pi-plus" /> {{ t('join_tournament') }}
      </RouterLink>
    </div>

    <!-- Stats -->
    <div class="grid grid-cols-2 lg:grid-cols-4 gap-3 mb-8 animate-slide-up" style="animation-delay:0.05s">
      <div v-for="stat in stats" :key="stat.label" class="card p-4">
        <p class="input-label mb-1 normal-case tracking-normal text-xs">{{ stat.label }}</p>
        <p class="text-3xl font-display font-bold" :style="{ color: stat.color || 'var(--text-primary)' }">{{ stat.value }}</p>
      </div>
    </div>

    <div v-if="store.loading" class="flex justify-center py-20">
      <i class="pi pi-spin pi-spinner text-3xl" style="color: var(--brand)" />
    </div>

    <template v-else>
      <!-- Owned -->
      <section v-if="store.owning.length" class="mb-8 animate-slide-up" style="animation-delay:0.1s">
        <div class="flex items-center justify-between mb-3">
          <h2 class="font-display text-base font-semibold" style="color: var(--text-secondary)">{{ t('organizing') }}</h2>
          <RouterLink to="/tournaments/new" class="btn-ghost btn-sm" style="color: var(--brand)">
            <i class="pi pi-plus text-xs" /> {{ t('new') }}
          </RouterLink>
        </div>
        <div class="grid gap-3 sm:grid-cols-2 xl:grid-cols-3">
          <TournamentCard v-for="t_ in store.owning" :key="t_.id" :tournament="t_" :owned="true" />
        </div>
      </section>

      <!-- Participating -->
      <section v-if="store.participating.length" class="mb-8 animate-slide-up" style="animation-delay:0.15s">
        <h2 class="font-display text-base font-semibold mb-3" style="color: var(--text-secondary)">{{ t('competing') }}</h2>
        <div class="grid gap-3 sm:grid-cols-2 xl:grid-cols-3">
          <TournamentCard v-for="t_ in store.participating" :key="t_.id" :tournament="t_" :owned="false" />
        </div>
      </section>

      <!-- Empty state -->
      <div v-if="!store.owning.length && !store.participating.length" class="text-center py-20 animate-fade-in">
        <div class="w-20 h-20 rounded-3xl flex items-center justify-center mx-auto mb-6"
          style="background: var(--bg-subtle); border: 1px solid var(--border);">
          <i class="pi pi-trophy text-4xl" style="color: var(--text-faint)" />
        </div>
        <h3 class="font-display text-xl font-semibold mb-2" style="color: var(--text-secondary)">{{ t('no_tournaments_title') }}</h3>
        <p class="text-sm mb-6" style="color: var(--text-muted)">{{ t('no_tournaments_body') }}</p>
        <div class="flex flex-wrap gap-2 justify-center">
        <RouterLink v-if="auth.isAdmin" to="/tournaments/new" class="btn-primary width-100">
          <i class="pi pi-plus" /> {{ t('create_tournament') }}
        </RouterLink>
        <RouterLink to="/tournaments/join" class="btn-primary">
          <i class="pi pi-plus" /> {{ t('join_tournament') }}
        </RouterLink>
        </div>
      </div>
    </template>

    <!-- Mobile FAB -->
    <RouterLink to="/tournaments/join"
      class="lg:hidden fixed bottom-20 right-4 z-20 w-14 h-14 rounded-2xl flex items-center justify-center shadow-lg"
      style="background: var(--brand); color: white; box-shadow: 0 8px 20px rgba(var(--brand-rgb), 0.4); animation: pulse-brand 2s ease-in-out infinite;">
      <i class="pi pi-plus text-xl" />
    </RouterLink>
  </div>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useTournamentStore } from '@/stores/tournament'
import { useI18n } from '@/i18n/index.js'
import TournamentCard from '@/components/TournamentCard.vue'

const auth = useAuthStore()
const store = useTournamentStore()
const { t } = useI18n()

const firstName = computed(() => auth.user?.visible_name?.split(' ')[0] || 'there')
const liveCount = computed(() => [...store.owning, ...store.participating].filter(x => x.state === 'live').length)
const stats = computed(() => [
  { label: t('organizing'), value: store.owning.length },
  { label: t('competing'), value: store.participating.length },
  { label: t('live'), value: liveCount.value, color: 'var(--brand)' },
  { label: t('total'), value: store.owning.length + store.participating.length },
])

onMounted(() => store.fetchAll())
</script>
