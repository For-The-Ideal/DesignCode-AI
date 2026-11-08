<template>
  <!-- ═══ 任务列表（展开时显示） ═══ -->
  <Transition name="card-fade">
    <div v-show="sidebarOpen" class="task-list-wrapper">
          <div class="task-list-header" v-if="showTitle">
            <span class="task-count">共 {{ tasks.length }} 个任务</span>
          </div>

          <div
            v-for="task in tasks"
            :key="task.id"
            class="task-group"
          >
            <!-- ── 紧凑条目（始终可见） ── -->
            <div
              class="task-item"
              :class="{ expanded: expandedId === task.id }"
              @click="onItemClick(task)"
            >
              <div class="task-item-top">
                <div class="task-item-title">{{ taskName(task) }}</div>
                <span class="task-status-badge" :class="'status--' + task.status">
                  <span class="status-dot"></span>
                  {{ statusLabel(task.status) }}
                </span>
              </div>
              <div v-if="task.status !== 'failed'" class="task-progress">
                <div class="task-progress-bar">
                  <div
                    class="task-progress-fill"
                    :class="{ animated: task.status === 'running' }"
                    :style="{ width: task.progress + '%' }"
                  ></div>
                </div>
                <span class="task-progress-num">{{ task.progress }}%</span>
              </div>
              <div v-if="task.status === 'failed'" class="task-error">
                <i class="fas fa-exclamation-circle"></i>
                {{ task.error || '生成失败' }}
              </div>
              <div class="task-item-meta">
                <span>平台：{{ configStore.platformLabel(task.platform) }}</span>
                <span>创建时间：{{ task.createdAt || task.time }}</span>
              </div>
            </div>

            <!-- ── 展开详情卡片 ── -->
            <Transition name="detail-expand">
              <div v-if="expandedId === task.id" class="task-detail">
                <!-- 头部：缩略图 + 标题/标签 + 时间 -->
                <div class="card-head">
                  <div
                    class="card-thumb"
                    :class="{ 'is-clickable': task.images?.[0]?.url }"
                    @click.stop="openViewer(task)"
                  >
                    <!-- 缩略图：点击打开大图（自建查看器，可叠加框选位置） -->
                    <img v-if="task.images?.[0]?.url" class="thumb-img" :src="task.images[0].url" alt="" />
                    <i v-else class="fa-solid fa-image"></i>
                  </div>
                  <div class="card-head-info">
                    <h3 class="card-title">{{ taskName(task) }}</h3>
                    <div class="card-tags">
                      <span v-for="(tag, i) in taskTags(task)" :key="i" class="card-tag">{{ tag }}</span>
                    </div>
                  </div>
                  <span class="card-time">{{ task.createdAt || task.time }}</span>
                </div>

                <!-- 基础配置：框架/平台/组件库/选项/页数节点数（列表项已有字段） -->
                <section class="detail-sec">
                  <div class="sec-title">
                    <i class="fas fa-sliders"></i> 基础配置
                  </div>
                  <div class="cfg-list">
                    <div class="cfg-row">
                      <span class="cfg-k">框架</span>
                      <span class="cfg-v">{{ task.framework ? configStore.frameworkLabel(task.framework) : '—' }}</span>
                    </div>
                    <div class="cfg-row">
                      <span class="cfg-k">平台</span>
                      <span class="cfg-v">{{ task.platform ? configStore.platformLabel(task.platform) : '—' }}</span>
                    </div>
                    <div class="cfg-row">
                      <span class="cfg-k">组件库</span>
                      <span class="cfg-v">{{ task.component_lib || '—' }}</span>
                    </div>
                    <div class="cfg-row">
                      <span class="cfg-k">生成选项</span>
                      <span class="cfg-v">{{ optionText(task) }}</span>
                    </div>
                    <div class="cfg-row" v-if="(task.advanced || []).length">
                      <span class="cfg-k">高级选项</span>
                      <span class="cfg-v">{{ advancedText(task) }}</span>
                    </div>
                    <div class="cfg-row">
                      <span class="cfg-k">规模</span>
                      <span class="cfg-v">{{ (task.images || []).length }} 页 · {{ treeNodeCount(task) }} 节点</span>
                    </div>
                  </div>
                </section>

                <!-- 结构树：逐层全部展开（tree 只在详情接口返回，展开时懒加载） -->
                <section class="detail-sec">
                  <div class="sec-title">
                    <i class="fas fa-sitemap"></i> 结构树
                    <span class="sec-extra" v-if="treePages(task).length">{{ treeNodeCount(task) }} 个节点</span>
                  </div>
                  <div v-if="!detailOf(task)" class="sec-tip">
                    <i class="fas fa-spinner fa-spin"></i> 加载结构数据…
                  </div>
                  <div v-else-if="!treePages(task).length" class="sec-tip">未使用结构树（按视觉解析生成）</div>
                  <div v-else class="tree-wrap">
                    <div v-for="p in treePages(task)" :key="p.page_id || p.sort_order" class="tree-page">
                      <div v-if="(task.images || []).length > 1" class="tree-page-name">{{ p.page_name || p.page_id }}</div>
                      <div class="tree-nodes">
                        <div
                          v-for="n in flatNodes(p.tree)"
                          :key="n.key"
                          class="tree-node"
                          :class="'tree-lv' + Math.min(n.level, 3)"
                          :style="{ paddingLeft: (6 + n.level * 12) + 'px' }"
                        >
                          <span class="tree-dot"></span>
                          <span class="tree-name" :title="n.raw || n.name">{{ n.name }}</span>
                          <span v-if="n.component" class="tree-comp">{{ n.component }}</span>
                        </div>
                      </div>
                    </div>
                  </div>
                </section>

                <!-- 时间轴步骤 -->
                <div v-if="task.steps && task.steps.length" class="timeline">
                  <div class="timeline-line"></div>
                  <div
                    v-for="(step, idx) in task.steps"
                    :key="idx"
                    class="timeline-step"
                    :class="stepPhase(task, idx)"
                  >
                    <div class="step-left">
                      <div class="step-indicator" :class="{ done: step.completed }">
                        <i :class="step.icon" class="step-icon"></i>
                      </div>
                      <div class="step-body">
                        <div class="step-title">{{ step.title }}</div>
                        <div class="step-time">{{ step.time }}</div>
                      </div>
                    </div>
                    <div class="step-dot" :class="{ done: step.completed }">
                      <i v-if="step.completed" class="fas fa-check"></i>
                    </div>
                  </div>
                </div>

                <!-- 查看详情按钮 -->
                <button class="view-detail-btn" @click.stop="goToDetail(task.id)" v-if="showDetail">
                  <i class="fas fa-external-link-alt"></i> 查看完整详情
                </button>
              </div>
            </Transition>
          </div>

          <!-- 空态 -->
          <div v-if="tasks.length === 0" class="task-empty">
            <i class="fas fa-inbox"></i>
            <p>暂无进行中的任务</p>
          </div>
    </div>
  </Transition>

  <!-- ═══ 放大预览：原图 + 框选位置（多页可翻页，Esc/点击空白关闭） ═══ -->
  <Teleport to="body">
    <Transition name="viewer-fade">
      <div v-if="viewer" class="img-viewer" @click.self="closeViewer">
        <button class="viewer-close" type="button" title="关闭" @click="closeViewer">
          <i class="fas fa-xmark"></i>
        </button>
        <button
          v-if="viewerPages.length > 1"
          class="viewer-nav viewer-nav--prev"
          type="button"
          title="上一页"
          @click.stop="stepViewer(-1)"
        >
          <i class="fas fa-chevron-left"></i>
        </button>
        <button
          v-if="viewerPages.length > 1"
          class="viewer-nav viewer-nav--next"
          type="button"
          title="下一页"
          @click.stop="stepViewer(1)"
        >
          <i class="fas fa-chevron-right"></i>
        </button>

        <div v-if="viewerPage" class="viewer-stage" :style="viewerStageStyle" @click.stop>
          <img class="viewer-img" :src="viewerPage.url" alt="" @load="onViewerLoad" />
          <span
            v-for="m in viewerMarks"
            :key="m.key"
            class="viewer-mark"
            :class="['viewer-mark--lv' + m.level, { 'viewer-mark--flush': m.flush }]"
            :style="m.style"
          >
            <i class="viewer-mark-tag">{{ m.text }}</i>
          </span>
        </div>

        <div class="viewer-bar">
          <span v-if="viewerPages.length > 1" class="viewer-idx">
            {{ (viewer?.index || 0) + 1 }} / {{ viewerPages.length }}
          </span>
          <span v-if="viewerPage?.title" class="viewer-name">{{ viewerPage.title }}</span>
          <span v-if="viewerMarks.length" class="viewer-count">{{ viewerMarks.length }} 个框选区域</span>
          <span v-else-if="viewerPage" class="viewer-count">该图未使用结构树，无框选位置</span>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useConfigStore } from '~/stores/config'
