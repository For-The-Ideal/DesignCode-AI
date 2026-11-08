<template>
  <div class="canvas-wrap">
    <!-- 视口：fit 时铺满不滚动；Ctrl/⌘+滚轮以光标为锚放大缩小（不小于 fit），空格或中键拖拽平移 -->
    <div ref="viewportRef" class="canvas-scroll" :class="{ 'is-space': spaceDown, 'is-panning': panning, 'is-pannable': zoom > 1 }" @pointerdown="onViewportDown">
      <div
        ref="frameRef"
        class="img-frame"
        :class="{ 'img-frame--select': mode === 'select', 'img-frame--phone': phoneShell }"
        :style="{ width: dispW + 'px', height: dispH + 'px' }"
        @pointerdown="onFrameDown"
        @dragstart.prevent
      >
        <img :key="store.image.url" :src="store.image.url" alt="设计稿" draggable="false" />

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

    <!-- 模式提示：固定在画布底部留白区（左下），不再悬浮遮挡图片 -->
    <div class="mode-banner">
      <i class="fas" :class="mode === 'select' ? 'fa-mouse-pointer' : 'fa-crosshairs'"></i>
      <span class="banner-text">{{ bannerText }}</span>
      <button v-if="mode === 'child'" class="mode-cancel" @click="emit('cancelMode')">退出子级</button>
    </div>

    <!-- 缩放控制：缩小 / 百分比(点击复位到适应窗口) / 放大；
         内嵌到多图卡片时由父组件在卡片底部栏渲染（zoomBar=false） -->
    <div v-if="zoomBar" class="zoom-bar">
      <button class="zoom-btn" :disabled="zoom <= MIN_ZOOM" @click="zoomOut"><i class="fas fa-minus"></i></button>
      <span class="zoom-pct" title="点击复位到适应窗口" @click="resetZoom">{{ zoomPct }}</span>
      <button class="zoom-btn" :disabled="zoom >= MAX_ZOOM" @click="zoomIn"><i class="fas fa-plus"></i></button>
    </div>

    <!-- 多图切页：贴在可编辑图片左右两侧（不占工具栏），垂直居中 -->
    <button
      v-if="showNav && dispW"
      class="pg-nav pg-nav--prev"
      type="button"
      :disabled="!hasPrev"
      title="上一页"
      @click="emit('prevPage')"
    >
      <i class="fas fa-chevron-left"></i>
    </button>
    <button
      v-if="showNav && dispW"
      class="pg-nav pg-nav--next"
      type="button"
      :disabled="!hasNext"
      title="下一页"
      @click="emit('nextPage')"
    >
      <i class="fas fa-chevron-right"></i>
    </button>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount, nextTick, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { useUiProjectStore } from '~/stores/uiProject'
import { regionLabel, nodeLabel } from '~/utils/uiTree.js'

const store = useUiProjectStore()
const props = defineProps({
  // 'draw' 默认框选（已有框不拦截指针，随处可画新框）；'select' 选中/拖动/缩放已有框；'child' 框子区域（挂 pendingParentKey）
  mode: { type: String, default: 'draw' },
  // 多图切页箭头：是否显示 + 首尾可用性（由父组件按当前页信息传入）
  showNav: { type: Boolean, default: false },
  // 单图：给画布加手机壳装饰（圆角/黑边/刘海/Home 条，纯视觉，不参与坐标与缩放的像素计算）
  phoneShell: { type: Boolean, default: false },
  hasPrev: { type: Boolean, default: false },
  hasNext: { type: Boolean, default: false },
  // 是否显示自带缩放条：多图卡片把它放到卡片底部栏，不浮在图上
  zoomBar: { type: Boolean, default: true },
})
const emit = defineEmits(['drawn', 'cancelMode', 'prevPage', 'nextPage'])

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

// 手势状态（draw / move / resize）：fit 重算与像素换算都要参考它，提前声明
let interaction = null // { mode, start:{px,py}, origin, node }

