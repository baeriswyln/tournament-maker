<template>
  <div class="p-4 lg:p-8 pb-24 lg:pb-8 max-w-2xl mx-auto">
    <div class="flex items-center gap-3 mb-8 animate-slide-up">
      <button @click="$router.back()" class="btn-ghost btn-sm p-2">
        <i class="pi pi-arrow-left" />
      </button>
      <h1 class="font-display text-2xl font-bold" style="color: var(--text-primary)">
        {{ isEditing ? t('edit_tournament_title') : t('new_tournament_title') }}
      </h1>
    </div>

    <form @submit.prevent="handleSubmit" class="space-y-4 animate-slide-up" style="animation-delay: 0.05s">

      <!-- Basic info -->
      <div class="card p-5 space-y-4">
        <h2 class="input-label">{{ t('basic_info') }}</h2>
        <div>
          <label class="input-label">{{ t('tournament_name') }} <span style="color: var(--brand)">*</span></label>
          <input v-model="form.name" type="text" class="input" :placeholder="t('tournament_name_placeholder')" required />
        </div>
        <div>
          <label class="input-label">{{ t('start_date') }}</label>
          <input v-model="form.start_date" type="datetime-local" class="input" />
        </div>
      </div>

      <!-- Format -->
      <div class="card p-5 space-y-4">
        <h2 class="input-label">{{ t('format_section') }}</h2>
        <div>
          <label class="input-label">{{ t('format_name') }}</label>
          <input v-model="form.format.name" type="text" class="input" :placeholder="t('format_name_placeholder')" />
        </div>
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="input-label">{{ t('groups') }}</label>
            <input v-model="form.format.nb_groups" type="text" class="input" placeholder="1" />
          </div>
          <div>
            <label class="input-label">{{ t('end_rounds') }}</label>
            <input v-model.number="form.format.nb_endrounds" type="number" min="0" class="input" placeholder="0" />
          </div>
        </div>
      </div>

      <!-- Visibility -->
      <div class="card p-5 space-y-4">
        <h2 class="input-label">{{ t('visibility_section') }}</h2>

        <label v-for="vis in visibilityOptions" :key="vis.key" class="flex items-center justify-between gap-4 cursor-pointer">
          <div>
            <p class="text-sm font-medium" style="color: var(--text-primary)">{{ vis.label }}</p>
            <p class="text-xs mt-0.5" style="color: var(--text-muted)">{{ vis.desc }}</p>
          </div>
          <div class="toggle-track" :class="{ on: form.visibility[vis.key] }"
            @click="form.visibility[vis.key] = !form.visibility[vis.key]">
            <div class="toggle-thumb" />
          </div>
        </label>
      </div>

      <!-- Status (edit only) -->
      <div v-if="isEditing" class="card p-5 space-y-3">
        <h2 class="input-label">{{ t('status_section') }}</h2>
        <div class="grid grid-cols-2 sm:grid-cols-4 gap-2">
          <button v-for="s in states" :key="s.value" type="button"
            @click="form.state = s.value"
            class="py-2.5 px-3 rounded-xl text-sm font-display font-semibold border transition-all"
            :style="form.state === s.value
              ? 'background: rgba(var(--brand-rgb), 0.1); border-color: var(--brand); color: var(--brand)'
              : 'background: var(--surface-raised); border-color: var(--border); color: var(--text-muted)'"
          >
            {{ s.label }}
          </button>
        </div>
        <p class="text-xs" style="color: var(--text-muted)">
          {{ t('status_note', { live: t('state_live'), open: t('state_open') }) }}
        </p>
      </div>

      <!-- Error -->
      <div v-if="error" class="text-sm rounded-xl px-4 py-3"
        style="color: var(--danger); background: rgba(239,68,68,0.08); border: 1px solid rgba(239,68,68,0.2)">
        {{ error }}
      </div>

      <!-- Actions -->
      <div class="flex gap-3 pt-2">
        <button type="button" @click="$router.back()" class="btn-secondary flex-1">{{ t('cancel') }}</button>
        <button type="submit" :disabled="loading" class="btn-primary flex-1">
          <i v-if="loading" class="pi pi-spin pi-spinner" />
          {{ isEditing ? t('save_changes') : t('new_tournament') }}
        </button>
      </div>
    </form>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useTournamentStore } from '@/stores/tournament'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/i18n/index.js'

const route = useRoute()
const router = useRouter()
const store = useTournamentStore()
const toast = useToast()
const { t } = useI18n()

const isEditing = !!route.params.id
const loading = ref(false)
const error = ref('')

const states = computed(() => [
  { value: 'draft',  label: t('state_draft') },
  { value: 'open',   label: t('state_open') },
  { value: 'live',   label: t('state_live') },
  { value: 'closed', label: t('state_closed') },
])

const visibilityOptions = computed(() => [
  { key: 'scoreboard', label: t('show_scoreboard'), desc: t('show_scoreboard_desc') },
  { key: 'scores',     label: t('show_scores'),     desc: t('show_scores_desc') },
])

const form = ref({
  name: '', start_date: '', nb_participants: null,
  visibility: { scoreboard: true, scores: true },
  format: { name: '', nb_groups: '', nb_endrounds: 0 },
  state: 'draft'
})

onMounted(async () => {
  if (isEditing) {
    const t_ = await store.fetchOne(route.params.id)
    form.value = {
      name: t_.name || '',
      start_date: t_.start_date ? new Date(t_.start_date).toISOString().slice(0, 16) : '',
      nb_participants: t_.nb_participants || null,
      visibility: { scoreboard: t_.visibility?.scoreboard ?? true, scores: t_.visibility?.scores ?? true },
      format: { name: t_.format?.name || '', nb_groups: t_.format?.nb_groups || '', nb_endrounds: t_.format?.nb_endrounds || 0 },
      state: t_.state || 'draft'
    }
  }
})

async function handleSubmit() {
  error.value = ''
  if (!form.value.name) { error.value = t('tournament_name_required'); return }
  loading.value = true
  try {
    const payload = { ...form.value }
    if (payload.start_date) payload.start_date = new Date(payload.start_date).toISOString()
    else delete payload.start_date
    if (!payload.nb_participants) delete payload.nb_participants

    if (isEditing) {
      await store.updateTournament(route.params.id, payload)
      toast.add({ severity: 'success', summary: t('updated'), life: 2000 })
      router.push(`/tournaments/${route.params.id}/manage`)
    } else {
      const { id } = await store.createTournament(payload)
      toast.add({ severity: 'success', summary: t('tournament_created'), life: 2000 })
      router.push(`/tournaments/${id}/manage`)
    }
  } catch (e) {
    error.value = e.response?.data?.message || t('error')
  } finally {
    loading.value = false
  }
}
</script>
