<template>
  <div class="app">
    <aside class="sidebar">
      <div class="logo">游戏管理器</div>
      <nav class="category-list">
        <div
          class="category-item"
          :class="{ active: !currentCategoryId && !seriesView }"
          @click="selectAll"
        >
          全部
        </div>
        <div
          v-for="c in visibleCategories"
          :key="c.Id"
          class="category-item"
          :class="{ active: currentCategoryId === c.Id }"
          @click="selectCategory(c.Id)"
        >
          <span class="cat-name">{{ c.Name }}</span>
          <span class="cat-num">{{ c.Num }}</span>
        </div>
      </nav>
    </aside>

    <main class="main">
      <header class="toolbar">
        <div class="title-area">
          <template v-if="seriesView">
            <button class="btn-back" @click="backToList">← 返回</button>
            <h1 class="series-title">{{ seriesView }}</h1>
          </template>
          <h1 v-else>游戏库</h1>
        </div>

        <div class="filters" v-if="!seriesView">
          <input
            v-model="keyword"
            class="filter-keyword"
            placeholder="搜索游戏名称..."
            @keyup.enter="loadGames"
          />
          <input
            v-model="seriesKeyword"
            class="filter-keyword"
            placeholder="搜索系列..."
            @keyup.enter="loadGames"
          />
          <input
            v-model="descKeyword"
            class="filter-keyword"
            placeholder="搜索描述..."
            @keyup.enter="loadGames"
          />
          <select v-model="playFilter" class="filter-select" @change="loadGames">
            <option :value="null">全部状态</option>
            <option :value="true">已玩</option>
            <option :value="false">未玩</option>
          </select>
          <input v-model="dateFrom" type="date" class="filter-date" @change="loadGames" />
          <span class="date-sep">至</span>
          <input v-model="dateTo" type="date" class="filter-date" @change="loadGames" />
          <button class="btn" @click="loadGames">搜索</button>
        </div>

        <button class="btn btn-primary btn-add" @click="showAdd = true">添加游戏</button>
      </header>

      <div class="content">
        <div v-if="seriesView" class="game-grid">
          <GameCard
            v-for="g in seriesGames"
            :key="g.Id"
            :game="g"
            @edit="openEdit"
            @delete="onDelete"
          />
          <div v-if="seriesGames.length === 0" class="empty">该系列暂无游戏</div>
        </div>
        <template v-else>
          <div v-if="loading" class="empty">加载中...</div>
          <div v-else-if="seriesCards.length === 0 && normalGames.length === 0" class="empty">
            暂无游戏数据
          </div>
          <div v-else class="game-grid">
            <template v-for="item in gridItems" :key="item.type === 'series' ? 'series-' + item.group.name : item.game.Id">
              <SeriesCard
                v-if="item.type === 'series'"
                :group="item.group"
                @open="openSeries"
              />
              <GameCard
                v-else
                :game="item.game"
                @edit="openEdit"
                @delete="onDelete"
              />
            </template>
          </div>
        </template>
      </div>
    </main>

    <AddDialog v-if="showAdd" @close="showAdd = false" @added="onAdded" />
    <EditDialog
      v-if="editing"
      :game="editing"
      :categories="categories"
      @close="editing = null"
      @saved="onSaved"
    />
    <DeleteConfirmDialog
      v-if="deleting"
      :game-name="deleting.NickName || deleting.Name"
      @cancel="deleting = null"
      @confirm="doDelete"
    />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { api } from './api'
import type { CategoryDTO, GameDTO, SeriesGroup } from './types'
import GameCard from './components/GameCard.vue'
import SeriesCard from './components/SeriesCard.vue'
import AddDialog from './components/AddDialog.vue'
import EditDialog from './components/EditDialog.vue'
import DeleteConfirmDialog from './components/DeleteConfirmDialog.vue'

const categories = ref<CategoryDTO[]>([])
const games = ref<GameDTO[]>([])
const loading = ref(false)

const currentCategoryId = ref('')
const seriesView = ref('')
const seriesGames = ref<GameDTO[]>([])