import { taskApi } from '~/api/task'
import { elementLabel, regionLabel } from '~/utils/uiTree.js'

const configStore = useConfigStore()
onMounted(() => configStore.load())

const props = defineProps({
  tasks: { type: Array, required: true },
  expandedId: { type: [String, null], default: null },
  sidebarOpen: { type: Boolean, default: true },
  collapsible: { type: Boolean, default: false },
  showDetail: { type: Boolean, default: true },
  showTitle: { type: Boolean, default: true },
})

const emit = defineEmits(['update:expandedId', 'navigate'])

// ═══ 展开详情：结构树 / 页面信息只在详情接口返回（列表项为瘦身数据） ═══
// 按任务懒加载并缓存，避免列表接口因回传整棵树而膨胀
const detailCache = ref({})
const detailLoading = ref({})

async function ensureDetail(taskId, retry = 1) {
  if (!taskId || detailCache.value[taskId] || detailLoading.value[taskId]) return
  detailLoading.value[taskId] = true
  try {
    // cancelDuplicate=false：详情页也会请求同一任务，关闭同键去重避免互相 abort
    const res = await taskApi.getTaskById(taskId, true, false)
    if (res?.code === 200 && res.data) detailCache.value[taskId] = res.data
  } catch {
    // 静默：失败不影响卡片其余信息展示（重新展开可重试）
  } finally {
    detailLoading.value[taskId] = false
    // 没拿到数据（并发取消/瞬时失败）且卡片或大图仍在使用时补拉一次
    if (!detailCache.value[taskId] && retry > 0 && (props.expandedId === taskId || viewer.value?.taskId === taskId)) {
      setTimeout(() => ensureDetail(taskId, retry - 1), 1200)
    }
  }
}

