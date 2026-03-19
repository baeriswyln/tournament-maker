<template>
  <div class="p-4 lg:p-8 pb-24 lg:pb-8 max-w-[100vw] lg:max-w-6xl mx-auto">
    <!-- Header -->
    <div class="flex items-center gap-3 mb-6 animate-slide-up">
      <button @click="$router.back()" class="btn-ghost btn-sm p-2">
        <i class="pi pi-arrow-left"/>
      </button>
      <div class="flex-1 min-w-0">
        <div class="flex items-center gap-2 flex-wrap">
          <h1 class="font-display text-xl lg:text-2xl font-bold truncate" style="color: var(--text-primary)">
            {{ store.current?.name || '…' }}
          </h1>
          <span :class="badgeClass">{{ stateLabel }}</span>
        </div>
        <p class="text-xs mt-0.5" style="color: var(--text-muted)">{{ t('admin_panel') }}</p>
      </div>
      <div class="hidden md:flex gap-2">
        <RouterLink :to="`/tournaments/${tournamentId}`" class="btn-ghost btn-sm">
          <i class="pi pi-eye"/> {{ t('view') }}
        </RouterLink>
        <RouterLink :to="`/tournaments/${tournamentId}/edit`" class="btn-secondary btn-sm">
          <i class="pi pi-pencil"/>
        </RouterLink>
      </div>
    </div>

    <!-- Mobile FAB -->
    <div class="flex gap-3 md:hidden fixed bottom-20 right-4 z-20 ">
      <RouterLink :to="`/tournaments/${tournamentId}`"
                  class="w-14 h-14 rounded-2xl flex items-center justify-center btn-secondary">
        <i class="pi pi-eye text-xl"/>
      </RouterLink>
      <RouterLink :to="`/tournaments/${tournamentId}/edit`"
                  class="z-20 w-14 h-14 rounded-2xl flex items-center justify-center btn-secondary">
        <i class="pi pi-pencil text-xl"/>
      </RouterLink>
    </div>


    <div v-if="loading" class="flex justify-center py-20">
      <i class="pi pi-spin pi-spinner text-3xl" style="color: var(--brand)"/>
    </div>

    <template v-else>
      <!-- Quick actions -->
      <div class="grid grid-cols-2 lg:grid-cols-4 gap-3 mb-6 animate-slide-up" style="animation-delay: 0.05s">

        <!-- Status -->
        <div class="card p-4 lg:col-span-2">
          <p class="input-label mb-3">{{ t('quick_status') }}</p>
          <div class="flex gap-2 flex-wrap">
            <button v-for="s in states" :key="s.value"
                    @click="changeState(s.value)"
                    :disabled="store.current?.state === s.value || stateLoading"
                    class="btn btn-sm rounded-lg border text-xs font-display font-semibold transition-all"
                    :style="store.current?.state === s.value
                ? 'background: rgba(var(--brand-rgb), 0.1); border-color: var(--brand); color: var(--brand); cursor: default'
                : 'background: var(--surface-raised); border-color: var(--border); color: var(--text-muted)'"
            >
              {{ s.label }}
            </button>
          </div>
        </div>

        <!-- QR invite -->
        <div class="card p-4 flex flex-col items-center justify-center gap-2 cursor-pointer transition-all"
             @click="openQrModal"
             @mouseenter="e => e.currentTarget.style.borderColor = 'rgba(var(--brand-rgb), 0.35)'"
             @mouseleave="e => e.currentTarget.style.borderColor = 'var(--border)'"
        >
          <i class="pi pi-qrcode text-2xl" style="color: var(--brand)"/>
          <p class="text-xs font-display font-medium" style="color: var(--text-secondary)">{{ t('invite_qr') }}</p>
        </div>

        <!-- Danger -->
        <div class="card p-4">
          <p class="input-label mb-3">{{ t('edit') }}</p>
          <div class="flex flex-col gap-1.5">
            <button @click="handleReset" class="btn-ghost btn-sm text-xs justify-start text-left"
                    style="color: var(--amber)">
              <i class="pi pi-refresh"/> {{ t('reset_encounters') }}
            </button>
            <button @click="handleDelete" class="btn-ghost btn-sm text-xs justify-start"
                    style="color: var(--danger)">
              <i class="pi pi-trash"/> {{ t('delete_tournament') }}
            </button>
          </div>
        </div>
      </div>

      <!-- Two-column: Participants + Encounters -->
      <div class="grid lg:grid-cols-2 gap-6">

        <!-- Participants -->
        <div class="animate-slide-up" style="animation-delay: 0.1s">
          <div class="flex items-center justify-between mb-3">
            <h2 class="font-display font-semibold" style="color: var(--text-secondary)">
              {{ t('participants_section') }}
              <span class="font-normal text-sm ml-1" style="color: var(--text-faint)">({{
                  store.participants.length
                }})</span>
            </h2>
            <button @click="showAddParticipant = !showAddParticipant" class="btn-ghost btn-sm text-xs"
                    style="color: var(--brand)">
              <i class="pi pi-plus"/> {{ t('add') }}
            </button>
          </div>

          <!-- Add participant panel -->
          <Transition name="slide-up">
            <div v-if="showAddParticipant" class="card p-4 mb-3 space-y-3">
              <div>
                <label class="input-label">{{ t('search_or_uuid') }}</label>
                <input v-model="addSearch" type="text" class="input py-2 text-sm"/>
              </div>
              <div v-if="filteredUsers.length" class="space-y-1 max-h-40 overflow-y-auto">
                <button v-for="u in filteredUsers" :key="u.uuid || u.id"
                        @click="toggleSelectUser(u.uuid || u.id)"
                        class="w-full flex items-center gap-2 px-3 py-2 rounded-lg text-left text-sm transition-all"
                        :style="selectedUserIds.includes(u.uuid || u.id)
                    ? 'background: rgba(var(--brand-rgb), 0.08); border: 1px solid rgba(var(--brand-rgb), 0.2); color: var(--brand)'
                    : 'background: var(--bg-subtle); color: var(--text-secondary); border: 1px solid transparent'"
                >
                  <div class="w-6 h-6 rounded font-bold text-xs flex items-center justify-center shrink-0"
                       style="background: var(--border); color: var(--text-muted)">
                    {{ u.visible_name?.[0]?.toUpperCase() }}
                  </div>
                  {{ u.visible_name }}
                  <span class="ml-auto text-xs" style="color: var(--text-faint)">{{ u.username }}</span>
                </button>
              </div>
              <div class="flex gap-2">
                <button @click="handleAddParticipants" class="btn-primary btn-sm flex-1 text-xs">{{
                    t('add_selected')
                  }}
                </button>
                <button @click="showAddParticipant = false" class="btn-secondary btn-sm px-3 text-xs">{{
                    t('cancel')
                  }}
                </button>
              </div>
            </div>
          </Transition>

          <div class="card overflow-hidden">
            <div v-if="!store.participants.length" class="p-8 text-center text-sm" style="color: var(--text-muted)">
              {{ t('no_participants') }}
            </div>
            <div v-for="(p, i) in store.participants" :key="p.id"
                 class="flex items-center gap-3 px-4 py-3"
                 :style="i > 0 ? 'border-top: 1px solid var(--border)' : ''">
              <div class="w-8 h-8 rounded-lg flex items-center justify-center font-display font-bold text-sm shrink-0"
                   style="background: var(--bg-subtle); color: var(--text-muted)">
                {{ p.visible_name?.[0]?.toUpperCase() || '?' }}
              </div>
              <div class="flex-1 min-w-0">
                <p class="text-sm font-medium truncate" style="color: var(--text-primary)">{{ p.visible_name }}</p>
                <p class="text-xs" style="color: var(--text-muted)">
                  {{ p.encounters?.finished || 0 }}/{{ p.encounters?.total || 0 }} · {{ p.encounters?.score || 0 }}
                  {{ t('pts') }}
                </p>
              </div>
              <button @click="handleRemoveParticipant(p.id)"
                      class="btn-ghost btn-sm p-1.5 rounded-lg"
                      style="color: rgba(239,68,68,0.6)"
                      @mouseenter="e => e.currentTarget.style.color = 'var(--danger)'"
                      @mouseleave="e => e.currentTarget.style.color = 'rgba(239,68,68,0.6)'">
                <i class="pi pi-times text-xs"/>
              </button>
            </div>
          </div>
        </div>

        <!-- Encounters -->
        <div class="animate-slide-up" style="animation-delay: 0.15s">
          <div class="flex items-center justify-between mb-3">
            <h2 class="font-display font-semibold" style="color: var(--text-secondary)">
              {{ t('encounters_section') }}
              <span class="font-normal text-sm ml-1" style="color: var(--text-faint)">({{
                  store.encounters.length
                }})</span>
            </h2>
            <button @click="showEncounterForm = !showEncounterForm" class="btn-ghost btn-sm text-xs"
                    style="color: var(--brand)">
              <i class="pi pi-plus"/> {{ t('add') }}
            </button>
          </div>

          <!-- Add encounter panel -->
          <Transition name="slide-up">
            <div v-if="showEncounterForm" class="card p-4 mb-3 space-y-3">
              <div class="flex gap-1 p-1 rounded-xl mb-6 overflow-x-auto animate-slide-up flex-wrap"
                   style="background: var(--bg-subtle); animation-delay: 0.05s">
                <button v-for="tab in availableTabs" :key="tab.id"
                        @click="activeEncounterTab = tab.id"
                        class="flex-1 min-w-fit py-2 px-3 rounded-lg text-sm font-display font-medium transition-all whitespace-nowrap"
                        :style="activeEncounterTab === tab.id
            ? 'background: var(--surface); color: var(--text-primary); box-shadow: 0 1px 4px rgba(0,0,0,0.08)'
            : 'color: var(--text-muted)'"
                >
                  {{ tab.label }}
                </button>
              </div>

              <div v-if="activeEncounterTab === 'generate'" key="generate" class="space-y-3 animate-slide-up">
                <h2>{{t('generate_encounters_description')}}</h2>

                <div class="flex gap-2 flex-wrap justify-between">
                  <div v-for="group in groups" class="card flex-1">
                    <h2 class="px-4 py-2">{{ group.name }}</h2>
                    <div v-for="(p, i) in group.participants" :key="p.id"
                         class="flex items-center gap-3 px-4 py-3"
                         :style="i > 0 ? 'border-top: 1px solid var(--border)' : ''">
                      <div
                          class="w-8 h-8 rounded-lg flex items-center justify-center font-display font-bold text-sm shrink-0"
                          style="background: var(--bg-subtle); color: var(--text-muted)">
                        {{ p.visible_name?.[0]?.toUpperCase() || '?' }}
                      </div>
                      <div class="flex-1 min-w-0">
                        <p class="text-sm font-medium truncate" style="color: var(--text-primary)">{{
                            p.visible_name
                          }}</p>
                      </div>
                    </div>
                  </div>
                </div>
                <button @click="generateGroups" class="btn-primary btn-sm text-xs w-full">{{ t('generate_groups') }}
                </button>

                <div>
                  <label class="input-label">{{ t('round') }}</label>
                  <input v-model="newEncounter.round" type="text" class="input py-2 text-sm" placeholder="R1, Final…"/>
                </div>

                <div class="flex gap-2">
                  <button @click="generateGroupEncounters"
                          :disabled="!groups.length"
                          class="btn-primary btn-sm flex-1 text-xs">{{ t('generate_encounters') }}
                  </button>
                  <button @click="showEncounterForm = false; groups = []" class="btn-secondary btn-sm px-3 text-xs">{{
                      t('cancel')
                    }}
                  </button>
                </div>
              </div>

              <div v-if="activeEncounterTab === 'manual'" key="manual" class="space-y-3 animate-slide-up">
                <div class="grid grid-cols-2 gap-3">
                  <div>
                    <label class="input-label">{{ t('player_a') }}</label>
                    <select v-model="newEncounter.partyA" class="input py-2 text-sm">
                      <option value="">{{ t('select_player') }}</option>
                      <option v-for="p in store.participants" :key="p.id" :value="p.id">{{ p.visible_name }}</option>
                    </select>
                  </div>
                  <div>
                    <label class="input-label">{{ t('player_b') }}</label>
                    <select v-model="newEncounter.partyB" class="input py-2 text-sm">
                      <option value="">{{ t('select_player') }}</option>
                      <option v-for="p in store.participants" :key="p.id" :value="p.id"
                              :disabled="p.id === newEncounter.partyA">{{ p.visible_name }}
                      </option>
                    </select>
                  </div>
                </div>
                <div>
                  <label class="input-label">{{ t('round') }}</label>
                  <input v-model="newEncounter.round" type="text" class="input py-2 text-sm" placeholder="R1, Final…"/>
                </div>
                <div class="flex gap-2">
                  <button @click="handleAddEncounter"
                          :disabled="!newEncounter.partyA || !newEncounter.partyB"
                          class="btn-primary btn-sm flex-1 text-xs">{{ t('add_encounter') }}
                  </button>
                  <button @click="showEncounterForm = false" class="btn-secondary btn-sm px-3 text-xs">{{
                      t('cancel')
                    }}
                  </button>
                </div>
              </div>
            </div>
          </Transition>

          <div class="space-y-2 max-h-[640px] overflow-y-auto">
            <div v-if="!store.encounters.length" class="card p-8 text-center text-sm" style="color: var(--text-muted)">
              {{ t('no_encounters') }}
            </div>
            <EncounterRow v-for="e in store.encounters"
                          :key="`${e.partyA}-${e.partyB}-${e.round}`"
                          :encounter="e" :participants="store.participants"
                          :tournament-id="tournamentId" :is-owner="true"
                          @score-changed="store.fetchEncounters(tournamentId)"/>
          </div>
        </div>
      </div>
    </template>

    <!-- QR Modal -->
    <Dialog v-model:visible="showQrModal" modal :header="t('invite_qr')" :style="{ width: '380px' }">
      <div class="flex flex-col items-center gap-4 py-2">
        <p class="text-sm text-center" style="color: var(--text-muted)">{{ t('share_qr_desc') }}</p>
        <div class="p-4 rounded-2xl shadow-xl" style="background: white">
          <canvas ref="qrCanvas" width="192" height="192"/>
        </div>
        <div class="w-full rounded-xl px-4 py-3 flex items-center gap-3"
             style="background: var(--bg-subtle); border: 1px solid var(--border)">
          <span class="input-label normal-case tracking-normal mb-0 text-xs shrink-0" style="color: var(--text-muted)">
            {{ t('tournament_id') }}
          </span>
          <code class="text-xs flex-1 truncate font-mono" style="color: var(--brand)">{{ tournamentId }}</code>
          <button @click="copyId" class="btn-ghost btn-sm p-1.5">
            <i class="pi pi-copy text-sm" style="color: var(--text-muted)"/>
          </button>
        </div>
      </div>
    </Dialog>

  </div>
