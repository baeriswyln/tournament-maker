<template>
  <div class="card overflow-hidden">
    <!-- Header -->
    <div class="grid grid-cols-14 sm:grid-cols-10 px-2 py-2.5"
      style="border-bottom: 1px solid var(--border); background: var(--bg-subtle);">
      <div class="col-span-2 sm:col-span-1 text-center input-label normal-case tracking-normal" style="font-size: 0.7rem; margin-bottom: 0">#</div>
      <div class="col-span-8 sm:col-span-5 input-label normal-case tracking-normal" style="font-size: 0.7rem; margin-bottom: 0">{{ t('players') }}</div>
      <div class="col-span-2 text-center input-label normal-case tracking-normal hidden sm:block" style="font-size: 0.7rem; margin-bottom: 0">W/D/L</div>
      <div class="col-span-2 sm:col-span-1 text-center input-label normal-case tracking-normal" style="font-size: 0.7rem; margin-bottom: 0">+/-</div>
      <div class="col-span-2 sm:col-span-1 text-center input-label normal-case tracking-normal font-bold" style="font-size: 0.7rem; margin-bottom: 0">{{ t('pts') }}</div>
    </div>

    <!-- Empty state -->
    <div v-if="!entries.length" class="p-10 text-center">
      <i class="pi pi-chart-bar text-3xl mb-3 block" style="color: var(--text-faint)" />
      <p class="text-sm" style="color: var(--text-muted)">{{ t('no_scores') }}</p>
    </div>

    <!-- Rows -->
    <div v-for="(p, i) in entries" :key="p.id">

      <!-- Stats row -->
      <div
        class="grid grid-cols-14 sm:grid-cols-10 px-2 py-3 items-center cursor-pointer transition-colors select-none"
        :style="{
          borderTop: i > 0 ? '1px solid var(--border)' : 'none',
          background: selectedId === p.id
            ? 'rgba(var(--brand-rgb), 0.06)'
            : i === 0 ? 'rgba(var(--brand-rgb), 0.04)' : 'transparent'
        }"
        @click="toggleParticipant(p.id)"
      >
        <!-- Rank -->
        <div class="col-span-2 sm:col-span-1 flex justify-center">
          <span v-if="i === 0" class="text-lg">🥇</span>
          <span v-else-if="i === 1" class="text-lg">🥈</span>
          <span v-else-if="i === 2" class="text-lg">🥉</span>
          <span v-else class="font-display font-bold text-sm" style="color: var(--text-faint)">{{ i + 1 }}</span>
        </div>

        <!-- Name -->
        <div class="col-span-8 sm:col-span-5 flex items-center gap-2 min-w-0">
          <div class="w-7 h-7 rounded-lg flex items-center justify-center text-xs font-display font-bold shrink-0"
            :style="selectedId === p.id
              ? 'background: rgba(var(--brand-rgb), 0.2); color: var(--brand)'
              : i === 0
                ? 'background: rgba(var(--brand-rgb), 0.15); color: var(--brand)'
                : 'background: var(--bg-subtle); color: var(--text-muted)'"
          >
            {{ p.visible_name?.[0]?.toUpperCase() || '?' }}
          </div>
          <span class="text-sm font-medium truncate"
            :style="{ color: i < 3 ? 'var(--text-primary)' : 'var(--text-secondary)', fontWeight: i === 0 ? '700' : '500' }">
            {{ p.visible_name || '?' }}
          </span>
          <i class="pi text-xs ml-auto shrink-0 transition-transform"
            :class="selectedId === p.id ? 'pi-chevron-up' : 'pi-chevron-down'"
            :style="{ color: selectedId === p.id ? 'var(--brand)' : 'var(--text-faint)' }" />
        </div>

        <!-- W/D/L -->
        <div class="col-span-2 text-center text-sm hidden sm:block" style="color: var(--text-muted)">
          <span style="color: var(--success-muted, #4ade80)">{{ p.won }}</span>
          <span class="mx-0.5">/</span>
          <span>{{ p.draw }}</span>
          <span class="mx-0.5">/</span>
          <span style="color: var(--danger-muted, #f87171)">{{ p.lost }}</span>
        </div>

        <!-- Score diff -->
        <div class="col-span-2 sm:col-span-1 text-center text-sm"
          :style="{ color: p.score_diff > 0 ? 'var(--success-muted, #4ade80)' : p.score_diff < 0 ? 'var(--danger-muted, #f87171)' : 'var(--text-muted)' }">
          {{ p.score_diff > 0 ? '+' : '' }}{{ p.score_diff }}
        </div>

        <!-- Points -->
        <div class="col-span-2 sm:col-span-1 text-center">
          <span class="font-display font-bold text-base"
            :style="{ color: i === 0 ? 'var(--brand)' : 'var(--text-primary)' }">
            {{ p.points }}
          </span>
        </div>
      </div>

      <!-- Expanded encounters -->
      <Transition name="expand">
        <div v-if="selectedId === p.id"
          class="overflow-hidden"
          style="border-top: 1px solid var(--border); background: var(--bg-subtle);">

          <div v-if="loadingId === p.id" class="flex justify-center py-5">
            <i class="pi pi-spin pi-spinner text-sm" style="color: var(--brand)" />
          </div>

          <div v-else class="overflow-y-auto" style="max-height: 240px;">

            <div v-if="!encounterCache[p.id]?.length" class="py-5 text-center text-sm" style="color: var(--text-muted)">
              {{ t('no_encounters') }}
            </div>

            <template v-else>
              <!-- Section label if both exist -->
              <template v-for="section in getSections(p.id)" :key="section.label">
                <div v-if="getSections(p.id).length > 1"
                  class="px-4 py-1.5 text-xs font-display font-semibold sticky top-0"
                  style="color: var(--text-muted); background: var(--bg-subtle); border-bottom: 1px solid var(--border)">
                  {{ section.label }}
                </div>
                <div v-for="(e, ei) in section.encounters" :key="`${e.partyA}-${e.partyB}-${e.round}`"
                  class="flex items-center gap-3 px-4 py-2.5"
                  :style="ei > 0 || getSections(p.id).length > 1 ? 'border-top: 1px solid var(--border)' : ''">

                  <!-- Round -->
                  <span class="text-xs font-display px-2 py-0.5 rounded-full shrink-0"
                    style="background: var(--border); color: var(--text-muted)">
                    {{ e.round }}
                  </span>

                  <!-- Opponent -->
                  <div class="flex items-center gap-1.5 flex-1 min-w-0">
                    <div class="w-5 h-5 rounded flex items-center justify-center text-xs font-bold shrink-0"
                      style="background: var(--surface); color: var(--text-muted)">
                      {{ getOpponentName(e, p.id)?.[0]?.toUpperCase() || '?' }}
                    </div>
                    <span class="text-sm truncate" style="color: var(--text-secondary)">
                      {{ getOpponentName(e, p.id) || t('unknown') }}
                    </span>
                  </div>

                  <!-- Score relative to this player -->
                  <div class="flex items-center gap-1 shrink-0 font-display font-bold text-sm" style="font-variant-numeric: tabular-nums">
                    <span class="inline-block text-center" style="min-width: 1.75rem" :style="{ color: scoreColor(e, p.id, true) }">{{ myScore(e, p.id) ?? '—' }}</span>
                    <span style="color: var(--text-faint)">–</span>
                    <span class="inline-block text-center" style="min-width: 1.75rem" :style="{ color: scoreColor(e, p.id, false) }">{{ oppScore(e, p.id) ?? '—' }}</span>
                  </div>

                  <!-- Result badge -->
                  <span class="text-xs font-display font-semibold py-0.5 rounded text-center shrink-0" style="min-width: 2rem"
                    :style="resultStyle(e, p.id)">
                    {{ resultLabel(e, p.id) }}
                  </span>
                </div>
              </template>
            </template>
          </div>
        </div>
      </Transition>

    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useI18n } from '@/i18n/index.js'
