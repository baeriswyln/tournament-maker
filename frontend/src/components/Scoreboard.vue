<template>
  <div class="card overflow-hidden">
    <!-- Header -->
    <div class="grid grid-cols-12 px-4 py-2.5"
      style="border-bottom: 1px solid var(--border); background: var(--bg-subtle);">
      <div class="col-span-1 text-center input-label normal-case tracking-normal" style="font-size: 0.7rem; margin-bottom: 0">#</div>
      <div class="col-span-5 lg:col-span-6 input-label normal-case tracking-normal" style="font-size: 0.7rem; margin-bottom: 0">{{ t('players') }}</div>
      <div class="col-span-2 text-center input-label normal-case tracking-normal hidden sm:block" style="font-size: 0.7rem; margin-bottom: 0">W/D/L</div>
      <div class="col-span-2 text-center input-label normal-case tracking-normal hidden sm:block" style="font-size: 0.7rem; margin-bottom: 0">+/-</div>
      <div class="col-span-4 sm:col-span-2 text-center input-label normal-case tracking-normal font-bold" style="font-size: 0.7rem; margin-bottom: 0">{{ t('pts') }}</div>
    </div>

    <!-- Empty state -->
    <div v-if="!entries.length" class="p-10 text-center">
      <i class="pi pi-chart-bar text-3xl mb-3 block" style="color: var(--text-faint)" />
      <p class="text-sm" style="color: var(--text-muted)">{{ t('no_scores') }}</p>
    </div>

    <!-- Rows -->
    <TransitionGroup name="scoreboard-row" tag="div">
      <div
        v-for="(p, i) in entries" :key="p.id"
        class="grid grid-cols-12 px-4 py-3 items-center"
        :style="{
          borderTop: i > 0 ? '1px solid var(--border)' : 'none',
          background: i === 0 ? 'rgba(var(--brand-rgb), 0.04)' : 'transparent'
        }"
      >
        <!-- Rank -->
        <div class="col-span-1 flex justify-center">
          <span v-if="i === 0" class="text-lg">🥇</span>
          <span v-else-if="i === 1" class="text-lg">🥈</span>
          <span v-else-if="i === 2" class="text-lg">🥉</span>
          <span v-else class="font-display font-bold text-sm" style="color: var(--text-faint)">{{ i + 1 }}</span>
        </div>

        <!-- Name -->
        <div class="col-span-5 lg:col-span-6 flex items-center gap-2 min-w-0">
          <div class="w-7 h-7 rounded-lg flex items-center justify-center text-xs font-display font-bold shrink-0"
            :style="i === 0
              ? 'background: rgba(var(--brand-rgb), 0.15); color: var(--brand)'
              : 'background: var(--bg-subtle); color: var(--text-muted)'"
          >
            {{ p.visible_name?.[0]?.toUpperCase() || '?' }}
          </div>
          <span class="text-sm font-medium truncate"
            :style="{ color: i < 3 ? 'var(--text-primary)' : 'var(--text-secondary)', fontWeight: i === 0 ? '700' : '500' }">
            {{ p.visible_name || '?' }}
          </span>
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
        <div class="col-span-2 text-center text-sm hidden sm:block"
          :style="{ color: p.score_diff > 0 ? 'var(--success-muted, #4ade80)' : p.score_diff < 0 ? 'var(--danger-muted, #f87171)' : 'var(--text-muted)' }">
          {{ p.score_diff > 0 ? '+' : '' }}{{ p.score_diff }}
        </div>

        <!-- Points -->
        <div class="col-span-4 sm:col-span-2 text-center">
          <span class="font-display font-bold text-base"
            :style="{ color: i === 0 ? 'var(--brand)' : 'var(--text-primary)' }">
            {{ p.points }}
          </span>
        </div>
      </div>
    </TransitionGroup>
  </div>
</template>

<script setup>
import { useI18n } from '@/i18n/index.js'

defineProps({ entries: { type: Array, default: () => [] } })
const { t } = useI18n()
</script>

<style scoped>
.scoreboard-row-move,
.scoreboard-row-enter-active,
.scoreboard-row-leave-active { transition: all 0.4s cubic-bezier(0.16, 1, 0.3, 1); }
.scoreboard-row-enter-from, .scoreboard-row-leave-to { opacity: 0; transform: translateX(-16px); }
.scoreboard-row-leave-active { position: absolute; width: 100%; }
</style>