<template>
  <Toast position="top-right" />
  <ConfirmDialog />
  <RouterView v-slot="{ Component, route }">
    <Transition name="page" mode="out-in">
      <component :is="Component" :key="route.path" />
    </Transition>
  </RouterView>
</template>

<script setup>
import { onMounted, onUnmounted } from 'vue'
import { RouterView } from 'vue-router'
import Toast from 'primevue/toast'
import ConfirmDialog from 'primevue/confirmdialog'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/i18n/index.js'

const toast = useToast()
const { t } = useI18n()

function onSessionExpired() {
  toast.add({ severity: 'warn', summary: t('session_expired_title'), detail: t('session_expired_detail'), life: 5000 })
}

onMounted(() => window.addEventListener('session-expired', onSessionExpired))
onUnmounted(() => window.removeEventListener('session-expired', onSessionExpired))
</script>
