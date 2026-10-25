<template>
  <div class="canvas-wrap">
    <!-- 视口：fit 时铺满不滚动；Ctrl/⌘+滚轮以光标为锚放大缩小（不小于 fit），空格或中键拖拽平移 -->
    <div ref="viewportRef" class="canvas-scroll" :class="{ 'is-space': spaceDown, 'is-panning': panning }" @pointerdown="onViewportDown">
      <div
        ref="frameRef"
        class="img-frame"
        :class="{ 'img-frame--select': mode === 'select' }"
        :style="{ width: dispW + 'px', height: dispH + 'px' }"
        @pointerdown="onFrameDown"
        @dragstart.prevent
      >
        <img :src="store.image.url" alt="设计稿" draggable="false" />

        <!-- 已有节点框（原图像素 -> 百分比叠加，缩放只在渲染层） -->
        <div
          v-for="n in store.nodeList"
          :key="n.nodeKey"
          class="node-box"
          :class="[
            'node-box--' + n.status,
            { 'node-box--active': store.selectedNodeKey === n.nodeKey, 'node-box--static': mode !== 'select' },
          ]"
          :style="boxStyle(n.bbox)"
          @pointerdown.stop="onBoxDown($event, n)"
        >
          <span class="node-tag">{{ label(n) }}</span>
          <span
            class="resize-handle"
            @pointerdown.stop="onHandleDown($event, n)"
          ></span>
        </div>

        <!-- 正在绘制的临时框 -->
        <div v-if="drawingBox" class="node-box node-box--drawing" :style="boxStyle(drawingBox)"></div>
      </div>
    </div>

    <div class="mode-banner">
      <i class="fas" :class="mode === 'select' ? 'fa-mouse-pointer' : 'fa-crosshairs'"></i>
      {{ bannerText }}
      <button v-if="mode === 'child'" class="mode-cancel" @click="emit('cancelMode')">退出子级</button>
    </div>

    <!-- 缩放控制：缩小 / 百分比(点击复位到适应窗口) / 放大 -->
    <div class="zoom-bar">
      <button class="zoom-btn" :disabled="zoom <= MIN_ZOOM" @click="zoomOut"><i class="fas fa-minus"></i></button>
      <span class="zoom-pct" title="点击复位到适应窗口" @click="resetZoom">{{ zoomPct }}</span>
      <button class="zoom-btn" :disabled="zoom >= MAX_ZOOM" @click="zoomIn"><i class="fas fa-plus"></i></button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount, nextTick, watch } from 'vue'
import { useUiProjectStore } from '~/stores/uiProject'

const store = useUiProjectStore()
const props = defineProps({
  // 'draw' 默认框选（已有框不拦截指针，随处可画新框）；'select' 选中/拖动/缩放已有框；'child' 框子区域（挂 pendingParentKey）
  mode: { type: String, default: 'draw' },
})
const emit = defineEmits(['drawn', 'cancelMode'])

const bannerText = computed(() => {
  if (props.mode === 'child') return '在大框内部拖拽框出子区域（可连续标多个）'
  if (props.mode === 'select') return '编辑模式：拖动已有框移动，右下角缩放'
  return '框选模式：直接在图上拖拽画新框（可覆盖已有框），松手后选类型'
})

const frameRef = ref(null)
const drawingBox = ref(null)

// ═══ 缩放 / 平移：以「适应视口」为 100% 基准，不小于 fit（MIN_ZOOM=1）═══
const viewportRef = ref(null)
const MIN_ZOOM = 1
const MAX_ZOOM = 4
const zoom = ref(1)
const fitW = ref(0)
const fitH = ref(0)
const spaceDown = ref(false)
const panning = ref(false)

const dispW = computed(() => (fitW.value ? Math.round(fitW.value * zoom.value) : 0))
const dispH = computed(() => (fitH.value ? Math.round(fitH.value * zoom.value) : 0))
const zoomPct = computed(() => Math.round(zoom.value * 100) + '%')

// 按容器可用尺寸对图片做 contain，求出 fit（100%）像素基准
function computeFit() {
  const vp = viewportRef.value
  const iw = store.image.width
  const ih = store.image.height
  if (!vp || !iw || !ih) return
  const vw = vp.clientWidth
  const vh = vp.clientHeight
  if (!vw || !vh) return
  const s = Math.min(vw / iw, vh / ih)
  fitW.value = Math.floor(iw * s)
  fitH.value = Math.floor(ih * s)
}

