<template>
  <div class="modal-mask" @click.self="close">
    <div class="modal edit-modal">
      <div class="modal-header">
        <h2>编辑游戏</h2>
        <button class="btn btn-small" @click="close">关闭</button>
      </div>

      <div class="edit-body">
        <div class="edit-form">
          <div class="icon-row">
            <img v-if="iconSrc" class="icon-preview" :src="iconSrc" alt="图标" />
            <div v-else class="icon-preview icon-empty">无图标</div>
            <div class="icon-input">
              <label>图标路径</label>
              <input v-model="form.IconPath" placeholder="可留空，支持本地路径或粘贴图片 base64" />
            </div>
          </div>

          <div class="field">
            <label>名称 *</label>
            <input v-model="form.Name" placeholder="游戏名称" />
          </div>
          <div class="field">
            <label>别名</label>
            <input v-model="form.NickName" placeholder="可留空" />
          </div>
          <div class="field">
            <label>系列</label>
            <input v-model="form.Series" placeholder="可留空" />
          </div>
          <div class="field">
            <label>描述</label>
            <textarea v-model="form.Description" rows="3" placeholder="可留空"></textarea>
          </div>
          <div class="field">
            <label>游戏路径</label>
            <input v-model="form.Path" placeholder="游戏主路径" />
          </div>
          <div class="field">
            <label>启动路径</label>
            <input v-model="form.StartPath" placeholder="可留空，留空则自动推断" />
          </div>
          <div class="field-row">
            <div class="field">
              <label>分类</label>
              <div class="category-combobox">
                <input
                  v-model="form.CategoryName"
                  placeholder="选择或输入新分类"
                  autocomplete="off"
                  @focus="openCategoryOptions"
                  @input="onCategoryInput"
                  @keydown="onCategoryKeydown"
                  @blur="closeCategoryOptions"
                />
                <button type="button" aria-label="展开分类" @mousedown.prevent="toggleCategoryOptions">
                  <svg viewBox="0 0 20 20" aria-hidden="true"><path d="m6 8 4 4 4-4" /></svg>
                </button>
                <div v-if="showCategoryOptions" class="category-options">
                  <div class="category-options-label">选择分类</div>
                  <button
                    v-for="category in filteredCategories"
                    :key="category.Id"
                    type="button"
                    :class="{ selected: category.Id === form.CategoryId }"
                    @mousedown.prevent="selectCategory(category)"
                  >
                    <span class="category-option-name">
                      <span class="category-check" aria-hidden="true">{{ category.Id === form.CategoryId ? '✓' : '' }}</span>
                      {{ category.Name }}
                    </span>
                    <small>{{ category.Num }} 款</small>
                  </button>
                  <div v-if="filteredCategories.length === 0 && form.CategoryName.trim()" class="category-create">
                    <span aria-hidden="true">+</span>
                    <span>保存时创建<strong>“{{ form.CategoryName.trim() }}”</strong></span>
                  </div>
                </div>
              </div>
            </div>
            <div class="field">
              <label>状态</label>
              <select v-model="form.IsPlay">
                <option :value="true">已玩</option>
                <option :value="false">未玩</option>
              </select>
            </div>
          </div>
        </div>

        <div class="edit-imgs">
          <div class="imgs-title">游戏截图（{{ displayImgs.length }}）</div>
          <div class="img-grid">
            <div v-for="(img, i) in displayImgs" :key="img + i" class="img-item">
              <img :src="img" alt="" @click="previewImg = img" />
              <button class="img-remove" @click="removeImg(img)">×</button>
            </div>
            <button class="img-add" type="button" @click="toggleImageAddPanel">
              <span>+ 添加截图</span>
            </button>
          </div>
          <div v-if="showImageAddPanel" class="base64-panel">
            <div class="base64-panel-title">添加截图</div>
            <label class="file-picker">
              <input type="file" accept="image/*" multiple hidden @change="onPickImgs" />
              <span>从本地选择图片</span>
            </label>
            <div class="add-divider"><span>或粘贴 Base64</span></div>
            <textarea
              id="screenshot-base64"
              v-model="base64Input"
              rows="4"
              placeholder="粘贴 data:image/...;base64,... 或纯 Base64 内容"
              @keydown="onBase64Keydown"
            ></textarea>
            <p v-if="base64Error" class="base64-error">{{ base64Error }}</p>
            <div class="base64-actions">
              <span>Ctrl + Enter 添加</span>
              <button class="btn btn-small" type="button" @click="showImageAddPanel = false">取消</button>
              <button class="btn btn-small btn-primary" type="button" @click="addBase64Img">添加截图</button>
            </div>
          </div>
        </div>
      </div>

      <div class="modal-footer">
        <div class="footer-left">
          <button class="btn" @click="openFolder">在资源管理器打开</button>
          <button class="btn btn-danger" @click="showDelete = true">删除</button>
        </div>
        <div class="footer-right">
          <button class="btn" @click="openGame">打开游戏</button>
          <button class="btn btn-primary" :disabled="saving || pendingImageReads > 0" @click="save">
            {{ pendingImageReads > 0 ? '正在读取图片...' : '保存' }}
          </button>
        </div>
      </div>
    </div>

    <div v-if="previewImg" class="img-lightbox" @click="previewImg = ''">
      <img :src="previewImg" alt="" @click.stop />
    </div>

    <DeleteConfirmDialog
      v-if="showDelete"
      :game-name="form.Name"
      @cancel="showDelete = false"
      @confirm="doDelete"
    />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { api } from '../api'
