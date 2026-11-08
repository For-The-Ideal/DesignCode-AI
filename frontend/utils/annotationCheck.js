/**
 * 标注完整性校验：AI 描述、开始生成共用同一套判定，调用方只负责提示文案。
 * 支持两种输入：
 *   - 编辑中的扁平列表（uiProject 的 nodeList，用 parentId 表达层级）
 *   - 提交树（images[].annotationTree，children 嵌套，节点内仍带 parentId）
 */

// 扁平化并按 nodeKey 去重（提交树里父子节点可能同时出现在数组与 children 中）
function flatten(nodes) {
  const seen = new Set()
  const out = []
  const walk = (list) => (list || []).forEach((n) => {
    if (!n) return
    const key = n.nodeKey || n.id
    if (key && seen.has(key)) return
    if (key) seen.add(key)
    out.push(n)
    walk(n.children)
  })
  walk(nodes)
  return out
}

/**
 * @param {Array} nodes 扁平节点列表或嵌套树
 * @returns {{
 *   total: number,           // 节点总数
 *   unnamedChildren: Array,  // 子级缺区域名（层级里会显示“未命名”）
 *   untypedTop: Array,       // 顶层未选类型（AI 描述与生成都依赖区域名）
 *   untypedElements: Array,  // 子级未选元素类型（软提示，不影响 ok）
 *   ok: boolean,             // 硬性校验是否通过
 * }}
 */
export function checkAnnotations(nodes) {
  const flat = flatten(nodes)
  const unnamedChildren = flat.filter((n) => n.parentId && !n.name)
  const untypedTop = flat.filter((n) => !n.parentId && !n.name)
  // 元素类型是新字段：老项目 / 未选也不拦截，由调用方决定是否软提示
  const untypedElements = flat.filter((n) => n.parentId && !n.component)
  return {
    total: flat.length,
    unnamedChildren,
    untypedTop,
    untypedElements,
    ok: !unnamedChildren.length && !untypedTop.length,
  }
}