const dispW = computed(() => (fitW.value ? Math.round(fitW.value * zoom.value) : 0))
const dispH = computed(() => (fitH.value ? Math.round(fitH.value * zoom.value) : 0))
const zoomPct = computed(() => Math.round(zoom.value * 100) + '%')
// 供父组件（卡片底部栏）控制的缩放状态与方法
const canZoomIn = computed(() => zoom.value < MAX_ZOOM)
const canZoomOut = computed(() => zoom.value > MIN_ZOOM)

// fit 内缩：给标注标签（上 -10px）与缩放手柄（右下 +7px）留出溢出空间，
// 它们贴边时不至于撑出滚动条或被裁掉
const FIT_GUTTER = 10

// 按容器可用尺寸对图片做 contain，求出 fit（100%）像素基准
// 必须用 offsetWidth/offsetHeight（不含滚动条）而不是 clientWidth/clientHeight：
// 用 client 尺寸会形成回路——装饰溢出 → 出滚动条 → 可用宽高变小 → fit 变小 → 图片变小
// → 溢出消失 → 滚动条消失 → fit 变大……自激振荡，表现就是框选时整张图一直抖
function computeFit() {
  const vp = viewportRef.value
  const iw = store.image.width
  const ih = store.image.height
  if (!vp || !iw || !ih) return
  const vw = vp.offsetWidth - FIT_GUTTER * 2
  const vh = vp.offsetHeight - FIT_GUTTER * 2
  if (vw <= 0 || vh <= 0) return
  const s = Math.min(vw / iw, vh / ih)
  const w = Math.floor(iw * s)
  const h = Math.floor(ih * s)
  if (w === fitW.value && h === fitH.value) return
  fitW.value = w
  fitH.value = h
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
let roRaf = 0
// 合并到下一帧再重算；手势进行中直接跳过——任何重排都会让画框的像素映射漂移，表现为抖动
function scheduleFit() {
  if (interaction || roRaf) return
  roRaf = requestAnimationFrame(() => { roRaf = 0; computeFit() })
}
onMounted(() => {
  computeFit()
  const vp = viewportRef.value
  if (vp) {
    vp.addEventListener('wheel', onWheel, { passive: false })
    if ('ResizeObserver' in window) { ro = new ResizeObserver(scheduleFit); ro.observe(vp) }
  }
  window.addEventListener('keydown', onKey)
  window.addEventListener('keyup', onKey)
})
onBeforeUnmount(() => {
  const vp = viewportRef.value
  if (vp) vp.removeEventListener('wheel', onWheel)
  if (ro) ro.disconnect()
  if (roRaf) cancelAnimationFrame(roRaf)
  unbindWindow()
  window.removeEventListener('keydown', onKey)
  window.removeEventListener('keyup', onKey)
  onPanUp()
})
watch(() => [store.image.width, store.image.height], scheduleFit)
watch(() => store.image.url, () => {
  zoom.value = 1
  // 切页/换图：放弃未完成的手势，避免旧图上的临时框残留到新图
  interaction = null
  drawingBox.value = null
  unbindWindow()
  releaseFrame()
})

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

// 子级框选的父级边界（原图像素）：仅 child 模式且父节点存在时限制，其余模式不限制
function parentBBox() {
  if (props.mode !== 'child' || !store.pendingParentKey) return null
  return store.nodes[store.pendingParentKey]?.bbox || null
}

// 父级显示名：优先区域类型/名称（如 banner → 横幅），未命名回退「父级」
function parentLabel() {
  if (!store.pendingParentKey) return '父级'
  return regionLabel(store.nodes[store.pendingParentKey]?.name) || '父级'
}

// 手势内的图片矩形快照：pointerdown 时定格，整个手势用同一基准换算像素。
// 即使期间面板/滚动条/尺寸有微变，画框、拖框也不会因基准漂移而跳
let frameSnap = null
function snapFrame() {
  const r = frameRef.value?.getBoundingClientRect()
  frameSnap = r ? { left: r.left, top: r.top, width: r.width, height: r.height } : null
}
function releaseFrame() { frameSnap = null }

// 客户端坐标 -> 原图像素点（含夹取）；子级模式下同时夹到父级 bbox，out 标记原始点是否越界
function toPixel(e) {
  const rect = frameSnap || frameRef.value.getBoundingClientRect()
  const rx = Math.min(Math.max((e.clientX - rect.left) / rect.width, 0), 1)
  const ry = Math.min(Math.max((e.clientY - rect.top) / rect.height, 0), 1)
  let px = rx * store.image.width
  let py = ry * store.image.height
  let out = false
  const pb = parentBBox()
  if (pb) {
    const cx = Math.min(Math.max(px, pb.x), pb.x + pb.w)
    const cy = Math.min(Math.max(py, pb.y), pb.y + pb.h)
    out = cx !== px || cy !== py
    px = cx
    py = cy
  }
  return { px, py, out }
}

// 顶层：区域类型（多块补序号）；子级：「元素类型 · 业务名」（画布标签与层级树同一套）
function label(n) {
  return nodeLabel(n, '未命名')
}

// ═══ 交互状态机：draw / move / resize ═══

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
  snapFrame() // 定格本次手势的坐标基准，之后所有 pointermove 都用它
  // draw/child 下空白处按下即画新框（已有框 pointer-events:none，不拦截）
  const { px, py, out } = toPixel(e)
  // 子级框选：起点在父级 bbox 之外时直接拦截并给出明确提醒
  if (out) {
    releaseFrame()
    ElMessage.warning(`子区域超出了「${parentLabel()}」区域边界，请重新框选`)
    return
  }
  interaction = { mode: 'draw', start: { px, py } }
  drawingBox.value = { x: px, y: py, w: 0, h: 0 }
  bindWindow()
}