</template>

<script setup>
import {ref, computed, onMounted, nextTick} from 'vue'
import {useRoute, RouterLink, useRouter} from 'vue-router'
import {useTournamentStore} from '@/stores/tournament'
import {useToast} from 'primevue/usetoast'
import {useConfirm} from 'primevue/useconfirm'
import Dialog from 'primevue/dialog'
import ConfirmDialog from 'primevue/confirmdialog'
import {useI18n} from '@/i18n/index.js'
import EncounterRow from '@/components/EncounterRow.vue'
import QRCode from 'qrcode'

const route = useRoute()
const router = useRouter()
const store = useTournamentStore()
const toast = useToast()
const confirm = useConfirm()
const {t} = useI18n()

const tournamentId = route.params.id
const loading = ref(true)
const stateLoading = ref(false)
const showAddParticipant = ref(false)
const showEncounterForm = ref(false)
const showQrModal = ref(false)
const addSearch = ref('')
const selectedUserIds = ref([])
const qrCanvas = ref(null)
const newEncounter = ref({partyA: '', partyB: '', round: 'R1'})
const activeEncounterTab = ref('generate')
const groups = ref([])

const states = computed(() => [
  {value: 'draft', label: t('state_draft')},
  {value: 'open', label: t('state_open')},
  {value: 'live', label: t('state_live')},
  {value: 'closed', label: t('state_closed')},
])

