<template>
  <DialogModel ref="dialogRef" width="1040px" customClass="anno-dialog" @close="emit('close')">
    <template #header>
      <div class="modal-header">
        <h3><i class="fas fa-vector-square"></i> 设计稿标注</h3>
        <!-- 多张图切换 tab -->
        <div v-if="images.length > 1" class="img-tabs">
          <button
            v-for="(img, idx) in images"
            :key="idx"
            class="img-tab"
            :class="{ active: currentIdx === idx }"
            @click="switchImage(idx)"
          >
            {{ idx + 1 }}
            <span v-if="drafts[idx]?.annotations?.length" class="tab-dot">{{ drafts[idx].annotations.length }}</span>
          </button>
        </div>
      </div>
    </template>

    <div class="anno-layout">
      <!-- ═══ 左：画布区（img + div overlay 画框，不依赖任何图形库）═══ -->
      <div class="canvas-col">
        <div class="canvas-wrap">
          <!-- img-frame 收缩包裹 img，其尺寸 = 图片实际渲染区域，框坐标直接按百分比落位 -->
          <div ref="frameRef" class="img-frame" @pointerdown="onFrameDown">
            <img :src="currentImage?.preview" alt="设计稿" draggable="false" />

            <!-- 已标注的框 -->
            <div
              v-for="(anno, j) in currentDraft.annotations"
              :key="anno.id"
              class="anno-box"
              :class="{ 'anno-box--active': selectedId === anno.id }"
              :style="boxStyle(anno)"
              @pointerdown.stop="onBoxDown($event, anno)"
            >
              <span class="anno-num">{{ j + 1 }}</span>
              <span
                class="anno-resize-handle"
                @pointerdown.stop="onHandleDown($event, anno)"
              ></span>
            </div>

            <!-- 正在绘制的临时框 -->
            <div v-if="drawingBox" class="anno-box anno-box--drawing" :style="boxStyle(drawingBox)"></div>
          </div>
        </div>
        <div class="canvas-toolbar">
          <span class="toolbar-hint">
            <i class="fas fa-mouse-pointer"></i>
            在图上拖拽画框 · 拖动框移动 · 右下角调整大小
          </span>
          <button class="tool-btn" :disabled="!undoStack.length" @click="undo">
            <i class="fas fa-rotate-left"></i> 撤销
          </button>
          <button class="tool-btn tool-btn--danger" :disabled="!currentDraft.annotations.length" @click="clearAll">
            <i class="fas fa-trash-can"></i> 清空
          </button>
        </div>
      </div>

      <!-- ═══ 右：标注列表 + 描述编辑 ═══ -->
      <div class="panel-col">
        <div class="panel-title">标注列表（{{ currentDraft.annotations.length }}）</div>

        <div class="anno-list">
          <div
            v-for="(anno, j) in currentDraft.annotations"
            :key="anno.id"
            class="anno-item"
            :class="{ active: selectedId === anno.id }"
            @click="selectedId = anno.id"
          >
            <span class="anno-item-num">{{ j + 1 }}</span>
            <div class="anno-item-text">
              <span class="anno-item-label">{{ anno.label || '未命名区域' }}</span>
              <span class="anno-item-desc">{{ anno.desc || '点击填写描述…' }}</span>
            </div>
            <button class="anno-item-del" @click.stop="removeAnno(anno.id)">
              <i class="fa-regular fa-trash-can"></i>
            </button>
          </div>
          <div v-if="!currentDraft.annotations.length" class="anno-empty">
            暂无标注，在左侧图片上拖拽即可框选区域
          </div>
        </div>

        <!-- 选中框的编辑区 -->
        <div v-if="selectedAnno" class="anno-editor">
          <input
            v-model="selectedAnno.label"
            class="editor-label-input"
            maxlength="20"
            placeholder="区域名称（如：导航栏、商品卡片）"
          />
          <textarea
            ref="descInputRef"
            v-model="selectedAnno.desc"
            class="editor-desc-input"
            maxlength="300"
            rows="3"
            placeholder="这个区域是做什么的？AI 会参考它来生成代码…"
          ></textarea>
          <div class="editor-footer">
            <span class="char-count">{{ (selectedAnno.desc || '').length }}/300</span>
          </div>
        </div>

        <!-- 页面整体描述 -->
        <div class="page-desc">
          <label class="panel-subtitle">页面整体说明（可选）</label>
          <textarea
            v-model="currentDraft.description"
            class="page-desc-input"
            maxlength="500"
            rows="2"
            placeholder="整张设计稿的用途和风格…"
          ></textarea>
        </div>
      </div>
    </div>

    <template #footer>
      <div class="modal-footer">
        <button class="btn-cancel" @click="closeDialog">取消</button>
        <button class="btn-save" @click="save">保存标注</button>
      </div>
    </template>
  </DialogModel>
