<template>
  <!-- 多图画板：页面自由摆放 + 页序连线；当前页卡片就地放大并内嵌标注画布（不再单开标注视图） -->
  <div class="wb-board">
    <div class="board-tip">
      <i class="fas fa-hand-pointer"></i>
      <span>拖卡片排版；点图片选中该页后才能在图上框选，选中页拖底部标题移动</span>
    </div>

    <div
      ref="viewportRef"
      class="board-viewport"
      :class="{ 'is-drop': dropActive }"
      @dragover.prevent="onDragOver"
      @dragleave="onDragLeave"
      @drop.prevent="onDrop"
    >
      <div ref="canvasRef" class="board-canvas" :style="{ width: world.w + 'px', height: world.h + 'px' }">
        <!-- 页序连线：按列表顺序连依次相邻的两张图 -->
        <svg class="board-links" :width="world.w" :height="world.h">
          <defs>
            <marker id="wb-board-arrow" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="7" markerHeight="7" orient="auto-start-reverse">
              <path d="M0,0 L10,5 L0,10 z" fill="rgba(0, 255, 255, 0.65)" />
            </marker>
          </defs>
          <path
            v-for="seg in arrows"
            :key="seg.key"
            :d="seg.d"
            class="link-path"
            marker-end="url(#wb-board-arrow)"
          />
        </svg>

        <div
          v-for="img in pages"
          :key="img.pageId"
          class="board-card"
          :class="{ active: isActive(img), dragging: img.pageId === drag.pageId }"
          :style="cardStyle(img)"
          @pointerdown="onCardDown($event, img)"
          @dblclick="onCardDblClick(img)"
        >
          <div class="card-media" :style="{ height: cardHeight(img) + 'px' }">
            <!-- 当前页：内嵌标注画布（无手机壳），底部工具条的框选/编辑直接作用在这张图上 -->
            <AnnotationCanvas
              v-if="isActive(img) && img.cosUrl"
              :ref="setCardCanvas"
              class="card-canvas"
              :mode="mode"
              :show-nav="false"
              :phone-shell="false"
              :zoom-bar="false"
              @drawn="(bbox) => emit('drawn', bbox)"
              @cancelMode="emit('cancelMode')"
            />
            <template v-else>
              <!-- 图片源与当前页画布保持一致（优先远程地址），切页时浏览器直接命中缓存，不闪 -->
              <img v-if="img.cosUrl || img.preview" class="card-still" :src="img.cosUrl || img.preview" alt="" draggable="false" />
              <i v-else class="fas fa-image"></i>
              <!-- 非当前页：把该页已保存的标注叠成只读框线 + 区域名标签，切到其它页时仍能对照它标过哪些范围 -->
              <div class="card-marks" aria-hidden="true">
                <span
                  v-for="m in marksOf(img)"
                  :key="m.id"
                  class="card-mark"
                  :class="[
                    'card-mark--' + (m.status || 'pending'),
                    { 'card-mark--child': !m.top, 'card-mark--flush': m.y * cardHeight(img) < 14 },
                  ]"
                  :style="markStyle(m)"
                >
                  <i v-if="markLabel(m)" class="card-mark-tag">{{ markLabel(m) }}</i>
                </span>
              </div>
            </template>
            <span class="card-no">P{{ pageNo(img) }}</span>
            <span v-if="img.uploading" class="card-mask"><i class="fas fa-spinner fa-spin"></i></span>
            <span v-else-if="img.uploadError" class="card-mask card-mask--err">上传失败</span>
          </div>
          <!-- 底部栏：入口标识 + 当前页缩放（图面不再浮控件，留给框选/编辑） -->
          <div class="card-foot" :class="{ 'is-handle': isActive(img) }" :title="isActive(img) ? '拖动这里移动当前页' : ''">
            <span v-if="isEntryPage(img)" class="card-entry">入口</span>
            <div v-if="isActive(img)" class="card-zoom" @pointerdown.stop @dblclick.stop>
              <button class="cz-btn" :disabled="!cardCanvas?.canZoomOut" title="缩小" @click="cardCanvas?.zoomOut()">
                <i class="fas fa-minus"></i>
              </button>
              <span class="cz-pct" title="点击复位到适应卡片" @click="cardCanvas?.resetZoom()">{{ cardCanvas?.zoomPct || '100%' }}</span>
              <button class="cz-btn" :disabled="!cardCanvas?.canZoomIn" title="放大" @click="cardCanvas?.zoomIn()">
                <i class="fas fa-plus"></i>
              </button>
            </div>
          </div>
        </div>

        <div v-if="!pages.length" class="board-empty">
          <i class="fa-regular fa-images"></i>
          <p>画布还没有页面</p>
          <span>从左侧页面列表拖一张图片进来开始排版</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, watch, nextTick, onMounted, onBeforeUnmount } from 'vue'
