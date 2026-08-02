<template>
  <div class="series-card" @click="open" @dblclick="open">
    <div class="series-icon">
      <img v-if="iconSrc" :src="iconSrc" alt="" />
      <div v-else class="series-placeholder">系列</div>
      <span class="series-badge">{{ group.games.length }} 款</span>
    </div>
    <div class="series-body">
      <div class="series-name" :title="group.name">{{ group.name }}</div>
      <div class="series-time" :title="`最新加入：${formatTime(latestTime)}`">
        {{ formatTime(latestTime) }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { api, formatTime } from '../api'
import type { SeriesGroup } from '../types'

const props = defineProps<{ group: SeriesGroup }>()
const emit = defineEmits<{ (e: 'open', name: string): void }>()

const iconSrc = ref('')

const latestTime = computed(() => props.group.latest.InsertTime)

function resolveIcon() {
  const p = props.group.latest.IconPath
  if (!p) return
  if (p.startsWith('data:')) {
    iconSrc.value = p
    return
  }
  api
    .iconToBase64(p)
    .then((src) => {
      if (src) iconSrc.value = src
    })
    .catch(() => {})
}

onMounted(resolveIcon)

function open() {
  emit('open', props.group.name)
}
</script>

<style scoped>
.series-card {
  background: var(--series-bg);
  border: 1px solid rgba(41, 128, 185, 0.45);
  border-left: 3px solid var(--series);
  border-radius: 10px;
  overflow: hidden;
  cursor: pointer;
  transition: transform 0.12s, box-shadow 0.12s;
  display: flex;
  flex-direction: column;
}
.series-card:hover {
  transform: translateY(-3px);
  box-shadow: 0 8px 20px rgba(41, 128, 185, 0.25);
  background: rgba(41, 128, 185, 0.26);
}

.series-icon {
  position: relative;
  height: 110px;
  background: rgba(27, 38, 59, 0.6);
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
}
.series-icon img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.series-placeholder {
  font-size: 20px;
  font-weight: 700;
  color: rgba(52, 152, 219, 0.9);
}

.series-badge {
  position: absolute;
  top: 8px;
  right: 8px;
  font-size: 11px;
  padding: 2px 8px;
  border-radius: 10px;
  background: var(--series);
  color: #fff;
}

.series-body {
  padding: 10px 12px 10px;
}
.series-name {
  font-size: 14px;
  font-weight: 600;
  color: #5fb3e8;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.series-time {
  margin-top: 4px;
  font-size: 12px;
  color: rgba(147, 164, 195, 0.9);
}
</style>
