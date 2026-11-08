import { defineStore } from 'pinia'
import { reactive, computed, toRefs } from 'vue'
import * as uiTree from '~/utils/uiTree.js'

/**
 * UI 图转代码 —— 单项目（一张图 / 一棵树）状态与逻辑 store
 *
 * 唯一状态源约定：
 * - 选中：selectedNodeKey（树与画布共享，不各存一份）
 * - 节点：nodes 扁平 map（nodeKey -> node），parentId 表达父子，O(1) 查找
 * - 坐标：全部原图像素，缩放只在渲染层
 *
 * 本 store 只管数据与逻辑，不发请求；analyze / 保存 / 生成 的网络调用由上层 api + 页面编排。
 */
export const useUiProjectStore = defineStore('uiProject', () => {
  const state = reactive({
    projectId: '',
    target: 'flutter', // flutter | vue3 | react（复用现有）
    globalConfig: {
      themeColor: '',
      fontFamily: '',
      spacingUnit: 8,
      borderRadius: 0,
      backgroundColor: '',
    },
    // 原图信息（COS 地址 + 像素宽高；宽高用于相对<->像素换算与位置软提示）
    // uploadId：当前树归属的主图唯一标识，写回时精确定位（两页同图/同 URL 也不会串页）
    image: { url: '', width: 0, height: 0, uploadId: '' },
    // 当前编辑主图的 uploadId：用于过滤归属该图的附图素材（标注引用素材）
    assetScopeId: '',
    // nodeKey -> node
    nodes: {},
    selectedNodeKey: null,
    // 框选子节点模式：进入画布框选前设置，框完挂到该父下
    pendingParentKey: null,
    // 已提示过「位置可能选错」的节点，避免同一节点重复弹
    positionHinted: {},
  })

  // ═══ Getters ═══
  const nodeList = computed(() => Object.values(state.nodes))
  const selectedNode = computed(() => state.nodes[state.selectedNodeKey] || null)
  const hasNodes = computed(() => Object.keys(state.nodes).length > 0)
  const tree = computed(() => uiTree.assembleTree(state.nodes))
  // 已占用的顶层枚举 -> 数量（占用即禁用）
  const enumUsage = computed(() => {
    const map = {}
    nodeList.value
      .filter((n) => !n.parentId)
      .forEach((n) => { if (n.name) map[n.name] = (map[n.name] || 0) + 1 })
    return map
  })
  // 生成前校验：未选类型 / 未确认 / 待补全的节点
  const unconfirmedNodes = computed(() =>
    nodeList.value.filter(
      (n) =>
        n.status === uiTree.NODE_STATUS.EMPTY ||
        n.status === uiTree.NODE_STATUS.AI_DRAFT ||
        n.status === uiTree.NODE_STATUS.PENDING,
    ),
  )

  // 某父节点的直接子节点数组（按 sortIndex）
  const childrenOf = (parentKey) =>
    nodeList.value
      .filter((n) => n.parentId === parentKey)
      .sort((a, b) => (a.sortIndex || 0) - (b.sortIndex || 0))

  // 内部：重算某父节点下子的排序 + layout（manual 不覆盖）
  function relayoutParent(parentKey) {
    if (!parentKey || !state.nodes[parentKey]) return
    const parent = state.nodes[parentKey]
    const kids = childrenOf(parentKey)
    uiTree.sortSiblingsByPosition(kids)
    parent.childrenIds = kids.map((k) => k.nodeKey)
    parent.nodeType = uiTree.NODE_TYPE.CONTAINER // 有子即容器
    if (parent.layoutSource !== uiTree.LAYOUT_SOURCE.MANUAL) {
      const inferred = uiTree.inferLayout(kids.map((k) => k.bbox), 20, parent.bbox)
      if (inferred) {
        parent.layout = inferred
        parent.layoutSource = uiTree.LAYOUT_SOURCE.AUTO
      }
    }
  }

  // 重排顶层（parentId=null）
  function relayoutRoots() {
    const roots = nodeList.value.filter((n) => !n.parentId)
    uiTree.sortSiblingsByPosition(roots)
  }

  // 载入 / 恢复的树可能没保存过 layout（旧数据、后端精简回传）：
  // 按各层自身子框位置补一次自动推断，否则右侧「排列方式」全空，
  // 生成侧也拿不到该层子区域的排布信息。
  // 已有值（AI / 保存过的自动 / 手动）不覆盖；手动清空过的（layoutSource=manual）也不补。
  function inferMissingLayouts() {
    nodeList.value.forEach((n) => {
      if (n.layout || n.layoutSource === uiTree.LAYOUT_SOURCE.MANUAL) return
      const kids = childrenOf(n.nodeKey)
      if (!kids.length) return
      const inferred = uiTree.inferLayout(kids.map((k) => k.bbox), 20, n.bbox)
      if (inferred) {
        n.layout = inferred
        n.layoutSource = uiTree.LAYOUT_SOURCE.AUTO
      }
    })
  }

  // ═══ Actions ═══

  // 载入 / 恢复项目（按 projectId 拉回后回填）
  function loadProject(payload = {}) {
    state.projectId = payload.projectId || ''
    if (payload.target) state.target = payload.target
    if (payload.globalConfig) Object.assign(state.globalConfig, payload.globalConfig)
    // 旧项目数据不含 uploadId，先置空避免残留到新项目
    state.image.uploadId = ''
    if (payload.image) Object.assign(state.image, payload.image)
    state.nodes = {}
    ;(payload.nodes || []).forEach((n) => { state.nodes[n.nodeKey] = n })
    state.selectedNodeKey = null
    state.pendingParentKey = null
    state.positionHinted = {}
    inferMissingLayouts()
  }

  function reset() {
    state.projectId = ''
    state.nodes = {}
    state.selectedNodeKey = null
    state.pendingParentKey = null
    state.positionHinted = {}
  }

  // 从已确认的嵌套树恢复扁平节点（重开编辑器继续编辑 / 刷新后恢复用）
  // 后端回传的嵌套树不含 parentId/sortIndex，必须按 children 重新推导：
  // 扁平 map 只认 parentId，缺失时层级树被拉平，再次提交会把子级全丢。
  function loadTree(tree) {
    state.nodes = {}
    const walk = (list, parentId) => {
      (list || []).forEach((n, i) => {
        const kids = n.children || []
        const node = { ...n, parentId: parentId || null, sortIndex: n.sortIndex ?? i }
        delete node.children
        node.childrenIds = kids.length ? kids.map((k) => k.nodeKey) : (n.childrenIds || [])
        state.nodes[node.nodeKey] = node
        walk(kids, node.nodeKey)
      })
    }
    walk(Array.isArray(tree) ? tree : Object.values(tree || {}), null)
    // 自愈：父级不存在的孤儿升为顶层，否则层级树会丢掉第一层（只剩一个悬空子级）
    repairOrphans()
    // 旧树没存 layout（或后端回传的精简数据没有）：按子框位置补齐各层排列方式
    inferMissingLayouts()
    state.selectedNodeKey = null
    state.pendingParentKey = null
    state.positionHinted = {}
  }

  function setImage(image = {}) {
    Object.assign(state.image, image)
  }

  // 打开编辑器时绑定当前主图（uploadId）：素材引用只能选该主图名下的附图
  function setAssetScope(uploadId) {
    state.assetScopeId = uploadId || ''
  }

  function selectNode(key) {
    state.selectedNodeKey = key || null
  }

  function setPendingParent(key) {
    state.pendingParentKey = key || null
  }

  // 新框落点检测：返回重叠冲突 + 建议父节点 + 位置软提示（不改动数据）
  // topOnly：只与顶层区域类型比对（子级框在父级内的包含关系不参与提醒）
  function detectNewBox(bbox, { excludeKey, topOnly = false } = {}) {
    const nodes = topOnly ? nodeList.value.filter((n) => !n.parentId) : nodeList.value
    const { conflicts, suggestedParentKey } = uiTree.detectAgainstNodes(bbox, nodes, { excludeKey })
    const hints = suggestedParentKey ? [] : uiTree.positionHint(state.nodes[excludeKey]?.name || '', bbox, state.image)
    return { conflicts, suggestedParentKey, positionHints: hints }
  }

  // 新增顶层占位节点：draw 模式画完即建（暂无类型，status=pending），由右侧面板选类型
  function addTopPlaceholder(bbox) {
    const node = uiTree.createNode({
      bbox,
      name: '',
      parentId: null,
      nodeType: uiTree.NODE_TYPE.CONTAINER,
    })
    node.variantIndex = 0
    node.status = uiTree.NODE_STATUS.PENDING // 待选类型
    state.nodes[node.nodeKey] = node
    relayoutRoots()
    state.selectedNodeKey = node.nodeKey
    return node
  }

  // 给顶层占位节点指派类型，返回位置软提示（供上层弹提示）
  function assignTopName(key, name) {
    const node = state.nodes[key]
    if (!node || node.parentId || !name) return { ok: false, hints: [] }
    // 同枚举多块：统计其它顶层已用数量作序号
    node.variantIndex = nodeList.value.filter(
      (n) => !n.parentId && n.name === name && n.nodeKey !== key,
    ).length
    node.name = name
    relayoutRoots()
    return { ok: true, hints: uiTree.positionHint(name, node.bbox, state.image) }
  }

  // 新增顶层节点：框完选枚举后调用
  function addTopNode({ bbox, name }) {
    const node = uiTree.createNode({
      bbox,
      name,
      parentId: null,
      nodeType: uiTree.NODE_TYPE.CONTAINER,
    })
    // 同枚举多块：记序号，用于展示 name-2 等
    node.variantIndex = enumUsage.value[name] || 0
    state.nodes[node.nodeKey] = node
    relayoutRoots()
    state.selectedNodeKey = node.nodeKey
    return node
  }

  // 给 ≥2 级节点指派元素类型（component）：只表达结构语义，不动区域名；
  // 传空 = 取消选择；叶子型却已有子框时只回软提示，不阻止（生成侧仍按容器处理）
  function assignElementType(key, type) {
    const node = state.nodes[key]
    if (!node || !node.parentId) return { ok: false, hints: [] }
    node.component = type || ''
    uiTree.markEdited(node)
    const hints = []
    if (node.component && !uiTree.isContainerElement(node.component) && childrenOf(key).length) {
      hints.push(`「${uiTree.elementLabel(node.component)}」通常不含子元素，已有子框会被保留并按容器生成`)
    }
    return { ok: true, hints }
  }

  // 新增子节点：优先用显式 parentKey，其次 pendingParentKey，最后建议父（包含判断）
  function addChildNode({ bbox, name = '', parentKey }) {
    const resolved = detectNewBox(bbox)
    const requested = parentKey || state.pendingParentKey || resolved.suggestedParentKey || null
    // 父级可能已被删除（框选进行中删了父节点/换了页）：不写悬空 parentId，退化为顶层占位
    const finalParent = requested && state.nodes[requested] ? requested : null
    const node = uiTree.createNode({
      bbox,
      name,
      parentId: finalParent,
      // 子节点默认叶子；退化为顶层时按容器处理（等待右栏选区域类型）
      nodeType: finalParent ? uiTree.NODE_TYPE.LEAF : uiTree.NODE_TYPE.CONTAINER,
    })
    node.status = finalParent ? uiTree.NODE_STATUS.EMPTY : uiTree.NODE_STATUS.PENDING
    // 挂到父下则成为父的子；若父由空容器变有子，relayout 里会转 container
    state.nodes[node.nodeKey] = node
    if (finalParent) relayoutParent(finalParent)
    else relayoutRoots()
    state.pendingParentKey = null
    state.selectedNodeKey = node.nodeKey
    return { node, suggestedParentKey: resolved.suggestedParentKey, conflicts: resolved.conflicts }
  }

  // 修复悬空 parentId：父节点已不存在的孤儿升为顶层（返回修复数量）。
  // 典型来源：删父节点 / 换页后旧的 pendingParentKey 仍被当作父级。
  function repairOrphans() {
    let fixed = 0
    nodeList.value.forEach((n) => {
      if (n.parentId && !state.nodes[n.parentId]) {
        n.parentId = null
        n.nodeType = uiTree.NODE_TYPE.CONTAINER
        fixed += 1
      }
    })
    if (fixed) relayoutRoots()
    return fixed
  }

  // 删除：cascade=true 连带子孙；false 把子节点提升到被删节点的父下
  function removeNode(key, cascade = true) {
    const node = state.nodes[key]
    if (!node) return
    const parentKey = node.parentId
    // 删除前先算好「下一个选中」：同层下一个 → 上一个 → 父级 → 任意剩余节点。
    // 否则选中被清空，右栏会从「节点」跳回「生成配置」，打断连续删除。
    const sameLevel = nodeList.value
      .filter((n) => (n.parentId || null) === (parentKey || null))
      .sort((a, b) => (a.sortIndex || 0) - (b.sortIndex || 0))
    const pos = sameLevel.findIndex((n) => n.nodeKey === key)
    const removed = new Set(cascade ? uiTree.collectDescendants(state.nodes, key) : [key])
    const rest = sameLevel.filter((n) => !removed.has(n.nodeKey))
    const nextKey = rest[pos]?.nodeKey || rest[pos - 1]?.nodeKey || null

    if (cascade) {
      uiTree.collectDescendants(state.nodes, key).forEach((k) => { delete state.nodes[k] })
    } else {
      childrenOf(key).forEach((child) => { child.parentId = parentKey })
      delete state.nodes[key]
    }
    if (parentKey) relayoutParent(parentKey)
    else relayoutRoots()
    // 被删节点如果正被当作「框选子区域的父级」，一并失效
    if (state.pendingParentKey && !state.nodes[state.pendingParentKey]) state.pendingParentKey = null
    if (state.selectedNodeKey === key || !state.nodes[state.selectedNodeKey]) {
      state.selectedNodeKey = nextKey || parentKey || nodeList.value[0]?.nodeKey || null
    }
  }

  // 拖拽改层级：把 key 移到 newParentKey 下（newParentKey=null 即升为顶层）
  function setParent(key, newParentKey) {
    const node = state.nodes[key]
    if (!node || key === newParentKey) return
    // 目标父级必须存在
    if (newParentKey && !state.nodes[newParentKey]) return
    // 不能挂到自己或自己的后代下
    if (newParentKey && uiTree.collectDescendants(state.nodes, key).includes(newParentKey)) return
    const oldParent = node.parentId
    node.parentId = newParentKey || null
    if (oldParent) relayoutParent(oldParent)
    else relayoutRoots()
    if (newParentKey) relayoutParent(newParentKey)
  }

  // 移动 / 缩放框：更新 bbox 后重排 + 重推父 layout
  function moveNode(key, bbox) {
    const node = state.nodes[key]
    if (!node) return
    node.bbox = uiTree.normalizeBbox(bbox)
    if (node.parentId) relayoutParent(node.parentId)
    else relayoutRoots()
  }

  // 改节点字段：描述类字段视为用户确认；layout 走 manual
  function updateNode(key, patch = {}) {
    const node = state.nodes[key]
    if (!node) return
    const descKeys = ['summary', 'styleHints', 'behavior', 'name', 'assetId', 'component']
    let touchedDesc = false
    Object.entries(patch).forEach(([k, v]) => {
      if (k === 'layout') { uiTree.setManualLayout(node, v); return }
      if (k === 'bbox') { node.bbox = uiTree.normalizeBbox(v); return }
      node[k] = v
      if (descKeys.includes(k)) touchedDesc = true
    })
    if (patch.bbox) { if (node.parentId) relayoutParent(node.parentId); else relayoutRoots() }
    if (touchedDesc) uiTree.markEdited(node)
  }

  // 提交生成用：组装嵌套树（做法 B，随请求体发送）
  function buildSubmitTree() {
    return uiTree.assembleTree(state.nodes)
  }

  return {
    ...toRefs(state),
    // getters
    nodeList,
    selectedNode,
    hasNodes,
    tree,
    enumUsage,
    unconfirmedNodes,
    childrenOf,
    // actions
    loadProject,
    reset,
    setImage,
    setAssetScope,
    selectNode,
    setPendingParent,
    detectNewBox,
    addTopPlaceholder,
    assignTopName,
    assignElementType,
    loadTree,
    addTopNode,
    addChildNode,
    removeNode,
    setParent,
    moveNode,
    repairOrphans,
    updateNode,
    buildSubmitTree,
  }
})
