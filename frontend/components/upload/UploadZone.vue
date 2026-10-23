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

  <!-- 预览卡片列表（竖向卡片网格） -->
  <div v-if="images.length" class="preview-area" :class="{ 'preview-area--locked': isSubmitting }">
    <div class="preview-grid">
      <div
        v-for="(img, i) in images"
        :key="i"
        class="preview-card"
        :class="{
          'preview-card--dragging': dragIdx === i,
          'preview-card--drag-over': dragOverIdx === i && dragIdx !== i,
          'preview-card--disabled': isNotReady(img),
        }"
        draggable="true"
        @click="openAnno(i)"
        @dragstart="onDragStart($event, i)"
        @dragover.prevent="onDragOver(i)"
        @dragleave="dragOverIdx = -1"
        @drop.prevent="onDrop(i)"
        @dragend="onDragEnd"
      >
        <!-- ── 图区 ── -->
        <div class="card-media">
          <!-- img-frame 收缩包裹图片渲染区域，标注框按 0~1 相对坐标直接叠放 -->
          <div class="img-frame">
            <img :src="img.preview" alt="设计稿" draggable="false" />
            <span
              v-for="(anno, j) in img.annotations || []"
              :key="anno.id"
              class="anno-box"
              :style="{
                left: anno.x * 100 + '%',
                top: anno.y * 100 + '%',
                width: anno.w * 100 + '%',
                height: anno.h * 100 + '%',
              }"
            >
              <span class="anno-num">{{ j + 1 }}</span>
            </span>
          </div>

          <!-- 上传状态 -->
          <div v-if="img.uploading" class="upload-status uploading">
            <i class="fas fa-spinner fa-spin"></i>
            <span class="uploading-text">上传中…</span>
          </div>
          <div v-else-if="img.uploadError" class="upload-status failed">
            <i class="fas fa-circle-exclamation"></i>
            <span class="failed-text">上传失败</span>
            <button class="retry-btn" @click.stop="onRetry(img)">
              <i class="fas fa-rotate-right"></i> 重试
            </button>
          </div>

          <!-- 悬浮工具条 -->
          <div class="card-toolbar">
            <span class="card-index">{{ String(i + 1).padStart(2, '0') }}</span>
            <span class="toolbar-spacer"></span>
            <button
              class="tool-icon"
              :disabled="isNotReady(img)"
              :title="isNotReady(img) ? '图片未上传完成，暂不可标注' : '标注'"
              @click.stop="openAnno(i)"
            >
              <i class="fa-regular fa-pen-to-square"></i>
            </button>
            <button class="tool-icon tool-icon--danger" title="删除" @click.stop="store.removeImage(i)">
              <i class="fa-regular fa-trash-can"></i>
            </button>
          </div>
        </div>

        <!-- ── 信息区 ── -->
        <div class="card-info">
          <div class="info-line">
            <span class="file-name">{{ img.file?.name || '设计稿' }}</span>
            <span class="file-size">{{ img.naturalWidth || '--' }} x {{ img.naturalHeight || '--' }}</span>
          </div>
          <!-- 上传状态优先展示，未就绪不提供标注入口 -->
          <div v-if="img.uploading" class="anno-summary empty">
            <i class="fas fa-spinner fa-spin"></i>
            上传中，完成后可标注
          </div>
          <div v-else-if="img.uploadError" class="anno-summary error">
            <i class="fas fa-exclamation-circle"></i>
            上传失败，暂不可标注
          </div>
          <div v-else-if="img.annotations?.length" class="anno-summary">
            <i class="fa-solid fa-location-target"></i>
            {{ img.annotations.length }} 个标注 · {{ annoLabels(img) }}
          </div>
          <div v-else class="anno-summary empty">
            <i class="fa-regular fa-square-dashed"></i>
            未标注，点击进入标注编辑器
          </div>
          <div v-if="img.description" class="page-desc-text">{{ img.description }}</div>
        </div>
      </div>
    </div>

    <div class="upload-actions">
      <span class="upload-count">{{ images.length }} 张已选择</span>
      <button class="clear-btn" @click="store.clearAll()">清空全部</button>
    </div>

    <div class="upload-tip">
      <i class="fa-regular fa-lightbulb"></i>
      <span class="tip-bold">提示：</span>拖拽卡片可调整页面顺序，框选区域并填写描述能帮 AI 更准地理解设计稿
    </div>
  </div>

  <!-- 标注编辑器 -->
  <AnnotationModal
    ref="annoModalRef"
    :images="images"
    :focusIndex="editingIndex"
    @save="onAnnoSave"
    @close="onAnnoClose"
  />
