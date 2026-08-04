<template>
  <div class="series-card" @click="open" @dblclick="open">
    <div class="series-icon">
      <img v-if="iconSrc" :src="iconSrc" alt="" />
      <div v-else class="series-placeholder">系列</div>
      <span class="series-badge">{{ group.games.length }} 款</span>
    </div>
    <div class="series-body">
      <div class="series-name" :title="group.name">{{ group.name }}</div>
      <div class="series-thumbnails">
        <template v-for="(g, idx) in group.games.slice(0, 5)" :key="g.Id">
          <div class="thumb-wrapper">
             <img v-if="thumbSrcs[g.Id]" :src="thumbSrcs[g.Id]" alt="" class="thumb-img"/>
             <div v-else class="thumb-placeholder">{{ (g.NickName || g.Name).charAt(0) }}</div>
          </div>
        </template>
        <div v-if="group.games.length > 5" class="thumb-more">+{{ group.games.length - 5 }}</div>
      </div>
      <div class="series-time" :title="`最新加入：${formatTime(latestTime)}`">
        {{ formatTime(latestTime) }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { api, formatTime } from '../api'
import type { SeriesGroup } from '../types'

const props = defineProps<{ group: SeriesGroup }>()
const emit = defineEmits<{ (e: 'open', name: string): void }>()

const iconSrc = ref('')
const thumbSrcs = ref<Record<string, string>>({})

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

function resolveThumbs() {
  const games = props.group.games.slice(0, 5)
  games.forEach(g => {
    const p = g.IconPath
    if (!p) return
    if (p.startsWith('data:')) {
      thumbSrcs.value[g.Id] = p
      return
    }
    api.iconToBase64(p).then(src => {
      if (src) thumbSrcs.value[g.Id] = src
    }).catch(() => {})
  })
}

onMounted(() => {
  resolveIcon()
  resolveThumbs()
})

watch(() => props.group, () => {
  thumbSrcs.value = {}
  resolveIcon()
  resolveThumbs()
}, { deep: true })

function open() {
  emit('open', props.group.name)
}
</script>

<style scoped>
.series-card {
  background: var(--series-bg);
  border: 1px solid rgba(41, 128, 185, 0.45);
  border-left: 3px solid var(--series);
  border-radius: 12px;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.16, 1, 0.3, 1);
  display: flex;
  flex-direction: column;
  position: relative;
  height: 330px;
}
.series-card::after {
  content: '';
  position: absolute;
  inset: 0;
  border-radius: 12px;
  padding: 2px;
  background: linear-gradient(135deg, rgba(41, 128, 185, 0.5), transparent);
  -webkit-mask: linear-gradient(#fff 0 0) content-box, linear-gradient(#fff 0 0);
  -webkit-mask-composite: xor;
  mask-composite: exclude;
  opacity: 0;
  transition: opacity 0.3s;
  pointer-events: none;
}
.series-card:hover {
  transform: translateY(-4px) scale(1.02);
  box-shadow: 0 12px 24px rgba(41, 128, 185, 0.2), 0 0 0 1px var(--series);
  background: rgba(41, 128, 185, 0.26);
}
.series-card:hover::after {
  opacity: 1;
}

.series-icon {
  position: relative;
  height: 140px;
  background: linear-gradient(to bottom, rgba(41, 128, 185, 0.1), transparent);
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  padding: 6px;
  flex-shrink: 0;
}
.series-icon::before {
  content: '';
  position: absolute;
  inset: 0;
  background: radial-gradient(circle at center, rgba(41, 128, 185, 0.2), transparent);
  opacity: 0;
  transition: opacity 0.3s;
}
.series-card:hover .series-icon::before {
  opacity: 1;
}
.series-icon img {
  width: 100%;
  height: 100%;
  object-fit: contain;
  transition: transform 0.5s cubic-bezier(0.16, 1, 0.3, 1);
  opacity: 0.9;
  border-radius: 6px;
}
.series-card:hover .series-icon img {
  transform: scale(1.08);
  opacity: 1;
}
.series-placeholder {
  font-size: 22px;
  font-weight: 800;
  color: var(--series);
  letter-spacing: 2px;
  text-shadow: 0 0 15px rgba(41, 128, 185, 0.5);
}

.series-badge {
  position: absolute;
  top: 10px;
  right: 10px;
  font-size: 11px;
  font-weight: 700;
  padding: 3px 8px;
  border-radius: 10px;
  background: var(--series);
  color: #fff;
  box-shadow: 0 4px 8px rgba(41, 128, 185, 0.3);
  z-index: 2;
}

.series-body {
  padding: 12px 14px;
  background: linear-gradient(to top, rgba(41, 128, 185, 0.05), transparent);
  flex: 1;
  display: flex;
  flex-direction: column;
}
.series-name {
  font-size: 15px;
  font-weight: 700;
  color: #5fb3e8;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-bottom: 8px;
}
.series-thumbnails {
  display: flex;
  gap: 4px;
  margin-bottom: 10px;
  flex-wrap: wrap;
  min-height: 24px;
}
.thumb-wrapper {
  width: 24px;
  height: 24px;
  border-radius: 4px;
  background: rgba(0,0,0,0.2);
  border: 1px solid rgba(41, 128, 185, 0.3);
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
}
.thumb-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.thumb-placeholder {
  font-size: 12px;
  color: var(--series);
  font-weight: bold;
}
.thumb-more {
  width: 24px;
  height: 24px;
  border-radius: 4px;
  background: rgba(41, 128, 185, 0.15);
  border: 1px dashed rgba(41, 128, 185, 0.4);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 11px;
  color: var(--series);
  font-weight: 600;
}
.series-time {
  font-size: 12px;
  color: var(--text-dim);
  display: flex;
  align-items: center;
  gap: 6px;
  margin-top: auto;
}
.series-time::before {
  content: '🕒';
  font-size: 10px;
  filter: hue-rotate(240deg);
}
</style>