import { ElMessage } from 'element-plus'
import AnnotationCanvas from '~/components/uiProject/AnnotationCanvas.client.vue'
import { useCodeStore } from '~/stores/code'
import { nodeLabel } from '~/utils/uiTree.js'

// draw | select | child：透传给当前页内嵌的标注画布（由底部工具条切换）
defineProps({
  mode: { type: String, default: 'draw' },
})
const store = useCodeStore()
const emit = defineEmits(['drawn', 'cancelMode', 'dropFiles'])

// 卡片几何：所有页同一规则（固定宽 + 图片比例高），与「是否当前页」无关。
// 切换页面时任何卡片都不会改变尺寸，画布视觉绝对静止。
const CARD_W = 300
const MAX_CARD_H = 860
// 卡片底部标题栏占位（world 计算用，与 .card-foot 固定高一致）
const FOOT_H = 34

const viewportRef = ref(null)
const canvasRef = ref(null)
// 当前页内嵌画布实例：底部栏的缩放按钮直接驱动它（同一时刻只有当前页有画布）。
// v-for 里的字符串 ref 会变成数组，用函数 ref 拿到单个实例
const cardCanvas = ref(null)
const setCardCanvas = (el) => { cardCanvas.value = el }
// 视口尺寸：仅用于给画布铺底，不参与卡片几何（切页仍不会改变任何尺寸）
const viewport = reactive({ w: 0, h: 0 })
let ro = null

// 卡片按图片比例取几何：宽固定；超高长图以高度上限反推宽度，图面不留空白
function pageBox(img) {
  const W = img.naturalWidth || img.treeWidth || 0
  const H = img.naturalHeight || img.treeHeight || 0
  if (!W || !H) return { w: CARD_W, h: Math.round(CARD_W * 1.4) }
  const h = Math.round((CARD_W * H) / W)
  if (h <= MAX_CARD_H) return { w: CARD_W, h }
  return { w: Math.max(140, Math.round((MAX_CARD_H * W) / H)), h: MAX_CARD_H }
}

const isActive = (img) => !!img.pageId && img.pageId === activePageId.value
const cardW = (img) => pageBox(img).w
const cardHeight = (img) => pageBox(img).h

// 只渲染已拖入画布的页；顺序仍以页面列表为准（P 序号即列表序号）
const pages = computed(() => {
  const all = store.images
  return all
    .map((img, i) => ({ img, no: i }))
    .filter(({ img }) => !!img.board)
    .sort((a, b) => a.no - b.no)
    .map(({ img }) => img)
})
const activePageId = computed(() => store.currentImage?.pageId || '')
const pageNo = (img) => store.images.findIndex((i) => i.pageId === img.pageId) + 1
// 入口页恒为页面列表第 1 张：列表排序/增删后，画布上的入口标识跟着换
const isEntryPage = (img) => store.images[0]?.pageId === img.pageId