// onMounted 才发请求（SSR 不执行）；切换展开项时按需拉取
onMounted(() => { if (props.expandedId) ensureDetail(props.expandedId) })
watch(() => props.expandedId, (id) => { if (id) ensureDetail(id) })

const detailOf = (task) => detailCache.value[task.id] || null

// 有结构树的页（tree 为空 = 该页按视觉解析生成）
const treePages = (task) => (detailOf(task)?.images || []).filter(i => Array.isArray(i.tree) && i.tree.length)

// 节点总数 + 逐层拍平（所有层级全部展示，不折叠）
// 界面全部中文：顶层是区域枚举（header/main…）走 regionLabel，
// 子级是元素类型（text/upload…）走 elementLabel；英文原值只放 title 悬停查看
const countNodes = (nodes) => (nodes || []).reduce((s, n) => s + 1 + countNodes(n.children), 0)
const treeNodeCount = (task) => treePages(task).reduce((s, p) => s + countNodes(p.tree), 0)
const flatNodes = (nodes, level = 0, out = []) => {
  (nodes || []).forEach((n) => {
    const isTop = level === 0
    const typeLabel = isTop ? regionLabel(n.name) : elementLabel(n.component)
    const name = isTop ? (typeLabel || '未命名区域') : (n.name || typeLabel || '未命名')
    // 子级右侧胶囊：类型与区域名不同才展示，避免「上传 上传」重复
    const component = !isTop && typeLabel && typeLabel !== name ? typeLabel : ''
    out.push({ key: n.nodeKey || `n${level}-${out.length}`, level, name, component, raw: n.name || '' })
    flatNodes(n.children, level + 1, out)
  })
  return out
}

// ═══ 框选位置（bbox 为原图像素，递归所有层级） ═══
// 框线标签同样中文化：顶层显示区域中文名，子级显示「元素类型 · 区域名」
const nodeText = (n, level = 0) => {
  if (!n) return ''
  if (level === 0) return regionLabel(n.name) || '未命名区域'
  const typeLabel = n.component ? elementLabel(n.component) : ''
  if (typeLabel && n.name && typeLabel !== n.name) return `${typeLabel} · ${n.name}`
  return n.name || typeLabel || ''
}
const bboxNodes = (nodes, level = 0, out = []) => {
  (nodes || []).forEach((n) => {
    const b = n.bbox || {}
    if (b.w > 0 && b.h > 0) {
      out.push({ x: b.x, y: b.y, w: b.w, h: b.h, level: Math.min(level, 3), label: nodeText(n, level) })
    }
    bboxNodes(n.children, level + 1, out)
  })
  return out
}

