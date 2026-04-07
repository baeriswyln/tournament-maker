<template>
  <div class="p-4 lg:p-8 pb-24 lg:pb-8 max-w-[100vw] sm:max-w-xl mx-auto">
    <h1 class="font-display text-2xl font-bold mb-8 animate-slide-up" style="color: var(--text-primary)">
      {{ t('profile') }}
    </h1>

    <!-- Avatar card -->
    <div class="card p-6 mb-4 flex items-center gap-4 animate-slide-up" style="animation-delay: 0.05s">
      <div class="w-16 h-16 rounded-2xl flex items-center justify-center font-display font-bold text-2xl shrink-0"
        style="background: rgba(var(--brand-rgb), 0.1); border: 1px solid rgba(var(--brand-rgb), 0.2); color: var(--brand)">
        {{ initial }}
      </div>
      <div class="flex-1 min-w-0">
        <h2 class="font-display font-bold text-xl" style="color: var(--text-primary)">{{ auth.user?.visible_name || '—' }}</h2>
        <p class="text-sm" style="color: var(--text-muted)">
          {{ auth.user?.username ? `@${auth.user.username}` : t('participant_only') }}
        </p>
        <button @click="copyId"
          class="text-xs font-mono mt-0.5 truncate max-w-full block transition-colors"
          style="color: var(--text-faint)"
          @mouseenter="e => e.currentTarget.style.color = 'var(--brand)'"
          @mouseleave="e => e.currentTarget.style.color = 'var(--text-faint)'">
          {{ shortId }}
        </button>
      </div>
    </div>

    <!-- Edit form -->
    <div class="card p-5 mb-4 space-y-4 animate-slide-up" style="animation-delay: 0.1s">
      <h3 class="input-label">{{ t('edit_profile') }}</h3>
      <div>
        <label class="input-label">{{ t('display_name') }}</label>
        <input v-model="form.visible_name" type="text" class="input" :placeholder="auth.user?.visible_name" />
      </div>
      <div v-if="auth.user?.username">
        <label class="input-label">{{ t('new_password') }}</label>
        <input v-model="form.password" type="password" class="input" :placeholder="t('new_password_placeholder')" autocomplete="new-password" />
      </div>
      <div v-if="saveError" class="text-sm rounded-xl px-4 py-3"
        style="color: var(--danger); background: rgba(239,68,68,0.08); border: 1px solid rgba(239,68,68,0.2)">
        {{ saveError }}
      </div>
      <button @click="handleSave" :disabled="saving" class="btn-primary w-full">
        <i v-if="saving" class="pi pi-spin pi-spinner" />
        {{ t('save') }}
      </button>
    </div>

    <!-- Account info -->
    <div class="card p-5 mb-4 animate-slide-up" style="animation-delay: 0.15s">
      <h3 class="input-label mb-4">{{ t('account') }}</h3>
      <div class="space-y-3">
        <div class="flex justify-between items-center">
          <span class="text-sm" style="color: var(--text-muted)">{{ t('account_type') }}</span>
          <span class="text-sm font-medium" style="color: var(--text-primary)">
            {{ auth.user?.username ? t('full_account') : t('participant_only') }}
          </span>
        </div>
        <div class="flex justify-between items-center">
          <span class="text-sm" style="color: var(--text-muted)">{{ t('user_id') }}</span>
          <button @click="copyId"
            class="flex items-center gap-1.5 text-xs font-mono transition-colors"
            style="color: var(--brand)">
            {{ shortId }} <i class="pi pi-copy text-xs" />
          </button>
        </div>
      </div>
    </div>

    <!-- Upgrade prompt -->
    <div v-if="!auth.user?.username" class="card p-5 mb-4 animate-slide-up"
      style="animation-delay: 0.2s; border-color: rgba(var(--brand-rgb), 0.25); background: rgba(var(--brand-rgb), 0.03)">
      <div class="flex gap-3">
        <i class="pi pi-info-circle text-lg mt-0.5" style="color: var(--brand)" />
        <div class="flex-1">
          <h4 class="font-display font-semibold mb-1" style="color: var(--text-primary)">{{ t('upgrade_title') }}</h4>
          <p class="text-sm mb-3" style="color: var(--text-muted)">{{ t('upgrade_body') }}</p>
          <button @click="showUpgrade = !showUpgrade" class="btn-primary btn-sm">{{ t('upgrade_account') }}</button>
        </div>
      </div>
      <Transition name="slide-up">
        <div v-if="showUpgrade" class="mt-4 pt-4 space-y-3" style="border-top: 1px solid rgba(var(--brand-rgb), 0.2)">
          <div>
            <label class="input-label">{{ t('username') }}</label>
            <input v-model="upgrade.username" type="text" class="input py-2 text-sm" autocomplete="username" />
          </div>
          <div>
            <label class="input-label">{{ t('password') }}</label>
            <input v-model="upgrade.password" type="password" class="input py-2 text-sm" autocomplete="new-password" />
          </div>
          <button @click="handleUpgrade" :disabled="saving" class="btn-primary btn-sm w-full">
            <i v-if="saving" class="pi pi-spin pi-spinner" />
            {{ t('upgrade_account') }}
          </button>
        </div>
      </Transition>
    </div>

    <!-- Sign out -->
    <button @click="handleLogout"
      class="w-full btn-secondary flex items-center justify-center gap-2 animate-slide-up"
      style="animation-delay: 0.2s; color: var(--danger); border-color: rgba(239,68,68,0.2)">
      <i class="pi pi-sign-out" />
      {{ t('sign_out') }}
    </button>

    <!-- About -->
    <div class="mt-10 text-center animate-slide-up" style="animation-delay: 0.25s">
      <div class="space-y-1.5">
        <p class="text-xs" style="color: var(--text-faint)">
          {{ t('version_label') }} · {{ t('version', {version, date: d(buildDate, 'datetime')}) }}
        </p>
        <p class="text-xs" style="color: var(--text-faint)">
          {{ t('built_with') }} Vue 3 · Vite · PrimeVue
        </p>
        <p class="text-xs mt-6" style="color: var(--text-faint)">© {{ new Date().getFullYear() }} - Made with ☕ by <a href="https://baeriswyl.dev" target="_blank"
            style="color: var(--brand); text-decoration: underline; text-underline-offset: 2px;">baeriswyl.dev</a></p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { usersApi } from '@/api'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/i18n/index.js'
