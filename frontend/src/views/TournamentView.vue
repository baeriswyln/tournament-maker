<template>
  <div class="p-4 lg:p-8 pb-24 lg:pb-8 max-w-[100vw] lg:max-w-5xl mx-auto">
    <!-- Back + header -->
    <div class="flex items-center gap-3 mb-6 animate-slide-up">
      <button @click="$router.push('/')" class="btn-ghost btn-sm p-2">
        <i class="pi pi-arrow-left" />
      </button>
      <div class="flex-1 min-w-0">
        <div class="flex items-center gap-2 flex-wrap">
          <h1 class="font-display text-xl lg:text-2xl font-bold truncate whitespace-break-spaces" style="color: var(--text-primary)">
            {{ store.current?.name || '…' }}
          </h1>
          <span :class="badgeClass">{{ stateLabel }}</span>
        </div>
        <p v-if="store.current?.start_date" class="text-xs mt-0.5" style="color: var(--text-muted)">
          {{ formatDate(store.current.start_date) }} · {{ store.current?.format?.name }}
        </p>
      </div>
      <RouterLink v-if="isOwner" :to="`/tournaments/${tournamentId}/manage`" class="btn-secondary btn-sm hidden md:inline-flex">
        <i class="pi pi-cog" /> {{ t('manage') }}
      </RouterLink>
    </div>

    <!-- Mobile FAB -->
    <RouterLink :to="`/tournaments/${tournamentId}/manage`" v-if="isOwner"
                class="md:hidden fixed bottom-20 right-4 z-20 w-14 h-14 rounded-2xl flex items-center justify-center
                bg-secondary btn-secondary">

      <i class="pi pi-cog text-xl" />
    </RouterLink>

    <div v-if="loading" class="flex justify-center py-20">
      <i class="pi pi-spin pi-spinner text-3xl" style="color: var(--brand)" />
    </div>

    <template v-else>
      <!-- Tabs -->
      <div class="flex gap-1 p-1 rounded-xl mb-6 overflow-x-auto animate-slide-up flex-wrap"
        style="background: var(--bg-subtle); animation-delay: 0.05s">
        <button v-for="tab in availableTabs" :key="tab.id"
          @click="activeTab = tab.id"
          class="flex-1 min-w-fit py-2 px-3 rounded-lg text-sm font-display font-medium transition-all whitespace-nowrap"
          :style="activeTab === tab.id
            ? 'background: var(--surface); color: var(--text-primary); box-shadow: 0 1px 4px rgba(0,0,0,0.08)'
            : 'color: var(--text-muted)'"
        >
          {{ tab.label }}
          <span v-if="tab.count !== undefined"
            class="ml-1.5 text-xs px-1.5 py-0.5 rounded-full"
            style="background: var(--border); color: var(--text-muted)">{{ tab.count }}</span>
        </button>
      </div>

      <!-- Tab content -->
      <Transition name="fade" mode="out-in">

        <!-- Upcoming -->
        <div v-if="activeTab === 'upcoming'" key="upcoming" class="space-y-3 animate-slide-up">
          <div v-if="!upcomingEncounters.length" class="card p-8 text-center">
            <i class="pi pi-clock text-3xl mb-3 block" style="color: var(--text-faint)" />
            <p class="text-sm" style="color: var(--text-muted)">{{ t('no_upcoming') }}</p>
          </div>
          <div v-if="upcomingEncounters.length" class="pb-3">
            <label class="input-label">{{ t('search_encounters') }}</label>
            <input v-model="upcomingEncounterSearch" type="text" class="input py-2 text-sm"/>
          </div>
          <div v-if="!filteredUpcomingEncounters.length" class="card p-8 text-center text-sm" style="color: var(--text-muted)">
            {{ t('no_encounters') }}
          </div>
          <div v-else-if="!filteredUpcomingEncounters.length" class="card p-8 text-center">
            <i class="pi pi-check-circle text-3xl mb-3 block" style="color: var(--text-faint)" />
            <p class="text-sm" style="color: var(--text-muted)">{{ t('encounter_not_found') }}</p>
          </div>
          <EncounterRow v-for="e in filteredUpcomingEncounters"
            :key="`${e.partyA}-${e.partyB}-${e.round}`"
            :encounter="e" :participants="store.participants"
            :current-user-id="auth.user?.id" :tournament-id="tournamentId"
            :is-owner="false" @score-changed="refreshEncounters" />
        </div>

        <!-- Past -->
        <div v-else-if="activeTab === 'past'" key="past" class="space-y-3 animate-slide-up">
          <div v-if="pastEncounters.length" class="pb-3">
            <label class="input-label">{{ t('search_encounters') }}</label>
            <input v-model="pastEncounterSearch" type="text" class="input py-2 text-sm"/>
          </div>
          <div v-if="!pastEncounters.length" class="card p-8 text-center">
            <i class="pi pi-check-circle text-3xl mb-3 block" style="color: var(--text-faint)" />
            <p class="text-sm" style="color: var(--text-muted)">{{ t('no_past') }}</p>
          </div>
          <div v-else-if="!filteredPastEncounters.length" class="card p-8 text-center">
            <i class="pi pi-check-circle text-3xl mb-3 block" style="color: var(--text-faint)" />
            <p class="text-sm" style="color: var(--text-muted)">{{ t('encounter_not_found') }}</p>
          </div>
          <EncounterRow v-for="e in filteredPastEncounters"
            :key="`${e.partyA}-${e.partyB}-${e.round}`"
            :encounter="e" :participants="store.participants"
            :current-user-id="auth.user?.id" :tournament-id="tournamentId"
            :is-owner="false" @score-changed="refreshEncounters" />
        </div>

        <!-- Scoreboard -->
        <div v-else-if="activeTab === 'scoreboard'" key="scoreboard" class="animate-slide-up">
          <Scoreboard :entries="store.scoreboard" :participants="store.participants" :tournament-id="tournamentId" />
        </div>

        <!-- Players -->
        <div v-else-if="activeTab === 'players'" key="players" class="animate-slide-up">
          <div class="card overflow-hidden">
            <div v-for="(p, i) in store.participants" :key="p.id"
              class="flex items-center gap-3 px-4 py-3"
              :style="i > 0 ? 'border-top: 1px solid var(--border)' : ''">
              <div class="w-8 h-8 rounded-lg flex items-center justify-center font-display font-bold text-sm shrink-0"
                style="background: var(--bg-subtle); color: var(--text-muted)">
                {{ p.visible_name?.[0]?.toUpperCase() || '?' }}
              </div>
              <div class="flex-1 min-w-0">
                <p class="text-sm font-medium truncate" style="color: var(--text-primary)">{{ p.visible_name }}</p>
              </div>
            </div>
          </div>
        </div>

      </Transition>
    </template>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, RouterLink } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useTournamentStore } from '@/stores/tournament'