// ═══ 放大预览：原图 + 框选位置（缩略图只展示原图，框选位置仅在放大后叠加） ═══
// el-image 的预览层在 body 上渲染，无法叠框，改用自建查看器（支持多页翻页/键盘）
const viewer = ref(null)                  // { taskId, index }
const viewerDims = ref({ w: 0, h: 0 })    // 老任务没存原图尺寸时用图片自身尺寸兜底
const vp = ref({ w: 0, h: 0 })

const viewerPages = computed(() => {
  const t = props.tasks.find((x) => x.id === viewer.value?.taskId)
  if (!t) return []
  const details = detailOf(t)?.images || []
  return (t.images || []).map((im, i) => {
    const d = details[i] || {}
    return { url: im.url, title: im.title || '', W: d.width || 0, H: d.height || 0, tree: d.tree || null }
  })
})
const viewerPage = computed(() => {
  const pages = viewerPages.value
  return pages[viewer.value?.index] || pages[0] || null
})
// 舞台按原图比例铺进视口：框线用百分比，舞台缩放时自动跟随
const viewerStageStyle = computed(() => {
  const p = viewerPage.value
  const W = p?.W || viewerDims.value.w || 0
  const H = p?.H || viewerDims.value.h || 0
  const { w: vw, h: vh } = vp.value
  if (!W || !H || !vw || !vh) return {}
  const s = Math.min((vw * 0.9) / W, (vh * 0.86) / H)
  return { width: Math.round(W * s) + 'px', height: Math.round(H * s) + 'px' }
})
const viewerMarks = computed(() => {
  const p = viewerPage.value
  const W = p?.W || viewerDims.value.w || 0
  const H = p?.H || viewerDims.value.h || 0
  if (!p || !W || !H || !Array.isArray(p.tree)) return []
  return bboxNodes(p.tree).map((b, i) => ({
    key: i,
    text: b.label,
    level: b.level,
    // 贴顶的框把标签放进框内，避免被视口裁掉
    flush: b.y / H < 0.03,
    style: {
      left: (b.x / W) * 100 + '%',
      top: (b.y / H) * 100 + '%',
      width: (b.w / W) * 100 + '%',
      height: (b.h / H) * 100 + '%',
    },
  }))
})

const onViewerResize = () => { vp.value = { w: window.innerWidth, h: window.innerHeight } }
const onViewerKey = (e) => {
  if (e.key === 'Escape') closeViewer()
  else if (e.key === 'ArrowLeft') stepViewer(-1)
  else if (e.key === 'ArrowRight') stepViewer(1)
}
function openViewer(task) {
  if (!task?.images?.[0]?.url) return
  viewer.value = { taskId: task.id, index: 0 }
  viewerDims.value = { w: 0, h: 0 }
  onViewerResize()
  window.addEventListener('resize', onViewerResize)
  window.addEventListener('keydown', onViewerKey)
  document.body.style.overflow = 'hidden'
  // 结构树只在详情接口返回：没缓存就补拉，拿到后框线自动出现
  ensureDetail(task.id)
}
function closeViewer() {
  viewer.value = null
  window.removeEventListener('resize', onViewerResize)
  window.removeEventListener('keydown', onViewerKey)
  document.body.style.overflow = ''
}
function stepViewer(d) {
  const v = viewer.value
  const n = viewerPages.value.length
  if (!v || n < 2) return
  v.index = (v.index + d + n) % n
  viewerDims.value = { w: 0, h: 0 }
}
const onViewerLoad = (e) => {
  const p = viewerPage.value
  if (p && (!p.W || !p.H)) {
    viewerDims.value = { w: e?.target?.naturalWidth || 0, h: e?.target?.naturalHeight || 0 }
  }
}
// 组件卸载时清掉全局监听与 body 滚动锁
onBeforeUnmount(() => { if (viewer.value) closeViewer() })

// 条目点击：展开该项并通知切换详情（失败/排队中也可进入，详情页内展示对应状态）
const onItemClick = (task) => {
  emit('update:expandedId', task.id)
  emit('navigate', task.id)
}

const goToDetail = (id) => {
  emit('navigate', id)
}

// 定位当前进度：返回 'phase-past' | 'phase-current' | 'phase-future'
const stepPhase = (task, idx) => {
  const firstPending = task.steps.findIndex(s => !s.completed)
  const currentIdx = firstPending === -1 ? task.steps.length - 1 : firstPending
  if (idx < currentIdx) return 'phase-past'
  if (idx === currentIdx) return 'phase-current'
  return 'phase-future'
}