const stateLabel = computed(() => {
  const map = {draft: t('state_draft'), open: t('state_open'), live: t('state_live'), closed: t('state_closed')}
  return map[store.current?.state] || store.current?.state
})
const badgeClass = computed(() => {
  const map = {draft: 'badge-draft', open: 'badge-open', live: 'badge-live', closed: 'badge-closed'}
  return `badge ${map[store.current?.state] || 'badge-draft'}`
})

const filteredUsers = computed(() => {
  const ids = store.participants.map(p => p.id)
  return store.users.filter(u => {
    const uid = u.uuid || u.id
    if (ids.includes(uid)) return false
    if (!addSearch.value) return true
    return (u.visible_name + u.username).toLowerCase().includes(addSearch.value.toLowerCase())
  })
})

const availableTabs = computed(() => {
  const tabs = [
    {id: 'generate', label: t('create_encounter_generate')},
    {id: 'manual', label: t('create_encounter_manual')},
  ]
  return tabs
})

async function openQrModal() {
  showQrModal.value = true
  await nextTick()
  setTimeout(generateQr, 80)
}

async function generateQr() {
  if (!qrCanvas.value) return
  try {
    const url = `${window.location.origin}/tournaments/${tournamentId}`
    await QRCode.toCanvas(qrCanvas.value, url, {
      width: 192, margin: 1,
      color: {dark: '#0A1628', light: '#FFFFFF'}
    })
  } catch (e) {
    console.warn('QR generation failed', e)
  }
}

