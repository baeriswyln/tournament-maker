<template>
  <div class="p-4 lg:p-8 pb-24 lg:pb-8 max-w-[100vw] md:max-w-2xl mx-auto">
    <!-- Scanner card -->
    <div class="card overflow-hidden mb-4 animate-slide-up" style="animation-delay: 0.05s">

      <!-- Viewfinder -->
      <div class="relative aspect-square w-full" style="background: #000">
        <video
            ref="videoEl"
            class="w-full h-full object-cover"
            :class="{ 'opacity-0': !cameraActive }"
            autoplay playsinline muted
        />

        <!-- Overlay: corner brackets -->
        <div v-if="cameraActive && !result" class="absolute inset-0 flex items-center justify-center pointer-events-none">
          <div class="relative w-56 h-56">
            <!-- TL -->
            <span class="absolute top-0 left-0 w-8 h-8 border-t-2 border-l-2 rounded-tl-lg" style="border-color: var(--brand)" />
            <!-- TR -->
            <span class="absolute top-0 right-0 w-8 h-8 border-t-2 border-r-2 rounded-tr-lg" style="border-color: var(--brand)" />
            <!-- BL -->
            <span class="absolute bottom-0 left-0 w-8 h-8 border-b-2 border-l-2 rounded-bl-lg" style="border-color: var(--brand)" />
            <!-- BR -->
            <span class="absolute bottom-0 right-0 w-8 h-8 border-b-2 border-r-2 rounded-br-lg" style="border-color: var(--brand)" />
            <!-- Scanning line -->
            <div class="scan-line absolute left-1 right-1" style="background: var(--brand)" />
          </div>
        </div>

        <!-- Success overlay -->
        <Transition name="fade">
          <div v-if="result" class="absolute inset-0 flex flex-col items-center justify-center gap-3"
               style="background: rgba(var(--brand-rgb), 0.92)">
            <div class="w-16 h-16 rounded-full bg-white/20 flex items-center justify-center">
              <i class="pi pi-check text-white text-3xl" />
            </div>
            <p class="text-white font-display font-bold text-lg">{{t('scanned')}}</p>
          </div>
        </Transition>

        <!-- Idle / permission state -->
        <div v-if="!cameraActive && !result" class="absolute inset-0 flex flex-col items-center justify-center gap-3"
             style="background: var(--bg-subtle)">
          <i class="pi pi-camera text-5xl" style="color: var(--text-faint)" />
          <p class="text-sm font-medium" style="color: var(--text-muted)">
            {{ statusMessage }}
          </p>
        </div>
      </div>

      <!-- Hidden canvas for jsQR frame processing -->
      <canvas ref="canvasEl" class="hidden" />

      <!-- Camera controls -->
      <div class="p-4 flex items-center gap-3 flex-wrap" style="border-top: 1px solid var(--border)">
        <button v-if="!cameraActive && !result" @click="startCamera" :disabled="starting"
                class="btn-primary flex-1">
          <i v-if="starting" class="pi pi-spin pi-spinner" />
          <i v-else class="pi pi-camera" />
          {{ starting ? t('starting') : t('scan_qr_code') }}
        </button>

        <button v-if="cameraActive && !result" @click="stopCamera" class="btn-secondary flex-1">
          <i class="pi pi-times" /> Stop
        </button>

        <template v-if="result">
          <button @click="reset" class="btn-secondary flex-1">
            <i class="pi pi-refresh" /> {{t('retry')}}
          </button>
        </template>

        <!-- Camera flip (mobile) -->
        <button v-if="cameraActive && !result && hasMultipleCameras" @click="flipCamera"
                class="btn-secondary btn-sm p-2.5">
          <i class="pi pi-sync" />
        </button>
      </div>

      <div class="p-4 flex items-center gap-3 flex-wrap">
        <label class="input-label">{{ t('join_with_code') }} <span style="color: var(--brand)">*</span></label>
        <input v-model="form.id" type="text" class="input" :placeholder="t('tournament_id')" required />
      </div>

    </div>

    <!-- Actions -->
    <div class="flex gap-3 pt-2">
      <button type="button" @click="$router.back()" class="btn-secondary flex-1">{{ t('cancel') }}</button>
      <button type="submit" :disabled="loading || !form.id" class="btn-primary flex-1" @click="submit">
        <i v-if="loading" class="pi pi-spin pi-spinner" />
        {{t('join')}}
      </button>
    </div>

    <!-- Error banner -->
    <Transition name="fade">
      <div v-if="errorMsg" class="mt-4 rounded-xl px-4 py-3 flex items-start gap-2 text-sm animate-slide-up"
           style="color: var(--danger); background: rgba(239,68,68,0.08); border: 1px solid rgba(239,68,68,0.2)">
        <i class="pi pi-exclamation-triangle mt-0.5 shrink-0" />
        {{ errorMsg }}
      </div>
    </Transition>
  </div>
</template>