// 任务标签列表：framework、platform、options 均回显字典 label
const taskTags = (task) => {
  const tags = []
  if (task.framework) tags.push(configStore.frameworkLabel(task.framework))
  if (task.platform)  tags.push(configStore.platformLabel(task.platform))
  ;(task.options || []).forEach(v => { if (v !== 'component') tags.push(configStore.optionLabel(v)) })
  if (task.component_lib) tags.push(task.component_lib)
  return tags
}

// 基础配置文本：选项字典回显（component 已单列组件库，不重复展示）
const optionText = (task) => (task.options || [])
  .filter(v => v !== 'component')
  .map(v => configStore.optionLabel(v))
  .filter(Boolean)
  .join('、') || '—'
const advancedText = (task) => (task.advanced || [])
  .map(v => configStore.advancedLabel(v))
  .filter(Boolean)
  .join('、') || '—'

// 模板保持简洁：任务名 / 状态 label 在此封装为短方法（数据源仍是配置字典）
const taskName    = (task) => configStore.taskName(task)
const statusLabel = (s) => configStore.statusLabel(s)
</script>

<style scoped>
/* ═══ 任务列表 ═══ */
.task-list-wrapper {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.task-list-header {
  display: flex;
  align-items: center;
  padding: 0 2px;
}
.task-count {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.35);
}

/* ── 任务分组（紧凑条目 + 展开详情） ── */
.task-group {
  display: flex;
  flex-direction: column;
}

/* ── 紧凑条目（可点击） ── */
.task-item {
  background: rgba(0, 255, 255, 0.02);
  border-radius: 12px;
  padding: 10px 14px;
  border: 1px solid rgba(0, 255, 255, 0.06);
  cursor: pointer;
  transition: all 0.2s;
}
.task-item:hover {
  border-color: rgba(0, 255, 255, 0.18);
  background: rgba(0, 255, 255, 0.05);
}
.task-item.expanded {
  border-color: rgba(0, 255, 255, 0.2);
  background: rgba(0, 255, 255, 0.06);
  border-radius: 12px 12px 0 0;
  border-bottom-color: transparent;
}

.task-item-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 6px;
}
.task-status-badge {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  padding: 2px 10px;
  border-radius: 20px;
  font-size: 11px;
  font-weight: 600;
}
.task-status-badge .status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
}

