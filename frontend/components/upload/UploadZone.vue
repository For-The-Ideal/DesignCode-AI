<template>
  <!-- 上传区 -->
  <div
    class="upload-zone"
    :class="{ 'upload-zone--disabled': disabled }"
    @click="onZoneClick"
    @dragover.prevent
    @drop.prevent="onZoneDrop"
  >
    <i class="fa-solid fa-cloud-arrow-up upload-icon"></i>
    <div class="upload-text">
      <span class="upload-text-bold">{{ uploadHint.title }}</span>
    </div>
    <div class="upload-hint">
      {{ uploadHint.desc }}
    </div>
  </div>

  <!-- 隐藏文件选择 -->
  <input type="file" ref="fileInput" multiple accept="image/*" style="display:none" @change="onFileChange" />

  <!-- 预览列表 -->
  <div v-if="images.length" class="preview-area" :class="{ 'preview-area--locked': isSubmitting }">
    <div class="preview-list">
      <div
        v-for="(img, i) in images"
        :key="i"
        class="preview-item"
        :class="{
          'preview-item--dragging': dragIdx === i,
          'preview-item--drag-over': dragOverIdx === i && dragIdx !== i,
        }"
        draggable="true"
        @dragstart="onDragStart($event, i)"
        @dragover.prevent="onDragOver(i)"
        @dragleave="dragOverIdx = -1"
        @drop.prevent="onDrop(i)"
        @dragend="onDragEnd"
      >
        <!-- 序号 -->
        <span class="item-index">{{ String(i + 1).padStart(2, '0') }}</span>
        <!-- 文件信息 -->
        <div class="item-file">
          <div class="file-thumb">
            <img :src="img.preview" alt="preview" />
            <div v-if="img.uploading" class="upload-status uploading">
              <i class="fas fa-spinner fa-spin"></i>
            </div>
            <div v-else-if="img.uploadError" class="upload-status failed">
              <i class="fas fa-exclamation-circle"></i>
            </div>
            <div v-else-if="img.cosUrl" class="upload-status success">
              <i class="fas fa-check-circle"></i>
            </div>
          </div>
          <div class="file-meta">
            <span class="file-name">{{ img.file?.name || '设计稿' }}</span>
            <span class="file-size">{{ img.naturalWidth || '--' }} x {{ img.naturalHeight || '--' }}</span>
          </div>
        </div>
        <!-- 描述输入 -->
        <div class="item-desc">
          <i class="fa-regular fa-file-lines desc-icon"></i>
          <span v-if="img.description" class="desc-text">{{ img.description }}</span>
          <span v-else class="desc-placeholder">添加图片描述</span>
          <button class="desc-edit-btn" @click.stop="onEditDesc(i)">
            <i class="fa-regular fa-pen-to-square"></i>
          </button>
        </div>
        <!-- 操作按钮 -->
        <div class="item-actions">
          <button class="action-btn" @click.stop="store.removeImage(i)">
            <i class="fa-regular fa-trash-can"></i>
          </button>
        </div>
      </div>
    </div>

    <div class="upload-actions">
      <span class="upload-count">{{ images.length }} 张已选择</span>
      <button class="clear-btn" @click="store.clearAll()">清空全部</button>
    </div>

    <div class="upload-tip">
      <i class="fa-regular fa-lightbulb"></i>
      <span class="tip-bold">提示：</span>拖拽可调整页面顺序，AI会按顺序理解页面关系
    </div>
  </div>

  <!-- 描述编辑弹窗 -->
  <DescEditorModal
    ref="descEditorRef"
    :images="images"
    :focusIndex="editingIndex"
    @save="onDescSave"
    @close="onDescClose"
  />
</template>

<script setup>
import { ref, computed } from 'vue'
import { storeToRefs } from 'pinia'
import DescEditorModal from '~/components/upload/DescEditorModal.vue'
import { useCodeStore } from '~/stores/code'

const store = useCodeStore()
const { images, isSubmitting, isMaxReached, isConcurrencyFull, uploadHint } = storeToRefs(store)

// ── 从 store 派生 ──
const disabled = computed(() => isMaxReached.value || isSubmitting.value || isConcurrencyFull.value)