</template>

<script setup>
import { ref, computed } from 'vue'
import { storeToRefs } from 'pinia'
import { ElMessage } from 'element-plus'
import AnnotationModal from '~/components/upload/AnnotationModal.vue'
import { useCodeStore } from '~/stores/code'
import { useUserStore } from '~/stores/user'
import { useCommonStore } from '~/stores/common'

const userStore = useUserStore()
const commonStore = useCommonStore()
const { isLogin } = storeToRefs(userStore)
const openLoginModal = () => commonStore.setLoginModalVisible(true)

const store = useCodeStore()
const { images, isSubmitting, isMaxReached, isConcurrencyFull, uploadHint } = storeToRefs(store)

// ── 从 store 派生 ──
const disabled = computed(() => {
  return isMaxReached.value || isSubmitting.value || isConcurrencyFull.value
})

// ── 文件选择 ──
const fileInput = ref(null)

const onZoneClick = () => {
  if (!isLogin.value) return openLoginModal()
  if (disabled.value) return
  fileInput.value?.click()
}

const onFileChange = (e) => {
  const files = Array.from(e.target.files)
  if (files.length) store.addImages(files)
  if (fileInput.value) fileInput.value.value = ''
}

const onZoneDrop = (e) => {
  if (!isLogin.value) return openLoginModal()
  if (disabled.value) return
  const files = Array.from(e.dataTransfer.files).filter(f => f.type.startsWith('image/'))
  if (files.length) store.addImages(files)
}

// ── 标注编辑 ──
const annoModalRef = ref(null)
const editingIndex = ref(0)

// 图片未就绪（上传中 / 失败 / 无 COS 地址）时禁止标注
const isNotReady = (img) => !!img.uploading || !!img.uploadError || !img.cosUrl

const openAnno = (idx) => {
  if (disabled.value) return
  const img = images.value[idx]
  if (!img) return
  if (isNotReady(img)) {
    if (img.uploadError) ElMessage.warning('图片上传失败，请先点击重试')
    return
  }
  editingIndex.value = idx
  annoModalRef.value?.open()
}

// 上传失败后重试
const onRetry = (img) => {
  if (!img.uploadId) return
  store.retryUpload(img.uploadId)
}

// 编辑器保存回写（drafts 与 images 索引一一对应）
const onAnnoSave = (drafts) => {
  drafts.forEach((draft, idx) => {
    const img = images.value[idx]
    if (!img) return
    img.description = draft.description || ''
    img.annotations = draft.annotations || []
  })
}

const onAnnoClose = () => {
  editingIndex.value = -1
}

// 卡片上展示标注的标签摘要
const annoLabels = (img) => {
  const labels = (img.annotations || []).map(a => a.label || '未命名').filter(Boolean)
  return labels.join('、')
}

// ── 拖拽排序 ──
const dragIdx = ref(-1)
const dragOverIdx = ref(-1)

const onDragStart = (e, idx) => {
  e.dataTransfer.effectAllowed = 'move'
  e.dataTransfer.setData('text/plain', String(idx))
  dragIdx.value = idx
}
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

/* ── 预览卡片区 ── */
.preview-area {
  margin-top: 16px;
}
.preview-area--locked {
  pointer-events: none;
  opacity: 0.45;
}

.preview-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
  gap: 12px;
}

.preview-card {
  position: relative;
  border-radius: 12px;
  background: rgba(0, 0, 0, 0.3);
  border: 1px solid rgba(0, 255, 255, 0.08);
  overflow: visible;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}
.preview-card:hover { border-color: rgba(0, 255, 255, 0.25); }
.preview-card--disabled { cursor: not-allowed; }
.preview-card--dragging {
  opacity: 0;
  transform: scale(0.8);
}
.preview-card--drag-over {
  border-color: transparent;
  transform: none;
}
.preview-card--drag-over::after {
  content: '';
  position: absolute;
  inset: -3px;
  border: 2px dashed rgba(0, 255, 255, 0.4);
  border-radius: 14px;
  pointer-events: none;
  animation: dropPulse 1s ease-in-out infinite;
}

@keyframes dropPulse {
  0%, 100% { border-color: rgba(0, 255, 255, 0.3); box-shadow: none; }
  50% { border-color: rgba(0, 255, 255, 0.7); box-shadow: 0 0 12px rgba(0, 255, 255, 0.15); }
}

/* ── 图区 ── */
.card-media {
  position: relative;
  height: 220px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.4);
  border-radius: 12px 12px 0 0;
  overflow: hidden;
}
/* 收缩包裹 img：框层与图片渲染区域严格一致 */
.img-frame {
  position: relative;
  display: inline-block;
  max-width: 100%;
  max-height: 100%;
}
.img-frame img {
  display: block;
  max-width: 100%;
  max-height: 100%;
  user-select: none;
}

