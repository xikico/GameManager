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
              <span>+ 添加图片</span>
            </label>
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
  padding: 20px 24px;
}
.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}
.modal-header h2 {
  font-size: 18px;
}

.edit-body {
  display: flex;
  gap: 20px;
}
.edit-form {
  flex: 1;
  min-width: 0;
}
.edit-imgs {
  width: 280px;
  flex-shrink: 0;
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
  border-radius: 8px;
  object-fit: contain;
  padding: 3px;
  border: 1px solid var(--border);
  background: var(--bg-soft);
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
  color: var(--text-dim);
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
  border-radius: 6px;
  overflow: hidden;
  background: var(--bg-soft);
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
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  color: var(--text-dim);
  cursor: pointer;
  transition: border-color 0.15s, color 0.15s;
}
.img-add:hover {
  border-color: var(--accent);
  color: var(--accent);
}

.modal-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 18px;
  padding-top: 14px;
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
