<template>
  <!-- 左栏：模式切换 + 页面列表（多图 = 页序列表可切页；单图 = 唯一主图卡） -->
  <aside class="wb-upload">
    <div class="mode-switch">
      <button class="mode-btn" :class="{ active: mode === 'single' }" :disabled="isSubmitting" @click="store.setMode('single')">
        <i class="fa-regular fa-image"></i> 单图
      </button>
      <button class="mode-btn" :class="{ active: mode === 'multi' }" :disabled="isSubmitting" @click="store.setMode('multi')">
        <i class="fa-solid fa-layer-group"></i> 多图
      </button>
    </div>

    <!-- 多图：入口页没有上层，页面关系从第 2 页开始逐页指定 -->
    <p v-if="mode === 'multi' && images.length" class="rel-hint">
      <i class="fas fa-circle-info"></i>
      <span><b>P1</b> 是入口页，没有上层页面。页面关系或页面组件从第 2 页开始配置：每一页指定它由哪一页进入。</span>
    </p>

    <!-- ═══ 多图：页面列表（点卡片切页，画布与区域树随之切换） ═══ -->
    <section v-if="mode === 'multi'" class="panel panel--pages">
      <header class="panel-head">
        <span class="panel-title"><i class="fa-regular fa-file-image"></i> 页面列表</span>
        <span class="panel-count">{{ images.length }}/{{ maxImages }}</span>
      </header>

      <div class="page-list">
        <button
          v-for="(img, i) in images"
          :key="img.pageId"
          class="page-item"
          :class="{
            active: img.pageId === currentImage?.pageId,
            'is-dragging': dragFrom === i,
            'is-drop-target': dragFrom > -1 && dropAt === i && dragFrom !== i,
            'is-unplaced': !img.board,
          }"
          :draggable="!isSubmitting"
          @click="onSelectPage(img)"
          @dblclick="onAnnotatePage(img)"
          @dragstart="onDragStart($event, i, img)"
          @dragover.prevent="onDragOver(i)"
          @drop.prevent="onDropItem($event, i)"
          @dragend="resetDrag"
        >
          <div class="page-item-media">
            <img v-if="img.preview || img.cosUrl" :src="img.preview || img.cosUrl" alt="" draggable="false" />
            <i v-else class="fas fa-image"></i>
            <span class="page-no">P{{ i + 1 }}</span>
            <span v-if="img.uploading" class="page-mask"><i class="fas fa-spinner fa-spin"></i></span>
            <span v-else-if="img.uploadError" class="page-mask page-mask--err" @click.stop="onRetryPage(img)">
              <i class="fas fa-rotate"></i> 重试
            </span>
            <!-- 未入画布的页：移入缩略图才浮出拖拽提示，不占常驻文案 -->
            <span v-else-if="!img.board" class="page-drag-hint">
              <i class="fas fa-up-down-left-right"></i> 拖入画布
            </span>
          </div>
          <!-- 缩略图下：入口标识 + 文件名靠左，文件大小靠右 -->
          <div class="page-item-info">
            <span class="page-item-title">
              <!-- 入口页恒为列表第 1 张：拖拽排序后自动跟着换，不单独维护标记 -->
              <span v-if="i === 0" class="page-entry">入口</span>
              <span class="page-item-name" :title="pageName(img)">{{ pageName(img) }}</span>
              <span v-if="fileSize(img)" class="page-item-size">{{ fileSize(img) }}</span>
            </span>
          </div>
        </button>

        <div v-if="!images.length" class="page-empty">
          <i class="fa-regular fa-images"></i>
          <p>暂无页面</p>
          <span>先在画布中上传第一张设计稿，再从这里添加页面</span>
        </div>
      </div>

      <!-- 整块随 v-if 隐藏，已保证有图；仅需处理提交中与页数达上限 -->
      <div class="file-actions" v-if="images.length">
        <button class="act-btn" :disabled="isSubmitting || images.length >= maxImages" @click="pickMore">
          <i class="fas fa-plus"></i> 添加页面
        </button>
        <button class="act-btn act-btn--danger" :disabled="isSubmitting" @click="onDeletePage">
          <i class="fa-regular fa-trash-can"></i> 删除当前页
        </button>
      </div>
    </section>

    <!-- ═══ 单图：唯一主图卡（保持原样） ═══ -->
    <section v-else class="panel panel--pages">
      <header class="panel-head">
        <span class="panel-title" :title="currentImage?.file?.name || '设计稿'">
          <i class="fa-regular fa-file-image"></i>
          {{ currentImage?.file?.name || '设计稿' }}
        </span>
      </header>

      <div v-if="!currentImage" class="page-empty">
        <i class="fa-regular fa-images"></i>
        <p>暂无上传页面</p>
        <span>上传设计稿后在此管理单页面</span>
      </div>

      <div v-else class="page-card">
        <div class="page-media">
          <img v-if="currentImage.preview || currentImage.cosUrl" :src="currentImage.preview || currentImage.cosUrl" alt="设计稿" draggable="false" />
          <i v-else class="fas fa-image"></i>
          <span v-if="currentImage.uploading" class="thumb-mask"><i class="fas fa-spinner fa-spin"></i></span>
          <span v-else-if="currentImage.uploadError" class="thumb-mask thumb-mask--err"><i class="fas fa-circle-exclamation"></i></span>
        </div>
        <div class="page-info">
          <div class="page-meta">
            <span class="page-dim">
              {{ currentImage.naturalWidth || currentImage.treeWidth || '--' }} × {{ currentImage.naturalHeight || currentImage.treeHeight || '--' }}
            </span>
            <span class="state-pill" :class="stateOf(currentImage)">{{ stateText[stateOf(currentImage)] }}</span>
          </div>
        </div>
      </div>

      <div v-if="currentImage" class="file-actions">
        <button class="act-btn" :disabled="isSubmitting" @click="onReplace"><i class="fas fa-rotate"></i> 更换</button>
        <button class="act-btn act-btn--danger" :disabled="isSubmitting" @click="onDeleteSingle">
          <i class="fa-regular fa-trash-can"></i> 删除
        </button>
      </div>
    </section>

    <input ref="fileInput" type="file" accept="image/*" :multiple="mode === 'multi'" hidden @change="onFileChange" />
  </aside>
