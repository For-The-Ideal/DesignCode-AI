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
// 注意：枚举值是全链路数据（AI prompt / 生成 DSL 映射 / 历史数据），只能保持英文；
// 界面上要显示中文时一律走 regionLabel()，不要改这里的值。
// 顺序即类型按钮的展示顺序（上→中→下→浮层→其他，便于快速找到）。
export const TOP_REGION_ENUM = [
  'header', 'search', 'banner', 'main', 'sidebar', 'footer', 'tabbar', 'modal', 'floating', 'other',
]

// 枚举 → 界面短标签（2-3 字，供类型按钮、标题、提示文案使用）
export const TOP_REGION_LABELS = {
  header: '顶部导航',
  search: '搜索栏',
  banner: '横幅',
  main: '主内容',
  sidebar: '侧边栏',
  footer: '底部栏',
  tabbar: '标签栏',
  modal: '弹窗',
  floating: '悬浮',
  other: '其他',
}

// 枚举 → 悬停释义（与后端 vision/region.txt 的判定口径保持一致）
export const TOP_REGION_TIPS = {
  header: '顶部导航/标题栏/状态栏',
  search: '搜索入口/搜索框区域（含顶部搜索条）',
  banner: '首屏大图轮播/运营横幅',
  main: '主体内容区（列表/卡片/表单）',
  sidebar: '侧边抽屉或竖向导航',
  footer: '底部操作条/版权栏（不含 Tab 导航）',
  tabbar: '底部 Tab 导航栏（主页面切换）',
  modal: '弹窗/浮层对话框',
  floating: '悬浮按钮/浮动元素',
  other: '无法归类的顶层块',
}

// 显示用：枚举转中文；子节点自由命名（非枚举）原样返回
export function regionLabel(name) {
  return TOP_REGION_LABELS[name] || name || ''
}

// 元素类型枚举（≥2 级节点用）：只描述「这块框是什么组件」，与顶层区域类型（在哪一块）分离。
// 值是全链路数据（提交 / 后端 DSL type 映射），只能保持英文；界面显示一律走 elementLabel()。
// 顺序即选择器展示顺序：容器型在前（可继续下钻），叶子型在后。
export const ELEMENT_CONTAINER_TYPES = [
  'group', 'card', 'list', 'listItem', 'grid', 'tabs', 'tab', 'swiper', 'form',
]
export const ELEMENT_LEAF_TYPES = [
  'title', 'text', 'image', 'icon', 'button', 'link', 'input', 'select',
  'switch', 'checkbox', 'avatar', 'badge', 'tag', 'divider',
]
export const ELEMENT_TYPE_ENUM = [...ELEMENT_CONTAINER_TYPES, ...ELEMENT_LEAF_TYPES]

// 枚举 → 界面短标签（供类型按钮 / 层级树 / 框线标签使用）
export const ELEMENT_TYPE_LABELS = {
  group: '分组',
  card: '卡片',
  list: '列表',
  listItem: '列表项',
  grid: '网格',
  tabs: '选项卡组',
  tab: '选项卡',
  swiper: '轮播',
  form: '表单',
  title: '标题',
  text: '文本',
  image: '图片',
  icon: '图标',
  button: '按钮',
  link: '链接',
  input: '输入框',
  select: '下拉',
  switch: '开关',
  checkbox: '勾选',
  avatar: '头像',
  badge: '角标',
  tag: '标签',
  divider: '分割线',
}

// 枚举 → 悬停释义
export const ELEMENT_TYPE_TIPS = {
  group: '通用分组容器，无特定语义',
  card: '卡片/数据卡（标题 + 内容的独立块）',
  list: '列表容器（同类项纵向排列）',
  listItem: '列表项/导航项/菜单行',
  grid: '宫格/多列网格容器',
  tabs: '选项卡组（含多个 tab）',
  tab: '单个选项卡',
  swiper: '轮播/横滑容器',
  form: '表单容器（输入项集合）',
  title: '标题文字',
  text: '正文/说明文字',
  image: '图片/插画/照片',
  icon: '图标',
  button: '按钮',
  link: '文字链接/可点文字',
  input: '单行输入框',
  select: '下拉选择/日期选择',
  switch: '开关',
  checkbox: '勾选框/单选',
  avatar: '头像',
  badge: '角标/未读数',
  tag: '标签/胶囊',
  divider: '分割线/间隔',
}