/* 状态颜色 */
.status--pending {
  background: rgba(59, 130, 246, 0.1);
  color: #60a5fa;
  border: 1px solid rgba(59, 130, 246, 0.2);
}
.status--pending .status-dot { background: #60a5fa; animation: pulse 2s ease-in-out infinite; }
.status--running {
  background: rgba(250, 204, 21, 0.1);
  color: #facc15;
  border: 1px solid rgba(250, 204, 21, 0.2);
}
.status--running .status-dot { background: #facc15; animation: pulse 1s ease-in-out infinite; }
.status--success {
  background: rgba(52, 211, 153, 0.1);
  color: #34d399;
  border: 1px solid rgba(52, 211, 153, 0.2);
}
.status--success .status-dot { background: #34d399; }
.status--failed {
  background: rgba(248, 113, 113, 0.1);
  color: #f87171;
  border: 1px solid rgba(248, 113, 113, 0.2);
}
.status--failed .status-dot { background: #f87171; }

.task-framework {
  font-size: 10px;
  color: rgba(255, 255, 255, 0.3);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}
.task-item-title {
  font-size: 13px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.75);
  flex: 1;
  min-width: 0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* 紧凑进度条 */
.task-progress {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 6px;
}
.task-progress-bar {
  flex: 1;
  height: 5px;
  background: rgba(255, 255, 255, 0.06);
  border-radius: 3px;
  overflow: hidden;
}
.task-progress-fill {
  height: 100%;
  border-radius: 3px;
  background: linear-gradient(90deg, #3b82f6, #22d3ee);
  transition: width 0.5s ease-out;
}
.task-progress-fill.animated {
  background: linear-gradient(90deg, #3b82f6, #22d3ee, #3b82f6);
  background-size: 200% 100%;
  animation: progressPulse 1.5s ease-in-out infinite;
}
@keyframes progressPulse {
  0%, 100% { background-position: 0% 50%; }
  50% { background-position: 100% 50%; }
}
.task-progress-num {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.4);
  font-weight: 600;
  flex-shrink: 0;
}
.task-error {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 11px;
  color: #f87171;
  margin-bottom: 6px;
}
.task-error i { font-size: 12px; flex-shrink: 0; }
.task-item-meta {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: rgba(255, 255, 255, 0.55);
}

/* ── 展开详情卡片 ── */
.task-detail {
  background: rgba(0, 255, 255, 0.03);
  border-radius: 0 0 14px 14px;
  padding: 16px 14px 14px;
  border: 1px solid rgba(0, 255, 255, 0.1);
  border-top: none;
  position: relative;
  overflow: hidden;
}

/* 卡片头部：缩略图 + 信息 + 时间 */
.card-head {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  margin-bottom: 14px;
}
/* 手机比例缩略图：比列表卡片更大，点开放大预览 */
.card-thumb {
  width: 76px;
  aspect-ratio: 9 / 16;
  border-radius: 10px;
  overflow: hidden;
  flex-shrink: 0;
  background: linear-gradient(135deg, rgba(59,130,246,0.12), rgba(34,211,238,0.08));
  display: flex;
  align-items: center;
  justify-content: center;
}
.card-thumb.is-clickable {
  cursor: zoom-in;
}
/* 缩略图：容器 overflow:hidden 裁切，cover 填满手机比例 */
.thumb-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}
.card-thumb i {
  font-size: 22px;
  color: rgba(148, 163, 184, 0.4);
}
.card-head-info {
  flex: 1;
  min-width: 0;
}
.card-title {
  font-size: 15px;
  font-weight: 700;
  color: #f1f5f9;
  line-height: 1.25;
  margin-bottom: 8px;
}
.card-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}
.card-tag {
  display: inline-block;
  font-size: 10px;
  color: rgba(255, 255, 255, 0.8);
  background: rgba(148, 163, 184, 0.08);
  padding: 2px 10px;
  border-radius: 20px;
  border: 1px solid rgba(148, 163, 184, 0.12);
}
.card-time {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.3);
  font-family: 'Fira Code', monospace;
  flex-shrink: 0;
  padding-top: 2px;
}

/* ── 详情分区：基础配置 / 结构树 ── */
.detail-sec { margin-bottom: 14px; }
.sec-title {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 8px;
  font-size: 12px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.55);
}
.sec-title i { font-size: 11px; color: rgba(34, 211, 238, 0.65); }
.sec-extra {
  margin-left: auto;
  font-size: 10px;
  font-weight: 400;
  color: rgba(255, 255, 255, 0.25);
}
.sec-tip {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 11px;
  color: rgba(255, 255, 255, 0.28);
}

/* 基础配置：键值行，值可换行展示全量 */
.cfg-list { display: flex; flex-direction: column; }
.cfg-row {
  display: flex;
  gap: 8px;
  padding: 4px 0;
  font-size: 11px;
  border-bottom: 1px dashed rgba(255, 255, 255, 0.05);
}
.cfg-row:last-child { border-bottom: none; }
.cfg-k { width: 58px; flex-shrink: 0; color: rgba(255, 255, 255, 0.3); }
.cfg-v { flex: 1; min-width: 0; color: rgba(255, 255, 255, 0.75); word-break: break-word; }

/* 结构树：全部层级逐行铺开，层级越深越淡 */
.tree-wrap { display: flex; flex-direction: column; gap: 8px; }
.tree-page-name {
  margin-bottom: 3px;
  font-size: 10px;
  font-family: 'Fira Code', monospace;
  color: rgba(34, 211, 238, 0.7);
}
.tree-nodes { display: flex; flex-direction: column; gap: 1px; }
.tree-node {
  display: flex;
  align-items: center;
  gap: 6px;
  padding-top: 3px;
  padding-right: 4px;
  padding-bottom: 3px;
  border-radius: 6px;
}
.tree-node:hover { background: rgba(34, 211, 238, 0.06); }
.tree-dot {
  width: 4px; height: 4px;
  border-radius: 50%;
  flex-shrink: 0;
  background: rgba(34, 211, 238, 0.55);
}
.tree-name {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.72);
  white-space: nowrap; overflow: hidden; text-overflow: ellipsis;
}
.tree-comp {
  margin-left: auto;
  flex-shrink: 0;
  padding: 1px 6px;
  border-radius: 4px;
  font-size: 10px;
  color: rgba(129, 140, 248, 0.75);
  background: rgba(129, 140, 248, 0.1);
}
/* 层级：顶层加粗提亮，越深越淡 */
.tree-lv0 .tree-name { color: #e2e8f0; font-weight: 600; }
.tree-lv0 .tree-dot { background: #22d3ee; box-shadow: 0 0 6px rgba(34, 211, 238, 0.5); }
.tree-lv1 .tree-name { color: rgba(255, 255, 255, 0.66); }
.tree-lv2 .tree-name { color: rgba(255, 255, 255, 0.5); }
.tree-lv3 .tree-name { color: rgba(255, 255, 255, 0.38); }

/* 时间轴 */
.timeline {
  position: relative;
  padding-left: 8px;
}
.timeline-line {
  position: absolute;
  left: 19px;
  top: 6px;
  bottom: 6px;
  width: 1px;
  background: rgba(255, 255, 255, 0.08);
}
.timeline-step {
  display: flex;
  align-items: center;
  gap: 10px;
  justify-content: space-between;
  margin-bottom: 8px;
  position: relative;
  z-index: 1;
}
.timeline-step:last-child { margin-bottom: 0; }

/* 左侧分组（indicator + body） */
.step-left {
  display: flex;
  align-items: center;
  gap: 10px;
  flex: 1;
  min-width: 0;
}

/* 步骤指示器（仅图标） */
.step-indicator {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  width: 22px;
  height: 22px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.08);
  transition: all 0.3s;
}
.step-icon {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.2);
  transition: all 0.3s;
}
.step-indicator.done {
  background: rgba(52, 211, 153, 0.1);
  border-color: rgba(52, 211, 153, 0.3);
}
.step-indicator.done .step-icon {
  color: rgba(52, 211, 153, 0.5);
}

