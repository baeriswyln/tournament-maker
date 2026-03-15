<template>
  <RouterLink :to="`/tournaments/${tournament.id}`"
    class="card block group transition-all duration-200"
    style="text-decoration: none;"
    :style="hoverStyle"
    @mouseenter="hovered = true" @mouseleave="hovered = false"
  >
    <div class="p-4">
      <div class="flex items-start justify-between gap-3 mb-3">
        <div class="flex-1 min-w-0">
          <h3 class="font-display font-semibold truncate transition-colors whitespace-break-spaces"
            :style="{ color: hovered ? 'var(--brand)' : 'var(--text-primary)' }">
            {{ tournament.name || t('unknown') }}
          </h3>
          <p class="text-xs mt-0.5 truncate" style="color: var(--text-muted)">
            {{ tournament.format?.name || '—' }}
          </p>
        </div>
        <span :class="badgeClass">{{ stateLabel }}</span>
      </div>

      <div class="flex items-center gap-4 text-xs" style="color: var(--text-faint)">
        <span class="flex items-center gap-1.5">
          <i class="pi pi-users" />
          {{ tournament.nb_participants || '?' }}
        </span>
        <span v-if="tournament.start_date" class="flex items-center gap-1.5">
          <i class="pi pi-calendar" />
          {{ formatDate(tournament.start_date) }}
        </span>
      </div>

      <div v-if="owned" class="mt-3 pt-3 flex items-center justify-between"
        style="border-top: 1px solid var(--border)">
        <span class="text-xs flex items-center gap-1" style="color: var(--amber)">
          <i class="pi pi-crown" />
          {{ t('owner') }}
        </span>
        <span @click.prevent.stop="$router.push(`/tournaments/${tournament.id}/manage`)"
          class="text-xs flex items-center gap-1 cursor-pointer transition-colors"
          style="color: var(--brand)"
          @mouseenter="e => e.target.style.color = 'var(--brand-light)'"
          @mouseleave="e => e.target.style.color = 'var(--brand)'"
        >
          {{ t('manage') }} <i class="pi pi-arrow-right text-xs" />
        </span>
      </div>
    </div>
  </RouterLink>
</template>

<script setup>
import { ref, computed } from 'vue'
import { RouterLink } from 'vue-router'
import { useI18n } from '@/i18n/index.js'
import { format } from 'date-fns'

const props = defineProps({
  tournament: { type: Object, required: true },
  owned: { type: Boolean, default: false }
})

const { t } = useI18n()
const hovered = ref(false)

const hoverStyle = computed(() =>
  hovered.value
    ? { borderColor: 'rgba(var(--brand-rgb), 0.35)', boxShadow: '0 4px 20px rgba(var(--brand-rgb), 0.08)' }
    : {}
)

const stateLabel = computed(() => {
  const map = { draft: t('state_draft'), open: t('state_open'), live: t('state_live'), closed: t('state_closed') }
  return map[props.tournament.state] || props.tournament.state
})

const badgeClass = computed(() => {
  const map = { draft: 'badge-draft', open: 'badge-open', live: 'badge-live', closed: 'badge-closed' }
  return `badge ${map[props.tournament.state] || 'badge-draft'}`
})

function formatDate(d) {
  try { return format(new Date(d), 'MMM d, yyyy') } catch { return '' }
}
</script>