// ── 文件选择 ──
const fileInput = ref(null)

const onZoneClick = () => {
  if (disabled.value) return
  fileInput.value?.click()
}

const onFileChange = (e) => {
  const files = Array.from(e.target.files)
  if (files.length) store.addImages(files)
  if (fileInput.value) fileInput.value.value = ''
}

const onZoneDrop = (e) => {
  if (disabled.value) return
  const files = Array.from(e.dataTransfer.files).filter(f => f.type.startsWith('image/'))
  if (files.length) store.addImages(files)
}

// ── 描述编辑 ──
const descEditorRef = ref(null)
const editingIndex = ref(-1)

const onEditDesc = (idx) => {
  editingIndex.value = idx
  descEditorRef.value?.open()
}

const onDescSave = (descriptions) => {
  images.value.forEach((img, idx) => {
    img.description = descriptions[idx] || ''
  })
  descEditorRef.value?.close()
}

const onDescClose = () => {
  editingIndex.value = -1
}

// ── 拖拽排序 ──
const dragIdx = ref(-1)
const dragOverIdx = ref(-1)

const onDragStart = (e, idx) => { dragIdx.value = idx }
const onDragOver = (idx) => { dragOverIdx.value = idx }
const onDrop = (idx) => {
  if (dragIdx.value < 0 || dragIdx.value === idx) return
  const arr = [...images.value]
  const [item] = arr.splice(dragIdx.value, 1)
  arr.splice(idx, 0, item)
  images.value = arr
  dragIdx.value = -1
  dragOverIdx.value = -1
}
const onDragEnd = () => {
  dragIdx.value = -1
  dragOverIdx.value = -1
}
</script>

<style scoped>
/* ── 上传区 ── */
.upload-zone {
  border: 2px dashed rgba(0, 255, 255, 0.1);
  border-radius: 16px;
  padding: 40px 16px;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s;
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}
.upload-zone:hover {
  border-color: rgba(0, 255, 255, 0.3);
  background: rgba(0, 255, 255, 0.02);
}

.upload-zone--disabled {
  opacity: 0.4;
  cursor: not-allowed;
  border-color: rgba(255, 255, 255, 0.08) !important;
}
.upload-zone--disabled:hover {
  background: transparent;
  border-color: rgba(255, 255, 255, 0.08) !important;
}

.upload-icon {
  font-size: 40px;
  color: rgba(0, 255, 255, 0.25);
  margin-bottom: 12px;
}

.upload-text {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.5);
  margin-bottom: 6px;
}
.upload-text-bold {
  color: rgba(255, 255, 255, 0.7);
  font-weight: 500;
}

.upload-hint {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.35);
}

/* ── 预览列表 ── */
.preview-area {
  margin-top: 16px;
}
.preview-area--locked {
  pointer-events: none;
  opacity: 0.45;
}
.preview-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.preview-item {
  position: relative;
  display: flex;
  align-items: center;
  border-radius: 10px;
  justify-content: space-between;
  overflow: visible;
  background: rgba(0, 0, 0, 0.3);
  border: 1px solid rgba(0, 255, 255, 0.08);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  cursor: grab;
  padding: 12px 15px;
}
.preview-item:hover { border-color: rgba(0, 255, 255, 0.15); }
.preview-item--dragging {
  opacity: 0;
  transform: scale(0.8);
}
.preview-item--drag-over {
  margin-top: 18px;
  border-color: transparent;
  box-shadow: none;
  transform: none;
}
.preview-item--drag-over::before {
  content: '';
  position: absolute;
  left: 0;
  right: 0;
  top: -14px;
  height: 4px;
  background: rgba(0, 255, 255, 0.3);
  border-radius: 2px;
  box-shadow: 0 0 8px rgba(0, 255, 255, 0.2);
}
.preview-item--drag-over::after {
  content: '';
  position: absolute;
  inset: -3px;
  border: 2px dashed rgba(0, 255, 255, 0.4);
  border-radius: 12px;
  pointer-events: none;
  animation: dropPulse 1s ease-in-out infinite;
}