const clampZoom = (z) => Math.min(Math.max(z, MIN_ZOOM), MAX_ZOOM)

// 以光标（默认视口中心）为锚点缩放，缩放后回算滚动位置保持锚点不动
function zoomAt(factor, clientX, clientY) {
  const vp = viewportRef.value
  if (!vp || !dispW.value) return
  const nz = clampZoom(zoom.value * factor)
  if (nz === zoom.value) return
  const rect = vp.getBoundingClientRect()
  const fr = frameRef.value.getBoundingClientRect()
  const cx = clientX ?? rect.left + rect.width / 2
  const cy = clientY ?? rect.top + rect.height / 2
  const fx = fr.width ? (cx - fr.left) / fr.width : 0.5
  const fy = fr.height ? (cy - fr.top) / fr.height : 0.5
  zoom.value = nz
  nextTick(() => {
    const f = frameRef.value
    if (!f) return
    vp.scrollLeft = f.offsetLeft + fx * f.offsetWidth - (cx - rect.left)
    vp.scrollTop = f.offsetTop + fy * f.offsetHeight - (cy - rect.top)
  })
}
const zoomIn = () => zoomAt(1.25)
const zoomOut = () => zoomAt(1 / 1.25)
function resetZoom() {
  zoom.value = 1
  nextTick(() => { const vp = viewportRef.value; if (vp) { vp.scrollLeft = 0; vp.scrollTop = 0 } })
}

function onWheel(e) {
  if (!e.ctrlKey && !e.metaKey) return // 非缩放手势交给原生滚动平移
  e.preventDefault()
  zoomAt(Math.exp(-e.deltaY * 0.0015), e.clientX, e.clientY)
}

// 空格 或 鼠标中键 拖拽平移（左键留给画框）
let panStart = null
function onViewportDown(e) {
  if (!(spaceDown.value || e.button === 1)) return
  e.preventDefault()
  const vp = viewportRef.value
  panStart = { x: e.clientX, y: e.clientY, sl: vp.scrollLeft, st: vp.scrollTop }
  panning.value = true
  window.addEventListener('pointermove', onPanMove)
  window.addEventListener('pointerup', onPanUp)
}
function onPanMove(e) {
  if (!panStart) return
  const vp = viewportRef.value
  vp.scrollLeft = panStart.sl - (e.clientX - panStart.x)
  vp.scrollTop = panStart.st - (e.clientY - panStart.y)
}
function onPanUp() {
  panStart = null
  panning.value = false
  window.removeEventListener('pointermove', onPanMove)
  window.removeEventListener('pointerup', onPanUp)
}

const isTyping = (el) => {
  const t = el && el.tagName
  return t === 'INPUT' || t === 'TEXTAREA' || !!(el && el.isContentEditable)
}
function onKey(e) {
  if (e.code !== 'Space' || isTyping(e.target)) return
  if (e.type === 'keydown') {
    if (!e.repeat) spaceDown.value = true
    e.preventDefault()
  } else {
    spaceDown.value = false
  }
}

let ro = null
onMounted(() => {
  computeFit()
  const vp = viewportRef.value
  if (vp) {
    vp.addEventListener('wheel', onWheel, { passive: false })
    if ('ResizeObserver' in window) { ro = new ResizeObserver(computeFit); ro.observe(vp) }
  }
  window.addEventListener('keydown', onKey)
  window.addEventListener('keyup', onKey)
})
onBeforeUnmount(() => {
  const vp = viewportRef.value
  if (vp) vp.removeEventListener('wheel', onWheel)
  if (ro) ro.disconnect()
  window.removeEventListener('keydown', onKey)
  window.removeEventListener('keyup', onKey)
  onPanUp()
})
watch(() => [store.image.width, store.image.height], computeFit)
watch(() => store.image.url, () => { zoom.value = 1 })

// 原图像素坐标 <-> 渲染百分比：框样式用相对值，数据存原图像素
function boxStyle(bbox) {
  const w = store.image.width || 1
  const h = store.image.height || 1
  return {
    left: (bbox.x / w) * 100 + '%',
    top: (bbox.y / h) * 100 + '%',
    width: (bbox.w / w) * 100 + '%',
    height: (bbox.h / h) * 100 + '%',
  }
}

