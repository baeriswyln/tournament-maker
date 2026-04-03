<template>
  <div class="min-h-screen flex items-center justify-center p-4 relative overflow-hidden"
    style="background: var(--bg);">
    <!-- Decorative blobs -->
    <div class="absolute inset-0 overflow-hidden pointer-events-none">
      <div class="absolute -top-40 -right-40 w-96 h-96 rounded-full blur-3xl opacity-30"
        style="background: radial-gradient(circle, rgba(var(--brand-rgb), 0.3), transparent 70%)" />
      <div class="absolute -bottom-40 -left-40 w-96 h-96 rounded-full blur-3xl opacity-20"
        style="background: radial-gradient(circle, rgba(var(--brand-rgb), 0.2), transparent 70%)" />
    </div>

    <!-- Language + Theme controls (top-right) -->
    <div class="absolute top-4 right-4 flex items-center gap-1">
      <ThemeToggle />
      <LanguageSwitcher />
    </div>

    <div class="w-full max-w-md relative animate-slide-up">
      <!-- Logo -->
      <div class="text-center mb-8">
        <div class="inline-flex items-center justify-center w-14 h-14 rounded-2xl mb-4"
          style="background: rgba(var(--brand-rgb), 0.1); border: 1px solid rgba(var(--brand-rgb), 0.25);">
          <svg class="w-7 h-7" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2" style="color: var(--brand)">
            <path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4M7.835 4.697a3.42 3.42 0 001.946-.806 3.42 3.42 0 014.438 0 3.42 3.42 0 001.946.806 3.42 3.42 0 013.138 3.138 3.42 3.42 0 00.806 1.946 3.42 3.42 0 010 4.438 3.42 3.42 0 00-.806 1.946 3.42 3.42 0 01-3.138 3.138 3.42 3.42 0 00-1.946.806 3.42 3.42 0 01-4.438 0 3.42 3.42 0 00-1.946-.806 3.42 3.42 0 01-3.138-3.138 3.42 3.42 0 00-.806-1.946 3.42 3.42 0 010-4.438 3.42 3.42 0 00.806-1.946 3.42 3.42 0 013.138-3.138z" />
          </svg>
        </div>
        <h1 class="font-display text-3xl font-bold" style="color: var(--text-primary)">{{ t('app_name') }}</h1>
        <p class="text-sm mt-1" style="color: var(--text-muted)">{{ t('version_label')}} {{ t('version', {version, date: d(buildDate, 'datetime')}) }}</p>
      </div>

      <!-- Tab switcher -->
      <div class="card-elevated p-1 flex gap-1 mb-4 rounded-2xl">
        <button v-for="tab in tabs" :key="tab.id"
          @click="activeTab = tab.id"
          class="flex-1 py-2.5 rounded-xl text-sm font-display font-semibold transition-all duration-200"
          :style="activeTab === tab.id
            ? 'background: var(--brand); color: white; box-shadow: 0 4px 12px rgba(var(--brand-rgb), 0.3)'
            : 'color: var(--text-muted)'"
        >
          {{ tab.label }}
        </button>
      </div>

      <!-- Login form -->
      <Transition name="fade" mode="out-in">
        <div v-if="activeTab === 'login'" key="login" class="card-elevated p-6 space-y-4">
          <div>
            <label class="input-label">{{ t('username') }}</label>
            <input v-model="loginForm.username" type="text" class="input"
              :placeholder="t('username_placeholder')" @keyup.enter="handleLogin" autocomplete="username" />
          </div>
          <div>
            <label class="input-label">{{ t('password') }}</label>
            <input v-model="loginForm.password" type="password" class="input"
              :placeholder="t('password_placeholder')" @keyup.enter="handleLogin" autocomplete="current-password" />
          </div>
          <div v-if="loginError" class="text-sm rounded-xl px-4 py-3"
            style="color: var(--danger); background: rgba(239,68,68,0.08); border: 1px solid rgba(239,68,68,0.2);">
            {{ loginError }}
          </div>
          <button @click="handleLogin" :disabled="loading" class="btn-primary w-full btn-lg mt-2">
            <i v-if="loading" class="pi pi-spin pi-spinner" />
            <span>{{ t('sign_in') }}</span>
          </button>
        </div>

        <!-- Register form -->
        <div v-else key="register" class="card-elevated p-6 space-y-4">
          <div>
            <label class="input-label">
              {{ t('display_name') }}
              <span style="color: var(--brand)">*</span>
            </label>
            <input v-model="registerForm.visible_name" type="text" class="input"
              :placeholder="t('display_name_placeholder')" />
          </div>

          <div class="relative flex items-center gap-3 py-1">
            <div class="flex-1 h-px" style="background: var(--border)" />
            <span class="text-xs px-2" style="color: var(--text-faint); background: var(--surface)">
              {{ t('optional_for_organizers') }}
            </span>
            <div class="flex-1 h-px" style="background: var(--border)" />
          </div>

          <div>
            <label class="input-label">
              {{ t('username') }}
              <span style="color: var(--text-faint)" class="normal-case ml-1">({{ t('optional') }})</span>
            </label>
            <input v-model="registerForm.username" type="text" class="input"
              :placeholder="t('username_placeholder')" autocomplete="username" />
          </div>
          <div>
            <label class="input-label">
              {{ t('password') }}
              <span style="color: var(--text-faint)" class="normal-case ml-1">({{ t('optional') }})</span>
            </label>
            <input v-model="registerForm.password" type="password" class="input"
              :placeholder="t('password_placeholder')" autocomplete="new-password" />
          </div>

          <div v-if="registerError" class="text-sm rounded-xl px-4 py-3"
            style="color: var(--danger); background: rgba(239,68,68,0.08); border: 1px solid rgba(239,68,68,0.2);">
            {{ registerError }}
          </div>

          <button @click="handleRegister" :disabled="loading" class="btn-primary w-full btn-lg">
            <i v-if="loading" class="pi pi-spin pi-spinner" />
            <span>{{ t('register') }}</span>
          </button>

          <p class="text-xs text-center" style="color: var(--text-faint)">
            {{ t('guest_account_note') }}
          </p>
        </div>
      </Transition>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/i18n/index.js'