</template>

<script setup>
import { ref, computed, watch, nextTick } from 'vue'
import DialogModel from '~/components/dialog/DialogModel.vue'

const props = defineProps({
  images: { type: Array, default: () => [] },
  focusIndex: { type: Number, default: 0 },
})
const emit = defineEmits(['save', 'close'])

const dialogRef = ref(null)
const frameRef = ref(null)
const descInputRef = ref(null)

const currentIdx = ref(0)
const selectedId = ref(null)

// 每张图片的临时编辑数据: [{ annotations: [...], description: '' }]
const drafts = ref([])
// 撤销快照栈（按当前图索引存）
const undoStack = ref([])

const currentImage = computed(() => props.images[currentIdx.value])
const currentDraft = computed(() => drafts.value[currentIdx.value] || { annotations: [], description: '' })
const selectedAnno = computed(() => currentDraft.value.annotations.find(a => a.id === selectedId.value) || null)

// ═══ 打开时深拷贝 store 数据，关闭/保存前不污染 ═══
watch(() => dialogRef.value?.visible, (val) => {
  if (val) {
    drafts.value = props.images.map(img => ({
      description: img.description || '',
      annotations: (img.annotations || []).map(a => ({ ...a })),
    }))
    currentIdx.value = props.focusIndex ?? 0
    selectedId.value = null
    undoStack.value = []
  }
})

const open = () => dialogRef.value?.open()
const close = () => dialogRef.value?.close()
const closeDialog = () => close()
defineExpose({ open, close })

const switchImage = (idx) => {
  if (idx === currentIdx.value) return
  currentIdx.value = idx
  selectedId.value = null
  undoStack.value = []
}

// ═══ 坐标工具：基于 .img-frame 渲染矩形换算 0~1 相对坐标 ═══
const toRel = (e) => {
  const rect = frameRef.value.getBoundingClientRect()
  const x = (e.clientX - rect.left) / rect.width
  const y = (e.clientY - rect.top) / rect.height
  return {
    x: Math.min(Math.max(x, 0), 1),
    y: Math.min(Math.max(y, 0), 1),
  }
}

const boxStyle = (box) => ({
  left: box.x * 100 + '%',
  top: box.y * 100 + '%',
  width: box.w * 100 + '%',
  height: box.h * 100 + '%',
})

let annoSeq = 0
const genId = () => `a${Date.now()}_${++annoSeq}`

const snapshot = () => {
  undoStack.value.push(JSON.stringify(currentDraft.value.annotations))
}

const undo = () => {
  const last = undoStack.value.pop()
  if (last !== undefined) {
    currentDraft.value.annotations = JSON.parse(last)
    selectedId.value = null
  }
}

const clearAll = () => {
  snapshot()
  currentDraft.value.annotations = []
  selectedId.value = null
}

// ═══ 交互状态机：draw / move / resize 三种模式共用 pointermove/up ═══
const drawingBox = ref(null)
let interaction = null // { mode, startRel, origin, anno }

const bindWindowPointer = () => {
  window.addEventListener('pointermove', onPointerMove)
  window.addEventListener('pointerup', onPointerUp, { once: true })
}

// 1) 在空白处按下 → 开始画新框
const onFrameDown = (e) => {
  if (e.button !== 0) return
  const start = toRel(e)
  interaction = { mode: 'draw', startRel: start }
  drawingBox.value = { x: start.x, y: start.y, w: 0, h: 0 }
  bindWindowPointer()
}

// 2) 在框上按下 → 移动
const onBoxDown = (e, anno) => {
  if (e.button !== 0) return
  selectedId.value = anno.id
  interaction = { mode: 'move', startRel: toRel(e), origin: { ...anno }, anno }
  snapshot()
  bindWindowPointer()
}

// 3) 在 resize handle 上按下 → 调整大小
const onHandleDown = (e, anno) => {
  if (e.button !== 0) return
  selectedId.value = anno.id
  interaction = { mode: 'resize', origin: { ...anno }, anno }
  snapshot()
  bindWindowPointer()
}

