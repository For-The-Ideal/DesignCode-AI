<template>
  <!-- 设计稿工作台：左「页面列表/素材」 / 中「画布或空态上传」 / 右「节点或生成配置」，三块均为独立组件 -->
  <div class="wb-view">
    <WorkbenchUpload class="wb-col wb-col--upload" @annotate="onAnnotate" />

    <div class="wb-col wb-col--canvas">
      <div class="wb-canvas-wrap" :class="{ 'is-board': showBoard }">
        <!-- 多图画板：当前页卡片就地放大并内嵌标注画布，框选/编辑直接在图上生效 -->
        <WorkbenchBoard
          v-if="showBoard"
          class="wb-canvas-slot"
          :mode="mode"
          @drawn="onDrawn"
          @cancelMode="exitMode"
          @dropFiles="onDropFiles"
        />

        <WorkbenchCanvas
          v-else
          class="wb-canvas-slot"
          :mode="mode"
          :scene="codeStore.mode"
          :has-image="hasImage"
          :uploading="!!currentImage?.uploading"
          :error="!!currentImage?.uploadError"
          :max-count="codeStore.maxImages"
          @drawn="onDrawn"
          @cancelMode="exitMode"
          @pick="pickImage"
          @drop="onDropFiles"
        />

        <!-- 底部工具条：与单图模式同一套（框选 / 编辑 / AI 描述），始终作用在当前页 -->
        <div v-if="hasImage" class="wb-toolbar-float">
          <WorkbenchToolbar
            :mode="mode"
            :parent-name="parentName"
            :has-nodes="uiStore.hasNodes"
            :ai-loading="aiLoading"
            @setMode="setMode"
            @exitChild="exitMode"
            @ai="onAiFill"
          />
        </div>
      </div>
    </div>

    <WorkbenchInspector
      class="wb-col wb-col--inspector"
      @addChild="onAddChild"
      @deleteNode="onDeleteNode"
      @select="onTreeSelect"
    />

    <!-- 空态上传入口（左栏不再放上传框）；单图只允许选一张 -->
    <input
      ref="pickInput"
      type="file"
      accept="image/*"
      :multiple="codeStore.mode === 'multi'"
      hidden
      @change="onPickChange"
    />

    <!-- AI 整页描述生成中：遮罩整个工作台（画布/工具栏/左栏/右侧检查器都不可操作） -->
    <AiLoadingOverlay :show="aiLoading" :sub="'共 ' + uiStore.nodeList.length + ' 个区域，生成后自动回填'" />
  </div>
</template>

<script setup>
import { computed, ref, watch, onMounted, onBeforeUnmount } from 'vue'
import { storeToRefs } from 'pinia'
import { ElMessage, ElMessageBox } from 'element-plus'
import WorkbenchUpload from './WorkbenchUpload.vue'
import WorkbenchBoard from './WorkbenchBoard.vue'
import WorkbenchToolbar from './WorkbenchToolbar.vue'
import WorkbenchCanvas from './WorkbenchCanvas.vue'
import WorkbenchInspector from './WorkbenchInspector.vue'
import { useCodeStore } from '~/stores/code'
import { useUiProjectStore } from '~/stores/uiProject'
import { useUserStore } from '~/stores/user'
import { useCommonStore } from '~/stores/common'
import { useRegionDescribe } from '~/composables/useRegionDescribe'
import AiLoadingOverlay from '~/components/common/AiLoadingOverlay.vue'
import { checkAnnotations } from '~/utils/annotationCheck'
import { nodeLabel, hintLabel, describeConflicts } from '~/utils/uiTree.js'

const codeStore = useCodeStore()
const uiStore = useUiProjectStore()
const userStore = useUserStore()
const commonStore = useCommonStore()
const { isLogin } = storeToRefs(userStore)
const { aiLoading, aiDescribe } = useRegionDescribe()

// 当前编辑主图：多图按 codeStore.activePageId 取页，单图恒为唯一主图
const currentImage = computed(() => codeStore.currentImage)
const hasImage = computed(() => !!currentImage.value?.cosUrl)

// ═══ 多图视图：中栏恒为画板，当前页卡片就地放大标注（不再单开一页标注视图） ═══
const showBoard = computed(() => codeStore.mode === 'multi' && codeStore.images.length > 0)

// 画布编辑模式：换页 / 换模式时必须先复位再载入新树，
// 否则旧的 pendingParentKey 会把新框挂到已不存在的父节点下（层级树丢第一层）
const mode = ref('draw')
const pendingParentKey = ref(null)

// 左栏双击/「标注」：只是把该页设为当前页，画板会把它放大成可标注的卡面
function onAnnotate(pageId) {
  if (pageId) codeStore.setActivePage(pageId)
}