// 客户端坐标 -> 原图像素点（含夹取）
function toPixel(e) {
  const rect = frameRef.value.getBoundingClientRect()
  const rx = Math.min(Math.max((e.clientX - rect.left) / rect.width, 0), 1)
  const ry = Math.min(Math.max((e.clientY - rect.top) / rect.height, 0), 1)
  return { px: rx * store.image.width, py: ry * store.image.height }
}

function label(n) {
  if (!n.name) return '未命名'
  return n.variantIndex > 0 ? `${n.name}-${n.variantIndex + 1}` : n.name
}

// ═══ 交互状态机：draw / move / resize ═══
let interaction = null // { mode, start:{px,py}, origin, node }

function bindWindow() {
  // 绑定前先清理，避免上次未收尾时 pointermove 重复绑定导致跳框
  unbindWindow()
  window.addEventListener('pointermove', onPointerMove)
  window.addEventListener('pointerup', onPointerUp)
  window.addEventListener('pointercancel', onPointerUp)
}

function unbindWindow() {
  window.removeEventListener('pointermove', onPointerMove)
  window.removeEventListener('pointerup', onPointerUp)
  window.removeEventListener('pointercancel', onPointerUp)
}

function onFrameDown(e) {
  if (e.button !== 0) return
  if (spaceDown.value) return // 空格平移中，不画新框
  if (props.mode === 'select') return // 编辑模式不画新框（交由拖动/缩放）
  e.preventDefault() // 阻止浏览器对图片的原生拖拽（避免出现拖整张图的鬼影）
  // draw/child 下空白处按下即画新框（已有框 pointer-events:none，不拦截）
  const { px, py } = toPixel(e)
  interaction = { mode: 'draw', start: { px, py } }
  drawingBox.value = { x: px, y: py, w: 0, h: 0 }
  bindWindow()
}

function onBoxDown(e, node) {
  if (e.button !== 0) return
  store.selectNode(node.nodeKey)
  if (props.mode !== 'select') return // 仅编辑模式可拖动已有框
  const { px, py } = toPixel(e)
  interaction = { mode: 'move', start: { px, py }, origin: { ...node.bbox }, node }
  bindWindow()
}

function onHandleDown(e, node) {
  if (e.button !== 0) return
  store.selectNode(node.nodeKey)
  if (props.mode !== 'select') return
  interaction = { mode: 'resize', origin: { ...node.bbox }, node }
  bindWindow()
}

function onPointerMove(e) {
  if (!interaction) return
  const { px, py } = toPixel(e)
  const W = store.image.width
  const H = store.image.height

  if (interaction.mode === 'draw') {
    const s = interaction.start
    drawingBox.value = {
      x: Math.min(s.px, px),
      y: Math.min(s.py, py),
      w: Math.abs(px - s.px),
      h: Math.abs(py - s.py),
    }
    return
  }

  if (interaction.mode === 'move') {
    const { origin, node, start } = interaction
    const dx = px - start.px
    const dy = py - start.py
    const nx = Math.min(Math.max(origin.x + dx, 0), W - origin.w)
    const ny = Math.min(Math.max(origin.y + dy, 0), H - origin.h)
    store.moveNode(node.nodeKey, { x: nx, y: ny, w: origin.w, h: origin.h })
    return
  }

  if (interaction.mode === 'resize') {
    const { origin, node } = interaction
    const nw = Math.min(Math.max(px - origin.x, 4), W - origin.x)
    const nh = Math.min(Math.max(py - origin.y, 4), H - origin.y)
    store.moveNode(node.nodeKey, { x: origin.x, y: origin.y, w: nw, h: nh })
  }
}

function onPointerUp() {
  unbindWindow()
  if (!interaction) return

  if (interaction.mode === 'draw') {
    const box = drawingBox.value
    drawingBox.value = null
    // 忽略误点（宽高 < 原图 1%）
    if (box && box.w > store.image.width * 0.01 && box.h > store.image.height * 0.01) {
      emit('drawn', { ...box })
    }
  }
  interaction = null
}
</script>