function copyId() {
  navigator.clipboard.writeText(tournamentId)
  toast.add({severity: 'success', summary: t('copied'), life: 1500})
}

function toggleSelectUser(id) {
  const idx = selectedUserIds.value.indexOf(id)
  if (idx >= 0) selectedUserIds.value.splice(idx, 1)
  else selectedUserIds.value.push(id)
}

async function changeState(state) {
  stateLoading.value = true
  try {
    await store.updateTournament(tournamentId, {state})
    toast.add({severity: 'success', summary: `${t('quick_status')}: ${state}`, life: 2000})
  } finally {
    stateLoading.value = false
  }
}

async function handleAddParticipants() {
  const ids = selectedUserIds.value.length
      ? selectedUserIds.value
      : addSearch.value ? [addSearch.value] : []
  if (!ids.length) return
  try {
    await store.addParticipants(tournamentId, ids)
    toast.add({severity: 'success', summary: t('added'), life: 2000})
    showAddParticipant.value = false
    selectedUserIds.value = []
    addSearch.value = ''
  } catch (e) {
    toast.add({severity: 'error', summary: t('error'), detail: e.response?.data?.message || '', life: 3000})
  }
}

function handleRemoveParticipant(id) {
  confirm.require({
    message: t('remove_participant'),
    header: t('confirm'),
    accept: async () => {
      await store.removeParticipant(tournamentId, id)
      toast.add({severity: 'info', summary: t('removed'), life: 2000})
    }
  })
}

