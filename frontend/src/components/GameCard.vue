<template>
  <div class="game-card" :class="{ played: game.IsPlay }" @dblclick="open">
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
        {{ formatTime(game.InsertTime) }}
      </div>
    </div>
    <div class="card-actions">
      <button class="btn btn-small btn-open" @click.stop="open">启动</button>
      <button class="btn btn-small" @click.stop="$emit('edit', game)">编辑</button>
      <button class="btn btn-small btn-del" @click.stop="$emit('delete', game)">删除</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { api, formatTime } from '../api'
import type { GameDTO } from '../types'

const props = defineProps<{ game: GameDTO }>()
const emit = defineEmits<{
  (e: 'edit', game: GameDTO): void
  (e: 'delete', game: GameDTO): void
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
  api.openGame(props.game).catch(() => {})
}
</script>

<style scoped>
.game-card {
  background: var(--card);
  border: 1px solid var(--border);
  border-radius: 12px;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.16, 1, 0.3, 1);
  display: flex;
  flex-direction: column;
  position: relative;
  isolation: isolate;
  height: 330px;
}
.game-card::after {
  content: '';
  position: absolute;
  inset: 0;
  border-radius: 12px;
  padding: 2px;
  background: linear-gradient(135deg, rgba(255,255,255,0.1), rgba(255,255,255,0));
  -webkit-mask: linear-gradient(#fff 0 0) content-box, linear-gradient(#fff 0 0);
  -webkit-mask-composite: xor;
  mask-composite: exclude;
  opacity: 0;
  transition: opacity 0.3s;
  pointer-events: none;
}
.game-card:hover {
  transform: translateY(-4px) scale(1.02);
  background: var(--card-hover);
  box-shadow: 0 12px 24px rgba(0, 0, 0, 0.25), 0 0 0 1px var(--accent);
}
.game-card:hover::after {
  opacity: 1;
}
.game-card.played {
  border-left: 3px solid var(--played);
}
.game-card:not(.played) {
  border-left: 3px solid var(--unplayed);
}

.card-icon {
  position: relative;
  height: 140px;
  background: linear-gradient(to bottom, rgba(0,0,0,0.15), transparent);
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  padding: 6px;
  flex-shrink: 0;
}
.card-icon::before {
  content: '';
  position: absolute;
  inset: 0;
  background: radial-gradient(circle at center, rgba(52, 152, 219, 0.1), transparent);
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
  border-radius: 8px;
}
.game-card:hover .card-icon img {
  transform: scale(1.1);
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
  font-weight: 600;
  padding: 3px 8px;
  border-radius: 10px;
  color: #fff;
  backdrop-filter: blur(8px);
  box-shadow: 0 2px 6px rgba(0,0,0,0.2);
  z-index: 2;
}
.play-badge.played {
  background: rgba(39, 174, 96, 0.85);
}
.play-badge.unplayed {
  background: rgba(127, 140, 141, 0.85);
}

.card-body {
  padding: 12px;
  flex: 1;
  display: flex;
  flex-direction: column;
  background: linear-gradient(to top, var(--card), transparent);
}
.card-name {
  font-size: 15px;
  font-weight: 700;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  color: var(--text);
  margin-bottom: 4px;
}
.card-desc {
  margin-top: 6px;
  height: 56px;
  overflow-y: auto;
  font-size: 12px;
  line-height: 1.5;
  color: var(--text-dim);
  word-break: break-all;
  padding-right: 6px;
}
.card-desc::-webkit-scrollbar {
  width: 4px;
}
.card-desc::-webkit-scrollbar-thumb {
  background: var(--border);
  border-radius: 2px;
}
.desc-text {
  white-space: pre-wrap;
}
.desc-empty {
  color: rgba(148, 163, 184, 0.5);
  font-style: italic;
}
.card-time {
  margin-top: 12px;
  font-size: 12px;
  color: var(--text-dim);
  flex-shrink: 0;
  display: flex;
  align-items: center;
  gap: 6px;
}
.card-time::before {
  content: '🕒';
  font-size: 10px;
}

.card-actions {
  display: flex;
  gap: 8px;
  padding: 12px 16px 16px;
  opacity: 0.8;
  transition: opacity 0.3s;
}
.game-card:hover .card-actions {
  opacity: 1;
}
.btn-small {
  padding: 8px 12px;
  font-size: 13px;
  flex: 1;
  border-radius: 8px;
  font-weight: 500;
}
.btn-open {
  background: var(--accent);
  color: #fff;
}
.btn-open:hover {
  background: var(--accent-hover);
  box-shadow: 0 4px 12px rgba(52, 152, 219, 0.3);
}
.btn-del {
  color: var(--danger);
  background: rgba(239, 68, 68, 0.1);
}
.btn-del:hover {
  background: var(--danger);
  color: #fff;
  box-shadow: 0 4px 12px rgba(239, 68, 68, 0.3);
}
</style>