const keyword = ref('')
const seriesKeyword = ref('')
const descKeyword = ref('')
const playFilter = ref<boolean | null>(null)
const dateFrom = ref('')
const dateTo = ref('')

const showAdd = ref(false)
const editing = ref<GameDTO | null>(null)
const deleting = ref<GameDTO | null>(null)

const allGamesInSeries = new Map<string, GameDTO[]>()

const visibleCategories = computed(() =>
  categories.value.filter((c) => c.Num > 0),
)

const seriesCards = computed<SeriesGroup[]>(() => {
  const map = new Map<string, GameDTO[]>()
  const normals: GameDTO[] = []
  for (const g of games.value) {
    const s = (g.Series || '').trim()
    if (s) {
      if (!map.has(s)) map.set(s, [])
      map.get(s)!.push(g)
    } else {
      normals.push(g)
    }
  }
  allGamesInSeries.clear()
  const cards: SeriesGroup[] = []
  for (const [name, gs] of map) {
    allGamesInSeries.set(name, gs)
    cards.push({
      name,
      games: gs,
      latest: gs.reduce((a, b) =>
        new Date(a.InsertTime) > new Date(b.InsertTime) ? a : b,
      ),
    })
  }
  return cards
})

const normalGames = computed<GameDTO[]>(() =>
  games.value.filter((g) => !(g.Series || '').trim()),
)

type GridItem = { type: 'series'; group: SeriesGroup } | { type: 'game'; game: GameDTO }

const gridItems = computed<GridItem[]>(() => {
  const items: GridItem[] = [
    ...seriesCards.value.map((s) => ({ type: 'series' as const, group: s })),
    ...normalGames.value.map((g) => ({ type: 'game' as const, game: g })),
  ]
  return items.sort((a, b) => {
    const ta =
      a.type === 'series' ? new Date(a.group.latest.InsertTime).getTime() : new Date(a.game.InsertTime).getTime()
    const tb =
      b.type === 'series' ? new Date(b.group.latest.InsertTime).getTime() : new Date(b.game.InsertTime).getTime()
    return tb - ta
  })
})

function buildCondition() {
  return {
    Name: keyword.value.trim(),
    Series: seriesKeyword.value.trim(),
    Description: descKeyword.value.trim(),
    Id: currentCategoryId.value,
    IsPlay: playFilter.value,
    InsertTimeStart: dateFrom.value ? localIso(dateFrom.value, 0) : null,
    InsertTimeEnd: dateTo.value ? localIso(dateTo.value, 1) : null,
  }
}

function localIso(dateStr: string, dayEnd: number): string {
  const [y, m, d] = dateStr.split('-').map(Number)
  const date = new Date(y, m - 1, d)
  if (dayEnd) date.setHours(23, 59, 59, 999)
  return date.toISOString()
}

async function loadGames() {
  loading.value = true
  try {
    games.value = await api.getGames(buildCondition())
  } catch (err: unknown) {
    alert('获取游戏列表失败：' + String(err))
  } finally {
    loading.value = false
  }
}

async function loadCategories() {
  try {
    categories.value = await api.getAllCategory()
  } catch (err: unknown) {
    alert('获取分类失败：' + String(err))
  }
}

function selectAll() {
  currentCategoryId.value = ''
  loadGames()
}

function selectCategory(id: string) {
  currentCategoryId.value = id
  loadGames()
}

function openSeries(name: string) {
  const gs = allGamesInSeries.get(name)
  if (!gs) return
  seriesView.value = name
  seriesGames.value = gs
}

function backToList() {
  seriesView.value = ''
  seriesGames.value = []
}

function openEdit(game: GameDTO) {
  editing.value = game
}

function onDelete(game: GameDTO) {
  deleting.value = game
}

function doDelete(delAll: boolean) {
  const game = deleting.value
  if (!game) return
  api
    .deleteGame(game, delAll)
    .then(() => {
      deleting.value = null
      alert('删除成功')
      onSaved()
    })
    .catch((err: unknown) => alert('删除失败：' + String(err)))
}