@keyframes dropPulse {
  0%, 100% { border-color: rgba(0, 255, 255, 0.3); box-shadow: none; }
  50% { border-color: rgba(0, 255, 255, 0.7); box-shadow: 0 0 12px rgba(0, 255, 255, 0.15); }
}
/* 序号区 */
.item-index {
  flex: 0 0 40px;
  font-size: 15px;
  font-weight: 700;
  color: #ffffff;
  text-align: center;
  line-height: 1;
}

/* 文件信息区 */
.item-file {
  flex: 0 0 220px;
  display: flex;
  align-items: center;
  gap: 12px;
}

.file-thumb {
  position: relative;
  flex-shrink: 0;
  width: 90px;
  height: 60px;
  border-radius: 6px;
  overflow: hidden;
  background: rgba(0, 0, 0, 0.4);
}
.file-thumb img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}
.file-thumb .upload-status {
  position: absolute;
  font-size: 8px;
  padding: 1px 4px;
  top: 1px;
  right: 1px;
  border-radius: 4px;
  display: flex;
  align-items: center;
  gap: 2px;
  z-index: 2;
}
.file-thumb .upload-status.uploading { background: rgba(0, 255, 255, 0.2); color: #00ffff; }
.file-thumb .upload-status.failed { background: rgba(255, 0, 0, 0.2); color: #ff4757; }
.file-thumb .upload-status.success { background: rgba(0, 200, 83, 0.2); color: #00c853; }
.file-thumb .upload-status i { font-size: 7px; }

.file-meta {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
  min-width: 0;
}
.file-meta .file-name {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.75);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 130px;
}
.file-meta .file-size {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.35);
}

/* 描述区 */
.item-desc {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
  padding: 6px 12px;
  background: rgba(0, 0, 0, 0.35);
  border-radius: 8px;
  max-width: 450px;
  border: 1px solid transparent;
  cursor: pointer;
  transition: all 0.2s;
}
.item-desc:hover {
  background: rgba(0, 0, 0, 0.45);
  border-color: rgba(0, 255, 255, 0.15);
}

.desc-icon {
  font-size: 12px;
  color: rgba(0, 255, 255, 0.35);
  flex-shrink: 0;
}

.desc-text {
  flex: 1;
  font-size: 12px;
  color: rgba(255, 255, 255, 0.55);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.desc-placeholder {
  flex: 1;
  font-size: 12px;
  color: rgba(255, 255, 255, 0.2);
  font-style: italic;
}

.desc-edit-btn {
  flex-shrink: 0;
  background: none;
  border: none;
  color: rgba(0, 255, 255, 0.3);
  font-size: 12px;
  cursor: pointer;
  padding: 2px 5px;
  border-radius: 4px;
  transition: all 0.2s;
  opacity: 0;
}
.item-desc:hover .desc-edit-btn { opacity: 1; }
.desc-edit-btn:hover {
  color: #00ffff;
  background: rgba(0, 255, 255, 0.1);
}

/* 操作区 */
.item-actions {
  flex: 0 0 40px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.action-btn {
  background: none;
  border: none;
  color: rgba(255, 255, 255, 0.25);
  font-size: 16px;
  cursor: pointer;
  padding: 6px;
  border-radius: 8px;
  transition: all 0.2s;
  line-height: 1;
}
.action-btn:hover {
  color: #ff4757;
  background: rgba(255, 71, 87, 0.1);
}
.action-btn:disabled {
  opacity: 0.3;
  cursor: not-allowed;
}

.upload-actions {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 4px 0;
}
.upload-count {
  font-size: 12px;
  color: rgba(255,255,255,0.5);
}
.clear-btn {
  background: none;
  border: none;
  color: #ff6b6b;
  font-size: 12px;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 6px;
}
.clear-btn:hover { background: rgba(255,0,0,0.1); }
.clear-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
  pointer-events: none;
}

.upload-tip {
  font-size: 14px;
  color: rgba(255,255,255,0.35);
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 6px;
}
.upload-tip i {
  color: rgba(0, 255, 255, 0.3);
  font-size: 13px;
}
.tip-bold {
  color: rgba(255,255,255,0.5);
  font-weight: 600;
}
</style>