// ═══ 空态上传：登录校验 → 选文件/拖拽 → store.addImages ═══
const pickInput = ref(null)

function pickImage() {
  if (!isLogin.value) return commonStore.setLoginModalVisible(true)
  if (codeStore.isSubmitting) return
  pickInput.value?.click()
}
function onPickChange(e) {
  const files = Array.from(e.target.files || [])
  if (files.length) codeStore.addImages(files)
  if (pickInput.value) pickInput.value.value = ''
}
function onDropFiles(files) {
  if (!isLogin.value) return commonStore.setLoginModalVisible(true)
  if (codeStore.isSubmitting) return
  codeStore.addImages(files)
}

// 载入树后默认选中第一个顶层区域：
// 从详情页「重新生成」恢复时右侧直接可编辑，不再停在「选中区域后可编辑描述」空态
function selectFirstTop() {
  if (uiStore.selectedNodeKey || !uiStore.hasNodes) return
  const first = uiStore.nodeList.find((n) => !n.parentId)
  if (first) uiStore.selectNode(first.nodeKey)
}

// ═══ 会话初始化：图片就绪即载入画布/树；切模式时「先写回旧图 → 再载入新图」 ═══
// 单图 / 多图是两套独立数据：旧模式的编辑必须先落盘，新模式的树再从自己的图上恢复；
// 目标模式还没图时清空右栏，不残留上一模式的结构。
// 当前已载入的图标识：uploadId 每张图唯一（两页同 URL 也能正确换树），旧数据回退 cosUrl
let loadedKey = ''
watch(
  () => [codeStore.mode, codeStore.activePageId, currentImage.value?.cosUrl, currentImage.value?.naturalWidth, currentImage.value?.treeWidth],
  () => {
    const img = currentImage.value
    // 此刻 image.url 仍指向旧图，先把 UI 里的结构写到它所属的那张图
    saveToImage()
    // 目标模式没有主图：清空结构与画布，右栏不显示上一模式的树
    if (!img?.cosUrl) {
      loadedKey = ''
      exitMode()
      if (uiStore.hasNodes || uiStore.image.url) {
        uiStore.loadTree([])
        uiStore.setImage({ url: '', width: 0, height: 0, uploadId: '' })
        uiStore.setAssetScope('')
      }
      return
    }
    const W = img.naturalWidth || img.treeWidth || 0
    const H = img.naturalHeight || img.treeHeight || 0
    const key = img.uploadId || img.cosUrl
    if (key !== loadedKey) {
      loadedKey = key
      // 换页 / 换模式：退出子级框选，pendingParentKey 不跨页携带
      exitMode()
      uiStore.loadTree(img.annotationTree || [])
      uiStore.setAssetScope(img.uploadId || '')
      selectFirstTop()
    } else if (!uiStore.hasNodes && img.annotationTree?.length) {
      uiStore.loadTree(img.annotationTree)
      selectFirstTop()
    }
    // 尺寸未变就不重复写入：避免同一张图因预览→远程地址等变化被重置缩放（画布闪烁）；
    // uploadId 一并同步，保证 saveToImage 始终写回正在编辑的那一页
    if (
      uiStore.image.url !== img.cosUrl ||
      uiStore.image.width !== W ||
      uiStore.image.height !== H ||
      uiStore.image.uploadId !== (img.uploadId || '')
    ) {
      uiStore.setImage({ url: img.cosUrl, width: W, height: H, uploadId: img.uploadId || '' })
    }
  },
  { immediate: true },
)

// ═══ 画布模式：draw（框选）/ select（编辑）/ child（为父级框子区域） ═══
// mode / pendingParentKey 已在上方声明（换页时需要先复位，见加载 watcher）
// 工具条只显示语义标签（横幅/卡片/按钮…），不展示用户填写的业务名（可能是公司/品牌词）
const parentName = computed(() => hintLabel(uiStore.nodes[pendingParentKey.value]))

function exitMode() {
  mode.value = 'draw'
  pendingParentKey.value = null
  uiStore.setPendingParent(null)
}

function setMode(m) {
  if (mode.value === m) return
  mode.value = m
  pendingParentKey.value = null
  uiStore.setPendingParent(null)
}