<style scoped>
.canvas-wrap {
  position: relative;
  height: 100%;
  background: rgba(0, 0, 0, 0.45);
  border: 1px solid rgba(0, 255, 255, 0.1);
  border-radius: 12px;
  overflow: hidden;
}
/* 视口：fit 时铺满不滚动；放大后 overflow:auto 出滚动条做平移。子项 margin:auto 实现「小则居中、大则可滚到四边」 */
.canvas-scroll {
  position: absolute;
  inset: 14px;
  display: flex;
  overflow: auto;
}
/* 平移态：空格或拖拽中，整块显示抓手（含已有框），优先于画框/编辑光标 */
.canvas-scroll.is-space .img-frame,
.canvas-scroll.is-space .node-box { cursor: grab; }
.canvas-scroll.is-panning .img-frame,
.canvas-scroll.is-panning .node-box { cursor: grabbing; }
/* 图片框：尺寸由 JS 按 fit×zoom 设像素，宽高即图片矩形，标注框百分比对齐 */
.img-frame {
  position: relative;
  margin: auto;
  flex: 0 0 auto;
  cursor: crosshair;
  touch-action: none;
}
/* 编辑模式：空白处不画框，用默认箭头（已有框另用 move） */
.img-frame--select { cursor: default; }
.img-frame img {
  display: block;
  width: 100%;
  height: 100%;
  object-fit: contain;
  user-select: none;
  pointer-events: none;
}

.node-box {
  position: absolute;
  border: 2px solid rgba(0, 255, 255, 0.8);
  background: rgba(0, 255, 255, 0.06);
  box-sizing: border-box;
  cursor: move;
  /* 深色外描边：白底/浅色设计稿上也能看清边框（不影响原有语义配色） */
  box-shadow: 0 0 0 1px rgba(0, 0, 0, 0.55);
}
.node-box--empty { border-color: rgba(255, 71, 87, 0.8); background: rgba(255, 71, 87, 0.06); }
.node-box--ai_draft { border-color: rgba(255, 207, 77, 0.85); background: rgba(255, 207, 77, 0.08); }
.node-box--confirmed { border-color: rgba(55, 224, 141, 0.85); background: rgba(55, 224, 141, 0.08); }
.node-box--pending { border-color: rgba(255, 255, 255, 0.45); }
.node-box--active {
  border-color: #ff00ff !important;
  box-shadow: 0 0 0 1px rgba(0, 0, 0, 0.55), 0 0 0 3px rgba(255, 0, 255, 0.35), 0 0 14px rgba(255, 0, 255, 0.3);
  z-index: 3;
}
.node-box--drawing {
  border-style: dashed;
  pointer-events: none;
}
/* 画框模式（child）下，已有框不拦截指针，允许在任意区域（含父框内）叠加新框 */
.node-box--static,
.node-box--static .resize-handle {
  pointer-events: none;
}

.node-tag {
  position: absolute;
  top: -10px;
  left: -2px;
  max-width: 120px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  padding: 0 6px;
  height: 16px;
  line-height: 16px;
  border-radius: 8px;
  background: #00ffff;
  color: #0a0a0f;
  font-size: 10px;
  font-weight: 700;
}

.resize-handle {
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

.mode-banner {
  position: absolute;
  top: 12px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 12px;
  border-radius: 20px;
  background: rgba(0, 0, 0, 0.7);
  border: 1px solid rgba(0, 255, 255, 0.35);
  color: #00ffff;
  font-size: 12px;
  z-index: 5;
}
.mode-cancel {
  margin-left: 4px;
  background: none;
  border: none;
  color: rgba(255, 255, 255, 0.6);
  cursor: pointer;
  font-size: 12px;
}
.mode-cancel:hover { color: #ff4757; }

/* 缩放控制条 */
.zoom-bar {
  position: absolute;
  right: 16px;
  bottom: 16px;
  z-index: 6;
  display: flex;
  align-items: center;
  gap: 2px;
  padding: 4px;
  border-radius: 10px;
  background: rgba(0, 0, 0, 0.7);
  border: 1px solid rgba(0, 255, 255, 0.25);
}
.zoom-btn {
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  border-radius: 7px;
  background: transparent;
  color: #00ffff;
  cursor: pointer;
  font-size: 12px;
  transition: background 0.2s;
}
.zoom-btn:hover:not(:disabled) { background: rgba(0, 255, 255, 0.12); }
.zoom-btn:disabled { color: rgba(255, 255, 255, 0.25); cursor: not-allowed; }
.zoom-pct {
  min-width: 48px;
  text-align: center;
  font-size: 12px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.8);
  cursor: pointer;
  user-select: none;
  font-variant-numeric: tabular-nums;
}
.zoom-pct:hover { color: #00ffff; }
</style>