// ═══ 非当前页的只读框线 ═══
// 每页数据独立：切页时当前页的树已 saveToImage 写回该页，所以这里读到的就是最新标注。
// 优先用提交格式的 annotations（已是 0~1 百分比，零换算）；
// 老数据只有 annotationTree 时按该图原尺寸现算，并按 tree 引用缓存，避免拖卡片时重复遍历。
const treeMarkCache = new WeakMap()
function marksOf(img) {
  const saved = img.annotations
  if (Array.isArray(saved) && saved.length) {
    return saved.filter((m) => m.w > 0 && m.h > 0)
  }
  const tree = img.annotationTree
  if (!Array.isArray(tree) || !tree.length) return []
  const W = img.naturalWidth || img.treeWidth || 0
  const H = img.naturalHeight || img.treeHeight || 0
  if (!W || !H) return []
  const cached = treeMarkCache.get(img)
  if (cached && cached.tree === tree) return cached.list
  const list = []
  const walk = (nodes) => (nodes || []).forEach((n) => {
    const b = n.bbox || {}
    if (b.w > 0 && b.h > 0) {
      list.push({
        id: n.nodeKey,
        label: nodeLabel(n, '未命名'),
        top: !n.parentId,
        status: n.status || 'pending',
        x: b.x / W,
        y: b.y / H,
        w: b.w / W,
        h: b.h / H,
      })
    }
    walk(n.children)
  })
  walk(tree)
  treeMarkCache.set(img, { tree, list })
  return list
}
const markStyle = (m) => ({
  left: m.x * 100 + '%',
  top: m.y * 100 + '%',
  width: m.w * 100 + '%',
  height: m.h * 100 + '%',
})
// 标签只显示有业务名的框：未选类型/未命名的框不刷一排「未命名」
const markLabel = (m) => (m.label && m.label !== '未命名' ? m.label : '')

// ═══ 世界尺寸：只覆盖卡片真实占用范围，不足一屏时铺满视口 ═══
// 不再预留大留白 —— 只有一张图时画布恰好等于视口，不会出现多余滚动条；
// 卡片几何只取决于图片本身，切页不会改变世界大小，滚动位置保持不动
const world = computed(() => {
  let w = 0
  let h = 0
  pages.value.forEach((img) => {
    w = Math.max(w, img.board.x + cardW(img))
    h = Math.max(h, img.board.y + cardHeight(img) + FOOT_H)
  })
  return { w: Math.max(w, viewport.w), h: Math.max(h, viewport.h) }
})

// ═══ 页序连线：相邻两页按相对位置选连线边（右→左 / 左→右 / 下→上）═══
const arrows = computed(() => {
  const list = pages.value
  const segs = []
  for (let i = 0; i < list.length - 1; i++) {
    const a = list[i]
    const b = list[i + 1]
    const aw = cardW(a)
    const bw = cardW(b)
    const ha = cardHeight(a)
    const hb = cardHeight(b)
    const ax = a.board.x
    const ay = a.board.y
    const bx = b.board.x
    const by = b.board.y
    let d = ''
    if (bx >= ax + aw) {
      const sx = ax + aw
      const sy = ay + ha / 2
      const ex = bx
      const ey = by + hb / 2
      const cx = Math.max(60, (ex - sx) / 2)
      d = `M${sx},${sy} C${sx + cx},${sy} ${ex - cx},${ey} ${ex},${ey}`
    } else if (bx + bw <= ax) {
      const sx = ax
      const sy = ay + ha / 2
      const ex = bx + bw
      const ey = by + hb / 2
      const cx = Math.max(60, (sx - ex) / 2)
      d = `M${sx},${sy} C${sx - cx},${sy} ${ex + cx},${ey} ${ex},${ey}`
    } else {
      const sx = ax + aw / 2
      const sy = ay + ha
      const ex = bx + bw / 2
      const ey = by
      const cy = Math.max(40, (ey - sy) / 2)
      d = `M${sx},${sy} C${sx},${sy + cy} ${ex},${ey - cy} ${ex},${ey}`
    }
    segs.push({ key: `${a.pageId}-${b.pageId}`, d })
  }
  return segs
})

// ═══ 卡片拖动：拖动中只改本地坐标（不写 store），松手才落库 ═══
const drag = reactive({ pageId: '', sx: 0, sy: 0, ox: 0, oy: 0, moved: false })
const dragPos = reactive({ x: 0, y: 0 })

const cardStyle = (img) => {
  const x = img.pageId === drag.pageId ? dragPos.x : img.board.x
  const y = img.pageId === drag.pageId ? dragPos.y : img.board.y
  return { left: `${x}px`, top: `${y}px`, width: `${cardW(img)}px` }
}

function onCardDown(e, img) {
  if (e.button !== 0) return
  // 当前页的图面就是标注区：只有底部标题栏可以拖动它，避免与框选/编辑抢指针
  if (isActive(img) && !e.target.closest('.card-foot')) return
  drag.pageId = img.pageId
  drag.sx = e.clientX
  drag.sy = e.clientY
  drag.ox = img.board.x
  drag.oy = img.board.y
  drag.moved = false
  dragPos.x = img.board.x
  dragPos.y = img.board.y
  window.addEventListener('pointermove', onDragMove)
  window.addEventListener('pointerup', onDragUp)
}