const onPointerMove = (e) => {
  if (!interaction) return
  const p = toRel(e)

  if (interaction.mode === 'draw') {
    const s = interaction.startRel
    drawingBox.value = {
      x: Math.min(s.x, p.x),
      y: Math.min(s.y, p.y),
      w: Math.abs(p.x - s.x),
      h: Math.abs(p.y - s.y),
    }
    return
  }

  if (interaction.mode === 'move') {
    const { origin, anno } = interaction
    const dx = p.x - interaction.startRel.x
    const dy = p.y - interaction.startRel.y
    anno.x = Math.min(Math.max(origin.x + dx, 0), 1 - anno.w)
    anno.y = Math.min(Math.max(origin.y + dy, 0), 1 - anno.h)
    return
  }

  if (interaction.mode === 'resize') {
    const { origin, anno } = interaction
    anno.w = Math.min(Math.max(p.x - origin.x, 0.01), 1 - origin.x)
    anno.h = Math.min(Math.max(p.y - origin.y, 0.01), 1 - origin.y)
  }
}

const onPointerUp = () => {
  window.removeEventListener('pointermove', onPointerMove)
  if (!interaction) return

  if (interaction.mode === 'draw') {
    const box = drawingBox.value
    drawingBox.value = null
    // 忽略过小的误点击（阈值 1%）
    if (box && box.w > 0.01 && box.h > 0.01) {
      snapshot()
      const anno = { id: genId(), x: box.x, y: box.y, w: box.w, h: box.h, label: '', desc: '' }
      currentDraft.value.annotations.push(anno)
      selectedId.value = anno.id
      // 画完直接聚焦描述输入，流程不断
      nextTick(() => descInputRef.value?.focus())
    }
  }

  interaction = null
}

const removeAnno = (id) => {
  snapshot()
  currentDraft.value.annotations = currentDraft.value.annotations.filter(a => a.id !== id)
  if (selectedId.value === id) selectedId.value = null
}

// ═══ 保存：把 drafts 抛给父组件写回 store ═══
const save = () => {
  emit('save', drafts.value.map(d => ({
    description: d.description,
    annotations: d.annotations.map(a => ({ ...a })),
  })))
  close()
}
</script>

<style scoped>
/* ═══ header ═══ */
.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding: 4px 4px 0;
}
.modal-header h3 {
  font-size: 17px;
  color: #00ffff;
  display: flex;
  align-items: center;
  gap: 8px;
  white-space: nowrap;
}
.img-tabs { display: flex; gap: 8px; }
.img-tab {
  position: relative;
  width: 32px;
  height: 32px;
  border-radius: 8px;
  border: 1px solid rgba(0, 255, 255, 0.2);
  background: rgba(0, 0, 0, 0.3);
  color: rgba(255, 255, 255, 0.55);
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
}
.img-tab.active {
  border-color: #00ffff;
  color: #00ffff;
  background: rgba(0, 255, 255, 0.08);
}
.tab-dot {
  position: absolute;
  top: -6px;
  right: -6px;
  min-width: 15px;
  height: 15px;
  border-radius: 8px;
  background: linear-gradient(135deg, #00ffff, #ff00ff);
  color: #0a0a0f;
  font-size: 9px;
  font-weight: 700;
  line-height: 15px;
}

/* ═══ 主体两栏 ═══ */
.anno-layout {
  display: flex;
  gap: 16px;
  padding: 4px 0 8px;
}

/* ── 左列画布 ── */
.canvas-col {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 10px;
}
.canvas-wrap {
  height: 58vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.45);
  border: 1px solid rgba(0, 255, 255, 0.1);
  border-radius: 12px;
  padding: 12px;
  overflow: hidden;
}
/* 收缩包裹 img：框层尺寸 = 图片实际渲染尺寸，contain 黑边不再干扰坐标 */
.img-frame {
  position: relative;
  display: inline-block;
  max-width: 100%;
  max-height: 100%;
  cursor: crosshair;
  touch-action: none;
}
.img-frame img {
  display: block;
  max-width: 100%;
  max-height: 100%;
  user-select: none;
  pointer-events: none;
}

/* ── 标注框 ── */
.anno-box {
  position: absolute;
  border: 2px solid rgba(0, 255, 255, 0.85);
  background: rgba(0, 255, 255, 0.08);
  cursor: move;
  touch-action: none;
}
.anno-box--active {
  border-color: #ff00ff;
  background: rgba(255, 0, 255, 0.1);
  box-shadow: 0 0 0 1px rgba(255, 0, 255, 0.3), 0 0 12px rgba(255, 0, 255, 0.25);
  z-index: 2;
}
.anno-box--drawing {
  cursor: crosshair;
  border-style: dashed;
  pointer-events: none;
}
.anno-num {
  position: absolute;
  top: -10px;
  left: -2px;
  min-width: 18px;
  height: 18px;
  border-radius: 9px;
  background: #00ffff;
  color: #0a0a0f;
  font-size: 11px;
  font-weight: 700;
  line-height: 18px;
  text-align: center;
  pointer-events: none;
}
.anno-resize-handle {
  position: absolute;
  right: -5px;
  bottom: -5px;
  width: 12px;
  height: 12px;
  border-radius: 3px;
  background: #00ffff;
  border: 2px solid #0a0a0f;
  cursor: nwse-resize;
}

