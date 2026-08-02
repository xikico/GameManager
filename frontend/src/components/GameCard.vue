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
  border-radius: 10px;
  overflow: hidden;
  cursor: pointer;
  transition: transform 0.12s, box-shadow 0.12s, border-color 0.12s;
  display: flex;
  flex-direction: column;
}
.game-card:hover {
  transform: translateY(-3px);
  background: var(--card-hover);
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.35);
  border-color: var(--accent);
}
.game-card.played {
  border-left: 3px solid var(--played);
}
.game-card:not(.played) {
  border-left: 3px solid var(--unplayed);
}

.card-icon {
  position: relative;
  height: 110px;
  background: var(--bg-soft);
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
}
.card-icon img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.icon-placeholder {
  font-size: 40px;
  font-weight: 700;
  color: var(--text-dim);
}

.play-badge {
  position: absolute;
  top: 8px;
  right: 8px;
  font-size: 11px;
  padding: 2px 8px;
  border-radius: 10px;
  color: #fff;
}
.play-badge.played {
  background: var(--played);
}
.play-badge.unplayed {
  background: var(--unplayed);
}

.card-body {
  padding: 10px 12px 6px;
  flex: 1;
}
.card-name {
  font-size: 14px;
  font-weight: 600;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.card-time {
  margin-top: 4px;
  font-size: 12px;
  color: var(--text-dim);
}

.card-actions {
  display: flex;
  gap: 6px;
  padding: 8px 12px 10px;
}
.btn-small {
  padding: 5px 8px;
  font-size: 12px;
  flex: 1;
}
.btn-open {
  background: var(--accent);
  color: #fff;
}
.btn-open:hover {
  background: var(--accent-hover);
}
.btn-del {
  color: var(--danger);
}
.btn-del:hover {
  background: var(--danger);
  color: #fff;
}
</style>