function onDragMove(e) {
  if (!drag.pageId) return
  const dx = e.clientX - drag.sx
  const dy = e.clientY - drag.sy
  if (!drag.moved && Math.hypot(dx, dy) < 4) return
  drag.moved = true
  dragPos.x = Math.max(0, drag.ox + dx)
  dragPos.y = Math.max(0, drag.oy + dy)
}

function onDragUp() {
  if (!drag.pageId) return
  const id = drag.pageId
  const moved = drag.moved
  drag.pageId = ''
  drag.moved = false
  window.removeEventListener('pointermove', onDragMove)
  window.removeEventListener('pointerup', onDragUp)
  if (moved) store.moveOnBoard(id, dragPos.x, dragPos.y)
  // 单击（未拖动）才切当前页：拖动中放大卡片会跟着手指长，排版时很跳
  else if (id !== activePageId.value) store.setActivePage(id)
}

// ═══ 从页面列表拖入画布 ═══
const dropActive = ref(false)

function onDragOver(e) {
  const dt = e.dataTransfer
  if (!dt) return
  // 只接管「页面列表拖入」与「图片文件拖入」
  if (!dt.types?.includes('application/x-dc-page') && !dt.types?.includes('Files')) return
  dt.dropEffect = 'copy'
  dropActive.value = true
}

function onDragLeave(e) {
  if (e.currentTarget === e.target || !e.currentTarget.contains(e.relatedTarget)) dropActive.value = false
}

function onDrop(e) {
  dropActive.value = false
  // 直接拖图片文件到画布：交给外层统一走登录校验 + 追加页面
  const files = Array.from(e.dataTransfer?.files || []).filter((f) => f.type.startsWith('image/'))
  if (files.length) {
    emit('dropFiles', files)
    return
  }
  const pageId = e.dataTransfer?.getData('application/x-dc-page') || e.dataTransfer?.getData('text/plain')
  const img = store.images.find((i) => i.pageId === pageId)
  if (!img) return
  const rect = canvasRef.value?.getBoundingClientRect()
  if (!rect) return
  // 以光标为卡片中心落点，避免每次都叠在左上角
  const x = Math.max(0, e.clientX - rect.left - cardW(img) / 2)
  const y = Math.max(0, e.clientY - rect.top - cardHeight(img) / 2)
  // 拖入画板默认只做排版（卡片只能拖动），不抢当前页；
  // 画板第一张 或 页面列表第一张（入口页）例外：拖入即选中，上来就能框选
  const firstOnBoard = !pages.value.length
  store.placeOnBoard(img.pageId, x, y)
  if (firstOnBoard || isEntryPage(img)) store.setActivePage(img.pageId)
}

// ═══ 标注：把该页设为当前页即可就地框选（不再单开标注视图）═══
function focusPage(img) {
  if (!img.cosUrl) {
    ElMessage.warning(img.uploadError ? '该页上传失败，请先在页面列表重试' : '该页还在上传中，请稍候')
    return
  }
  store.setActivePage(img.pageId)
}
function onCardDblClick(img) {
  if (!isActive(img)) focusPage(img)
}

// ═══ 默认展示第一页：进入多图/首次有图时把列表第一页放到画布 ═══
watch(
  () => [store.mode, store.images.length],
  () => store.ensureFirstOnBoard(),
  { immediate: true },
)

onMounted(() => {
  store.ensureFirstOnBoard()
  // 量视口内容区（clientWidth 不含滚动条）：画布铺底尺寸恰好等于可视区，
  // 内容不足一屏时不会凭空多出滚动条
  const vp = viewportRef.value
  if (vp) {
    const sync = () => { viewport.w = vp.clientWidth; viewport.h = vp.clientHeight }
    sync()
    if ('ResizeObserver' in window) { ro = new ResizeObserver(sync); ro.observe(vp) }
  }
  // 进入时只滚一次到当前页；后续切页不自动滚动，避免画面跳来跳去
  const start = pages.value.find((p) => p.pageId === store.currentImage?.pageId) || pages.value[0]
  if (start?.board) nextTick(() => scrollToPage(start.pageId))
})