import type { CategoryDTO, GameDTO } from '../types'
import { normalizeImageDataUrl } from '../utils/imageDataUrl'
import DeleteConfirmDialog from './DeleteConfirmDialog.vue'

const props = defineProps<{ game: GameDTO; categories: CategoryDTO[] }>()
const emit = defineEmits<{
  (e: 'close'): void
  (e: 'saved'): void
}>()

const saving = ref(false)
const iconSrc = ref('')
const oldImgs = ref<string[]>([])
const newImgs = ref<string[]>([])
const removedImgs = ref<string[]>([])
const previewImg = ref('')
const showDelete = ref(false)
const showImageAddPanel = ref(false)
const base64Input = ref('')
const base64Error = ref('')
const pendingImageReads = ref(0)
const showCategoryOptions = ref(false)
const categoryFilter = ref('')

const form = reactive({
  IconPath: props.game.IconPath,
  Name: props.game.Name,
  NickName: props.game.NickName,
  Series: props.game.Series,
  Description: props.game.Description,
  Path: props.game.Path,
  StartPath: props.game.StartPath,
  CategoryId: props.game.Category.Id,
  CategoryName: props.game.Category.Name,
  IsPlay: props.game.IsPlay,
})

const filteredCategories = computed(() => {
  const keyword = categoryFilter.value.trim().toLocaleLowerCase()
  if (!keyword) return props.categories
  return props.categories.filter((category) => category.Name.toLocaleLowerCase().includes(keyword))
})

const displayImgs = computed(() => [
  ...oldImgs.value.filter((i) => !removedImgs.value.includes(i)),
  ...newImgs.value,
])

onMounted(async () => {
  try {
    oldImgs.value = await api.getGameImgs(props.game.Id)
  } catch {
    oldImgs.value = []
  }
  if (props.game.IconPath) {
    if (props.game.IconPath.startsWith('data:')) {
      iconSrc.value = props.game.IconPath
    } else {
      iconSrc.value = await api.iconToBase64(props.game.IconPath)
    }
  }
})

async function onPickImgs(e: Event) {
  const files = (e.target as HTMLInputElement).files
  if (!files) return
  const input = e.target as HTMLInputElement
  pendingImageReads.value += files.length
  base64Error.value = ''
  try {
    const images = await Promise.all(Array.from(files, readImageFile))
    for (const image of images) {
      if (!displayImgs.value.includes(image)) newImgs.value.push(image)
    }
  } catch (err: unknown) {
    base64Error.value = err instanceof Error ? err.message : '读取图片失败'
  } finally {
    pendingImageReads.value -= files.length
    input.value = ''
  }
}

