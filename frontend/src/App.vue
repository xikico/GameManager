<template>
  <div class="app">
    <aside class="sidebar">
      <div class="brand">
        <span>
          <strong>游戏管理器</strong>
          <small>PERSONAL GAME LIBRARY</small>
        </span>
      </div>
      <div class="nav-label">游戏分区</div>
      <nav class="category-list" aria-label="游戏分区">
        <button
          v-for="c in visibleCategories"
          :key="c.Id"
          class="category-item"
          :class="{ active: currentCategoryId === c.Id }"
          :aria-current="currentCategoryId === c.Id ? 'page' : undefined"
          :title="c.Name"
          @click="selectCategory(c.Id)"
        >
          <span class="cat-name">{{ c.Name }}</span>
          <span class="cat-num">{{ c.Num }}</span>
        </button>
      </nav>
      <div class="sidebar-footer">
        <button class="sidebar-action" @click="showSettings = true">
          <svg viewBox="0 0 24 24" aria-hidden="true"><path d="M12 15.5a3.5 3.5 0 1 0 0-7 3.5 3.5 0 0 0 0 7Z"/><path d="M19.4 15a1.7 1.7 0 0 0 .34 1.88l.06.06-2.86 2.86-.06-.06A1.7 1.7 0 0 0 15 19.4a1.7 1.7 0 0 0-1 .6 1.7 1.7 0 0 0-.4 1.1V21H9.5v-.1A1.7 1.7 0 0 0 8.4 19.4a1.7 1.7 0 0 0-1.88.34l-.06.06-2.86-2.86.06-.06A1.7 1.7 0 0 0 4 15a1.7 1.7 0 0 0-.6-1 1.7 1.7 0 0 0-1.1-.4H2V9.5h.3A1.7 1.7 0 0 0 4 8.4a1.7 1.7 0 0 0-.34-1.88l-.06-.06L6.46 3.6l.06.06A1.7 1.7 0 0 0 8.4 4 1.7 1.7 0 0 0 9.5 2.3V2h4.1v.3A1.7 1.7 0 0 0 15 4a1.7 1.7 0 0 0 1.88-.34l.06-.06 2.86 2.86-.06.06A1.7 1.7 0 0 0 19.4 8.4a1.7 1.7 0 0 0 1.7 1.1h.3v4.1h-.3a1.7 1.7 0 0 0-1.7 1.4Z"/></svg>
          <span>设置</span>
          <kbd>Ctrl ,</kbd>
        </button>
      </div>
    </aside>

    <main class="main">
      <header class="page-header">
        <div class="header-main">
          <div class="title-area">
          <template v-if="seriesView">
              <button class="btn-back" @click="backToList" aria-label="返回分区">←</button>
              <div>
                <span class="breadcrumb">{{ currentCategory?.Name }} / 系列</span>
                <h1 class="series-title">{{ seriesView }}</h1>
              </div>
          </template>
            <div v-else>
              <span class="breadcrumb">游戏库 / 分区</span>
              <div class="title-line">
                <h1>{{ currentCategory?.Name || '游戏库' }}</h1>
                <span class="result-count">{{ games.length }} 款游戏</span>
              </div>
            </div>
          </div>

          <button class="btn btn-primary btn-add" @click="showAdd = true">
            <span class="plus" aria-hidden="true">+</span>
            添加游戏
          </button>
        </div>

        <form v-if="!seriesView" class="filter-panel" @submit.prevent="loadGames">
          <label class="search-box">
            <svg viewBox="0 0 24 24" aria-hidden="true"><circle cx="11" cy="11" r="7"/><path d="m20 20-4-4"/></svg>
            <input v-model="keyword" placeholder="搜索游戏名称" aria-label="搜索游戏名称" />
          </label>
          <select v-model="playFilter" class="filter-select" aria-label="游玩状态">
            <option :value="null">全部状态</option>
            <option :value="true">已玩</option>
            <option :value="false">未玩</option>
          </select>
          <button type="button" class="btn btn-filter" :class="{ active: showAdvanced }" @click="showAdvanced = !showAdvanced">
            高级筛选
            <span v-if="advancedFilterCount" class="filter-count">{{ advancedFilterCount }}</span>
          </button>
          <button type="submit" class="btn btn-search" :disabled="loading">搜索</button>
          <button v-if="hasFilters" type="button" class="btn-clear" @click="clearFilters">清除</button>

          <div v-if="showAdvanced" class="advanced-filters">
            <label><span>系列</span><input v-model="seriesKeyword" placeholder="输入系列名称" /></label>
            <label><span>描述</span><input v-model="descKeyword" placeholder="搜索描述内容" /></label>
            <label><span>加入日期从</span><input v-model="dateFrom" type="date" /></label>
            <label><span>到</span><input v-model="dateTo" type="date" /></label>
          </div>
        </form>
      </header>

      <div ref="contentEl" class="content" :aria-busy="loading">
        <div v-if="loading && games.length" class="refresh-indicator" role="status">
          <span></span>正在更新
        </div>
        <div v-if="seriesView" class="game-grid">
          <GameCard
            v-for="g in seriesGames"
            :key="g.Id"
            :game="g"
            @edit="openEdit"
            @launch-error="showToast('启动失败：' + $event)"
          />
          <div v-if="seriesGames.length === 0" class="empty-state">
            <div class="empty-mark">S</div><h2>该系列暂无游戏</h2><p>游戏可能已被移动到其他系列。</p>
            <button class="btn" @click="backToList">返回分区</button>
          </div>
        </div>
        <template v-else>
          <div v-if="loading && games.length === 0" class="empty-state loading-state" role="status">
            <div class="loader"></div><h2>正在加载游戏库</h2><p>正在整理当前分区的内容...</p>
          </div>
          <div v-else-if="seriesCards.length === 0 && normalGames.length === 0" class="empty-state">
            <div class="empty-mark">G</div>
            <h2>{{ hasFilters ? '没有符合条件的游戏' : '这个分区还是空的' }}</h2>
            <p>{{ hasFilters ? '尝试减少筛选条件，或清除筛选重新查看。' : '添加第一款游戏，开始整理你的收藏。' }}</p>
            <button v-if="hasFilters" class="btn" @click="clearFilters">清除筛选</button>
            <button v-else class="btn btn-primary" @click="showAdd = true">添加游戏</button>
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
                @launch-error="showToast('启动失败：' + $event)"
              />
            </template>
          </div>
        </template>
      </div>
    </main>

    <Transition name="toast">
      <div v-if="toast" class="toast" role="status">{{ toast }}</div>
    </Transition>

    <AddDialog v-if="showAdd" @close="showAdd = false" @added="onAdded" />
    <SettingsDialog v-if="showSettings" @close="showSettings = false" />
    <EditDialog
      v-if="editing"
      :game="editing"
      :categories="categories"
      @close="editing = null"
      @saved="onSaved"
    />
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref } from 'vue'
import { api } from './api'
import type { CategoryDTO, GameDTO, SeriesGroup } from './types'
import GameCard from './components/GameCard.vue'
import SeriesCard from './components/SeriesCard.vue'
import AddDialog from './components/AddDialog.vue'
import EditDialog from './components/EditDialog.vue'
import SettingsDialog from './components/SettingsDialog.vue'

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
const showAdvanced = ref(false)
const toast = ref('')
let loadSequence = 0
let toastTimer: ReturnType<typeof setTimeout> | undefined