</template>

<script setup>
import { ref } from 'vue'
import { storeToRefs } from 'pinia'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useCodeStore } from '~/stores/code'
import { useUiProjectStore } from '~/stores/uiProject'
import { checkAnnotations } from '~/utils/annotationCheck'

const store = useCodeStore()
const uiStore = useUiProjectStore()
const { images, mode, isSubmitting, maxImages, currentImage } = storeToRefs(store)

const stateOf = (img) => (img.uploading ? 'uploading' : img.uploadError ? 'error' : 'ready')
const stateText = { uploading: '上传中', error: '上传失败', ready: '已就绪' }

// ═══ 多图页面列表 ═══
const pageName = (img) => {
  const i = images.value.indexOf(img)
  return img.pageName || (img.file?.name || '').replace(/\.[^.]+$/, '') || `P${i + 1}`
}

// 文件大小：B / KB / MB，最多一位小数（860KB、1.2MB）；取不到时留空
const fileSize = (img) => {
  const bytes = img.file?.size
  if (!bytes) return ''
  if (bytes < 1024) return `${bytes}B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1).replace(/\.0$/, '')}KB`
  return `${(bytes / 1024 / 1024).toFixed(1).replace(/\.0$/, '')}MB`
}

// 树里的区域总数（含子区域），用于在页卡上快速判断是否已标注
function nodeCount(img) {
  let n = 0
  const walk = (list) => (list || []).forEach((x) => { n += 1; walk(x.children) })
  walk(img.annotationTree)
  return n
}

function onSelectPage(img) {
  if (img.pageId !== currentImage.value?.pageId) store.setActivePage(img.pageId)
}

// 双击页卡：直接进入该页标注（列表里的页还没拖入画布时也能编辑）
const emit = defineEmits(['annotate'])
function onAnnotatePage(img) {
  if (!img.cosUrl) {
    ElMessage.warning(img.uploadError ? '该页上传失败，请先重试' : '该页还在上传中，请稍候')
    return
  }
  store.setActivePage(img.pageId)
  emit('annotate', img.pageId)
}

