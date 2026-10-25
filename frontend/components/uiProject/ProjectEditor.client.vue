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

          <!-- AI 顶层识别未响应时的整页 loading（遮挡主体区，工具条仍可关闭） -->
          <div v-if="analyzing" class="ed-loading">
            <div class="loading-card">
              <i class="fas fa-spinner fa-spin loading-icon"></i>
              <span class="loading-title">AI 正在识别图片布局…</span>
              <span class="loading-sub">识别完成后框选顶层将自动填充描述，稍候即可开始</span>
            </div>
          </div>
        </div>
      </div>

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

const store = useUiProjectStore()

defineProps({
  // 顶层识别中：由上层（UploadZone）驱动，仅用于工具条提示，不阻塞框选
  analyzing: { type: Boolean, default: false },
})

// 「确认标注」通过后上抛已组装的树；由上层写回图片缩略标注并关闭编辑器。
// 真正的下单（创建任务/扣积分/SSE）由 ConfigPanel「开始生成」另行触发，编辑器不碰网络。
// close：编辑器关闭（工具条「关闭」或确认后 close）时上抛，供上层中止仍在途的识别请求。
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
    // draw 模式：先建占位顶层节点并选中，类型由右侧面板内联选择（不在图上弹层）
    store.addTopPlaceholder(bbox)
  }
}

// ═══ 子级：直接创建并挂父（保持 child 模式，可连续标注多个子区域）═══
function createChild(bbox) {
  const res = store.addChildNode({ bbox, parentKey: pendingParentKey.value })
  applyConflicts({ conflicts: res.conflicts }, () => {
    if (res.suggestedParentKey && res.node.parentId === res.suggestedParentKey) {
      const p = store.nodes[res.suggestedParentKey]
      ElMessage.info(`已自动挂到「${p?.name || '未命名'}」下，可继续框选或点“退出”`)
    }
  })
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

// ═══ 确认标注：只负责「校验未标注 + 上抛已确认树」，不提交 ═══
function onGenerate() {
  // 硬拦：顶层框还没选类型（无 name）——无法判断区域身份，不能进生成
  const noType = store.nodeList.filter((n) => !n.parentId && !n.name)
  if (noType.length) {
    ElMessage.warning(`有 ${noType.length} 个顶层区域还没选类型，请先在右侧点选`)
    store.selectNode(noType[0].nodeKey) // 自动定位第一个未完成项
    return
  }
  const tree = store.buildSubmitTree()
  // 软提示：仍有「待补全 / AI 初稿未确认」的节点（不拦截，结构完整即可生成）
  const pending = store.unconfirmedNodes
  if (pending.length) {
    ElMessageBox.confirm(
      `还有 ${pending.length} 个区域未填描述或为 AI 初稿，可直接确认（生成时按结构 + 已有描述处理）。是否确认标注？`,
      '确认标注',
      { confirmButtonText: '确认标注', cancelButtonText: '返回检查', type: 'warning' },
    )
      .then(() => emit('confirm', tree))
      .catch(() => {})
  } else {
    emit('confirm', tree)
  }
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

/* 顶层识别 loading 遮罩：覆盖整个主体区（画布+树+表单） */
.ed-loading {
  position: absolute;
  inset: 0;
  z-index: 20;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(5, 6, 12, 0.55);
  backdrop-filter: blur(2px);
}
.loading-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
  padding: 26px 34px;
  border-radius: 14px;
  background: rgba(12, 14, 22, 0.92);
  border: 1px solid rgba(0, 255, 255, 0.28);
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.55);
}
.loading-icon { font-size: 30px; color: #00ffff; }
.loading-title { font-size: 15px; font-weight: 600; color: rgba(255, 255, 255, 0.88); }
.loading-sub { font-size: 12px; color: rgba(255, 255, 255, 0.42); text-align: center; }
</style>

<style>
/* 全屏标注编辑器打开时：锁定背后页面滚动（避免滚轮穿透），并把 Element 提示/确认浮层抬到遮罩（z-index 3000）之上 */
body.ui-editor-open { overflow: hidden; }
body.ui-editor-open .el-message { z-index: 4000 !important; }
body.ui-editor-open .el-overlay { z-index: 3900 !important; }
</style>
