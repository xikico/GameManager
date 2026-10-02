<template>
  <div class="modal-mask" @click.self="close">
    <div class="modal settings-modal">
      <div class="modal-header">
        <div>
          <p class="eyebrow">PREFERENCES</p>
          <h2>设置</h2>
        </div>
        <button class="btn btn-small" @click="close">关闭</button>
      </div>

      <div class="settings-body">
        <label class="setting-row">
          <span class="setting-copy">
            <strong>剪贴板图片转 Base64</strong>
            <small>检测剪贴板中的图片，并自动替换为可粘贴的 Base64 文本。</small>
          </span>
          <input v-model="settings.ClipboardImageDetectionEnabled" class="switch-input" type="checkbox" />
          <span class="switch" aria-hidden="true"></span>
        </label>
        <p v-if="errorMsg" class="error-msg">{{ errorMsg }}</p>
      </div>

      <div class="modal-footer">
        <span class="save-note">设置保存在本机数据库中</span>
        <button class="btn btn-primary" :disabled="saving || loading" @click="save">
          {{ saving ? '保存中...' : '保存设置' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { api } from '../api'

const emit = defineEmits<{ (e: 'close'): void }>()

const loading = ref(true)
const saving = ref(false)
const errorMsg = ref('')
const settings = reactive({ ClipboardImageDetectionEnabled: true })

onMounted(async () => {
  try {
    Object.assign(settings, await api.getSettings())
  } catch (err: unknown) {
    errorMsg.value = '读取设置失败：' + String(err)
  } finally {
    loading.value = false
  }
})

async function save() {
  saving.value = true
  errorMsg.value = ''
  try {
    await api.updateSettings({ ...settings })
    close()
  } catch (err: unknown) {
    errorMsg.value = '保存设置失败：' + String(err)
  } finally {
    saving.value = false
  }
}

function close() {
  emit('close')
}
</script>

<style scoped>
.settings-modal {
  width: 520px;
  max-width: 92vw;
  padding: 24px 26px 22px;
}
.modal-header,
.modal-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.modal-header {
  margin-bottom: 18px;
}
.modal-header h2 {
  font-family: 'Segoe UI Variable Display', 'Microsoft YaHei UI', sans-serif;
  font-size: 23px;
  font-weight: 720;
  letter-spacing: -0.4px;
}
.eyebrow {
  color: #6ec2f2;
  font-size: 10px;
  font-weight: 800;
  letter-spacing: 1.8px;
  margin-bottom: 3px;
}
.settings-body {
  padding: 5px 0 18px;
}
.setting-row {
  display: grid;
  grid-template-columns: 1fr auto;
  align-items: center;
  gap: 22px;
  padding: 19px;
  border: 1px solid var(--border-strong);
  border-radius: 14px;
  background: linear-gradient(135deg, rgba(19,35,55,.86), rgba(9,22,37,.82));
  box-shadow: inset 0 1px rgba(255,255,255,.025),0 10px 26px rgba(0,5,14,.12);
  cursor: pointer;
}
.setting-copy {
  display: flex;
  flex-direction: column;
  gap: 6px;
}
.setting-copy strong {
  font-size: 14px;
  letter-spacing: .1px;
}
.setting-copy small,
.save-note {
  color: var(--text-dim);
  font-size: 12px;
  line-height: 1.6;
}
.switch-input {
  position: absolute;
  opacity: 0;
  pointer-events: none;
}
.switch {
  width: 46px;
  height: 26px;
  border-radius: 99px;
  background: rgba(81,103,132,.5);
  border: 1px solid rgba(150,181,216,.13);
  position: relative;
  transition: 0.2s ease;
}
.switch::after {
  content: '';
  position: absolute;
  width: 20px;
  height: 20px;
  top: 3px;
  left: 3px;
  border-radius: 50%;
  background: linear-gradient(145deg,#fff,#dbe9f4);
  box-shadow: 0 3px 8px rgba(0,0,0,.38),inset 0 1px rgba(255,255,255,.8);
  transition: 0.2s ease;
}
.switch-input:checked + .switch {
  background: linear-gradient(135deg,#50ade8,#277eb5);
  box-shadow: 0 0 0 3px rgba(76,168,232,.1),0 6px 16px rgba(26,112,168,.28);
}
.switch-input:checked + .switch::after {
  transform: translateX(20px);
}
.error-msg {
  color: var(--danger);
  font-size: 13px;
  margin-top: 12px;
}
.modal-footer {
  border-top: 1px solid var(--border);
  padding-top: 16px;
}
</style>
