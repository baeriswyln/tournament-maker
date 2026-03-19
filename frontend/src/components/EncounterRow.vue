<template>
  <div class="card p-4">
    <!-- Round + state -->
    <div class="flex items-center justify-between mb-3">
      <span class="input-label normal-case tracking-normal" style="color: var(--text-muted); margin-bottom: 0">
        {{ encounter.round || t('round') }}
      </span>
      <span :class="encounterBadgeClass">{{ encounterStateLabel }}</span>
    </div>

    <!-- Match layout -->
    <div class="flex items-center gap-2">
      <!-- Party A -->
      <div class="flex-1 flex items-center gap-2 min-w-0"
        :style="{ opacity: encounter.state === 'final' && encounter.scoreA < encounter.scoreB ? 0.5 : 1 }">
        <div class="w-7 h-7 rounded-lg flex items-center justify-center text-xs font-display font-bold shrink-0"
          style="background: var(--bg-subtle); color: var(--text-muted)">
          {{ getParticipantName(encounter.partyA)?.[0]?.toUpperCase() || '?' }}
        </div>
        <span class="text-sm font-medium truncate"
          :style="{
            color: encounter.state === 'final' && encounter.scoreA > encounter.scoreB ? 'var(--brand)' : 'var(--text-primary)',
            fontWeight: encounter.state === 'final' && encounter.scoreA > encounter.scoreB ? '700' : '500'
          }">
          {{ getParticipantName(encounter.partyA) || t('unknown') }}
          <span v-if="encounter.partyA === currentUserId" class="text-xs ml-1" style="color: var(--brand)">({{ t('you') }})</span>
        </span>
      </div>

      <!-- Scores -->
      <div class="flex items-center gap-1.5 shrink-0">
        <div class="score-box text-sm" :style="scoreBoxStyle(encounter.scoreA, encounter.scoreB)">
          {{ encounter.scoreA ?? '—' }}
        </div>
        <span class="text-xs font-bold" style="color: var(--border-strong)">vs</span>
        <div class="score-box text-sm" :style="scoreBoxStyle(encounter.scoreB, encounter.scoreA)">
          {{ encounter.scoreB ?? '—' }}
        </div>
      </div>

      <!-- Party B -->
      <div class="flex-1 flex items-center gap-2 justify-end min-w-0"
        :style="{ opacity: encounter.state === 'final' && encounter.scoreB < encounter.scoreA ? 0.5 : 1 }">
        <span class="text-sm font-medium truncate text-right"
          :style="{
            color: encounter.state === 'final' && encounter.scoreB > encounter.scoreA ? 'var(--brand)' : 'var(--text-primary)',
            fontWeight: encounter.state === 'final' && encounter.scoreB > encounter.scoreA ? '700' : '500'
          }">
          <span v-if="encounter.partyB === currentUserId" class="text-xs mr-1" style="color: var(--brand)">({{ t('you') }})</span>
          {{ getParticipantName(encounter.partyB) || t('unknown') }}
        </span>
        <div class="w-7 h-7 rounded-lg flex items-center justify-center text-xs font-display font-bold shrink-0"
          style="background: var(--bg-subtle); color: var(--text-muted)">
          {{ getParticipantName(encounter.partyB)?.[0]?.toUpperCase() || '?' }}
        </div>
      </div>
    </div>

    <!-- Score entry -->
    <div v-if="canEnterScore || canConfirmVote || isOwner" class="mt-4 pt-3" style="border-top: 1px solid var(--border)">
      <div v-if="!showScoreForm" class="flex justify-center flex-wrap">
        <button v-if="isOwner || encounter.state === 'planned'" @click="showScoreForm = true" class="btn-ghost btn-sm text-xs" style="color: var(--brand)">
          <i class="pi pi-pencil text-xs" /> {{ t('enter_score') }}
        </button>
        <button v-if="!isOwner  && canConfirmVote" @click="submitScore" class="btn-ghost btn-sm text-xs" style="color: var(--brand)">
          <i class="pi pi-check text-xs" /> {{ t('confirm_score') }}
        </button>
        <button v-if="isOwner && encounter.state !== 'planned'" @click="deleteScore" class="btn-ghost btn-sm text-xs text-yellow-500">
          <i class="pi pi-refresh text-xs" /> {{ t('reset_score') }}
        </button>
        <button v-if="isOwner" @click="deleteEncounter" class="btn-ghost btn-sm text-xs text-red-400">
          <i class="pi pi-trash text-xs" /> {{ t('delete_encounter') }}
        </button>
      </div>
      <div v-else class="space-y-3">
        <div class="flex items-end gap-2">
          <div class="flex-1">
            <label class="input-label text-xs">{{ getParticipantName(encounter.partyA) }}</label>
            <input v-model.number="scoreForm.scoreA" type="number" min="0" class="input py-2 text-center" />
          </div>
          <div class="pb-2.5 font-bold" style="color: var(--text-faint)">–</div>
          <div class="flex-1">
            <label class="input-label text-xs">{{ getParticipantName(encounter.partyB) }}</label>
            <input v-model.number="scoreForm.scoreB" type="number" min="0" class="input py-2 text-center" />
          </div>
        </div>
        <div class="flex gap-2">
          <button @click="submitScore" :disabled="submitting" class="btn-primary btn-sm flex-1">
            <i v-if="submitting" class="pi pi-spin pi-spinner" />
            {{ isOwner ? t('set_final_score') : t('propose_score') }}
          </button>
          <button @click="showScoreForm = false" class="btn-secondary btn-sm px-3">
            <i class="pi pi-times" />
          </button>
        </div>
        <p v-if="!isOwner" class="text-xs text-center" style="color: var(--text-muted)">
          {{ t('opponent_must_agree') }}
        </p>
      </div>
    </div>

    <!-- Voted notice -->
    <div v-if="encounter.state === 'voted' && !showScoreForm && !canConfirmVote"
      class="mt-3 pt-3 flex items-center gap-2 text-xs"
      style="border-top: 1px solid var(--border); color: var(--amber)">
      <i class="pi pi-clock" />
      {{ isOwner ? t('score_proposed_notice_with_name', {name: getParticipantName(encounter.voter_id)}) : t('score_proposed_notice') }}
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useToast } from 'primevue/usetoast'
import { encountersApi } from '@/api'
import { useI18n } from '@/i18n/index.js'

