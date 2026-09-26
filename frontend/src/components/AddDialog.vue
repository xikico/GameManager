<template>
  <div class="modal-mask" @click.self="close">
    <div class="modal add-modal">
      <div class="modal-header">
        <h2>添加游戏</h2>
        <button class="btn btn-small" @click="close">关闭</button>
      </div>

      <div class="add-body">
        <div class="field">
          <label>文件夹绝对路径 / 压缩包路径</label>
          <input v-model="path" placeholder="例如：D:\Games\MyGame 或 D:\Downloads\game.7z" />
        </div>
        <div class="field">
          <label>压缩包密码（非压缩包可留空）</label>
          <input v-model="password" type="password" placeholder="请输入压缩包密码" />
        </div>
        <p v-if="errorMsg" class="error-msg">{{ errorMsg }}</p>
      </div>

      <div class="modal-footer">
        <button class="btn" @click="close">取消</button>
        <button class="btn btn-primary" :disabled="busy" @click="add">确定</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { api } from '../api'

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'added'): void
}>()

const path = ref('')
const password = ref('')
const busy = ref(false)
const errorMsg = ref('')

async function add() {
  const p = path.value.trim()
  if (!p) {
    errorMsg.value = '请输入路径'
    return
  }
  busy.value = true
  errorMsg.value = ''
  try {
    for (;;) {
      const result = await api.addGame(p, password.value)
      if (result.is_many) {
        if (!confirm('该文件夹包含多个子游戏，是否将其作为一个分类进行批量添加？')) {
          break
        }
        await api.confirmAddMany(p)
        alert('批量添加成功')
        emit('added')
        emit('close')
        return
      }
      if (result.need_pass) {
        errorMsg.value = '需要密码或密码错误，请重新输入'
        break
      }
      alert('添加游戏成功')
      emit('added')
      emit('close')
      return
    }
  } catch (err: unknown) {
    errorMsg.value = '添加失败：' + String(err)
  } finally {
    busy.value = false
  }
}

function close() {
  emit('close')
}
</script>

<style scoped>
.add-modal {
  width: 460px;
  padding: 23px 26px;
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
.field {
  margin-bottom: 14px;
}
.field label {
  display: block;
  font-size: 12px;
  color: var(--text-dim);
  margin-bottom: 5px;
}
.field input {
  width: 100%;
}
.error-msg {
  color: var(--danger);
  font-size: 13px;
  margin-bottom: 10px;
}
.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 16px;
  padding-top: 14px;
  border-top: 1px solid var(--border);
}
</style>