function scrollToPage(pageId) {
  const vp = viewportRef.value
  const img = store.images.find((i) => i.pageId === pageId)
  if (!vp || !img?.board) return
  // 已经在视口里就不动滚动：进入时画面保持原样，不被推一下
  const inView =
    img.board.x >= vp.scrollLeft &&
    img.board.x + cardW(img) <= vp.scrollLeft + vp.clientWidth &&
    img.board.y >= vp.scrollTop &&
    img.board.y + cardHeight(img) <= vp.scrollTop + vp.clientHeight
  if (inView) return
  vp.scrollTo({ left: Math.max(0, img.board.x - 16), top: Math.max(0, img.board.y - 12), behavior: 'auto' })
}

onBeforeUnmount(() => {
  if (ro) ro.disconnect()
  window.removeEventListener('pointermove', onDragMove)
  window.removeEventListener('pointerup', onDragUp)
})
</script>

<style scoped>
.wb-board {
  position: relative;
  height: 100%;
  min-height: 0;
  border-radius: 14px;
  border: 1px solid rgba(0, 255, 255, 0.08);
  background: rgba(10, 14, 23, 0.7);
  overflow: hidden;
}

.board-tip {
  position: absolute;
  top: 10px;
  left: 14px;
  z-index: 30;
  display: inline-flex;
  align-items: center;
  gap: 7px;
  padding: 5px 10px;
  border-radius: 8px;
  border: 1px solid rgba(0, 255, 255, 0.12);
  background: rgba(8, 12, 20, 0.82);
  font-size: 11px;
  color: rgba(255, 255, 255, 0.45);
  pointer-events: none;
}
.board-tip i { color: rgba(0, 255, 255, 0.6); }

/* 画板：默认滚动（不做缩放），拖入时高亮边框 */
.board-viewport {
  position: absolute;
  inset: 0;
  overflow: auto;
  padding: 0;
}
.board-viewport.is-drop {
  box-shadow: inset 0 0 0 2px rgba(0, 255, 255, 0.45);
}

.board-canvas {
  position: relative;
  /* 网格底纹：帮助判断页面排版与对齐 */
  background-image:
    radial-gradient(rgba(255, 255, 255, 0.06) 1px, transparent 1px);
  background-size: 24px 24px;
}

.board-links {
  position: absolute;
  inset: 0;
  pointer-events: none;
  overflow: visible;
}
.link-path {
  fill: none;
  stroke: rgba(0, 255, 255, 0.45);
  stroke-width: 1.6;
  stroke-dasharray: 6 5;
}

/* 页面卡片 */
.board-card {
  position: absolute;
  display: flex;
  flex-direction: column;
  border-radius: 10px;
  border: 1px solid rgba(255, 255, 255, 0.12);
  background: rgba(12, 17, 26, 0.95);
  box-shadow: 0 10px 26px rgba(0, 0, 0, 0.4);
  cursor: grab;
  user-select: none;
  transition: border-color 0.15s, box-shadow 0.15s;
}
.board-card:hover { border-color: rgba(0, 255, 255, 0.4); }
/* 当前页高亮：放大后要盖住周围卡片 */
.board-card.active {
  border-color: rgba(0, 255, 255, 0.75);
  box-shadow: 0 0 0 1px rgba(0, 255, 255, 0.3), 0 12px 30px rgba(0, 0, 0, 0.5);
  z-index: 5;
}
.board-card.dragging { cursor: grabbing; opacity: 0.92; }

.card-media {
  position: relative;
  border-radius: 9px 9px 0 0;
  overflow: hidden;
  background: rgba(0, 0, 0, 0.4);
}
.card-still { width: 100%; height: 100%; object-fit: contain; pointer-events: none; }
.card-media > i { position: absolute; inset: 0; margin: auto; width: 20px; height: 20px; color: rgba(255, 255, 255, 0.2); }