.step-dot {
  width: 22px; height: 22px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  background: rgba(10, 10, 15, 0.95);
  border: 2px solid rgba(255, 255, 255, 0.1);
  font-size: 10px;
  color: transparent;
  transition: all 0.3s;
}
.step-dot.done {
  border-color: rgba(52, 211, 153, 0.5);
  background: rgba(52, 211, 153, 0.1);
  color: #34d399;
}
.step-body {
  min-width: 0;
  padding-top: 2px;
}
.step-title {
  font-size: 12px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.6);
  transition: color 0.3s;
}
.step-time {
  font-size: 10px;
  color: rgba(255, 255, 255, 0.25);
  font-family: 'Fira Code', monospace;
  transition: color 0.3s;
}

/* 阶段色：已完成 → 白色 */
.phase-past .step-title { color: rgba(255, 255, 255, 0.9); }
.phase-past .step-time  { color: rgba(255, 255, 255, 0.4); }

/* 阶段色：当前进行中 → 蓝色高亮 */
.phase-current .step-title { color: #60a5fa; }
.phase-current .step-time  { color: rgba(96, 165, 250, 0.6); }

/* 阶段色：未开始 → 保持默认灰 */
/* .phase-future 使用默认色 */

/* ═══ 卡片淡入过渡 ═══ */
.card-fade-enter-active {
  transition: opacity 0.35s ease, transform 0.35s ease;
}
.card-fade-leave-active {
  transition: opacity 0.15s ease;
}
.card-fade-enter-from {
  opacity: 0;
  transform: translateY(8px);
}
.card-fade-leave-to { opacity: 0; }

/* ═══ 详情展开过渡 ═══ */
.detail-expand-enter-active {
  transition: all 0.3s ease-out;
}
.detail-expand-leave-active {
  transition: all 0.2s ease-in;
}
.detail-expand-enter-from {
  opacity: 0;
  max-height: 0;
  transform: translateY(-8px);
}
.detail-expand-enter-to {
  opacity: 1;
  max-height: 500px;
  transform: translateY(0);
}
.detail-expand-leave-from {
  opacity: 1;
  max-height: 500px;
}
.detail-expand-leave-to {
  opacity: 0;
  max-height: 0;
}

/* 空态 */
.task-empty {
  text-align: center;
  padding: 32px 16px;
  color: rgba(255, 255, 255, 0.15);
}
.task-empty i {
  font-size: 32px;
  margin-bottom: 8px;
  display: block;
}
.task-empty p { font-size: 12px; }

/* 查看详情按钮 */
.view-detail-btn {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  width: 100%;
  padding: 6px 12px;
  margin-top: 12px;
  border: 1px solid rgba(0,255,255,0.12);
  border-radius: 8px;
  background: rgba(0,255,255,0.04);
  color: rgba(0,255,255,0.55);
  font-size: 11px;
  cursor: pointer;
  transition: all 0.2s;
  font-family: inherit;
  justify-content: center;
}
.view-detail-btn:hover {
  background: rgba(0,255,255,0.08);
  border-color: rgba(0,255,255,0.25);
  color: #00d8ff;
}

/* ═══ 放大预览（原图 + 框选位置） ═══ */
.viewer-fade-enter-active,
.viewer-fade-leave-active { transition: opacity 0.18s ease; }
.viewer-fade-enter-from,
.viewer-fade-leave-to { opacity: 0; }
.img-viewer {
  position: fixed;
  inset: 0;
  z-index: 3000;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 28px;
  background: rgba(5, 8, 12, 0.88);
  backdrop-filter: blur(3px);
}
.viewer-stage {
  position: relative;
  border-radius: 10px;
  box-shadow: 0 24px 64px rgba(0, 0, 0, 0.6);
}
.viewer-img {
  display: block;
  width: 100%;
  height: auto;
  max-height: 86vh;
  border-radius: 10px;
}
.viewer-close {
  position: absolute;
  top: 18px;
  right: 20px;
  width: 34px;
  height: 34px;
  border-radius: 50%;
  border: 1px solid rgba(255, 255, 255, 0.12);
  background: rgba(15, 22, 32, 0.7);
  color: rgba(255, 255, 255, 0.6);
  font-size: 14px;
  cursor: pointer;
  transition: all 0.18s;
}
.viewer-close:hover { color: #e2e8f0; border-color: rgba(0, 255, 255, 0.3); }
.viewer-nav {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  width: 40px;
  height: 40px;
  border-radius: 50%;
  border: 1px solid rgba(255, 255, 255, 0.12);
  background: rgba(15, 22, 32, 0.7);
  color: rgba(255, 255, 255, 0.6);
  font-size: 14px;
  cursor: pointer;
  transition: all 0.18s;
}
.viewer-nav:hover { color: #22d3ee; border-color: rgba(0, 255, 255, 0.3); }
.viewer-nav--prev { left: 18px; }
.viewer-nav--next { right: 18px; }
.viewer-bar {
  position: absolute;
  bottom: 18px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  align-items: center;
  gap: 12px;
  max-width: 80vw;
  padding: 6px 14px;
  border-radius: 999px;
  background: rgba(15, 22, 32, 0.82);
  border: 1px solid rgba(255, 255, 255, 0.06);
  font-size: 12px;
  color: rgba(255, 255, 255, 0.55);
}
.viewer-idx { font-variant-numeric: tabular-nums; color: #22d3ee; font-weight: 600; }
.viewer-name {
  max-width: 260px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.viewer-count { font-size: 11px; color: rgba(255, 255, 255, 0.35); }
/* 框线：与详情页一致按层级配色，标签贴框上沿（贴顶时收进框内） */
.viewer-mark {
  position: absolute;
  box-sizing: border-box;
  border: 1.5px solid rgba(34, 211, 238, 0.85);
  background: rgba(34, 211, 238, 0.1);
  border-radius: 3px;
}
.viewer-mark--lv1 { border-color: rgba(129, 140, 248, 0.85); background: rgba(129, 140, 248, 0.1); }
.viewer-mark--lv2 { border-color: rgba(52, 211, 153, 0.8); background: rgba(52, 211, 153, 0.09); }
.viewer-mark--lv3 { border-color: rgba(251, 191, 36, 0.75); background: rgba(251, 191, 36, 0.08); }
.viewer-mark-tag {
  position: absolute;
  top: -16px;
  left: -1px;
  max-width: 160px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  padding: 0 5px;
  height: 15px;
  line-height: 15px;
  border-radius: 3px;
  font-size: 10px;
  font-style: normal;
  font-weight: 700;
  color: #06121a;
  background: rgba(34, 211, 238, 0.92);
}
.viewer-mark--lv1 .viewer-mark-tag { background: rgba(129, 140, 248, 0.92); color: #0a0a14; }
.viewer-mark--lv2 .viewer-mark-tag { background: rgba(52, 211, 153, 0.92); color: #04140d; }
.viewer-mark--lv3 .viewer-mark-tag { background: rgba(251, 191, 36, 0.92); color: #1a1204; }
.viewer-mark--flush .viewer-mark-tag { top: 1px; }
</style>
