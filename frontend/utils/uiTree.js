/**
 * UI 图转代码 —— 树结构纯逻辑模块（渲染无关）
 *
 * 约定（与流程文档 / 已定规范对齐）：
 * - 所有坐标全链路使用「原图像素」：bbox = { x, y, w, h }，缩放只在渲染层做，不进这里。
 * - layout 三层来源优先级：manual（用户手改）> auto（子节点坐标推断）> ai（顶层建议）。
 *   有真实子框时以子框为准；顶层尚未下钻时先沿用 AI 建议；用户手动改过后不再被任何自动逻辑覆盖。
 * - 节点以「扁平 map + parentId」存储，O(1) 查找；树按需组装。
 * - 本模块只做纯函数，不触网、不依赖 Vue/Pinia，便于单测与复用。
 */

// 顶层区域枚举：AI 返回与用户框选顶层时都只能取其中之一
export const TOP_REGION_ENUM = [
  'header', 'sidebar', 'main', 'footer', 'banner', 'modal', 'floating', 'other',
]

// layout 合法值
export const LAYOUT_VALUES = ['row', 'column', 'grid']

// 节点类型
export const NODE_TYPE = { CONTAINER: 'container', LEAF: 'leaf' }

// 节点状态
export const NODE_STATUS = {
  PENDING: 'pending',   // 待描述
  AI_DRAFT: 'ai_draft', // AI 初稿
  CONFIRMED: 'confirmed', // 已确认（用户改过）
  EMPTY: 'empty',       // 待补全
}

// layout 来源
export const LAYOUT_SOURCE = { AI: 'ai', AUTO: 'auto', MANUAL: 'manual' }

// 重叠检测阈值：>0.8 强提示（拦截），0.3~0.8 软提示
export const OVERLAP_STRONG = 0.8
export const OVERLAP_SOFT = 0.3

let keySeq = 0
// 前端生成的唯一 nodeKey，贯穿匹配 / 补全 / 生成全流程
export function genNodeKey() {
  keySeq += 1
  return `n${Date.now().toString(36)}_${keySeq.toString(36)}`
}

// 归一化 bbox，保证数值类型与非负宽高
export function normalizeBbox(bbox = {}) {
  return {
    x: Number(bbox.x) || 0,
    y: Number(bbox.y) || 0,
    w: Math.max(0, Number(bbox.w) || 0),
    h: Math.max(0, Number(bbox.h) || 0),
  }
}

/**
 * 创建一个节点（默认顶层容器、状态待补全）。四组字段一次性给全，
 * 后续各处再按流程改写。
 */
export function createNode({
  bbox,
  name = '',
  parentId = null,
  nodeType = NODE_TYPE.CONTAINER,
} = {}) {
  return {
    // 结构组
    nodeKey: genNodeKey(),
    parentId,
    nodeType,
    bbox: normalizeBbox(bbox),
    sortIndex: 0,
    childrenIds: [],
    // 排列组
    layout: '',
    layoutSource: '',
    // 描述组
    name,
    summary: '',
    styleHints: '',
    behavior: '',
    // 状态组
    status: NODE_STATUS.EMPTY,
    aiFilled: false,
  }
}

export function bboxArea(b) {
  return Math.max(0, b.w) * Math.max(0, b.h)
}

export function bboxCenter(b) {
  return { x: b.x + b.w / 2, y: b.y + b.h / 2 }
}

/**
 * 重叠比 = 交集面积 ÷ 较小框面积。返回 0~1。
 */
export function overlapRatio(a, b) {
  const ix = Math.max(0, Math.min(a.x + a.w, b.x + b.w) - Math.max(a.x, b.x))
  const iy = Math.max(0, Math.min(a.y + a.h, b.y + b.h) - Math.max(a.y, b.y))
  const inter = ix * iy
  const minArea = Math.min(bboxArea(a), bboxArea(b))
  return minArea > 0 ? inter / minArea : 0
}