function readImageFile(file: File): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = () => {
      try {
        resolve(normalizeImageDataUrl(String(reader.result)))
      } catch (err) {
        reject(err)
      }
    }
    reader.onerror = () => reject(new Error(`读取图片“${file.name}”失败`))
    reader.readAsDataURL(file)
  })
}

function toggleImageAddPanel() {
  showImageAddPanel.value = !showImageAddPanel.value
  base64Error.value = ''
}

function onBase64Keydown(event: KeyboardEvent) {
  if ((event.ctrlKey || event.metaKey) && event.key === 'Enter') {
    event.preventDefault()
    event.stopPropagation()
    addBase64Img()
  }
}

function addBase64Img() {
  base64Error.value = ''
  try {
    const image = normalizeImageDataUrl(base64Input.value)
    if (!displayImgs.value.includes(image)) newImgs.value.push(image)
    base64Input.value = ''
    showImageAddPanel.value = false
  } catch (err: unknown) {
    base64Error.value = err instanceof Error ? err.message : 'Base64 图片格式无效'
  }
}

function removeImg(img: string) {
  if (newImgs.value.includes(img)) {
    newImgs.value = newImgs.value.filter((i) => i !== img)
  } else if (!removedImgs.value.includes(img)) {
    removedImgs.value.push(img)
  }
}

function onCategoryInput() {
  categoryFilter.value = form.CategoryName
  const existing = props.categories.find(
    (category) => category.Name.toLocaleLowerCase() === form.CategoryName.trim().toLocaleLowerCase(),
  )
  form.CategoryId = existing?.Id ?? ''
  showCategoryOptions.value = true
}

function openCategoryOptions() {
  categoryFilter.value = ''
  showCategoryOptions.value = true
}

function toggleCategoryOptions() {
  if (showCategoryOptions.value) {
    showCategoryOptions.value = false
    return
  }
  openCategoryOptions()
}

function selectCategory(category: CategoryDTO) {
  form.CategoryId = category.Id
  form.CategoryName = category.Name
  categoryFilter.value = ''
  showCategoryOptions.value = false
}

function closeCategoryOptions() {
  window.setTimeout(() => { showCategoryOptions.value = false }, 100)
}

function onCategoryKeydown(event: KeyboardEvent) {
  if (event.key === 'Escape') showCategoryOptions.value = false
  if (event.key === 'Enter' && filteredCategories.value.length === 1) {
    event.preventDefault()
    selectCategory(filteredCategories.value[0])
  }
}

async function save() {
  if (!form.Name.trim()) {
    alert('名称不能为空')
    return
  }
  if (!form.CategoryName.trim()) {
    alert('分类不能为空')
    return
  }
  saving.value = true
  try {
    await api.editGame({
      Id: props.game.Id,
      IconPath: form.IconPath.trim(),
      Name: form.Name.trim(),
      NickName: form.NickName.trim(),
      Series: form.Series.trim(),
      Description: form.Description.trim(),
      Path: form.Path.trim(),
      StartPath: form.StartPath.trim(),
      CategoryId: form.CategoryId,
      CategoryName: form.CategoryName.trim(),
      IsPlay: form.IsPlay,
      Imgs: displayImgs.value,
    })
    emit('saved')
    emit('close')
  } catch (err: unknown) {
    alert('编辑失败：' + String(err))
  } finally {
    saving.value = false
  }
}

function openFolder() {
  api.openFolder(form.Path).catch(() => {})
}

function openGame() {
  api
    .openGame({
      ...props.game,
      Path: form.Path,
      StartPath: form.StartPath,
    })
    .catch(() => {})
}

function doDelete(delAll: boolean) {
  showDelete.value = false
  api
    .deleteGame(props.game, delAll)
    .then(() => {
      alert('删除成功')
      emit('saved')
      emit('close')
    })
    .catch((err: unknown) => alert('删除失败：' + String(err)))
}

function close() {
  emit('close')
}
</script>