function onBoxDown(e, node) {
  if (e.button !== 0) return
  store.selectNode(node.nodeKey)
  if (props.mode !== 'select') return // 仅编辑模式可拖动已有框
  snapFrame()
  const { px, py } = toPixel(e)
  interaction = { mode: 'move', start: { px, py }, origin: { ...node.bbox }, node }
  bindWindow()
}

function onHandleDown(e, node) {
  if (e.button !== 0) return
  store.selectNode(node.nodeKey)
  if (props.mode !== 'select') return
  snapFrame()
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

function onPointerUp(e) {
  unbindWindow()
  if (!interaction) {
    releaseFrame()
    return
  }

  if (interaction.mode === 'draw') {
    const box = drawingBox.value
    drawingBox.value = null
    // 子级框选：松手点或最终框超出父级 bbox 均拒绝提交，并提示具体区域类型（含 0.5px 容差）
    const pb = parentBBox()
    if (pb) {
      const endOut = e && e.type !== 'pointercancel' ? toPixel(e).out : false
      const EPS = 0.5
      const boxOut = !!box && (box.x < pb.x - EPS || box.y < pb.y - EPS
        || box.x + box.w > pb.x + pb.w + EPS || box.y + box.h > pb.y + pb.h + EPS)
      if (endOut || boxOut) {
        ElMessage.warning(`子区域超出了「${parentLabel()}」区域边界，请重新框选`)
        interaction = null
        releaseFrame()
        return
      }
    }
    // 忽略误点（宽高 < 原图 1%）
    if (box && box.w > store.image.width * 0.01 && box.h > store.image.height * 0.01) {
      emit('drawn', { ...box })
    }
  }
  interaction = null
  releaseFrame()
}

// 内嵌到卡片时隐藏自带缩放条，由卡片底部栏调用这些状态与方法
defineExpose({ zoom, zoomPct, zoomIn, zoomOut, resetZoom, canZoomIn, canZoomOut })
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
/* 视口：底部留出提示条/缩放条的高度，fit 时铺满不滚动；放大后 overflow:auto 出滚动条做平移。子项 margin:auto 实现「小则居中、大则可滚到四边」 */
.canvas-scroll {
  position: absolute;
  inset: 14px 14px 60px;
  display: flex;
  /* 未放大（fit）时禁止滚动：标注标签/手柄贴边不会撑出滚动条，
     从而避免「滚动条 ↔ fit 尺寸」互相触发造成框选抖动；放大后才可平移 */
  overflow: hidden;
}
.canvas-scroll.is-pannable { overflow: auto; }
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
  /* 切页时 key 变化重挂载，重放淡入，避免新图解码期间硬闪一下 */
  animation: img-fade-in 0.22s ease;
}
@keyframes img-fade-in {
  from { opacity: 0; }
  to { opacity: 1; }
}