// 显示用：元素类型转中文；未知值原样返回
export function elementLabel(type) {
  return ELEMENT_TYPE_LABELS[type] || type || ''
}

/**
 * 提示语专用标签：只用「区域类型 / 元素类型」的语义词，不拼接用户填写的业务名。
 * 业务名可能是公司名/品牌词（设计稿里的真实文案），不能进系统提示。
 * 顶层只认枚举（枚举外的脏数据返回空），子级只认元素类型；空 = 没有可用的语义标签。
 */
export function hintLabel(node) {
  if (!node) return ''
  if (!node.parentId) return TOP_REGION_LABELS[node.name] || ''
  return elementLabel(node.component) || ''
}

/**
 * 重叠冲突 → 语义化描述（按类型聚合计数，不暴露节点业务名）：
 * conflicts + nodes map => 「卡片」区域、2 个「按钮」区域
 */
export function describeConflicts(conflicts, nodesByKey = {}) {
  const UNTYPED = '未标注类型的区域'
  const counts = {}
  ;(conflicts || []).forEach((c) => {
    const label = hintLabel(nodesByKey[c.nodeKey]) || UNTYPED
    counts[label] = (counts[label] || 0) + 1
  })
  return Object.entries(counts)
    .map(([label, n]) => {
      if (label === UNTYPED) return n > 1 ? `${n} 个${UNTYPED}` : UNTYPED
      return n > 1 ? `${n} 个「${label}」区域` : `「${label}」区域`
    })
    .join('、')
}

// 容器型元素可含子元素；叶子型仍允许框子级，只是给软提示
export function isContainerElement(type) {
  return ELEMENT_CONTAINER_TYPES.includes(type)
}

/**
 * 节点展示名（层级树 / 画布标签 / 白板框线统一用）：
 * - 顶层：区域类型中文 + 多块序号；
 * - 子级：「元素类型 · 业务名」；只选了类型显示类型，都没有则用 fallback。
 */
export function nodeLabel(n, fallback = '未命名区域') {
  if (!n) return fallback
  if (!n.parentId) {
    if (!n.name) return fallback
    const text = regionLabel(n.name)
    return n.variantIndex > 0 ? `${text}-${n.variantIndex + 1}` : text
  }
  const type = elementLabel(n.component)
  if (type && n.name) return `${type} · ${n.name}`
  return type || n.name || fallback
}

// layout 合法值：row/column/grid 为流式排列；stack 为叠放（子元素共用同一区域、按坐标叠加）
export const LAYOUT_VALUES = ['row', 'column', 'grid', 'stack']

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
    component: '', // 元素类型（≥2 级）：group/card/button/...，空=未选
    summary: '',
    styleHints: '',
    behavior: '',
    // 引用组：关联的附图素材（上传区 assets 里的 assetId，空串=未关联）
    assetId: '',
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
 * - 有子框几乎铺满父框（如背景图，≥80% 父框面积）→ stack（叠放）：
 *   其余元素浮在它上面，按中心点判行列会把浮层误判成竖排/网格；
 * - 只有 1 个子节点 → column
 * - 横向分散、纵向不分散 → row；反之 → column；两向都分散 → grid
 * @param {Array<{x,y,w,h}>} childBboxes
 * @param {number} thresholdPx 原图像素阈值，可调
 * @param {{x,y,w,h}} parentBBox 父框（传了才会识别叠放）
 */
export function inferLayout(childBboxes, thresholdPx = 20, parentBBox = null) {
  const list = (childBboxes || []).filter((b) => b && b.w > 0 && b.h > 0)
  if (list.length === 0) return ''
  if (parentBBox && list.length > 1) {
    const pa = bboxArea(parentBBox)
    if (pa > 0 && list.some((b) => bboxArea(b) >= pa * 0.8)) return 'stack'
  }
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

// 用户改过任意描述字段 → 视为已确认，AI 初稿标记清除
export function markEdited(node) {
  node.status = NODE_STATUS.CONFIRMED
  node.aiFilled = false
  return node
}

// 用户手动设置 layout → 来源 manual，后续不再被自动覆盖；传空 = 显式取消（不再给生成侧排列暗示）
export function setManualLayout(node, layout) {
  if (!layout || LAYOUT_VALUES.includes(layout)) {
    node.layout = layout || ''
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
  if (name === 'tabbar' && ry < 0.6) hints.push('不在底部，是否选错？')
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