const showAdd = ref(false)
const showSettings = ref(false)
const editing = ref<GameDTO | null>(null)
const contentEl = ref<HTMLElement | null>(null)

const allGamesInSeries = new Map<string, GameDTO[]>()

const visibleCategories = computed(() =>
  categories.value.filter((c) => c.Num > 0),
)

const currentCategory = computed(() =>
  categories.value.find((category) => category.Id === currentCategoryId.value),
)

const advancedFilterCount = computed(() =>
  [seriesKeyword.value, descKeyword.value, dateFrom.value, dateTo.value].filter(Boolean).length,
)

const hasFilters = computed(() => Boolean(
  keyword.value.trim() || advancedFilterCount.value || playFilter.value !== null,
))

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
  const sequence = ++loadSequence
  loading.value = true
  try {
    const result = await api.getGames(buildCondition())
    if (sequence === loadSequence) games.value = result
  } catch (err: unknown) {
    showToast('获取游戏列表失败，请稍后重试')
  } finally {
    if (sequence === loadSequence) loading.value = false
  }
}

async function loadCategories() {
  try {
    categories.value = await api.getAllCategory()
  } catch (err: unknown) {
    showToast('获取分区失败：' + String(err))
  }
}

function selectCategory(id: string) {
  if (currentCategoryId.value === id) return
  currentCategoryId.value = id
  seriesView.value = ''
  loadGames()
}

