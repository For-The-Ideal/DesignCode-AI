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
    image: { url: '', width: 0, height: 0 },
    // 第一次顶层 AI 返回（后端已规范化），供框选顶层时按 name 匹配自动填充
    topRegions: [],
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
      const inferred = uiTree.inferLayout(kids.map((k) => k.bbox))
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

  // ═══ Actions ═══

  // 载入 / 恢复项目（按 projectId 拉回后回填）
  function loadProject(payload = {}) {
    state.projectId = payload.projectId || ''
    if (payload.target) state.target = payload.target
    if (payload.globalConfig) Object.assign(state.globalConfig, payload.globalConfig)
    if (payload.image) Object.assign(state.image, payload.image)
    state.topRegions = Array.isArray(payload.topRegions) ? payload.topRegions : []
    state.nodes = {}
    ;(payload.nodes || []).forEach((n) => { state.nodes[n.nodeKey] = n })
    state.selectedNodeKey = null
    state.pendingParentKey = null
    state.positionHinted = {}
  }

  function reset() {
    state.projectId = ''
    state.topRegions = []
    state.nodes = {}
    state.selectedNodeKey = null
    state.pendingParentKey = null
    state.positionHinted = {}
  }

  // 从已确认的嵌套树恢复扁平节点（重开编辑器继续编辑 / 刷新后恢复用）
  function loadTree(tree) {
    state.nodes = {}
    const walk = (list) => {
      (list || []).forEach((n) => {
        const kids = n.children || []
        const node = { ...n }
        delete node.children
        node.childrenIds = n.childrenIds || kids.map((k) => k.nodeKey)
        state.nodes[node.nodeKey] = node
        walk(kids)
      })
    }
    walk(Array.isArray(tree) ? tree : Object.values(tree || {}))
    state.selectedNodeKey = null
    state.pendingParentKey = null
  }

  function setImage(image = {}) {
    Object.assign(state.image, image)
  }

  function setTopRegions(regions) {
    state.topRegions = Array.isArray(regions) ? regions : []
  }

  function selectNode(key) {
    state.selectedNodeKey = key || null
  }

  function setPendingParent(key) {
    state.pendingParentKey = key || null
  }

  // 新框落点检测：返回重叠冲突 + 建议父节点 + 位置软提示（不改动数据）
  function detectNewBox(bbox, { excludeKey } = {}) {
    const { conflicts, suggestedParentKey } = uiTree.detectAgainstNodes(bbox, nodeList.value, { excludeKey })
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

  // 给顶层占位节点指派类型：匹配 topRegions 自动填充，返回位置软提示（供上层弹提示）
  function assignTopName(key, name) {
    const node = state.nodes[key]
    if (!node || node.parentId || !name) return { ok: false, hints: [] }
    // 同枚举多块：统计其它顶层已用数量作序号
    node.variantIndex = nodeList.value.filter(
      (n) => !n.parentId && n.name === name && n.nodeKey !== key,
    ).length
    node.name = name
    const region = uiTree.matchTopRegion(name, state.topRegions)
    uiTree.applyRegionToNode(node, region)
    if (node.variantIndex > 0 && region && node.summary) {
      node.summary = `与第 ${node.variantIndex + 1} 块同属 ${name}，请细化分工：` + node.summary
    }
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
    // 按 name 精确匹配顶层 AI 区域自动填充
    const region = uiTree.matchTopRegion(name, state.topRegions)
    uiTree.applyRegionToNode(node, region)
    if (node.variantIndex > 0 && region && node.summary) {
      node.summary = `与第 ${node.variantIndex + 1} 块同属 ${name}，请细化分工：` + node.summary
    }
    state.nodes[node.nodeKey] = node
    relayoutRoots()
    state.selectedNodeKey = node.nodeKey
    return node
  }

  // 新增子节点：优先用显式 parentKey，其次 pendingParentKey，最后建议父（包含判断）
  function addChildNode({ bbox, name = '', parentKey }) {
    const resolved = detectNewBox(bbox)
    const finalParent = parentKey || state.pendingParentKey || resolved.suggestedParentKey || null
    const node = uiTree.createNode({
      bbox,
      name,
      parentId: finalParent,
      nodeType: uiTree.NODE_TYPE.LEAF, // 子节点默认叶子
    })
    node.status = uiTree.NODE_STATUS.EMPTY // 子节点默认无描述
    // 挂到父下则成为父的子；若父由空容器变有子，relayout 里会转 container
    state.nodes[node.nodeKey] = node
    if (finalParent) relayoutParent(finalParent)
    else relayoutRoots()
    state.pendingParentKey = null
    state.selectedNodeKey = node.nodeKey
    return { node, suggestedParentKey: resolved.suggestedParentKey, conflicts: resolved.conflicts }
  }

  // 删除：cascade=true 连带子孙；false 把子节点提升到被删节点的父下
  function removeNode(key, cascade = true) {
    const node = state.nodes[key]
    if (!node) return
    const parentKey = node.parentId
    if (cascade) {
      uiTree.collectDescendants(state.nodes, key).forEach((k) => { delete state.nodes[k] })
    } else {
      childrenOf(key).forEach((child) => { child.parentId = parentKey })
      delete state.nodes[key]
    }
    if (parentKey) relayoutParent(parentKey)
    else relayoutRoots()
    if (state.selectedNodeKey === key || !state.nodes[state.selectedNodeKey]) {
      state.selectedNodeKey = null
    }
  }

  // 拖拽改层级：把 key 移到 newParentKey 下（newParentKey=null 即升为顶层）
  function setParent(key, newParentKey) {
    const node = state.nodes[key]
    if (!node || key === newParentKey) return
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
    const descKeys = ['summary', 'styleHints', 'behavior', 'name']
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
    setTopRegions,
    selectNode,
    setPendingParent,
    detectNewBox,
    addTopPlaceholder,
    assignTopName,
    loadTree,
    addTopNode,
    addChildNode,
    removeNode,
    setParent,
    moveNode,
    updateNode,
    buildSubmitTree,
  }
})