// inner 是否完全落在 outer 内部（用于自动建议父节点）
export function isContained(inner, outer) {
  return (
    inner.x >= outer.x &&
    inner.y >= outer.y &&
    inner.x + inner.w <= outer.x + outer.w &&
    inner.y + inner.h <= outer.y + outer.h
  )
}

/**
 * 坐标推断 layout：取所有直接子节点中心点，按 x/y 方向跨度对比阈值判断散列方向。
 * - 只有 1 个子节点 → column
 * - 横向分散、纵向不分散 → row；反之 → column；两向都分散 → grid
 * @param {Array<{x,y,w,h}>} childBboxes
 * @param {number} thresholdPx 原图像素阈值，可调
 */
export function inferLayout(childBboxes, thresholdPx = 20) {
  const list = (childBboxes || []).filter((b) => b && b.w > 0 && b.h > 0)
  if (list.length === 0) return ''
  if (list.length === 1) return 'column'
  const centers = list.map(bboxCenter)
  const xs = centers.map((c) => c.x)
  const ys = centers.map((c) => c.y)
  const xSpan = Math.max(...xs) - Math.min(...xs)
  const ySpan = Math.max(...ys) - Math.min(...ys)
  const xSpread = xSpan > thresholdPx
  const ySpread = ySpan > thresholdPx
  if (xSpread && !ySpread) return 'row'
  if (!xSpread && ySpread) return 'column'
  if (xSpread && ySpread) return 'grid'
  return 'column'
}

/**
 * 兄弟节点按视觉位置排序：先比中心 y（行带容差内视为同排），再比 x。
 * 就地写回每个节点的 sortIndex，并返回排序后的数组。
 */
export function sortSiblingsByPosition(nodes, bandPx = 40) {
  const arr = [...nodes].sort((a, b) => {
    const ca = bboxCenter(a.bbox)
    const cb = bboxCenter(b.bbox)
    if (Math.abs(ca.y - cb.y) > bandPx) return ca.y - cb.y
    if (Math.abs(ca.x - cb.x) > bandPx) return ca.x - cb.x
    return ca.y - cb.y
  })
  arr.forEach((n, i) => { n.sortIndex = i })
  return arr
}

// 顶层按 name 精确匹配第一次 AI 返回的 region（不做模糊语义匹配）
export function matchTopRegion(name, topRegions) {
  if (!Array.isArray(topRegions)) return null
  return topRegions.find((r) => r && r.name === name) || null
}

/**
 * 用匹配到的顶层 region 自动填充节点描述。
 * - 命中：填 summary/styleHints/behavior；layoutHint 合法则 → layout（来源 ai）；状态 ai_draft。
 * - region 为空或 layoutHint 非法：排列方式不替用户做主，清成未选中（重选类型时不残留上一次的建议）。
 */
export function applyRegionToNode(node, region) {
  if (!region) {
    node.layout = ''
    node.layoutSource = ''
    node.status = NODE_STATUS.EMPTY
    node.aiFilled = false
    return node
  }
  node.summary = region.summary || ''
  node.styleHints = region.styleHints || ''
  node.behavior = region.behavior || ''
  if (region.layoutHint && LAYOUT_VALUES.includes(region.layoutHint)) {
    node.layout = region.layoutHint
    node.layoutSource = LAYOUT_SOURCE.AI
  } else if (node.layoutSource !== LAYOUT_SOURCE.MANUAL) {
    // AI 没给（或换了个没给的类型）就回到未选中，但不覆盖用户手动选过的
    node.layout = ''
    node.layoutSource = ''
  }
  const hasDesc = !!(node.summary || node.styleHints)
  node.status = hasDesc ? NODE_STATUS.AI_DRAFT : NODE_STATUS.EMPTY
  node.aiFilled = hasDesc
  return node
}

