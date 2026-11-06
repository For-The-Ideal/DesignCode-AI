<template>
  <Teleport to="body">
    <div v-if="visible" class="editor-mask">
      <div class="editor">
        <!-- ═══ 顶部工具条 ═══ -->
        <div class="ed-toolbar">
          <div class="tb-left">
            <i class="fa-solid fa-diagram-project"></i>
            <span class="tb-title">设计稿结构编辑</span>
          </div>
          <div class="tb-mid">
            <div class="mode-seg">
              <button class="seg" :class="{ active: mode === 'draw' }" @click="setMode('draw')">
                <i class="fas fa-draw-polygon"></i> 框选
              </button>
              <button class="seg" :class="{ active: mode === 'select' }" @click="setMode('select')">
                <i class="fas fa-mouse-pointer"></i> 编辑
              </button>
            </div>
            <span v-if="mode === 'child'" class="tb-hint">
              子级模式：为「{{ parentName }}」框选
              <button class="tb-exit" @click="exitMode">退出</button>
            </span>
          </div>
          <div class="tb-right">
            <span class="tb-count">节点 {{ store.nodeList.length }} · 待补 {{ store.unconfirmedNodes.length }}</span>
            <button class="tb-btn" :class="{ 'tb-btn--loading': aiLoading }" :disabled="!store.hasNodes || aiLoading" @click="aiDescribe">
              <i class="fas" :class="aiLoading ? 'fa-spinner fa-spin' : 'fa-magic'"></i>
              {{ aiLoading ? '生成中…' : 'AI 描述' }}
            </button>
            <button class="tb-btn tb-btn--primary" :disabled="!store.hasNodes" @click="onGenerate">确认标注</button>
            <button class="tb-btn" @click="close">关闭</button>
          </div>
        </div>

        <!-- ═══ 主体三栏 ═══ -->
        <div class="ed-body">
          <div class="ed-canvas">
            <AnnotationCanvas :mode="mode" @drawn="onDrawn" @cancelMode="exitMode" />
          </div>
          <div class="ed-side">
            <div class="side-tree">
              <HierarchyTree @addChild="onAddChild" @deleteNode="onDeleteNode" @select="onTreeSelect" />
            </div>
            <div class="side-insp">
              <NodeInspector />
            </div>
          </div>
        </div>
      </div>

      <!-- AI 整页描述生成中：全屏 loading，遮罩整个编辑器避免误操作 -->
      <transition name="ai-fade">
        <div v-if="aiLoading" class="ai-page-loading">
          <div class="ai-loading-box">
            <div class="ai-orbit">
              <span class="ai-orbit-ring"></span>
              <span class="ai-orbit-core"><i class="fas fa-magic"></i></span>
            </div>
            <p class="ai-loading-text">AI 正在生成整页区域描述…</p>
            <p class="ai-loading-sub">共 {{ store.nodeList.length }} 个区域，生成后自动回填</p>
          </div>
        </div>
      </transition>

      <!-- 顶层类型改由右侧 NodeInspector 内联胶囊选择，不再在图上弹层 -->
    </div>
  </Teleport>
</template>

<script setup>
import { ref, computed, watch, onBeforeUnmount } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import AnnotationCanvas from './AnnotationCanvas.client.vue'
import HierarchyTree from './HierarchyTree.client.vue'
import NodeInspector from './NodeInspector.client.vue'
import { useUiProjectStore } from '~/stores/uiProject'
import { useCodeStore } from '~/stores/code'
import { useRegionDescribe } from '~/composables/useRegionDescribe'

const store = useUiProjectStore()
const codeStore = useCodeStore()
// AI 描述是整页操作（一次生成所有区域描述并全部回填），放在顶部工具栏
const { aiLoading, aiDescribe } = useRegionDescribe()

// 「确认标注」通过后上抛已组装的树；由上层写回图片缩略标注并关闭编辑器。
// 真正的下单（创建任务/扣积分/SSE）由 ConfigPanel「开始生成」另行触发，编辑器不碰网络。
// close：编辑器关闭（工具条「关闭」或确认后 close）时上抛。
const emit = defineEmits(['confirm', 'close'])

const visible = ref(false)
const mode = ref('draw') // draw（默认，直接框选）| select（拖动/缩放已有框）| child
const pendingParentKey = ref(null)