import { useI18n } from '@/i18n/index.js'
import EncounterRow from '@/components/EncounterRow.vue'
import Scoreboard from '@/components/Scoreboard.vue'
import { format } from 'date-fns'
import {useToast} from "primevue/usetoast";
import router from "../router/index.js";
import { useApiError } from '@/composables/useApiError.js'

const route = useRoute()
const auth = useAuthStore()
const toast = useToast()
const store = useTournamentStore()
const { t } = useI18n()
const { resolveApiError } = useApiError()
const tournamentId = route.params.id
const activeTab = ref('upcoming')
const loading = ref(true)
const upcomingEncounterSearch = ref('')
const pastEncounterSearch = ref('')

const isOwner = computed(() =>
  store.owning?.some(x => x.id === tournamentId) || store.current?.owner_id === auth.user?.id
)
const upcomingEncounters = computed(() => store.myEncounters.filter(e => e.state !== 'final'))
const pastEncounters = computed(() => store.myEncounters.filter(e => e.state === 'final'))
const showScoreboard = computed(() => store.current?.visibility?.scoreboard !== false)

const stateLabel = computed(() => {
  const map = { draft: t('state_draft'), open: t('state_open'), live: t('state_live'), closed: t('state_closed') }
  return map[store.current?.state] || store.current?.state
})
const badgeClass = computed(() => {
  const map = { draft:'badge-draft', open:'badge-open', live:'badge-live', closed:'badge-closed' }
  return `badge ${map[store.current?.state] || 'badge-draft'}`
})

const availableTabs = computed(() => {
  const tabs = [
    { id: 'upcoming', label: t('upcoming'), count: upcomingEncounters.value.length },
    { id: 'past', label: t('past'), count: pastEncounters.value.length },
  ]
  if (showScoreboard.value) tabs.push({ id: 'scoreboard', label: t('scoreboard') })
  tabs.push({ id: 'players', label: t('players'), count: store.participants.length })
  return tabs
})

function getParticipantName(id) {
  return store.participants.find(p => p.id === id)?.visible_name.toLowerCase()
}

const filteredUpcomingEncounters = computed(() => {
  const terms = upcomingEncounterSearch.value.trim().toLowerCase().split(' ')

  if (!terms.length) return upcomingEncounters.value
  return upcomingEncounters.value.filter(e => terms.every(t =>
      getParticipantName(e.partyA).includes(t) ||
      getParticipantName(e.partyB).includes(t) ||
      e.round.toLowerCase().includes(t)
  ))
})

const filteredPastEncounters = computed(() => {
  const terms = pastEncounterSearch.value.trim().toLowerCase().split(' ')

  if (!terms.length) return pastEncounters.value
  return pastEncounters.value.filter(e => terms.every(t =>
      getParticipantName(e.partyA).includes(t) ||
      getParticipantName(e.partyB).includes(t) ||
      e.round.toLowerCase().includes(t)
  ))
})

function formatDate(d) {
  try { return format(new Date(d), 'MMMM d, yyyy') } catch { return '' }
}

async function refreshEncounters() {
  await Promise.all([
    store.fetchEncounters(tournamentId, auth.user?.id),
    store.fetchParticipants(tournamentId),
    store.fetchScoreboard(tournamentId),
  ])
}

onMounted(async () => {
  loading.value = true
  try {
    await Promise.all([
      store.fetchOne(tournamentId),
      store.fetchParticipants(tournamentId),
      store.fetchEncounters(tournamentId, auth.user?.id),
      store.fetchScoreboard(tournamentId),
    ])
  } catch (e) {
    const { title, detail } = resolveApiError(e)
    toast.add({ severity: 'error', summary: title, detail, life: 4000 })
    router.push('/')
  }
  loading.value = false
  if (!upcomingEncounters.value.length) {
    activeTab.value = showScoreboard.value ? 'scoreboard' : 'players'
  }
})
</script>