import ThemeToggle from '@/components/ThemeToggle.vue'
import LanguageSwitcher from '@/components/LanguageSwitcher.vue'
import { useApiError } from '@/composables/useApiError.js'

const router = useRouter()
const route = useRoute()
const auth = useAuthStore()
const toast = useToast()
const { t, d } = useI18n()
const { resolveApiError } = useApiError()

const version = __APP_VERSION__
const buildDate = new Date(__BUILD_DATE__)

const tabs = computed(() => [
  { id: 'login', label: t('sign_in') },
  { id: 'register', label: t('register') }
])
const activeTab = ref('login')
const loading = ref(false)
const loginError = ref('')
const registerError = ref('')
const loginForm = ref({ username: '', password: '' })
const registerForm = ref({ visible_name: '', username: '', password: '' })

async function handleLogin() {
  loginError.value = ''
  if (!loginForm.value.username || !loginForm.value.password) {
    loginError.value = t('username_required'); return
  }
  loading.value = true
  try {
    await auth.login(loginForm.value.username, loginForm.value.password)
    router.push(route.query.redirect || '/')
  } catch (e) {
    const { title, detail } = resolveApiError(e)
    loginError.value = detail || title
  } finally {
    loading.value = false
  }
}

async function handleRegister() {
  registerError.value = ''
  if (!registerForm.value.visible_name) {
    registerError.value = t('display_name_required'); return
  }
  loading.value = true
  try {
    const payload = { visible_name: registerForm.value.visible_name }
    if (registerForm.value.username) payload.username = registerForm.value.username
    if (registerForm.value.password) payload.password = registerForm.value.password
    await auth.register(payload)
    if (registerForm.value.username && registerForm.value.password && !auth.isAuthenticated) {
      await auth.login(registerForm.value.username, registerForm.value.password)
    }
    toast.add({ severity: 'success', summary: t('success'), detail: t('register'), life: 3000 })
    router.push('/')
  } catch (e) {
    const { title, detail } = resolveApiError(e)
    registerError.value = detail || title
  } finally {
    loading.value = false
  }
}
</script>