// ═══ 框选落点：child 模式直接挂父；draw 模式画在已有区域内部时按子区域自动挂到最小包含它的节点 ═══
function onDrawn(bbox) {
  if (mode.value === 'child') {
    const res = uiStore.addChildNode({ bbox, parentKey: pendingParentKey.value })
    if (res.suggestedParentKey && res.node.parentId === res.suggestedParentKey) {
      const label = hintLabel(uiStore.nodes[res.suggestedParentKey])
      ElMessage.info(label ? `已自动挂到「${label}」下，可继续框选或点“退出”` : '已自动挂到所选父区域下，可继续框选或点“退出”')
    }
    return
  }
  // draw 模式：完全画在已有区域内部的框，逻辑上就是子区域 → 自动挂到最小包含节点，不再弹重叠确认
  const inner = uiStore.detectNewBox(bbox)
  if (inner.suggestedParentKey) {
    const res = uiStore.addChildNode({ bbox, parentKey: inner.suggestedParentKey })
    const label = hintLabel(uiStore.nodes[res.node.parentId])
    ElMessage.info(label ? `已在「${label}」下创建子区域，可在右侧选元素类型` : '已在包含它的区域下创建子区域，可在右侧选元素类型')
    return
  }
  // 不在任何区域内部：按顶层区域处理，只在顶层之间做重叠提醒
  const det = uiStore.detectNewBox(bbox, { topOnly: true })
  applyConflicts(det, () => uiStore.addTopPlaceholder(bbox))
}

// 重叠处理：high 拦截给保留/取消，soft 仅提示；冲突方只描述类型（不拼接业务名）
function applyConflicts(det, proceed) {
  const high = (det.conflicts || []).filter((c) => c.level === 'high')
  if (high.length) {
    ElMessageBox.confirm(
      `新框与 ${describeConflicts(high, uiStore.nodes)} 高度重叠，是否仍作为独立区域保留？`,
      '重叠提醒',
      { confirmButtonText: '仍保留', cancelButtonText: '放弃该框', type: 'warning' },
    ).then(() => proceed()).catch(() => exitMode())
    return
  }
  const soft = (det.conflicts || []).filter((c) => c.level === 'soft')
  if (soft.length) ElMessage.warning(`与 ${describeConflicts(soft, uiStore.nodes)} 可能重叠，请确认`)
  proceed()
}

// ═══ 层级树事件 ═══
function onAddChild(parentKey) {
  pendingParentKey.value = parentKey
  uiStore.setPendingParent(parentKey)
  mode.value = 'child'
}

// key 是否为 ancestorKey 自身或其后代（沿 parentId 上溯）
function isDescendantOrSelf(key, ancestorKey) {
  let cur = key
  while (cur) {
    if (cur === ancestorKey) return true
    cur = uiStore.nodes[cur]?.parentId || null
  }
  return false
}

function onTreeSelect(key) {
  if (mode.value === 'child' && pendingParentKey.value && !isDescendantOrSelf(key, pendingParentKey.value)) {
    exitMode()
  }
  uiStore.selectNode(key)
}

function onDeleteNode(key) {
  const node = uiStore.nodes[key]
  if (!node) return
  const kids = uiStore.childrenOf(key)
  if (kids.length) {
    ElMessage.warning(`该区域下还有 ${kids.length} 个子区域，请先删除其子区域`)
    return
  }
  // 正在为这个节点框子区域：先退出子级模式，避免后续框挂到已删除父下
  if (pendingParentKey.value === key) exitMode()
  uiStore.removeNode(key, true)
}

// ═══ 保存：树写回图片（沿用现有 annotationTree / annotations 结构，不改数据格式） ═══
// 写回目标按「正在编辑的图」定位（跨单图 / 多图两个数据桶）：
// 切模式瞬间 currentImage 已指向新模式的图，但 UI 里的树仍属于旧图，
// 直接用 currentImage 会把旧结构写到新图上。
// 优先用 uploadId 精确匹配（每张图唯一）：两页即使同图/同 URL 也不会串页；旧数据回退 url。
function editingImage() {
  const all = [...codeStore.modes.single.images, ...codeStore.modes.multi.images]
  if (uiStore.image.uploadId) {
    const hit = all.find((i) => i.uploadId === uiStore.image.uploadId)
    if (hit) return hit
  }
  const url = uiStore.image.url
  if (!url) return null
  return all.find((i) => i.cosUrl && i.cosUrl === url) || null
}
function saveToImage() {
  const img = editingImage()
  if (!img) return
  const tree = uiStore.buildSubmitTree()
  const W = uiStore.image.width
  const H = uiStore.image.height
  img.annotationTree = tree
  img.annotations = toAnnotations(tree, W, H)
  if (W && H) {
    img.treeWidth = W
    img.treeHeight = H
  }
}
function toAnnotations(tree, W, H) {
  const list = []
  const walk = (nodes) => (nodes || []).forEach((n) => {
    list.push({
      id: n.nodeKey,
      label: nodeLabel(n, '未命名'),
      top: !n.parentId,
      status: n.status,
      x: W ? n.bbox.x / W : 0,
      y: H ? n.bbox.y / H : 0,
      w: W ? n.bbox.w / W : 0,
      h: H ? n.bbox.h / H : 0,
    })
    walk(n.children)
  })
  walk(tree)
  return list
}

