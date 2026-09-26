<template>
  <article class="game-card">
    <div class="card-icon">
      <img v-if="iconSrc" :src="iconSrc" alt="" />
      <div v-else class="icon-placeholder">{{ displayName.charAt(0) }}</div>
      <span class="play-badge" :class="game.IsPlay ? 'played' : 'unplayed'">
        {{ game.IsPlay ? '已玩' : '未玩' }}
      </span>
    </div>
    <div class="card-body">
      <div class="card-name" :title="displayName">{{ displayName }}</div>
      <div class="card-desc">
        <span v-if="game.Description" class="desc-text">{{ game.Description }}</span>
        <span v-else class="desc-empty">暂无简介</span>
      </div>
      <div class="card-time" :title="`加入时间：${formatTime(game.InsertTime)}`">
        <svg viewBox="0 0 24 24" aria-hidden="true"><circle cx="12" cy="12" r="9"/><path d="M12 7v5l3 2"/></svg>
        {{ formatTime(game.InsertTime) }} 加入
      </div>
    </div>
    <div class="card-actions">
      <button class="btn btn-small btn-open" @click="open">启动游戏</button>
      <button class="btn btn-small btn-edit" @click="$emit('edit', game)" aria-label="编辑游戏">编辑</button>
    </div>
  </article>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { api, formatTime } from '../api'
import type { GameDTO } from '../types'

const props = defineProps<{ game: GameDTO }>()
const emit = defineEmits<{
  (e: 'edit', game: GameDTO): void
  (e: 'launch-error', message: string): void
}>()

const iconSrc = ref('')

const displayName = computed(() => props.game.NickName || props.game.Name)

function resolveIcon() {
  const p = props.game.IconPath
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
  api.openGame(props.game).catch((err: unknown) => emit('launch-error', String(err)))
}
</script>

<style scoped>
.game-card {
  background: linear-gradient(155deg,rgba(22,38,59,.94),rgba(12,24,40,.94));
  border: 1px solid var(--border);
  border-radius: 15px;
  overflow: hidden;
  transition: transform .25s cubic-bezier(0.16, 1, 0.3, 1), box-shadow .25s, border-color .25s;
  display: flex;
  flex-direction: column;
  position: relative;
  isolation: isolate;
  box-shadow: 0 10px 28px rgba(0,5,14,.16), inset 0 1px rgba(255,255,255,.025);
  height: 334px;
}
.game-card::after {
  content: '';
  position: absolute;
  inset: 0;
  border-radius: 15px;
  padding: 1px;
  background: linear-gradient(135deg, rgba(151,212,251,.3), rgba(255,255,255,0) 38%, rgba(76,168,232,.08));
  -webkit-mask: linear-gradient(#fff 0 0) content-box, linear-gradient(#fff 0 0);
  -webkit-mask-composite: xor;
  mask-composite: exclude;
  opacity: 0;
  transition: opacity 0.3s;
  pointer-events: none;
}
.game-card:hover {
  transform: translateY(-3px);
  border-color: rgba(103,185,238,.34);
  box-shadow: 0 20px 44px rgba(0,5,14,.32), 0 0 0 1px rgba(76,168,232,.08), inset 0 1px rgba(255,255,255,.04);
}
.game-card:hover::after {
  opacity: 1;
}
.card-icon {
  position: relative;
  height: 140px;
  background: radial-gradient(ellipse at 50% 38%,rgba(80,148,194,.11),transparent 63%),linear-gradient(180deg,rgba(2,9,18,.28),transparent);
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  padding: 10px;
  flex-shrink: 0;
}
.card-icon::before {
  content: '';
  position: absolute;
  inset: 0;
  background: radial-gradient(circle at center, rgba(76,168,232,.14), transparent 60%);
  opacity: 0;
  transition: opacity 0.3s;
}
.game-card:hover .card-icon::before {
  opacity: 1;
}
.card-icon img {
  width: 100%;
  height: 100%;
  object-fit: contain;
  transition: transform 0.5s cubic-bezier(0.16, 1, 0.3, 1);
  border-radius: 10px;
  filter: drop-shadow(0 10px 16px rgba(0,4,12,.24));
}
.game-card:hover .card-icon img {
  transform: scale(1.045);
}
.icon-placeholder {
  font-size: 64px;
  font-weight: 800;
  color: var(--border);
  transition: transform 0.5s cubic-bezier(0.16, 1, 0.3, 1), color 0.3s;
}
.game-card:hover .icon-placeholder {
  transform: scale(1.1);
  color: var(--accent);
}

.play-badge {
  position: absolute;
  top: 10px;
  right: 10px;
  font-size: 11px;
  font-weight: 700;
  padding: 4px 9px;
  border-radius: 99px;
  color: #fff;
  backdrop-filter: blur(8px);
  border: 1px solid rgba(255,255,255,.12);
  box-shadow: 0 5px 14px rgba(0,0,0,.2), inset 0 1px rgba(255,255,255,.1);
  z-index: 2;
}
.play-badge.played {
  background: rgba(45, 151, 98, 0.88);
}
.play-badge.unplayed {
  background: rgba(67, 82, 102, 0.9);
}

.card-body {
  padding: 13px 15px 10px;
  flex: 1;
  display: flex;
  flex-direction: column;
  background: linear-gradient(180deg,transparent,rgba(5,13,24,.14));
}
.card-name {
  font-size: 15px;
  font-weight: 720;
  letter-spacing: .1px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  color: var(--text);
  margin-bottom: 4px;
}
.card-desc {
  margin-top: 6px;
  height: 54px;
  overflow: hidden;
  font-size: 12px;
  line-height: 1.5;
  color: #8395ae;
  overflow-wrap: anywhere;
}
.desc-text {
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
.desc-empty {
  color: rgba(148, 163, 184, 0.5);
  font-style: italic;
}
.card-time {
  margin-top: 12px;
  font-size: 12px;
  color: #667c99;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  gap: 6px;
}
.card-time svg { width: 13px; fill: none; stroke: currentColor; stroke-width: 1.7; }

.card-actions {
  display: flex;
  gap: 8px;
  padding: 12px 15px 21px;
  border-top: 1px solid rgba(145,175,209,.07);
}
.btn-small {
  padding: 8px 12px;
  font-size: 13px;
  flex: 0 0 auto;
  border-radius: 8px;
  font-weight: 500;
}
.btn-open {
  flex: 1;
  background: linear-gradient(135deg,#47a5e4,#277daf);
  border: 1px solid rgba(128,207,251,.17);
  box-shadow: inset 0 1px rgba(255,255,255,.15),0 6px 16px rgba(20,99,150,.18);
  color: #fff;
}
.btn-open:hover {
  background: linear-gradient(135deg,#5cb4ea,#318bc1);
  box-shadow: inset 0 1px rgba(255,255,255,.2),0 8px 20px rgba(20,99,150,.3);
}
.btn-edit { width: 58px; background: rgba(255,255,255,.015); border: 1px solid var(--border); color: var(--text-dim); }
.btn-edit:hover { color: var(--text); background: rgba(145,175,209,.08); }
</style>