// 用户改过任意描述字段 → 视为已确认，AI 初稿标记清除
export function markEdited(node) {
  node.status = NODE_STATUS.CONFIRMED
  node.aiFilled = false
  return node
}

// 用户手动设置 layout → 来源 manual，后续不再被自动覆盖
export function setManualLayout(node, layout) {
  if (LAYOUT_VALUES.includes(layout)) {
    node.layout = layout
    node.layoutSource = LAYOUT_SOURCE.MANUAL
  }
  return node
}

/**
 * 新框与已有节点比对：返回重叠冲突列表 + 建议父节点（完全包含新框的最小节点）。
 */
export function detectAgainstNodes(bbox, nodes, { excludeKey } = {}) {
  const conflicts = []
  let suggestedParentKey = null
  let smallestArea = Infinity
  for (const n of nodes) {
    if (n.nodeKey === excludeKey) continue
    const r = overlapRatio(bbox, n.bbox)
    if (r > OVERLAP_STRONG) conflicts.push({ nodeKey: n.nodeKey, name: n.name, level: 'high', ratio: r })
    else if (r >= OVERLAP_SOFT) conflicts.push({ nodeKey: n.nodeKey, name: n.name, level: 'soft', ratio: r })
    if (isContained(bbox, n.bbox)) {
      const a = bboxArea(n.bbox)
      if (a < smallestArea) { smallestArea = a; suggestedParentKey = n.nodeKey }
    }
  }
  return { conflicts, suggestedParentKey }
}

/**
 * 位置合理性软提示：按框中心相对图片宽高比例判断枚举是否可能选错。只提示，不阻止。
 * @param {string} name 顶层枚举
 * @param {{x,y,w,h}} bbox
 * @param {{width,height}} image 原图尺寸
 */
export function positionHint(name, bbox, image) {
  const c = bboxCenter(bbox)
  const ry = image && image.height ? c.y / image.height : 0
  const rx = image && image.width ? c.x / image.width : 0
  const hints = []
  if (name === 'header' && ry > 0.75) hints.push('位于底部，是否选错？')
  if (name === 'footer' && ry < 0.25) hints.push('位于上部，是否选错？')
  if (name === 'sidebar' && rx > 0.3 && rx < 0.7) hints.push('位于中部，是否选错？')
  return hints
}

// 某节点的全部子孙 nodeKey（含自身），用于级联删除
export function collectDescendants(nodeMap, key) {
  const result = []
  const walk = (k) => {
    if (!nodeMap[k]) return
    result.push(k)
    ;(nodeMap[k].childrenIds || []).forEach(walk)
  }
  walk(key)
  return result
}

/**
 * 由扁平 map 组装成嵌套树（提交生成用，做法 B）。
 * 重算每个容器的 childrenIds（按 parentId 归集 + sortIndex 排序），根节点排在顶层。
 * 返回的每个节点是浅拷贝并附带 children 数组，不改动原 map。
 */
export function assembleTree(nodeMap) {
  const nodes = Object.values(nodeMap || {})
  const byKey = {}
  nodes.forEach((n) => { byKey[n.nodeKey] = n })

  // 以 parentId 为准重算父子与排序
  const childrenIds = {}
  const roots = []
  nodes.forEach((n) => { childrenIds[n.nodeKey] = [] })
  nodes.forEach((n) => {
    if (n.parentId && byKey[n.parentId]) childrenIds[n.parentId].push(n.nodeKey)
    else roots.push(n.nodeKey)
  })
  const bySort = (a, b) => (byKey[a].sortIndex || 0) - (byKey[b].sortIndex || 0)
  Object.keys(childrenIds).forEach((k) => childrenIds[k].sort(bySort))
  roots.sort(bySort)

  const build = (key) => {
    const n = byKey[key]
    const kids = childrenIds[key] || []
    return { ...n, childrenIds: kids, children: kids.map(build) }
  }
  return roots.map(build)
}
