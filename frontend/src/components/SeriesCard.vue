<template>
  <article class="series-card">
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
        最近更新于 {{ formatTime(latestTime) }}
      </div>
    </div>
    <button class="series-open" @click="open">查看系列 <span>→</span></button>
  </article>
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
  background: linear-gradient(150deg,rgba(19,55,80,.64),rgba(10,27,43,.92));
  border: 1px solid rgba(78,158,210,.3);
  border-radius: 15px;
  overflow: hidden;
  transition: transform .25s cubic-bezier(0.16, 1, 0.3, 1), box-shadow .25s, background .25s;
  display: flex;
  flex-direction: column;
  position: relative;
  height: 330px;
  box-shadow: 0 12px 30px rgba(0,5,14,.2),inset 0 1px rgba(128,207,251,.035);
}
.series-card::after {
  content: '';
  position: absolute;
  inset: 0;
  border-radius: 15px;
  padding: 1px;
  background: linear-gradient(135deg,rgba(105,195,245,.5),transparent 42%,rgba(41,128,185,.15));
  -webkit-mask: linear-gradient(#fff 0 0) content-box, linear-gradient(#fff 0 0);
  -webkit-mask-composite: xor;
  mask-composite: exclude;
  opacity: 0;
  transition: opacity 0.3s;
  pointer-events: none;
}
.series-card:hover {
  transform: translateY(-3px);
  border-color: rgba(105,195,245,.48);
  box-shadow: 0 22px 46px rgba(0,7,17,.34),0 0 0 1px rgba(68,158,216,.1);
  background: linear-gradient(150deg,rgba(24,68,98,.7),rgba(11,31,49,.96));
}
.series-card:hover::after {
  opacity: 1;
}

.series-icon {
  position: relative;
  height: 140px;
  background: radial-gradient(ellipse at center,rgba(70,162,219,.18),transparent 62%),linear-gradient(180deg,rgba(3,14,25,.25),transparent);
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  padding: 10px;
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
  transform: scale(1.045);
  opacity: 1;
}
.series-placeholder {
  font-size: 22px;
  font-weight: 800;
  color: var(--series);
  letter-spacing: 2px;
  text-shadow: 0 0 24px rgba(62,148,207,.38);
}

.series-badge {
  position: absolute;
  top: 10px;
  right: 10px;
  font-size: 11px;
  font-weight: 700;
  padding: 3px 8px;
  border-radius: 99px;
  background: rgba(35,117,169,.88);
  border: 1px solid rgba(153,220,255,.16);
  color: #fff;
  box-shadow: 0 4px 8px rgba(41, 128, 185, 0.3);
  z-index: 2;
}

.series-body {
  padding: 13px 15px;
  background: linear-gradient(180deg,transparent,rgba(3,15,26,.18));
  flex: 1;
  display: flex;
  flex-direction: column;
}
.series-name {
  font-size: 15px;
  font-weight: 700;
  color: #82cdf6;
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
  border-radius: 6px;
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
.series-open { margin: 0 14px 14px; padding: 9px 11px; border-radius: 9px; background: linear-gradient(90deg,rgba(48,138,194,.2),rgba(48,138,194,.1)); color: #82cdf6; font-size: 12px; font-weight: 700; text-align: left; border: 1px solid rgba(82,174,229,.24); box-shadow: inset 0 1px rgba(255,255,255,.025); }
.series-open span { float: right; transition: transform .2s; }
.series-open:hover { background: rgba(41,128,185,.27); }
.series-open:hover span { transform: translateX(3px); }
</style>