const parentName = computed(() => {
  const p = store.nodes[pendingParentKey.value]
  return p ? (p.name || '未命名') : ''
})

function open() { visible.value = true }
function close() { visible.value = false; exitMode(); emit('close') }
defineExpose({ open, close })

// 编辑器全屏遮罩 z-index=3000，会盖住 Element 默认层级的提示/确认浮层；
// 打开时给 body 加类，用作用域限定的样式把 ElMessage / ElMessageBox 抬到其上（不污染全站）。
watch(visible, (v) => {
  if (typeof document !== 'undefined') document.body.classList.toggle('ui-editor-open', v)
})
onBeforeUnmount(() => {
  if (typeof document !== 'undefined') document.body.classList.remove('ui-editor-open')
})

function exitMode() {
  mode.value = 'draw'
  pendingParentKey.value = null
  store.setPendingParent(null)
}

// 切换顶层框选 / 编辑模式（退出子级）
function setMode(m) {
  if (mode.value === m) return
  mode.value = m
  pendingParentKey.value = null
  store.setPendingParent(null)
}

// ── 树里点 +：进入子框选模式 ──
function onAddChild(parentKey) {
  pendingParentKey.value = parentKey
  store.setPendingParent(parentKey)
  mode.value = 'child'
}

// key 是否为 ancestorKey 自身或其后代（沿 parentId 上溯）
function isDescendantOrSelf(key, ancestorKey) {
  let cur = key
  while (cur) {
    if (cur === ancestorKey) return true
    cur = store.nodes[cur]?.parentId || null
  }
  return false
}

// 树中手动选中：若处于子级模式且选中了当前父子树以外的节点（如另一个顶层），退出整个子级标注
function onTreeSelect(key) {
  if (mode.value === 'child' && pendingParentKey.value && !isDescendantOrSelf(key, pendingParentKey.value)) {
    exitMode()
  }
  store.selectNode(key)
}

function onDeleteNode(key) {
  const node = store.nodes[key]
  if (!node) return
  const kids = store.childrenOf(key)
  // 父节点下还有子区域时不允许删除，先删完子区域
  if (kids.length) {
    ElMessage.warning(`「${node.name || '未命名区域'}」下还有 ${kids.length} 个子区域，请先删除其子区域`)
    return
  }
  store.removeNode(key, true)
}

// ── 画布框完一个矩形 ──
function onDrawn(bbox) {
  if (mode.value === 'child') {
    createChild(bbox)
  } else {
    // draw 模式：只在顶层区域类型之间做重叠提醒（子级框在父级内属正常包含，不提示）
    const det = store.detectNewBox(bbox, { topOnly: true })
    applyConflicts(det, () => {
      // 先建占位顶层节点并选中，类型由右侧面板内联选择（不在图上弹层）
      store.addTopPlaceholder(bbox)
    })
  }
}

// ═══ 子级：直接创建并挂父（保持 child 模式，可连续标注多个子区域）═══
// 子级与父级是包含关系，不做重叠提醒（避免把父节点当冲突）
function createChild(bbox) {
  const res = store.addChildNode({ bbox, parentKey: pendingParentKey.value })
  if (res.suggestedParentKey && res.node.parentId === res.suggestedParentKey) {
    const p = store.nodes[res.suggestedParentKey]
    ElMessage.info(`已自动挂到「${p?.name || '未命名'}」下，可继续框选或点“退出”`)
  }
}

// 重叠处理：high 拦截给保留/取消，soft 仅提示
function applyConflicts(det, proceed) {
  const high = (det.conflicts || []).filter(c => c.level === 'high')
  if (high.length) {
    const names = high.map(c => c.name || '未命名').join('、')
    ElMessageBox.confirm(
      `新框与「${names}」高度重叠，是否仍要保留这一块？`,
      '重叠提醒',
      { confirmButtonText: '仍保留', cancelButtonText: '放弃该框', type: 'warning' }
    ).then(() => proceed()).catch(() => { exitMode() })
    return
  }
  const soft = (det.conflicts || []).filter(c => c.level === 'soft')
  if (soft.length) ElMessage.warning(`与「${soft.map(c => c.name || '未命名').join('、')}」可能重叠，请确认`)
  proceed()
}