function clearFilters() {
  keyword.value = ''
  seriesKeyword.value = ''
  descKeyword.value = ''
  playFilter.value = null
  dateFrom.value = ''
  dateTo.value = ''
  loadGames()
}

function showToast(message: string) {
  toast.value = message
  if (toastTimer) clearTimeout(toastTimer)
  toastTimer = setTimeout(() => { toast.value = '' }, 2600)
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

function onAdded() {
  refreshLibrary(false)
}

async function onSaved() {
  await refreshLibrary(true)
}

function ensureSelectedCategory() {
  if (!visibleCategories.value.some((category) => category.Id === currentCategoryId.value)) {
    currentCategoryId.value = visibleCategories.value[0]?.Id ?? ''
  }
}

async function refreshLibrary(preserveScroll: boolean) {
  const scrollTop = preserveScroll ? contentEl.value?.scrollTop ?? 0 : 0
  await loadCategories()
  ensureSelectedCategory()
  await loadGames()
  if (seriesView.value) {
    seriesGames.value = games.value.filter((game) => game.Series.trim() === seriesView.value)
  }
  if (preserveScroll) {
    await nextTick()
    if (contentEl.value) contentEl.value.scrollTop = scrollTop
  }
}

onMounted(async () => {
  window.addEventListener('keydown', onGlobalKeydown)
  await loadCategories()
  ensureSelectedCategory()
  await loadGames()
})

onBeforeUnmount(() => window.removeEventListener('keydown', onGlobalKeydown))

function onGlobalKeydown(event: KeyboardEvent) {
  if (event.ctrlKey && event.key === ',') {
    event.preventDefault()
    showSettings.value = true
  }
}
</script>

<style scoped>
.app { display: flex; height: 100vh; background: var(--bg); }
.sidebar { width: 218px; flex-shrink: 0; background: linear-gradient(180deg,rgba(15,27,44,.98),rgba(10,20,33,.98)); border-right: 1px solid var(--border); display: flex; flex-direction: column; box-shadow: 12px 0 38px rgba(0,5,14,.22),inset -1px 0 rgba(255,255,255,.018); z-index: 10; }
.brand { min-height: 82px; padding: 19px 18px 17px; display: flex; align-items: center; }
.brand strong,.brand small { display: block; }
.brand strong { font-size: 17px; font-weight: 720; letter-spacing: .3px; }
.brand small { color: #617796; font-size: 8px; font-weight: 700; letter-spacing: 1.35px; margin-top: 5px; }
.nav-label { padding: 15px 18px 8px; color: #6f83a7; font-size: 10px; font-weight: 800; letter-spacing: 1.4px; border-top: 1px solid rgba(255,255,255,.035); }
.category-list { padding: 4px 10px 16px; flex: 1; overflow-y: auto; }
.category-item { width: 100%; display: flex; justify-content: space-between; align-items: center; padding: 10px 11px; border-radius: 10px; font-size: 13px; font-weight: 600; color: var(--text-dim); background: transparent; border: 1px solid transparent; margin-bottom: 3px; transition: background .18s,color .18s,border-color .18s; }
.category-item:hover { background: var(--panel-hover); color: var(--text); }
.category-item.active { background: linear-gradient(90deg,rgba(76,168,232,.16),rgba(76,168,232,.06)); border-color: rgba(76,168,232,.13); color: #78c6f4; box-shadow: inset 2px 0 #59b6ee,0 5px 18px rgba(0,8,20,.12); }
.cat-name { overflow: hidden; text-overflow: ellipsis; white-space: nowrap; text-align: left; }
.cat-num { min-width: 27px; flex-shrink: 0; text-align: center; font-size: 12px; background: rgba(255,255,255,.06); padding: 2px 7px; border-radius: 99px; }
.category-item.active .cat-num { background: var(--accent); color: #fff; }
.sidebar-footer { padding: 10px; border-top: 1px solid var(--border); }
.sidebar-action { width: 100%; display: flex; align-items: center; gap: 10px; padding: 10px 11px; border-radius: 9px; background: transparent; color: var(--text-dim); font-size: 13px; }
.sidebar-action:hover { background: var(--panel-hover); color: var(--text); }
.sidebar-action svg { width: 18px; fill: none; stroke: currentColor; stroke-width: 1.7; }
.sidebar-action kbd { margin-left: auto; color: #627698; font: 10px inherit; }
.main { flex: 1; display: flex; flex-direction: column; min-width: 0; position: relative; }
.main::before { content: ''; position: absolute; inset: 0 0 auto; height: 340px; background: radial-gradient(ellipse at 32% -20%,rgba(54,149,211,.13),transparent 58%),linear-gradient(180deg,rgba(255,255,255,.012),transparent); pointer-events: none; }
.page-header { padding: 19px 26px 15px; border-bottom: 1px solid var(--border); background: linear-gradient(180deg,rgba(14,26,42,.88),rgba(12,23,38,.76)); backdrop-filter: blur(22px) saturate(1.2); box-shadow: 0 10px 36px rgba(0,5,14,.14),inset 0 -1px rgba(255,255,255,.015); position: relative; z-index: 2; }
.header-main { display: flex; align-items: center; justify-content: space-between; gap: 20px; min-height: 48px; }
.title-area { display: flex; align-items: center; gap: 12px; min-width: 0; }
.breadcrumb { display: block; color: #7185a7; font-size: 10px; font-weight: 700; letter-spacing: 1px; margin-bottom: 3px; }
.title-line { display: flex; align-items: baseline; gap: 12px; }
.title-area h1 { font-family: 'Segoe UI Variable Display','Microsoft YaHei UI',sans-serif; font-size: 25px; font-weight: 720; color: var(--text); letter-spacing: -.6px; text-shadow: 0 1px 20px rgba(183,222,255,.06); }
.result-count { color: #7489a6; font-size: 11px; padding: 3px 8px; border: 1px solid var(--border); border-radius: 99px; background: rgba(255,255,255,.018); }
.series-title { color: #5fb3e8 !important; }
.btn-back { width: 34px; height: 34px; font-size: 18px; background: var(--panel); border: 1px solid var(--border); border-radius: 8px; color: var(--text-dim); transition: background .2s,color .2s,transform .2s; }
.btn-back:hover { background: var(--panel-hover); color: var(--text); transform: translateX(-2px); }
.btn-add { flex-shrink: 0; padding: 10px 17px; box-shadow: inset 0 1px rgba(255,255,255,.2),0 9px 24px rgba(23,107,163,.26); }
.plus { font-size: 19px; line-height: 12px; font-weight: 300; }
.filter-panel { display: flex; align-items: center; gap: 9px; flex-wrap: wrap; margin-top: 15px; padding-top: 14px; border-top: 1px solid rgba(153,185,221,.09); }
.search-box { width: min(360px,38vw); position: relative; display: flex; align-items: center; }
.search-box svg { position: absolute; left: 12px; width: 17px; fill: none; stroke: #7185a7; stroke-width: 1.8; pointer-events: none; }
.search-box input { width: 100%; padding-left: 38px; background: rgba(4,12,23,.46); box-shadow: inset 0 1px 7px rgba(0,5,14,.16); }
.filter-select { width: 116px; background: rgba(4,12,23,.46); }
.btn-filter,.btn-search { padding: 9px 14px; }
.btn-filter { border: 1px solid var(--border); background: transparent; }
.btn-filter.active { color: var(--accent); border-color: rgba(52,152,219,.5); background: rgba(52,152,219,.08); }
.filter-count { min-width: 18px; height: 18px; border-radius: 9px; display: grid; place-items: center; background: var(--accent); color: #fff; font-size: 10px; }
.btn-search { background: var(--panel-hover); }
.btn-clear { padding: 8px 5px; background: transparent; color: var(--text-dim); font-size: 12px; }
.btn-clear:hover { color: var(--text); }
.advanced-filters { width: 100%; display: grid; grid-template-columns: 1fr 1fr 160px 160px; gap: 10px; padding: 14px; background: rgba(4,12,23,.3); border: 1px solid var(--border); border-radius: 13px; box-shadow: inset 0 1px rgba(255,255,255,.018); }
.advanced-filters label { display: flex; flex-direction: column; gap: 5px; }
.advanced-filters label span { color: var(--text-dim); font-size: 10px; font-weight: 700; }
.advanced-filters input { width: 100%; padding: 8px 11px; background: var(--bg); }
.content { flex: 1; overflow-y: auto; padding: 26px; position: relative; z-index: 1; }
.refresh-indicator { position: sticky; top: 0; z-index: 5; display: flex; align-items: center; gap: 8px; width: max-content; margin: -12px auto 10px; padding: 6px 11px; border-radius: 99px; color: var(--text-dim); background: rgba(24,35,54,.92); border: 1px solid var(--border); font-size: 11px; box-shadow: 0 5px 15px rgba(0,0,0,.18); }
.refresh-indicator span { width: 7px; height: 7px; border-radius: 50%; background: var(--accent); animation: pulse 1s infinite; }
.game-grid { display: grid; grid-template-columns: repeat(auto-fill,minmax(224px,1fr)); gap: 20px; animation: fadeIn .45s cubic-bezier(.16,1,.3,1); align-items: start; }
.empty-state { min-height: 55vh; grid-column: 1/-1; display: flex; flex-direction: column; align-items: center; justify-content: center; text-align: center; color: var(--text-dim); }
.empty-mark { width: 62px; height: 62px; display: grid; place-items: center; border-radius: 20px; margin-bottom: 17px; color: var(--accent); font-size: 24px; font-weight: 900; border: 1px solid rgba(52,152,219,.3); background: rgba(52,152,219,.08); box-shadow: inset 0 0 25px rgba(52,152,219,.08); }
.empty-state h2 { color: var(--text); font-size: 18px; margin-bottom: 7px; }
.empty-state p { font-size: 13px; margin-bottom: 18px; }
.loader { width: 34px; height: 34px; border-radius: 50%; margin-bottom: 18px; border: 2px solid var(--border); border-top-color: var(--accent); animation: spin .8s linear infinite; }
.toast { position: fixed; right: 24px; bottom: 24px; z-index: 300; padding: 12px 16px; max-width: 360px; border-radius: 12px; color: var(--text); background: rgba(16,29,46,.94); backdrop-filter: blur(16px); border: 1px solid var(--border-strong); box-shadow: 0 18px 45px rgba(0,3,10,.4),inset 0 1px rgba(255,255,255,.04); font-size: 13px; }
.toast-enter-active,.toast-leave-active { transition: opacity .2s,transform .2s; }
.toast-enter-from,.toast-leave-to { opacity: 0; transform: translateY(8px); }
@keyframes fadeIn { from { opacity: 0; transform: translateY(10px); } to { opacity: 1; transform: translateY(0); } }
@keyframes spin { to { transform: rotate(360deg); } }
@keyframes pulse { 50% { opacity: .35; } }
@media (max-width:1100px) { .advanced-filters { grid-template-columns: 1fr 1fr; } .search-box { width: min(320px,45vw); } }
@media (max-width:800px) { .sidebar { width: 168px; } .brand { padding-inline: 12px; } .brand small,.sidebar-action kbd { display: none; } .page-header { padding: 15px 16px 12px; } .search-box { width: 100%; } .filter-panel { align-items: stretch; } .advanced-filters { grid-template-columns: 1fr; } .content { padding: 16px; } }
</style>