import { encountersApi } from '@/api'

const props = defineProps({
  entries: { type: Array, default: () => [] },
  participants: { type: Array, default: () => [] },
  tournamentId: { type: String, default: null },
})
const { t } = useI18n()

const selectedId = ref(null)
const loadingId = ref(null)
const encounterCache = ref({})

async function toggleParticipant(id) {
  if (selectedId.value === id) {
    selectedId.value = null
    return
  }
  selectedId.value = id
  if (encounterCache.value[id]) return
  loadingId.value = id
  try {
    const { data } = await encountersApi.list(props.tournamentId, id)
    encounterCache.value = { ...encounterCache.value, [id]: data }
  } finally {
    loadingId.value = null
  }
}

function getParticipantName(id) {
  return props.participants.find(p => p.id === id)?.visible_name
}

function getOpponentName(encounter, participantId) {
  const oppId = encounter.partyA === participantId ? encounter.partyB : encounter.partyA
  return getParticipantName(oppId)
}

function myScore(encounter, participantId) {
  return encounter.partyA === participantId ? encounter.scoreA : encounter.scoreB
}

function oppScore(encounter, participantId) {
  return encounter.partyA === participantId ? encounter.scoreB : encounter.scoreA
}

function getResult(encounter, participantId) {
  if (encounter.state !== 'final') return encounter.state
  const me = myScore(encounter, participantId)
  const opp = oppScore(encounter, participantId)
  if (me == null || opp == null) return 'planned'
  if (me > opp) return 'win'
  if (me < opp) return 'loss'
  return 'draw'
}

function resultLabel(encounter, participantId) {
  const r = getResult(encounter, participantId)
  const map = { win: 'W', loss: 'L', draw: 'D', planned: t('state_planned'), voted: t('state_voted') }
  return map[r] ?? r
}

function resultStyle(encounter, participantId) {
  const r = getResult(encounter, participantId)
  if (r === 'win')  return 'background: rgba(74,222,128,0.12); color: #4ade80'
  if (r === 'loss') return 'background: rgba(248,113,113,0.12); color: #f87171'
  if (r === 'draw') return 'background: var(--border); color: var(--text-muted)'
  return 'background: var(--border); color: var(--text-faint)'
}

function scoreColor(encounter, participantId, isMe) {
  if (encounter.state !== 'final') return 'var(--text-muted)'
  const me = myScore(encounter, participantId)
  const opp = oppScore(encounter, participantId)
  if (me == null || opp == null) return 'var(--text-muted)'
  const winning = isMe ? me > opp : opp > me
  const losing  = isMe ? me < opp : opp < me
  if (winning) return 'var(--brand)'
  if (losing)  return 'var(--text-faint)'
  return 'var(--text-secondary)'
}

function getSections(participantId) {
  const all = encounterCache.value[participantId] || []
  const upcoming = all.filter(e => e.state !== 'final')
  const past = all.filter(e => e.state === 'final')
  const sections = []
  if (past.length)     sections.push({ label: t('past'), encounters: past })
  if (upcoming.length) sections.push({ label: t('upcoming'), encounters: upcoming })
  return sections
}
</script>

<style scoped>
.expand-enter-active, .expand-leave-active {
  transition: max-height 0.3s cubic-bezier(0.16, 1, 0.3, 1), opacity 0.25s ease;
  max-height: 300px;
  opacity: 1;
  overflow: hidden;
}
.expand-enter-from, .expand-leave-to {
  max-height: 0;
  opacity: 0;
}
</style>