import { useApiError } from '@/composables/useApiError.js'

const auth = useAuthStore()
const router = useRouter()
const toast = useToast()
const { t, d } = useI18n()
const { resolveApiError } = useApiError()

const version = __APP_VERSION__
const buildDate = new Date(__BUILD_DATE__)

const saving = ref(false)
const saveError = ref('')
const showUpgrade = ref(false)
const form = ref({ visible_name: '', password: '' })
const upgrade = ref({ username: '', password: '' })

const initial = computed(() => (auth.user?.visible_name?.[0] || '?').toUpperCase())
const shortId = computed(() => {
  const id = auth.user?.id || ''
  return id ? id.slice(0, 8) + '…' : '—'
})

function copyId() {
  navigator.clipboard.writeText(auth.user?.id || '')
  toast.add({ severity: 'success', summary: t('copied'), life: 1500 })
}

async function handleSave() {
  saveError.value = ''
  const payload = {}
  if (form.value.visible_name) payload.visible_name = form.value.visible_name
  if (form.value.password) payload.password = form.value.password
  if (!Object.keys(payload).length) return
  saving.value = true
  try {
    await usersApi.update(auth.user.id, payload)
    if (payload.visible_name) auth.setUser({ ...auth.user, visible_name: payload.visible_name })
    form.value = { visible_name: '', password: '' }
    toast.add({ severity: 'success', summary: t('profile_updated'), life: 2000 })
  } catch (e) {
    const { title, detail } = resolveApiError(e)
    saveError.value = detail || title
  } finally {
    saving.value = false
  }
}

async function handleUpgrade() {
  if (!upgrade.value.username || !upgrade.value.password) return
  saving.value = true
  try {
    await usersApi.create({
      id: auth.user.id,
      visible_name: auth.user.visible_name,
      username: upgrade.value.username,
      password: upgrade.value.password
    })
    auth.setUser({ ...auth.user, username: upgrade.value.username })
    showUpgrade.value = false
    toast.add({ severity: 'success', summary: t('upgrade_success'), life: 2000 })
  } catch (e) {
    const { title, detail } = resolveApiError(e)
    toast.add({ severity: 'error', summary: title, detail, life: 4000 })
  } finally {
    saving.value = false
  }
}

function handleLogout() {
  auth.logout()
  router.push('/login')
}
</script>