/* 非当前页的只读框线：仅提示该页已标注的范围，不拦截指针、不抢当前页注意力 */
.card-marks {
  position: absolute;
  inset: 0;
  z-index: 2;
  pointer-events: none;
}
.card-mark {
  position: absolute;
  box-sizing: border-box;
  border: 1px solid rgba(0, 255, 255, 0.5);
  border-radius: 2px;
}
/* 颜色沿用标注状态语义，整体压低亮度 */
.card-mark--empty { border-color: rgba(255, 107, 122, 0.55); }
.card-mark--pending { border-color: rgba(150, 165, 185, 0.5); }
.card-mark--ai_draft { border-color: rgba(230, 170, 60, 0.5); }
.card-mark--confirmed { border-color: rgba(55, 224, 141, 0.5); }

/* 区域名标签：贴在框线上沿；子级用淡底白字区分层级；
   框贴顶时改放框内（card-mark--flush），不然会被卡片裁掉 */
.card-mark-tag {
  position: absolute;
  top: -14px;
  left: -1px;
  max-width: 130px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  padding: 0 4px;
  height: 13px;
  line-height: 13px;
  border-radius: 3px;
  font-size: 9px;
  font-style: normal;
  font-weight: 700;
  color: #0a0a0f;
  background: rgba(0, 255, 255, 0.72);
}
.card-mark--child .card-mark-tag {
  color: rgba(255, 255, 255, 0.92);
  background: rgba(0, 255, 255, 0.3);
}
.card-mark--flush .card-mark-tag { top: 1px; }

/* 当前页内嵌标注画布：铺满图面，压缩底部留白并隐藏模式提示条 */
.card-media .card-canvas {
  position: absolute;
  inset: 0;
  z-index: 1;
  border: none;
  border-radius: 0;
  background: transparent;
}
.card-media :deep(.canvas-scroll) { inset: 8px; }
.card-media :deep(.mode-banner) { display: none; }

.card-no {
  position: absolute;
  top: 5px;
  left: 5px;
  z-index: 3;
  pointer-events: none;
  padding: 1px 6px;
  border-radius: 5px;
  font-size: 10px;
  font-weight: 700;
  color: #00ffff;
  background: rgba(0, 0, 0, 0.62);
}
.card-entry {
  flex-shrink: 0;
  padding: 1px 7px;
  border-radius: 5px;
  font-size: 9px;
  font-weight: 700;
  color: #37e08d;
  background: rgba(55, 224, 141, 0.12);
  border: 1px solid rgba(55, 224, 141, 0.35);
}
.card-mask {
  position: absolute;
  inset: 0;
  z-index: 4;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 11px;
  color: #00ffff;
  background: rgba(0, 0, 0, 0.6);
}
.card-mask--err { color: #ff6b7a; }

.card-foot {
  display: flex;
  align-items: center;
  gap: 8px;
  height: 34px;
  padding: 0 7px;
  border-top: 1px solid rgba(255, 255, 255, 0.06);
  box-sizing: border-box;
}
/* 当前页：图面留给框选/编辑，卡片拖动收到底部标题栏 */
.card-foot.is-handle { cursor: grab; }
.card-foot.is-handle:active { cursor: grabbing; }
/* 底部栏缩放：与画布内缩放条同一套状态，由卡片底部开关驱动 */
.card-zoom {
  margin-left: auto;
  display: flex;
  align-items: center;
  gap: 2px;
}
.cz-btn {
  width: 22px;
  height: 22px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 6px;
  border: 1px solid rgba(0, 255, 255, 0.22);
  background: rgba(0, 255, 255, 0.06);
  color: rgba(0, 255, 255, 0.85);
  font-size: 9px;
  cursor: pointer;
}
.cz-btn:hover:not(:disabled) { background: rgba(0, 255, 255, 0.16); }
.cz-btn:disabled { opacity: 0.3; cursor: not-allowed; }
.cz-pct {
  min-width: 34px;
  text-align: center;
  font-size: 10px;
  color: rgba(255, 255, 255, 0.6);
  cursor: pointer;
}
.cz-pct:hover { color: #00ffff; }

.board-empty {
  position: absolute;
  inset: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  color: rgba(255, 255, 255, 0.3);
  pointer-events: none;
}
.board-empty i { font-size: 26px; color: rgba(0, 255, 255, 0.3); }
.board-empty p { margin: 0; font-size: 13px; font-weight: 600; color: rgba(255, 255, 255, 0.55); }
.board-empty span { font-size: 11px; }
</style>