// ═══ 列表拖动：上下拖 = 改页序；拖到右侧画板 = 放入画布 ═══
const dragFrom = ref(-1)
const dropAt = ref(-1)

function onDragStart(e, i, img) {
  if (isSubmitting.value) {
    e.preventDefault()
    return
  }
  dragFrom.value = i
  dropAt.value = i
  e.dataTransfer.effectAllowed = 'copyMove'
  // 画板用这个自定义 key 识别「从页面列表拖入」
  e.dataTransfer.setData('application/x-dc-page', img.pageId)
  e.dataTransfer.setData('text/plain', img.pageId)
}

function onDragOver(i) {
  if (dragFrom.value < 0) return
  dropAt.value = i
}

function onDropItem(e, i) {
  // 只在列表内部拖动时重排；从别处拖入不处理
  if (dragFrom.value < 0) {
    resetDrag()
    return
  }
  store.moveImage(dragFrom.value, i)
  resetDrag()
}

function resetDrag() {
  dragFrom.value = -1
  dropAt.value = -1
}

const onRetryPage = (img) => {
  if (img.uploadId) store.retryUpload(img.uploadId)
}

// ── 文件选择：多图=添加页面（可多选），单图=更换主图 ──
const fileInput = ref(null)
const pendingReplace = ref(false)

// 更换/删除会丢掉当前页的框选节点与附图素材，先告知影响面再让用户决定
const impactText = () => {
  const parts = []
  const nodeCountVal = checkAnnotations(uiStore.nodeList).total
  if (nodeCountVal) parts.push(`${nodeCountVal} 个框选区域`)
  const assetCount = store.assets.filter((a) => a.ownerId === currentImage.value?.uploadId).length
  if (assetCount) parts.push(`${assetCount} 张附图素材`)
  return parts.length ? `将清空当前页的${parts.join('、')}，` : ''
}

async function onFileChange(e) {
  const files = Array.from(e.target.files)
  if (files.length && pendingReplace.value) {
    images.value.forEach((img) => store.clearAssetsOf(img.uploadId))
    store.clearImages()
    uiStore.loadTree([])
  }
  if (files.length) {
    const before = images.value.length
    await store.addImages(files)
    // 多图：只在还没有页面时定位新页；后续追加不抢当前页（否则画布会闪一下）
    const firstNew = images.value[before]
    if (mode.value === 'multi' && before === 0 && firstNew) store.setActivePage(firstNew.pageId)
  }
  pendingReplace.value = false
  if (fileInput.value) fileInput.value.value = ''
}

// 多图「添加页面」：受 maxImages 限制，多选时超出部分由 store 丢弃
function pickMore() {
  if (isSubmitting.value || images.value.length >= maxImages.value) return
  pendingReplace.value = false
  fileInput.value?.click()
}

// 多图「删除当前页」：关系边与素材由 store 一并清理
async function onDeletePage() {
  if (isSubmitting.value || !currentImage.value) return
  const img = currentImage.value
  const impact = impactText()
  try {
    await ElMessageBox.confirm(
      impact ? `${impact}删除后不可恢复。` : `确定删除「${pageName(img)}」？`,
      '删除页面',
      { type: 'warning', confirmButtonText: '继续', cancelButtonText: '取消' },
    )
  } catch {
    return
  }
  const idx = images.value.findIndex((i) => i.pageId === img.pageId)
  if (idx > -1) store.removeImage(idx)
}

// ── 单图：更换主图 ──
async function onReplace() {
  if (isSubmitting.value || !currentImage.value) return
  const impact = impactText()
  if (impact) {
    try {
      await ElMessageBox.confirm(`${impact}更换主图后需重新框选。`, '更换主图', {
        type: 'warning',
        confirmButtonText: '继续',
        cancelButtonText: '取消',
      })
    } catch {
      return
    }
  }
  pendingReplace.value = true
  fileInput.value?.click()
}

