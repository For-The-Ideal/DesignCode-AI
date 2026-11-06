import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { useUiProjectStore } from '~/stores/uiProject'
import { useCodeStore } from '~/stores/code'
import { uiProjectApi } from '~/api/uiProject'

/**
 * 整页 AI 描述：提交整页区域树，AI 返回每个节点的描述，按 nodeKey 全部回填。
 * 属于全局操作（一次生成所有区域），由编辑器顶部工具栏调用。
 */
export function useRegionDescribe() {
  const store = useUiProjectStore()
  const codeStore = useCodeStore()
  const aiLoading = ref(false)

  // 整页区域树：全部顶层节点的完整子树（任意层级），isFocus 标记当前选中区域
  function buildRegionTree(focusNode) {
    const toNode = (cur, depth) => {
      const asset = cur.assetId ? codeStore.assets.find((a) => a.assetId === cur.assetId) : null
      return {
        nodeKey: cur.nodeKey,
        name: cur.name || '',
        summary: cur.summary || '',
        styleHints: cur.styleHints || '',
        layout: cur.layout || '',
        behavior: cur.behavior || '',
        bbox: cur.bbox ? { x: cur.bbox.x, y: cur.bbox.y, w: cur.bbox.w, h: cur.bbox.h } : null,
        // 该区域关联的附图（无引用时为 null）
        asset: asset && asset.cosUrl ? { label: asset.label || '', url: asset.cosUrl } : null,
        isFocus: !!focusNode && cur.nodeKey === focusNode.nodeKey,
        // 子级继续递归（任意层级），无子级为空数组
        children: depth >= 8 ? [] : store.childrenOf(cur.nodeKey).map((c) => toNode(c, depth + 1)),
      }
    }
    return store.nodeList
      .filter((n) => !n.parentId)
      .sort((a, b) => (a.sortIndex || 0) - (b.sortIndex || 0))
      .map((n) => toNode(n, 0))
  }

  async function aiDescribe() {
    if (aiLoading.value) return
    const imageUrl = store.image.url
    if (!imageUrl) {
      ElMessage.warning('当前主图缺少图片地址，无法生成描述')
      return
    }
    if (!store.nodeList.length) {
      ElMessage.warning('还没有标注区域，请先框选')
      return
    }
    aiLoading.value = true
    try {
      // 整页层级结构：所有顶层区域的完整子树
      const tree = buildRegionTree(store.selectedNode)
      const res = await uiProjectApi.describeRegion({ imageUrl, tree })
      // 响应树与请求同构，每个节点 summary 已由后端填入 AI 描述
      const treeRes = res?.code === 200 ? res?.data?.tree || [] : []
      if (!treeRes.length) throw new Error(res?.message || 'AI 未返回区域描述')
      let total = 0
      let filled = 0
      const walk = (nodes) => nodes.forEach((n) => {
        total++
        const k = n?.nodeKey
        const desc = n?.summary
        // 只回填仍存在且描述非空的节点（生成期间节点可能被删）
        if (k && desc && store.nodes[k]) {
          store.updateNode(k, { summary: desc })
          filled++
        }
        walk(n.children || [])
      })
      walk(treeRes)
      if (!filled) throw new Error('AI 返回的节点均无法匹配，请重试')
      if (filled < total) ElMessage.warning(`已回填 ${filled}/${total} 个区域描述`)
      else ElMessage.success(`已生成 ${filled} 个区域描述`)
    } catch (e) {
      ElMessage.error(e?.message || 'AI 描述生成失败')
    } finally {
      aiLoading.value = false
    }
  }

  return { aiLoading, aiDescribe }
}
