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
              <select v-model="form.CategoryId">
                <option v-for="c in categories" :key="c.Id" :value="c.Id">{{ c.Name }}</option>
              </select>
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
            <label class="img-add">
              <input type="file" accept="image/*" multiple hidden @change="onPickImgs" />
              <span>+ 选择图片</span>
            </label>
            <button class="img-add" type="button" @click="showBase64Input = !showBase64Input">
              &lt;/&gt; Base64
            </button>
          </div>
          <div v-if="showBase64Input" class="base64-panel">
            <label for="screenshot-base64">图片 Base64</label>
            <textarea
              id="screenshot-base64"
              v-model="base64Input"
              rows="4"
              placeholder="粘贴 data:image/...;base64,... 或纯 Base64 内容"
              @keydown.ctrl.enter="addBase64Img"
            ></textarea>
            <p v-if="base64Error" class="base64-error">{{ base64Error }}</p>
            <div class="base64-actions">
              <span>Ctrl + Enter 快速添加</span>
              <button class="btn btn-small" type="button" @click="showBase64Input = false">取消</button>
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
          <button class="btn btn-primary" :disabled="saving" @click="save">保存</button>
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
const showBase64Input = ref(false)
const base64Input = ref('')
const base64Error = ref('')

const form = reactive({
  IconPath: props.game.IconPath,
  Name: props.game.Name,
  NickName: props.game.NickName,
  Series: props.game.Series,
  Description: props.game.Description,
  Path: props.game.Path,
  StartPath: props.game.StartPath,
  CategoryId: props.game.Category.Id,
  IsPlay: props.game.IsPlay,
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

function onPickImgs(e: Event) {
  const files = (e.target as HTMLInputElement).files
  if (!files) return
  for (const file of Array.from(files)) {
    const reader = new FileReader()
    reader.onload = () => {
      const url = reader.result as string
      if (!newImgs.value.includes(url)) newImgs.value.push(url)
    }
    reader.readAsDataURL(file)
  }
  ;(e.target as HTMLInputElement).value = ''
}

function addBase64Img() {
  base64Error.value = ''
  try {
    const image = normalizeBase64Image(base64Input.value)
    if (!displayImgs.value.includes(image)) newImgs.value.push(image)
    base64Input.value = ''
    showBase64Input.value = false
  } catch (err: unknown) {
    base64Error.value = err instanceof Error ? err.message : 'Base64 图片格式无效'
  }
}

function normalizeBase64Image(value: string): string {
  const input = value.trim()
  if (!input) throw new Error('请粘贴图片 Base64 内容')

  const dataUrl = input.match(/^data:(image\/[\w.+-]+);base64,([\s\S]+)$/i)
  const mime = dataUrl?.[1].toLowerCase()
  const payload = (dataUrl?.[2] ?? input).replace(/\s/g, '')
  if (!payload || !/^[A-Za-z0-9+/]*={0,2}$/.test(payload)) {
    throw new Error('Base64 内容包含无效字符')
  }

  let binary: string
  try {
    binary = atob(payload)
  } catch {
    throw new Error('无法解析 Base64 内容')
  }
  if (!binary) throw new Error('Base64 图片内容为空')

  const detectedMime = detectImageMime(binary)
  if (!mime && !detectedMime) {
    throw new Error('无法识别图片格式，请使用 PNG、JPEG、GIF、WebP、BMP 或 ICO')
  }
  return `data:${mime ?? detectedMime};base64,${payload}`
}

function detectImageMime(binary: string): string | null {
  const byte = (index: number) => binary.charCodeAt(index)
  const text = (start: number, end: number) => binary.slice(start, end)
  if (byte(0) === 0x89 && text(1, 4) === 'PNG') return 'image/png'
  if (byte(0) === 0xff && byte(1) === 0xd8 && byte(2) === 0xff) return 'image/jpeg'
  if (text(0, 4) === 'GIF8') return 'image/gif'
  if (text(0, 4) === 'RIFF' && text(8, 12) === 'WEBP') return 'image/webp'
  if (text(0, 2) === 'BM') return 'image/bmp'
  if (byte(0) === 0 && byte(1) === 0 && byte(2) === 1 && byte(3) === 0) return 'image/x-icon'
  return null
}

function removeImg(img: string) {
  if (newImgs.value.includes(img)) {
    newImgs.value = newImgs.value.filter((i) => i !== img)
  } else if (!removedImgs.value.includes(img)) {
    removedImgs.value.push(img)
  }
}

async function save() {
  if (!form.Name.trim()) {
    alert('名称不能为空')
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