// ── 单图：删除主图 ──
async function onDeleteSingle() {
  if (isSubmitting.value) return
  const img = currentImage.value
  if (!img) return
  const impact = impactText()
  try {
    await ElMessageBox.confirm(
      impact ? `${impact}删除后不可恢复。` : '确定删除当前主图？',
      '删除主图',
      { type: 'warning', confirmButtonText: '继续', cancelButtonText: '取消' },
    )
  } catch {
    return
  }
  images.value.forEach((i) => store.clearAssetsOf(i.uploadId))
  store.clearImages()
  uiStore.loadTree([])
}
</script>

<style scoped>
.wb-upload {
  display: flex;
  flex-direction: column;
  gap: 10px;
  min-height: 0;
  height: 100%;
  overflow-y: auto;
  padding: 10px;
  border-radius: 14px;
  border: 1px solid rgba(0, 255, 255, 0.08);
  background: rgba(10, 14, 23, 0.7);
}

.mode-switch {
  display: flex;
  gap: 6px;
  padding: 3px;
  border-radius: 10px;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.08);
}
.mode-btn {
  flex: 1;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 6px 0;
  border-radius: 7px;
  border: none;
  background: transparent;
  color: rgba(255, 255, 255, 0.5);
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s;
}
.mode-btn:hover { color: rgba(255, 255, 255, 0.85); }
.mode-btn.active { background: rgba(0, 255, 255, 0.12); color: #00ffff; }
.mode-btn:disabled { opacity: 0.5; cursor: not-allowed; }

/* 列表上方的页面关系说明（仅多图） */
.rel-hint {
  flex-shrink: 0;
  display: flex;
  align-items: flex-start;
  gap: 7px;
  margin: 0;
  padding: 7px 9px;
  border-radius: 8px;
  border: 1px solid rgba(0, 255, 255, 0.12);
  background: rgba(0, 255, 255, 0.04);
  font-size: 11px;
  line-height: 1.6;
  color: rgba(255, 255, 255, 0.42);
}
.rel-hint i { margin-top: 2px; color: rgba(0, 255, 255, 0.5); }
.rel-hint b { color: rgba(255, 255, 255, 0.75); }

.panel {
  display: flex;
  flex-direction: column;
  min-height: 0;
  border-radius: 12px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  background: rgba(255, 255, 255, 0.02);
}
.panel--pages {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
}
.panel-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 9px 10px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
  font-size: 12px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.7);
}
.panel-head i { margin-right: 5px; color: rgba(0, 255, 255, 0.6); }
.panel-count { font-size: 11px; font-weight: 500; color: rgba(0, 255, 255, 0.7); }

.page-empty {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 28px 12px;
  text-align: center;
}
.page-empty i { font-size: 26px; color: rgba(255, 255, 255, 0.18); }
.page-empty p { margin: 0; font-size: 13px; font-weight: 600; color: rgba(255, 255, 255, 0.55); }
.page-empty span { font-size: 11px; line-height: 1.6; color: rgba(255, 255, 255, 0.3); }

