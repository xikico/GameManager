<template>
  <div class="modal-mask" @click.self="cancel">
    <div class="modal del-modal">
      <div class="modal-header">
        <h2>确认删除</h2>
      </div>

      <div class="del-body">
        <p class="del-tip">确定要删除「{{ gameName }}」吗？</p>
        <label class="del-option">
          <input v-model="delAll" type="checkbox" />
          <span>同时删除源文件（文件夹）</span>
        </label>
      </div>

      <div class="modal-footer">
        <button class="btn" @click="cancel">取消</button>
        <button class="btn btn-danger" :disabled="busy" @click="confirm">
          {{ delAll ? '删除（含源文件）' : '确定' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

defineProps<{ gameName: string }>()

const emit = defineEmits<{
  (e: 'confirm', delAll: boolean): void
  (e: 'cancel'): void
}>()

const delAll = ref(false)
const busy = ref(false)

async function confirm() {
  busy.value = true
  emit('confirm', delAll.value)
}

function cancel() {
  emit('cancel')
}
</script>

<style scoped>
.del-modal {
  width: 400px;
  padding: 20px 24px;
}
.modal-header {
  margin-bottom: 14px;
}
.modal-header h2 {
  font-size: 17px;
  color: var(--danger);
}
.del-tip {
  font-size: 14px;
  margin-bottom: 16px;
}
.del-option {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  cursor: pointer;
  color: var(--text-dim);
}
.del-option input {
  accent-color: var(--danger);
  width: 15px;
  height: 15px;
  cursor: pointer;
}
.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 20px;
  padding-top: 14px;
  border-top: 1px solid var(--border);
}
</style>
