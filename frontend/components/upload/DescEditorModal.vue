<!-- components/DescEditorModal.vue -->
<template>
  <DialogModel ref="dialogRef" width="600px" customClass="desc-dialog" @close="emit('close')">
    <div class="modal-header">
      <h3><i class="fas fa-edit"></i> 设计稿描述</h3>
      <button class="close-modal" @click="dialogRef.close()">&times;</button>
    </div>

    <div class="modal-body">
      <div class="desc-list">
        <div
          v-for="(img, idx) in images"
          :key="img.id"
          class="desc-item"
          :class="{ focused: focusIndex === idx }"
        >
          <div class="desc-item-preview">
            <img :src="img.preview" alt="预览" />
            <span class="desc-item-index">{{ idx + 1 }}</span>
          </div>
          <div class="desc-item-input">
            <label>{{ img.type || '设计稿' }}</label>
            <textarea
              :ref="el => setTextareaRef(el, idx)"
              v-model="tempDescriptions[idx]"
              @focus="changeFocusIndex(idx)"
              :placeholder="`描述第 ${idx + 1} 张设计稿的用途和关键元素...`"
              rows="3"
            ></textarea>
            <div class="desc-item-footer">
              <span class="char-count">{{ (tempDescriptions[idx] || '').length }}/500</span>
              <span class="desc-tip">帮助 AI 理解设计意图</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="modal-footer">
      <button class="btn-cancel" @click="dialogRef.close()">取消</button>
      <button class="btn-save" @click="save">保存描述</button>
    </div>
  </DialogModel>
</template>

<script setup>
import { ref, watch, nextTick } from 'vue'
import DialogModel from '~/components/dialog/DialogModel.vue'

const props = defineProps({
  images: { type: Array, default: () => [] },
  focusIndex: { type: Number, default: null }
})

const emit = defineEmits(['save', 'close', 'update:focusIndex'])

const dialogRef = ref(null)
const tempDescriptions = ref([])
const textareaRefs = ref({})

const open = () => { dialogRef.value.open() }
const close = () => { dialogRef.value.close() }

defineExpose({ open, close })

const setTextareaRef = (el, idx) => {
  if (el) textareaRefs.value[idx] = el
}

// 监听 dialog 打开状态，初始化数据
watch(() => dialogRef.value?.visible, (newVal) => {
  if (newVal) {
    tempDescriptions.value = props.images.map(img => img.description || '')
    nextTick(() => {
      if (props.focusIndex !== null && textareaRefs.value[props.focusIndex]) {
        textareaRefs.value[props.focusIndex].focus()
      }
    })
  }
})

// 监听 focusIndex 变化，切换聚焦
watch(() => props.focusIndex, async (newIdx) => {
  if (newIdx !== null && textareaRefs.value[newIdx]) {
    await nextTick()
    textareaRefs.value[newIdx].focus()
  }
})

// 保存描述
const save = () => {
  // 将临时描述保存到图片数组中
  const updatedDescriptions = [...tempDescriptions.value]
  emit('save', updatedDescriptions)
  close()
}

// 切换聚焦索引
const changeFocusIndex = (idx) => {
  emit('update:focusIndex', idx)
}
</script>

<style scoped>
.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  background: linear-gradient(135deg, rgba(0, 255, 255, 0.1), rgba(255, 0, 255, 0.1));
  border-bottom: 1px solid rgba(0, 255, 255, 0.2);
}

.modal-header h3 {
  font-size: 18px;
  color: #00ffff;
  display: flex;
  align-items: center;
  gap: 8px;
}

.close-modal {
  background: none;
  border: none;
  color: #888;
  font-size: 24px;
  cursor: pointer;
  transition: color 0.3s;
}

.close-modal:hover {
  color: #fff;
}

.modal-body {
  padding: 20px;
  max-height: 60vh;
  overflow-y: auto;
}

.desc-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.desc-item {
  display: flex;
  gap: 16px;
  padding: 12px;
  background: rgba(0, 0, 0, 0.3);
  border-radius: 16px;
  border: 1px solid rgba(0, 255, 255, 0.15);
  transition: all 0.3s;
}

.desc-item.focused {
  border-color: #00ffff;
  background: rgba(0, 255, 255, 0.05);
  box-shadow: 0 0 15px rgba(0, 255, 255, 0.1);
}

.desc-item-preview {
  position: relative;
  width: 70px;
  height: 70px;
  flex-shrink: 0;
}

.desc-item-preview img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 8px;
}

.desc-item-index {
  position: absolute;
  top: -6px;
  left: -6px;
  width: 20px;
  height: 20px;
  background: linear-gradient(135deg, #00ffff, #ff00ff);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 10px;
  font-weight: bold;
  color: #0a0a0f;
}

.desc-item-input {
  flex: 1;
}

.desc-item-input label {
  font-size: 12px;
  color: #00ffff;
  display: block;
  margin-bottom: 6px;
}

.desc-item-input textarea {
  width: 100%;
  background: rgba(0, 0, 0, 0.5);
  border: 1px solid rgba(0, 255, 255, 0.3);
  border-radius: 12px;
  padding: 10px;
  color: #fff;
  font-size: 13px;
  font-family: inherit;
  resize: vertical;
  outline: none;
  transition: all 0.3s;
}

.desc-item-input textarea:focus {
  border-color: #00ffff;
  box-shadow: 0 0 10px rgba(0, 255, 255, 0.1);
}

.desc-item-footer {
  display: flex;
  justify-content: space-between;
  margin-top: 6px;
  font-size: 10px;
  color: #666;
}

.char-count {
  font-family: monospace;
}

.desc-tip {
  color: #00ffff;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px 20px;
  border-top: 1px solid rgba(0, 255, 255, 0.2);
}

.btn-cancel, .btn-save {
  padding: 8px 24px;
  border-radius: 40px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s;
}

.btn-cancel {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.2);
  color: #aaa;
}

.btn-cancel:hover {
  background: rgba(255, 255, 255, 0.1);
}

.btn-save {
  background: linear-gradient(135deg, #00ffff, #ff00ff);
  border: none;
  color: white;
}

.btn-save:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 15px rgba(0, 255, 255, 0.3);
}
</style>