// 编辑即写回当前主图：框选 / 拖拽 / 改属性 / 删除 都同步到 annotationTree。
// 开始生成的完整性校验（createTask）读的就是 img.annotationTree，
// 只在「AI 描述」时写回会导致未选类型的框也能直接提交。
watch(() => uiStore.nodeList, () => saveToImage(), { deep: true })

// 工作台会因「多图模式已有主图」被页面卸载：卸载前补一次写回，
// 避免排队中（如 300ms debounce）的最后一次编辑丢失。
onBeforeUnmount(() => saveToImage())

// 会话里可能残留旧脏数据（父级已删、子级悬空）：挂载时自愈一次，
// 否则层级树会把悬空子级当作根节点显示，看起来「第一层没了」
onMounted(() => {
  const fixed = uiStore.repairOrphans()
  if (fixed) ElMessage.info(`已把 ${fixed} 个父级丢失的区域恢复为顶层，请重新选择区域类型`)
})

async function onAiFill() {
  const img = currentImage.value
  if (!img) return

  // 先校验（与开始生成同一套判定），定位到第一个问题节点
  const check = checkAnnotations(uiStore.nodeList)
  if (check.unnamedChildren.length) {
    // 按父级的语义标签分组（「卡片」下 2 个），不拼接用户业务名
    const groups = new Map()
    check.unnamedChildren.forEach((n) => {
      const label = hintLabel(uiStore.nodes[n.parentId]) || '未标注类型的区域'
      groups.set(label, (groups.get(label) || 0) + 1)
    })
    const detail = [...groups.entries()]
      .slice(0, 3)
      .map(([label, count]) => (label === '未标注类型的区域' ? `${label}下 ${count} 个` : `「${label}」下 ${count} 个`))
      .join('、')
    uiStore.selectNode(check.unnamedChildren[0].nodeKey)
    ElMessage.warning(`还有 ${check.unnamedChildren.length} 个子区域没填区域名（${detail}），已选中第一个；树里标红的「缺名」就是它们`)
    return
  }
  if (check.untypedTop.length) {
    uiStore.selectNode(check.untypedTop[0].nodeKey)
    ElMessage.warning(`还有 ${check.untypedTop.length} 个区域没选类型，请在右侧选择后再生成`)
    return
  }
  // 元素类型是软校验：老项目未选也不拦截，只提示生成侧按容器语义推断
  if (check.untypedElements.length) {
    ElMessage.info(`有 ${check.untypedElements.length} 个子区域未选元素类型，AI 将按容器语义推断`)
  }

  // 3) 先保存当前编辑，避免 AI 请求失败时丢失标注
  saveToImage()

  const ok = await aiDescribe()
  if (ok) {
    saveToImage()
    ElMessage.success('AI 描述已回填，标注已保存')
  } else {
    ElMessage.warning('AI 描述未完成，当前标注已保存')
  }
}
</script>

<style scoped>
.wb-view {
  position: relative;
  display: flex;
  gap: 12px;
  flex: 1;
  min-height: 0;
}

.wb-col { min-height: 0; }
.wb-col--upload { flex: 0 0 224px; }
.wb-col--canvas {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
}
.wb-col--inspector { flex: 0 0 340px; }

/* 画布包裹层：工具条浮层相对它定位，不被下方附图区遮挡 */
.wb-canvas-wrap {
  position: relative;
  flex: 1;
  min-height: 0;
}

/* 工具条浮层：放在画布底部预留区（不遮挡设计稿、不占布局高度），右侧让出缩放条 */
.wb-toolbar-float {
  position: absolute;
  bottom: 10px;
  left: 14px;
  z-index: 20;
  max-width: calc(100% - 170px);
}
/* 多图画板：页面加宽后底部工具条仍浮在画布底部（与单图一致） */
.wb-canvas-slot {
  height: 100%;
}
/* 模式提示已由底部工具条承担，隐藏画布自带的提示条，避免重叠 */
.wb-col--canvas :deep(.mode-banner) {
  display: none;
}
/* 多图：工具条改成底部常驻条，不压在可框选的当前页图面上 */
.wb-canvas-wrap.is-board {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.wb-canvas-wrap.is-board .wb-canvas-slot {
  flex: 1;
  min-height: 0;
  height: auto;
}
.wb-canvas-wrap.is-board .wb-toolbar-float {
  position: static;
  flex-shrink: 0;
  max-width: none;
}

@media (max-width: 1440px) {
  .wb-col--upload { flex-basis: 208px; }
  .wb-col--inspector { flex-basis: 300px; }
}
</style>