<style scoped>
.edit-modal {
  width: 860px;
  max-width: 94vw;
  padding: 22px 26px;
}
.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}
.modal-header h2 {
  font-family: 'Segoe UI Variable Display', 'Microsoft YaHei UI', sans-serif;
  font-size: 20px;
  font-weight: 720;
}

.edit-body {
  display: flex;
  gap: 24px;
}
.edit-form {
  flex: 1;
  min-width: 0;
}
.edit-imgs {
  width: 280px;
  flex-shrink: 0;
  padding: 14px;
  border: 1px solid var(--border);
  border-radius: 14px;
  background: rgba(4,12,23,.22);
}

.icon-row {
  display: flex;
  gap: 12px;
  margin-bottom: 14px;
  align-items: center;
}
.icon-preview {
  width: 64px;
  height: 64px;
  border-radius: 11px;
  object-fit: contain;
  padding: 3px;
  border: 1px solid var(--border);
  background: rgba(7,16,28,.65);
  box-shadow: inset 0 1px rgba(255,255,255,.025);
  flex-shrink: 0;
}
.icon-empty {
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  color: var(--text-dim);
}
.icon-input {
  flex: 1;
  min-width: 0;
}

.field {
  margin-bottom: 12px;
}
.field label {
  display: block;
  font-size: 12px;
  color: var(--text-dim);
  margin-bottom: 5px;
}
.field input,
.field textarea,
.field select {
  width: 100%;
}
.field-row {
  display: flex;
  gap: 12px;
}
.field-row .field {
  flex: 1;
}
.category-combobox {
  position: relative;
}
.category-combobox > input {
  padding-right: 38px;
}
.category-combobox > button {
  position: absolute;
  top: 1px;
  right: 1px;
  width: 36px;
  height: calc(100% - 2px);
  border-radius: 0 9px 9px 0;
  color: #748aa8;
  background: transparent;
  display: grid;
  place-items: center;
  border-left: 1px solid rgba(143,176,216,.1);
}
.category-combobox > button:hover {
  color: var(--accent);
  background: rgba(76,168,232,.08);
}
.category-combobox > button svg {
  width: 18px;
  height: 18px;
  fill: none;
  stroke: currentColor;
  stroke-width: 1.7;
  stroke-linecap: round;
  stroke-linejoin: round;
}
.category-options {
  position: absolute;
  top: calc(100% + 6px);
  left: 0;
  right: 0;
  z-index: 30;
  max-height: 230px;
  overflow-y: auto;
  padding: 7px;
  border: 1px solid var(--border-strong);
  border-radius: 12px;
  background: linear-gradient(155deg,rgba(19,35,55,.99),rgba(10,23,39,.99));
  box-shadow: 0 20px 46px rgba(0,4,12,.52),inset 0 1px rgba(255,255,255,.035);
  backdrop-filter: blur(20px) saturate(1.15);
}
.category-options-label {
  padding: 4px 9px 8px;
  color: #607895;
  font-size: 9px;
  font-weight: 800;
  letter-spacing: 1.3px;
  text-transform: uppercase;
}
.category-options > button {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
  min-height: 38px;
  padding: 7px 9px;
  border-radius: 8px;
  color: var(--text-dim);
  border: 1px solid transparent;
  background: transparent;
  font-size: 12px;
  font-weight: 600;
  text-align: left;
  transition: color .15s,background .15s,border-color .15s;
}
.category-options > button:hover {
  color: var(--text);
  background: rgba(76,168,232,.09);
  border-color: rgba(76,168,232,.1);
}
.category-options > button.selected {
  color: #87cef6;
  background: linear-gradient(90deg,rgba(76,168,232,.16),rgba(76,168,232,.06));
  border-color: rgba(76,168,232,.13);
}
.category-option-name {
  display: flex;
  align-items: center;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.category-check {
  width: 17px;
  flex-shrink: 0;
  color: var(--accent);
  font-size: 11px;
}
.category-options small {
  padding: 2px 7px;
  border-radius: 99px;
  color: #647b98;
  background: rgba(255,255,255,.035);
  font-size: 9px;
  font-weight: 700;
}
.category-create {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px;
  border: 1px dashed rgba(76,168,232,.22);
  border-radius: 8px;
  color: #8da2bb;
  background: rgba(76,168,232,.045);
  font-size: 11px;
  line-height: 1.4;
}
.category-create > span:first-child {
  width: 22px;
  height: 22px;
  display: grid;
  place-items: center;
  flex-shrink: 0;
  border-radius: 7px;
  color: #7bc8f4;
  background: rgba(76,168,232,.12);
  font-size: 16px;
}
.category-create strong {
  display: block;
  margin-top: 1px;
  color: #80cbf5;
  font-weight: 700;
}

.imgs-title {
  font-size: 13px;
  color: #9aacc3;
  font-weight: 700;
  margin-bottom: 8px;
}
.img-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 8px;
}
.img-item {
  position: relative;
  aspect-ratio: 4 / 3;
  border-radius: 8px;
  overflow: hidden;
  background: rgba(7,16,28,.68);
  border: 1px solid var(--border);
}
.img-item img {
  width: 100%;
  height: 100%;
  object-fit: contain;
  padding: 2px;
  cursor: zoom-in;
}
.img-remove {
  position: absolute;
  top: 3px;
  right: 3px;
  width: 18px;
  height: 18px;
  border-radius: 50%;
  background: rgba(231, 76, 60, 0.9);
  color: #fff;
  font-size: 12px;
  line-height: 1;
}
.img-add {
  aspect-ratio: 4 / 3;
  border: 1px dashed var(--border);
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  color: var(--text-dim);
  cursor: pointer;
  transition: border-color 0.15s, color 0.15s;
  background: transparent;
  font-family: inherit;
}
.img-add:hover {
  border-color: var(--accent);
  color: var(--accent);
}
.base64-panel {
  margin-top: 12px;
  padding: 12px;
  border: 1px solid var(--border);
  border-radius: 10px;
  background: rgba(3, 11, 21, .38);
}
.base64-panel-title {
  margin-bottom: 10px;
  color: #9aacc3;
  font-size: 12px;
  font-weight: 700;
}
.file-picker {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 34px;
  padding: 7px 10px;
  border: 1px solid var(--border);
  border-radius: 8px;
  color: var(--text-dim);
  background: rgba(255,255,255,.018);
  cursor: pointer;
  font-size: 12px;
  transition: border-color .15s, color .15s, background .15s;
}
.file-picker:hover {
  border-color: var(--accent);
  color: var(--accent);
  background: rgba(76,168,232,.07);
}
.add-divider {
  display: flex;
  align-items: center;
  gap: 8px;
  margin: 11px 0 8px;
  color: #617794;
  font-size: 10px;
}
.add-divider::before,
.add-divider::after {
  content: '';
  height: 1px;
  flex: 1;
  background: var(--border);
}
.base64-panel > label {
  display: block;
  margin-bottom: 6px;
  color: #9aacc3;
  font-size: 11px;
  font-weight: 700;
}
.base64-panel textarea {
  width: 100%;
  min-height: 88px;
  padding: 9px 10px;
  resize: vertical;
  font-family: Consolas, monospace;
  font-size: 11px;
  line-height: 1.45;
  word-break: break-all;
}
.base64-error {
  margin-top: 7px;
  color: var(--danger);
  font-size: 11px;
  line-height: 1.4;
}
.base64-actions {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 7px;
  margin-top: 9px;
}
.base64-actions > span {
  margin-right: auto;
  color: #617794;
  font-size: 9px;
}
.base64-actions .btn-small {
  padding: 6px 9px;
  font-size: 11px;
}

.modal-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 18px;
  padding-top: 16px;
  border-top: 1px solid var(--border);
}
.footer-left,
.footer-right {
  display: flex;
  gap: 10px;
}

.img-lightbox {
  position: fixed;
  inset: 0;
  background: rgba(10, 16, 26, 0.85);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 200;
  cursor: zoom-out;
}
.img-lightbox img {
  max-width: 92vw;
  max-height: 92vh;
  border-radius: 8px;
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.6);
}
</style>