// ═══ 确认标注：允许带待补保存，不提交 ═══
// 待补（未选类型 / AI 初稿 / 描述为空）不拦截：先存进度，进度卡会标黄；
// 真正的完整性校验在「开始生成」时执行（未选类型的区域不能进生成）。
function onGenerate() {
  // 子级必须有区域名：缺名会在层级结构里出现「未命名」内容，不符合数据结构
  const unnamed = store.nodeList.filter((n) => n.parentId && !n.name)
  if (unnamed.length) {
    // 按父级分组提示：xxx 一层缺少几个区域名
    const groups = new Map()
    unnamed.forEach((n) => {
      const pk = store.nodes[n.parentId]?.name || '未命名区域'
      groups.set(pk, (groups.get(pk) || 0) + 1)
    })
    const text = [...groups.entries()]
      .slice(0, 3)
      .map(([name, count]) => `「${name}」一层缺少 ${count} 个区域名`)
      .join('；')
    // 选中第一个缺名子级，右栏直接补名
    store.selectNode(unnamed[0].nodeKey)
    ElMessage.warning(`${text}${groups.size > 3 ? ' 等' : ''}，层级结构会出现「未命名」内容，请补齐后再确认`)
    return
  }
  const pending = store.unconfirmedNodes
  const noType = pending.filter((n) => !n.parentId && !n.name)
  const tree = store.buildSubmitTree()
  if (!pending.length) {
    offerCleanAssets(tree)
    return
  }
  const extra = noType.length ? `，其中 ${noType.length} 个还没选类型` : ''
  ElMessageBox.confirm(
    `还有 ${pending.length} 个区域描述未完善${extra}，可先保存稍后继续；开始生成前需要全部补齐。`,
    '确认标注',
    { confirmButtonText: '确认保存', cancelButtonText: '返回检查', type: 'warning' },
  )
    .then(() => offerCleanAssets(tree))
    .catch(() => {})
}

// 确认前提醒：当前主图名下没有被任何区域引用的附图可一键清空（清空后素材列表自动同步）；
// 点「继续」→ 清空未用附图并完成确认；点「取消」→ 什么都不动（不清理、不确认）。
async function offerCleanAssets(tree) {
  const scopeId = store.assetScopeId
  const used = new Set(store.nodeList.map((n) => n.assetId).filter(Boolean))
  const unused = scopeId
    ? codeStore.assets.filter((a) => a.ownerId === scopeId && !used.has(a.assetId))
    : []
  // 没有未使用附图：直接完成确认
  if (!unused.length) {
    emit('confirm', tree)
    return
  }
  try {
    await ElMessageBox.confirm(
      `当前主图有 ${unused.length} 张附图未被任何区域使用，继续将清空这些附图，取消则不做任何操作。`,
      '未使用的附图',
      { confirmButtonText: '继续', cancelButtonText: '取消', type: 'warning' },
    )
  } catch {
    // 取消：保持原样，不清理素材、不完成确认
    return
  }
  unused.forEach((a) => codeStore.removeAsset(a.assetId))
  emit('confirm', tree)
}
</script>