/* ── 画布工具条 ── */
.canvas-toolbar {
  display: flex;
  align-items: center;
  gap: 10px;
}
.toolbar-hint {
  flex: 1;
  font-size: 12px;
  color: rgba(255, 255, 255, 0.35);
  display: flex;
  align-items: center;
  gap: 6px;
}
.toolbar-hint i { color: rgba(0, 255, 255, 0.4); font-size: 11px; }
.tool-btn {
  background: rgba(0, 0, 0, 0.3);
  border: 1px solid rgba(0, 255, 255, 0.15);
  color: rgba(255, 255, 255, 0.6);
  font-size: 12px;
  padding: 5px 12px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
}
.tool-btn:hover:not(:disabled) { color: #00ffff; border-color: rgba(0, 255, 255, 0.4); }
.tool-btn:disabled { opacity: 0.35; cursor: not-allowed; }
.tool-btn--danger:hover:not(:disabled) { color: #ff4757; border-color: rgba(255, 71, 87, 0.4); }

/* ── 右列面板 ── */
.panel-col {
  flex: 0 0 320px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  max-height: calc(58vh + 40px);
}
.panel-title, .panel-subtitle {
  font-size: 13px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.6);
}
.panel-subtitle { margin-bottom: 6px; display: block; }

.anno-list {
  flex: 1;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 8px;
  min-height: 120px;
}
.anno-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  border-radius: 10px;
  background: rgba(0, 0, 0, 0.3);
  border: 1px solid rgba(0, 255, 255, 0.08);
  cursor: pointer;
  transition: all 0.2s;
}
.anno-item:hover { border-color: rgba(0, 255, 255, 0.25); }
.anno-item.active {
  border-color: #ff00ff;
  background: rgba(255, 0, 255, 0.05);
}
.anno-item-num {
  flex-shrink: 0;
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: linear-gradient(135deg, #00ffff, #ff00ff);
  color: #0a0a0f;
  font-size: 11px;
  font-weight: 700;
  line-height: 20px;
  text-align: center;
}
.anno-item-text {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 2px;
}
.anno-item-label {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.8);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.anno-item-desc {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.35);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.anno-item-del {
  flex-shrink: 0;
  background: none;
  border: none;
  color: rgba(255, 255, 255, 0.25);
  font-size: 13px;
  cursor: pointer;
  padding: 4px;
}
.anno-item-del:hover { color: #ff4757; }
.anno-empty {
  padding: 24px 12px;
  text-align: center;
  font-size: 12px;
  color: rgba(255, 255, 255, 0.25);
  border: 1px dashed rgba(0, 255, 255, 0.12);
  border-radius: 10px;
}

/* ── 选中框编辑区 ── */
.anno-editor {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 12px;
  border-radius: 10px;
  background: rgba(0, 255, 255, 0.04);
  border: 1px solid rgba(0, 255, 255, 0.15);
}
.editor-label-input, .editor-desc-input, .page-desc-input {
  width: 100%;
  background: rgba(0, 0, 0, 0.5);
  border: 1px solid rgba(0, 255, 255, 0.25);
  border-radius: 8px;
  padding: 8px 10px;
  color: #fff;
  font-size: 13px;
  font-family: inherit;
  outline: none;
  transition: all 0.2s;
}
.editor-label-input:focus, .editor-desc-input:focus, .page-desc-input:focus {
  border-color: #00ffff;
  box-shadow: 0 0 8px rgba(0, 255, 255, 0.1);
}
.editor-desc-input { resize: vertical; }
.editor-footer { display: flex; justify-content: flex-end; }
.char-count { font-size: 10px; color: rgba(255, 255, 255, 0.3); font-family: monospace; }

.page-desc-input { resize: vertical; }

/* ── footer ── */
.modal-footer { display: flex; justify-content: flex-end; gap: 12px; }
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
.btn-cancel:hover { background: rgba(255, 255, 255, 0.1); }
.btn-save {
  background: linear-gradient(135deg, #00ffff, #ff00ff);
  border: none;
  color: white;
}
.btn-save:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 15px rgba(0, 255, 255, 0.3);
}

/* ── 窄屏降级为上下布局 ── */
@media (max-width: 900px) {
  .anno-layout { flex-direction: column; }
  .panel-col { flex: none; max-height: none; }
}
</style>

<style>
/* DialogModel 内部结构样式（append-to-body 后 scoped 失效） */
.anno-dialog .el-dialog__body {
  padding: 12px 20px 0;
}
</style>