function onAdded() {
  loadGames()
  loadCategories()
}

function onSaved() {
  loadGames()
  loadCategories()
}

onMounted(() => {
  loadCategories()
  loadGames()
})
</script>

<style scoped>
.app {
  display: flex;
  height: 100vh;
  background: var(--bg);
}

.sidebar {
  width: 200px;
  flex-shrink: 0;
  background: var(--panel);
  border-right: 1px solid var(--border);
  display: flex;
  flex-direction: column;
  overflow-y: auto;
  box-shadow: 2px 0 10px rgba(0, 0, 0, 0.1);
  z-index: 10;
}

.logo {
  padding: 20px 16px;
  font-size: 18px;
  font-weight: 800;
  color: var(--accent);
  border-bottom: 1px solid var(--border);
  letter-spacing: 0.5px;
  display: flex;
  align-items: center;
  gap: 12px;
}
.logo::before {
  content: '';
  display: block;
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background: var(--accent);
  box-shadow: 0 0 8px var(--accent);
}

.category-list {
  padding: 16px 12px;
  flex: 1;
}

.category-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 14px;
  border-radius: 8px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  color: var(--text-dim);
  margin-bottom: 4px;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  overflow: hidden;
}
.category-item::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 3px;
  background: var(--accent);
  transform: scaleY(0);
  transition: transform 0.2s ease;
  border-radius: 0 4px 4px 0;
}
.category-item:hover {
  background: var(--panel-hover);
  color: var(--text);
  transform: translateX(4px);
}
.category-item.active {
  background: var(--series-bg);
  color: var(--accent);
}
.category-item.active::before {
  transform: scaleY(1);
}
.category-item .cat-num {
  font-size: 12px;
  background: var(--border);
  padding: 2px 8px;
  border-radius: 12px;
  transition: all 0.2s;
}
.category-item.active .cat-num {
  background: var(--accent);
  color: white;
}

.main {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
  position: relative;
}

.main::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 200px;
  background: radial-gradient(circle at 50% 0, rgba(52, 152, 219, 0.05) 0%, transparent 70%);
  pointer-events: none;
  z-index: 0;
}

.toolbar {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 14px 20px;
  border-bottom: 1px solid var(--border);
  flex-wrap: wrap;
  background: rgba(36, 51, 77, 0.8);
  backdrop-filter: blur(12px);
  position: relative;
  z-index: 2;
}

.title-area {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-shrink: 0;
}
.title-area h1 {
  font-size: 22px;
  font-weight: 700;
  color: var(--text);
  letter-spacing: -0.5px;
}
.series-title {
  color: var(--series) !important;
  display: flex;
  align-items: center;
  gap: 8px;
}
.series-title::before {
  content: '';
  display: block;
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--series);
}
.btn-back {
  padding: 6px 14px;
  font-size: 13px;
  background: var(--panel);
  border: 1px solid var(--border);
  border-radius: 6px;
  color: var(--text-dim);
  transition: all 0.2s;
}
.btn-back:hover {
  background: var(--panel-hover);
  color: var(--text);
  transform: translateX(-2px);
}

.filters {
  display: flex;
  align-items: center;
  gap: 12px;
  flex: 1;
  flex-wrap: wrap;
}
.filter-keyword {
  width: 220px;
}
.filter-select {
  width: 120px;
}
.filter-date {
  width: 140px;
}
.date-sep {
  color: var(--text-dim);
  font-size: 13px;
}

.btn-add {
  flex-shrink: 0;
  box-shadow: 0 4px 12px rgba(52, 152, 219, 0.2);
}

.content {
  flex: 1;
  overflow-y: auto;
  padding: 24px;
  position: relative;
  z-index: 1;
}

.game-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(210px, 1fr));
  gap: 16px;
  animation: fadeIn 0.4s ease-out;
  align-items: start;
}

.empty {
  text-align: center;
  color: var(--text-dim);
  padding: 100px 0;
  font-size: 16px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
}
.empty::before {
  content: '🎮';
  font-size: 48px;
  opacity: 0.5;
  filter: grayscale(1);
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>