/* 单图手机壳：纯装饰（黑边用 box-shadow 不占布局），标注坐标与缩放逻辑不变 */
.img-frame--phone {
  border-radius: 26px;
  box-shadow:
    0 0 0 7px #0d0d0d,
    0 0 0 9px #1f1f1f,
    0 18px 40px rgba(0, 0, 0, 0.55);
}
.img-frame--phone > img { border-radius: 26px; }
.img-frame--phone::before {
  content: '';
  position: absolute;
  top: 0;
  left: 50%;
  transform: translateX(-50%);
  width: 32%;
  height: 3%;
  min-height: 12px;
  background: #0d0d0d;
  border-radius: 0 0 14px 14px;
  pointer-events: none;
}
.img-frame--phone::after {
  content: '';
  position: absolute;
  bottom: 1.6%;
  left: 50%;
  transform: translateX(-50%);
  width: 30%;
  height: 1.1%;
  min-height: 3px;
  border-radius: 2px;
  background: rgba(0, 0, 0, 0.45);
  pointer-events: none;
}

.node-box {
  position: absolute;
  border: 2px solid rgba(0, 150, 170, 1);
  background: rgba(0, 150, 170, 0.08);
  box-sizing: border-box;
  cursor: move;
  /* 深色外描边：白底/浅色设计稿上也能看清边框（不影响原有语义配色） */
  box-shadow: 0 0 0 1px rgba(0, 0, 0, 0.55);
}
.node-box--empty { border-color: rgba(215, 45, 60, 1); background: rgba(215, 45, 60, 0.08); }
.node-box--ai_draft { border-color: rgba(190, 130, 0, 1); background: rgba(190, 130, 0, 0.1); }
.node-box--confirmed { border-color: rgba(15, 145, 85, 1); background: rgba(15, 145, 85, 0.1); }
.node-box--pending { border-color: rgba(100, 110, 125, 1); }
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
.node-box--static {
  pointer-events: none;
}
/* 不可编辑时手柄既不显示也不占滚动区域：贴边时不会多撑出 7px 造成滚动条 */
.node-box--static .resize-handle {
  display: none;
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
  left: 16px;
  bottom: 14px;
  z-index: 5;
  display: flex;
  align-items: center;
  gap: 8px;
  height: 36px; /* 与右下角缩放条同高同底，保持同一水平线 */
  max-width: calc(100% - 160px); /* 右侧给缩放条留位 */
  padding: 0 14px;
  box-sizing: border-box;
  border-radius: 18px;
  background: rgba(0, 0, 0, 0.55);
  border: 1px solid rgba(0, 255, 255, 0.25);
  color: #00ffff;
  font-size: 12px;
}
/* 文本过长时省略，不撑破底部单行布局 */
.banner-text {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
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
  bottom: 14px;
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

/* 多图切页：固定在图片编辑区域最左/最右，留一点内边距 */
.pg-nav {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  z-index: 6;
  width: 36px;
  height: 64px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 10px;
  border: 1px solid rgba(0, 255, 255, 0.28);
  background: rgba(10, 14, 23, 0.72);
  color: #00ffff;
  font-size: 15px;
  cursor: pointer;
  backdrop-filter: blur(4px);
  transition: all 0.18s ease;
}
.pg-nav:hover:not(:disabled) {
  border-color: rgba(0, 255, 255, 0.6);
  background: rgba(0, 255, 255, 0.16);
  box-shadow: 0 0 18px rgba(0, 255, 255, 0.16);
}
.pg-nav:disabled { opacity: 0.22; cursor: not-allowed; }
.pg-nav--prev { left: 12px; }
.pg-nav--next { right: 12px; }
</style>