/* ═══ 多图页面列表 ═══ */
.page-list { flex: 1; min-height: 0; display: flex; flex-direction: column; gap: 8px; padding: 8px; }
.page-item {
  display: flex;
  flex-direction: column;
  align-items: stretch;
  gap: 7px;
  padding: 6px;
  border-radius: 10px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  background: rgba(255, 255, 255, 0.02);
  cursor: pointer;
  text-align: left;
  transition: all 0.18s;
}
.page-item:hover { border-color: rgba(0, 255, 255, 0.25); }
.page-item.active {
  border-color: rgba(0, 255, 255, 0.55);
  background: rgba(0, 255, 255, 0.07);
  box-shadow: inset 0 0 0 1px rgba(0, 255, 255, 0.18);
}
/* 拖动排序：被拖行半透明，经过行用顶部高亮线标出插入位置 */
.page-item.is-dragging { opacity: 0.4; }
.page-item.is-drop-target { box-shadow: inset 0 2px 0 rgba(0, 255, 255, 0.75); }
/* 未放进画布的页：虚线边框提示；悬停缩略图时浮出「拖入画布」蒙层 */
.page-item.is-unplaced { border-style: dashed; }
.page-drag-hint {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 5px;
  font-size: 11px;
  font-weight: 600;
  color: rgba(0, 255, 255, 0.9);
  background: rgba(0, 0, 0, 0.58);
  opacity: 0;
  transition: opacity 0.18s;
  pointer-events: none;
}
.page-item-media:hover .page-drag-hint { opacity: 1; }
.page-item-media {
  position: relative;
  width: 100%;
  height: 112px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 7px;
  overflow: hidden;
  background: rgba(0, 0, 0, 0.4);
}
/* 上图下文：手机竖屏稿宽高比大，用 contain 保留整页，避免只截出中间一条 */
.page-item-media img { width: 100%; height: 100%; object-fit: contain; }
.page-item-media > i { font-size: 18px; color: rgba(255, 255, 255, 0.2); }
.page-no {
  position: absolute;
  top: 3px;
  left: 3px;
  padding: 1px 5px;
  border-radius: 5px;
  font-size: 10px;
  font-weight: 700;
  color: #00ffff;
  background: rgba(0, 0, 0, 0.6);
}
.page-entry {
  flex-shrink: 0;
  padding: 1px 6px;
  border-radius: 5px;
  font-size: 9px;
  font-weight: 700;
  color: #37e08d;
  background: rgba(55, 224, 141, 0.12);
  border: 1px solid rgba(55, 224, 141, 0.35);
}
.page-mask {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
  font-size: 10px;
  color: #00ffff;
  background: rgba(0, 0, 0, 0.6);
}
.page-mask--err { color: #ff6b7a; }
.page-item-info {
  flex: none;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 3px;
  padding: 0 2px 2px;
}
.page-item-title {
  display: flex;
  align-items: center;
  gap: 6px;
  min-width: 0;
}
.page-item-name {
  flex: 1;
  min-width: 0;
  font-size: 12px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.85);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.page-item-size {
  flex-shrink: 0;
  margin-left: auto;
  font-size: 10px;
  color: rgba(255, 255, 255, 0.38);
  font-variant-numeric: tabular-nums;
}

/* ═══ 单图主图卡 ═══ */
.page-card {
  flex-shrink: 0;
  margin: 8px;
  border-radius: 10px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  background: rgba(255, 255, 255, 0.03);
  overflow: hidden;
}
.page-media {
  position: relative;
  height: clamp(260px, 44vh, 420px);
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.35);
}
.page-media img { width: 100%; height: 100%; object-fit: contain; }
.page-media > i { font-size: 24px; color: rgba(255, 255, 255, 0.2); }
.thumb-mask {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.55);
  color: #00ffff;
  font-size: 14px;
}
.thumb-mask--err { color: #ff6b7a; }
.page-info { padding: 9px 10px 10px; }
.page-meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}
.page-dim { font-size: 11px; color: rgba(255, 255, 255, 0.35); }

.state-pill {
  flex-shrink: 0;
  padding: 2px 8px;
  border-radius: 999px;
  font-size: 10px;
  font-weight: 600;
}
.state-pill.uploading { background: rgba(0, 255, 255, 0.12); color: #00ffff; }
.state-pill.error { background: rgba(255, 71, 87, 0.15); color: #ff6b7a; }
.state-pill.ready { background: rgba(55, 224, 141, 0.12); color: #37e08d; }

.panel-title {
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* 换/删：图片下方的两个按钮 */
.file-actions { display: flex; flex-shrink: 0; gap: 8px; padding: 0 8px 10px; }
.act-btn {
  flex: 1;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 7px 0;
  border-radius: 8px;
  border: 1px solid rgba(255, 255, 255, 0.12);
  background: rgba(0, 0, 0, 0.3);
  color: rgba(255, 255, 255, 0.6);
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s;
}
.act-btn:hover:not(:disabled) { border-color: rgba(0, 255, 255, 0.4); color: #00ffff; }
.act-btn--danger:hover:not(:disabled) { border-color: rgba(255, 71, 87, 0.5); color: #ff6b7a; }
.act-btn:disabled { opacity: 0.4; cursor: not-allowed; }
</style>