<script setup>
import { ref, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/i18n/index.js'
import { useTournamentStore } from '@/stores/tournament'
import {useAuthStore} from "@/stores/auth.js";

const router = useRouter()
const toast = useToast()
const auth = useAuthStore()
const store = useTournamentStore()
const { t } = useI18n()

// ── Refs ──────────────────────────────────────────────────
const videoEl = ref(null)
const canvasEl = ref(null)
const cameraActive = ref(false)
const starting = ref(false)
const result = ref('')       // raw QR value
const errorMsg = ref('')
const statusMessage = ref(t('camera_not_started'))
const hasMultipleCameras = ref(false)
const facingMode = ref('environment') // prefer rear camera on mobile

let stream = null
let rafId = null
let detector = null

const loading = ref(false)
const form = ref({
  id: ''
})

async function submit() {
  errorMsg.value = ''
  loading.value = true
  try {
    const payload = {id: auth.user.id}
    await store.addParticipants(form.value.id, payload)
    toast.add({severity: 'success', summary: t('joined'), life: 2000})
    router.push(`/tournaments/${form.value.id}`)
  } catch (e) {
    errorMsg.value = e.response?.data?.message || t('error')
    toast.add({severity: 'error', summary: errorMsg.value, life: 2000})
  } finally {
    loading.value = false
  }
}

// ── BarcodeDetector feature-detect ───────────────────────
function getBarcodeDetector() {
  if ('BarcodeDetector' in window) {
    try {
      return new window.BarcodeDetector({ formats: ['qr_code'] })
    } catch { /* fall through to jsQR */ }
  }
  return null
}

// ── Camera lifecycle ──────────────────────────────────────
async function startCamera() {
  errorMsg.value = ''
  starting.value = true
  statusMessage.value = t('camera_request_permission')
  try {
    // Check for multiple cameras
    const devices = await navigator.mediaDevices.enumerateDevices()
    const videoInputs = devices.filter(d => d.kind === 'videoinput')
    hasMultipleCameras.value = videoInputs.length > 1

    stream = await navigator.mediaDevices.getUserMedia({
      video: { facingMode: facingMode.value, width: { ideal: 1280 }, height: { ideal: 1280 } },
      audio: false
    })
    videoEl.value.srcObject = stream
    await videoEl.value.play()
    cameraActive.value = true

    detector = getBarcodeDetector()
    scheduleFrame()
  } catch (e) {
    if (e.name === 'NotAllowedError') {
      errorMsg.value = 'Camera permission denied. Please allow camera access in your browser settings.'
    } else if (e.name === 'NotFoundError') {
      errorMsg.value = 'No camera found on this device.'
    } else {
      errorMsg.value = `Camera error: ${e.message}`
    }
    statusMessage.value = t('camera_not_available')
  } finally {
    starting.value = false
  }
}

function stopCamera() {
  cancelAnimationFrame(rafId)
  stream?.getTracks().forEach(t => t.stop())
  stream = null
  cameraActive.value = false
  statusMessage.value = t('camera_stopped')
}

async function flipCamera() {
  facingMode.value = facingMode.value === 'environment' ? 'user' : 'environment'
  stopCamera()
  await startCamera()
}

// ── Frame scanning loop ───────────────────────────────────
function scheduleFrame() {
  rafId = requestAnimationFrame(scanFrame)
}

async function scanFrame() {
  if (!cameraActive.value || result.value) return

  const video = videoEl.value
  if (!video || video.readyState < 2) { scheduleFrame(); return }

  try {
    if (detector) {
      // Native BarcodeDetector — fast, off-main-thread on some browsers
      const codes = await detector.detect(video)
      if (codes.length > 0) {
        handleDetected(codes[0].rawValue)
        return
      }
    } else {
      // jsQR fallback — draw frame to canvas and decode
      await scanWithJsQR(video)
    }
  } catch { /* keep scanning */ }

  scheduleFrame()
}

async function scanWithJsQR(video) {
  const canvas = canvasEl.value
  const ctx = canvas.getContext('2d', { willReadFrequently: true })
  canvas.width = video.videoWidth
  canvas.height = video.videoHeight
  ctx.drawImage(video, 0, 0)

  // Lazy-load jsQR only when needed
  if (!window._jsQR) {
    try {
      const mod = await import('https://cdn.jsdelivr.net/npm/jsqr@1.4.0/dist/jsQR.js')
      window._jsQR = mod.default ?? window.jsQR
    } catch {
      errorMsg.value = 'QR decoding library failed to load. Try pasting the ID manually.'
      stopCamera()
      return
    }
  }
  const imageData = ctx.getImageData(0, 0, canvas.width, canvas.height)
  const code = window._jsQR(imageData.data, imageData.width, imageData.height, {
    inversionAttempts: 'dontInvert'
  })
  if (code) handleDetected(code.data)
}

// ── Value extraction ──────────────────────────────────────
// Accepts either a bare UUID, or a URL like https://…/tournaments/<uuid>
const UUID_RE = /[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}/i

function extractId(raw) {
  const match = raw.match(UUID_RE)
  return match ? match[0] : raw.trim()
}

function handleDetected(raw) {
  result.value = raw
  form.value.id = extractId(raw)
  stopCamera()
}

function onManualInput() {
  // Also run extraction when user pastes a full URL
  id.value = extractId(id.value)
}

// ── Navigation ────────────────────────────────────────────
function confirm() {
  if (!id.value) return
  stopCamera()
  // Navigate to the tournament — uses the extracted UUID
  router.push(`/tournaments/${id.value}`)
}

function reset() {
  result.value = ''
  id.value = ''
  errorMsg.value = ''
  startCamera()
}

// ── Cleanup ───────────────────────────────────────────────
onUnmounted(() => {
  stopCamera()
})
</script>

<style scoped>
/* Scanning sweep line */
.scan-line {
  height: 2px;
  border-radius: 1px;
  opacity: 0.85;
  animation: sweep 2s ease-in-out infinite;
}
@keyframes sweep {
  0%   { top: 8px;   opacity: 0; }
  10%  { opacity: 1; }
  90%  { opacity: 1; }
  100% { top: calc(100% - 8px); opacity: 0; }
}

.fade-enter-active, .fade-leave-active { transition: opacity 0.2s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>