// -- encounter generation --------------------
function shuffle(array) {
  // source: https://stackoverflow.com/a/2450976
  let currentIndex = array.length;

  // While there remain elements to shuffle...
  while (currentIndex != 0) {

    // Pick a remaining element...
    let randomIndex = Math.floor(Math.random() * currentIndex);
    currentIndex--;

    // And swap it with the current element.
    [array[currentIndex], array[randomIndex]] = [
      array[randomIndex], array[currentIndex]];
  }
}

function generateGroups() {
  const nb_participants = store.participants.length
  const nb_groups = store.current.format.nb_groups
  const group_assignments = []
  for (let i = 0; i < nb_participants; i++) {
    group_assignments.push(i % nb_groups);
  }
  shuffle(group_assignments)

  groups.value = []
  for(let gid = 1; gid <= nb_groups; gid++) {
    groups.value.push({
      "name": `${t('group')} ${gid}`,
      "participants": store.participants.filter((e, idx) => gid-1 === group_assignments[idx])
    })
  }
}

async function generateGroupEncounters() {
  if (!groups.value.length || groups.value.some(g => !g.participants.length)) return

  const encounters = []
  for (let i = 0; i < groups.value.length; i++) {
    for (let pi1 = 0; pi1 < groups.value[i].participants.length-1; pi1++) {
      for (let pi2 = pi1+1; pi2 < groups.value[i].participants.length; pi2++) {
        encounters.push({
          partyA: groups.value[i].participants[pi1].id,
          partyB: groups.value[i].participants[pi2].id,
          round: newEncounter.value.round
        })
      }
    }
  }

  try {
    await store.createEncounters(tournamentId, encounters)
    toast.add({severity: 'success', summary: t('encounter_added'), life: 2000})
  } catch (e) {
    toast.add({severity: 'error', summary: t('error'), detail: e.response?.data?.message || '', life: 3000})
  }
}

async function handleAddEncounter() {
  if (!newEncounter.value.partyA || !newEncounter.value.partyB) return
  try {
    await store.createEncounters(tournamentId, [{...newEncounter.value}])
    toast.add({severity: 'success', summary: t('encounter_added'), life: 2000})
    newEncounter.value = {partyA: '', partyB: '', round: 'R1'}
    showEncounterForm.value = false
  } catch (e) {
    toast.add({severity: 'error', summary: t('error'), detail: e.response?.data?.message || '', life: 3000})
  }
}

function handleReset() {
  confirm.require({
    message: t('reset_confirm'),
    header: t('reset_title'),
    accept: async () => {
      await store.resetTournament(tournamentId)
      await store.fetchEncounters(tournamentId)
      toast.add({severity: 'warn', summary: t('tournament_reset'), life: 2000})
    }
  })
}

function handleDelete() {
  console.log('delete')
  confirm.require({
    message: t('delete_confirm'),
    header: t('delete_title'),
    accept: async () => {
      await store.deleteTournament(tournamentId)
      toast.add({severity: 'warn', summary: t('tournament_deleted'), life: 2000})
      router.push('/')
    }
  })
}

onMounted(async () => {
  loading.value = true
  try {
    await Promise.all([
      store.fetchOne(tournamentId),
      store.fetchParticipants(tournamentId),
      store.fetchEncounters(tournamentId),
      store.fetchUsers()
    ])
  } catch (e) {
    toast.add({severity: 'error', summary: t('forbidden'), life: 2000})
    router.push('/')
  }
  loading.value = false
})
</script>