<style scoped>
.editor-mask {
  position: fixed;
  inset: 0;
  z-index: 3000;
  background: rgba(5, 6, 12, 0.9);
  backdrop-filter: blur(4px);
  display: flex;
}
.editor {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 0;
}
.ed-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding: 12px 20px;
  border-bottom: 1px solid rgba(0, 255, 255, 0.12);
}
.tb-left { display: flex; align-items: center; gap: 8px; color: #00ffff; }
.tb-title { font-size: 15px; font-weight: 600; color: rgba(255, 255, 255, 0.85); }
.tb-mid { display: flex; align-items: center; gap: 10px; }
.tb-hint { font-size: 12px; color: rgba(255, 207, 77, 0.9); display: flex; align-items: center; gap: 8px; }
.mode-seg { display: inline-flex; gap: 4px; padding: 3px; border-radius: 8px; background: rgba(255, 255, 255, 0.05); border: 1px solid rgba(255, 255, 255, 0.08); }
.seg {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 4px 12px;
  border-radius: 6px;
  border: none;
  background: transparent;
  color: rgba(255, 255, 255, 0.55);
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s;
}
.seg:hover { color: rgba(255, 255, 255, 0.85); }
.seg.active { background: rgba(0, 255, 255, 0.12); color: #00ffff; }
.tb-exit { background: none; border: none; color: #ff6b7a; cursor: pointer; font-size: 12px; }
.tb-right { display: flex; align-items: center; gap: 12px; }
.tb-count { font-size: 12px; color: rgba(255, 255, 255, 0.4); }
.tb-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 14px;
  border-radius: 8px;
  border: 1px solid rgba(0, 255, 255, 0.2);
  background: rgba(0, 0, 0, 0.3);
  color: rgba(255, 255, 255, 0.7);
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
}
.tb-btn:hover:not(:disabled) { border-color: rgba(0, 255, 255, 0.5); color: #00ffff; }
.tb-btn.active { border-color: #00ffff; color: #00ffff; background: rgba(0, 255, 255, 0.1); }
.tb-btn--primary { background: linear-gradient(135deg, #00ffff, #ff00ff); color: #fff; border: none; }
.tb-btn:disabled { opacity: 0.4; cursor: not-allowed; }

/* AI 描述 loading：图标旋转 + 扫光，禁用态不置灰以保留动画可见性 */
.tb-btn.tb-btn--loading:disabled {
  position: relative;
  overflow: hidden;
  opacity: 1;
  cursor: progress;
  border-color: rgba(0, 255, 255, 0.7);
  color: #00ffff;
}
.tb-btn--loading::after {
  content: '';
  position: absolute;
  inset: 0;
  background: linear-gradient(100deg, transparent 20%, rgba(0, 255, 255, 0.25) 50%, transparent 80%);
  transform: translateX(-100%);
  animation: tb-sweep 1.1s linear infinite;
  pointer-events: none;
}
@keyframes tb-sweep { to { transform: translateX(100%); } }

/* ═══ AI 整页生成 loading：全屏遮罩 + 旋转光环 ═══ */
.ai-page-loading {
  position: absolute;
  inset: 0;
  z-index: 100;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(5, 6, 12, 0.6);
  backdrop-filter: blur(3px);
}
.ai-loading-box {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 14px;
  padding: 32px 48px;
  border-radius: 16px;
  background: rgba(10, 14, 22, 0.85);
  border: 1px solid rgba(0, 255, 255, 0.25);
  box-shadow: 0 0 40px rgba(0, 255, 255, 0.12);
}
.ai-orbit { position: relative; width: 64px; height: 64px; }
.ai-orbit-ring {
  position: absolute;
  inset: 0;
  border-radius: 50%;
  border: 2px solid rgba(0, 255, 255, 0.15);
  border-top-color: #00ffff;
  animation: ai-spin 0.9s linear infinite;
}
.ai-orbit-core {
  position: absolute;
  inset: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  background: rgba(0, 255, 255, 0.1);
  color: #00ffff;
  font-size: 18px;
  animation: ai-pulse 1.4s ease-in-out infinite;
}
.ai-loading-text { margin: 0; font-size: 15px; font-weight: 600; color: rgba(255, 255, 255, 0.9); }
.ai-loading-sub { margin: 0; font-size: 12px; color: rgba(0, 255, 255, 0.6); }
@keyframes ai-spin { to { transform: rotate(360deg); } }
@keyframes ai-pulse { 0%, 100% { transform: scale(1); opacity: 0.85; } 50% { transform: scale(1.08); opacity: 1; } }
.ai-fade-enter-active, .ai-fade-leave-active { transition: opacity 0.2s ease; }
.ai-fade-enter-from, .ai-fade-leave-to { opacity: 0; }

.ed-body { flex: 1; display: flex; min-height: 0; position: relative; }
.ed-canvas { flex: 1; min-width: 0; padding: 16px; }
.ed-side {
  flex: 0 0 400px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 16px 16px 16px 0;
  min-height: 0;
}
.side-tree { flex: 1; min-height: 0; overflow: hidden; }
.side-insp { flex: 0 0 auto; }
</style>

<style>
/* 全屏标注编辑器打开时：锁定背后页面滚动（避免滚轮穿透），并把 Element 提示/确认浮层抬到遮罩（z-index 3000）之上 */
body.ui-editor-open { overflow: hidden; }
body.ui-editor-open .el-message { z-index: 4000 !important; }
body.ui-editor-open .el-overlay { z-index: 3900 !important; }
</style>