const props = defineProps({
  encounter: { type: Object, required: true },
  participants: { type: Array, default: () => [] },
  currentUserId: { type: String, default: null },
  tournamentId: { type: String, required: true },
  isOwner: { type: Boolean, default: false }
})
const emit = defineEmits(['score-changed'])
const toast = useToast()
const { t } = useI18n()

const showScoreForm = ref(false)
const submitting = ref(false)
const scoreForm = ref({ scoreA: props.encounter.scoreA ?? 0, scoreB: props.encounter.scoreB ?? 0 })

const isInvolved = computed(() =>
  props.currentUserId === props.encounter.partyA || props.currentUserId === props.encounter.partyB
)
const canEnterScore = computed(() =>
  props.encounter.state === 'planned' && (props.isOwner || isInvolved.value)
)

const canConfirmVote = computed(() =>
    props.encounter.state === 'voted' && !props.isOwner && props.currentUserId !== props.encounter.voter_id
)

const encounterStateLabel = computed(() => {
  const map = { final: t('state_final'), voted: t('state_voted'), planned: t('state_planned') }
  return map[props.encounter.state] || props.encounter.state
})
const encounterBadgeClass = computed(() => {
  const map = { final: 'badge-final', voted: 'badge-voted', planned: 'badge-planned' }
  return `badge ${map[props.encounter.state] || 'badge-planned'}`
})

function getParticipantName(id) {
  return props.participants.find(p => p.id === id)?.visible_name
}

function scoreBoxStyle(a, b) {
  if (props.encounter.state !== 'final' || a == null) {
    return { background: 'var(--bg-subtle)', color: 'var(--text-muted)' }
  }
  if (a > b) return { background: 'rgba(var(--brand-rgb), 0.12)', color: 'var(--brand)', fontWeight: '700' }
  if (a < b) return { background: 'var(--bg-subtle)', color: 'var(--text-faint)' }
  return { background: 'var(--bg-subtle)', color: 'var(--text-secondary)' }
}

async function submitScore() {
  submitting.value = true
  try {
    await encountersApi.updateScore(props.tournamentId, {
      partyA: props.encounter.partyA,
      partyB: props.encounter.partyB,
      round: props.encounter.round,
      scoreA: scoreForm.value.scoreA,
      scoreB: scoreForm.value.scoreB,
      state: props.isOwner ? 'final' : 'voted'
    })
    toast.add({ severity: 'success', summary: t('score_saved'), life: 2000 })
    showScoreForm.value = false
    emit('score-changed')
  } catch (e) {
    toast.add({ severity: 'error', summary: t('error'), detail: e.response?.data?.message || '', life: 3000 })
  } finally {
    submitting.value = false
  }
}

async function deleteScore() {
  submitting.value = true
  try {
    await encountersApi.deleteScore(props.tournamentId, {
      partyA: props.encounter.partyA,
      partyB: props.encounter.partyB,
      round: props.encounter.round,
    })
    toast.add({ severity: 'success', summary: t('score_reset'), life: 2000 })
    showScoreForm.value = false
    emit('score-changed')
  } catch (e) {
    toast.add({ severity: 'error', summary: t('error'), detail: e.response?.data?.message || '', life: 3000 })
  } finally {
    submitting.value = false
  }
}

async function deleteEncounter() {
  submitting.value = true
  try {
    await encountersApi.deleteScore(props.tournamentId, {
      partyA: props.encounter.partyA,
      partyB: props.encounter.partyB,
      round: props.encounter.round,
      'delete': true
    })
    toast.add({ severity: 'success', summary: t('encounter_deleted'), life: 2000 })
    showScoreForm.value = false
    emit('score-changed')
  } catch (e) {
    toast.add({ severity: 'error', summary: t('error'), detail: e.response?.data?.message || '', life: 3000 })
  } finally {
    submitting.value = false
  }
}
</script>