/* 标注框（缩略展示，不可交互） */
.anno-box {
  position: absolute;
  border: 1.5px solid rgba(0, 255, 255, 0.9);
  background: rgba(0, 255, 255, 0.08);
  pointer-events: none;
}
.anno-num {
  position: absolute;
  top: -8px;
  left: -2px;
  min-width: 14px;
  height: 14px;
  border-radius: 7px;
  background: #00ffff;
  color: #0a0a0f;
  font-size: 9px;
  font-weight: 700;
  line-height: 14px;
  text-align: center;
}

/* 上传状态 */
.upload-status {
  position: absolute;
  z-index: 3;
  display: flex;
  align-items: center;
  justify-content: center;
}
/* 上传中：整图遮罩 + 居中 loading */
.upload-status.uploading {
  inset: 0;
  flex-direction: column;
  gap: 8px;
  background: rgba(0, 0, 0, 0.45);
  color: #00ffff;
  font-size: 26px;
}
.upload-status.uploading .uploading-text {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.7);
  letter-spacing: 1px;
}
/* 失败：整图遮罩 + 重试入口 */
.upload-status.failed {
  inset: 0;
  flex-direction: column;
  gap: 8px;
  background: rgba(0, 0, 0, 0.55);
  color: #ff4757;
  font-size: 22px;
}
.upload-status.failed .failed-text {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.75);
  letter-spacing: 1px;
}
.retry-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-top: 2px;
  padding: 5px 14px;
  border-radius: 16px;
  border: 1px solid rgba(255, 71, 87, 0.5);
  background: rgba(255, 71, 87, 0.12);
  color: #ff6b7a;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s;
}
.retry-btn:hover {
  background: rgba(255, 71, 87, 0.25);
  color: #fff;
}

/* 悬浮工具条 */
.card-toolbar {
  position: absolute;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 6px 8px;
  background: linear-gradient(transparent, rgba(0, 0, 0, 0.65));
  opacity: 0;
  transition: opacity 0.2s;
  z-index: 3;
}
.preview-card:hover .card-toolbar { opacity: 1; }
.card-index {
  font-size: 13px;
  font-weight: 700;
  color: rgba(255, 255, 255, 0.85);
}
.toolbar-spacer { flex: 1; }
.tool-icon {
  background: rgba(0, 0, 0, 0.4);
  border: 1px solid rgba(0, 255, 255, 0.2);
  color: rgba(0, 255, 255, 0.7);
  font-size: 12px;
  width: 26px;
  height: 26px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
}
.tool-icon:hover {
  color: #00ffff;
  background: rgba(0, 255, 255, 0.12);
  border-color: rgba(0, 255, 255, 0.5);
}
.tool-icon--danger { color: rgba(255, 255, 255, 0.5); border-color: rgba(255, 255, 255, 0.15); }
.tool-icon:disabled {
  opacity: 0.35;
  cursor: not-allowed;
  color: rgba(255, 255, 255, 0.4);
  border-color: rgba(255, 255, 255, 0.1);
  background: rgba(0, 0, 0, 0.3);
}
.tool-icon--danger:hover {
  color: #ff4757;
  background: rgba(255, 71, 87, 0.12);
  border-color: rgba(255, 71, 87, 0.4);
}

/* ── 信息区 ── */
.card-info {
  padding: 10px 12px 12px;
  display: flex;
  flex-direction: column;
  gap: 6px;
}
.info-line {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  min-width: 0;
}
.file-name {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.75);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.file-size {
  flex-shrink: 0;
  font-size: 11px;
  color: rgba(255, 255, 255, 0.35);
}

.anno-summary {
  font-size: 12px;
  color: rgba(0, 255, 255, 0.55);
  display: flex;
  align-items: center;
  gap: 6px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.anno-summary i { font-size: 11px; flex-shrink: 0; }
.anno-summary.empty { color: rgba(255, 255, 255, 0.25); }
.anno-summary.error { color: rgba(255, 71, 87, 0.75); }

.page-desc-text {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.4);
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

/* ── 底部操作与提示 ── */
.upload-actions {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 0 4px;
}
.upload-count {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.5);
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
.clear-btn:hover { background: rgba(255, 0, 0, 0.1); }
.clear-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
  pointer-events: none;
}

.upload-tip {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.35);
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
  color: rgba(255, 255, 255, 0.5);
  font-weight: 600;
}
</style>
