import { ElMessage } from 'element-plus'

// 后端 target（小写）→ 工作台框架字典值（提交时 targetMap 的逆映射）
const TARGET_LABEL = { flutter: 'Flutter', react: 'React', vue3: 'Vue' }

// page_links.trigger → 关系类型 kind（与 relation.js 的 KIND_TRIGGER 互逆）
const TRIGGER_KIND = { node_click: 'route', bottom_tab: 'tab', back: 'back', auto: 'auto' }

// 多图画板横向排布：固定卡宽 300 + 间距 48，恢复后不叠在一起
const BOARD_GAP_X = 348
const boardPos = (i) => ({ x: 24 + i * BOARD_GAP_X, y: 24 })

/**
 * 「重新生成」数据恢复：把任务详情（TaskStatusResponse）回填到工作台。
 *
 * 设计要点：
 * - 图片以 cosUrl 作为恢复后的预览地址，不依赖本地 File；
 * - 树挂到 annotationTree（工作台 watch 会 loadTree 到画布/层级树）；
 * - 素材按 attached_assets 归属到各页，assetId 与树的 assetId 引用保持一致；
 * - 页面关系把后端边表还原成关系板数据（back 边保留，提交时会去重）；
 * - 整桶替换目标模式，避免与残留的手工内容混排。
 *
 * 返回恢复后的模式（single | multi），无图时返回空串。
 */
export const createRestoreActions = (ctx) => {
  const { state, modeMax } = ctx

  const restoreFromTask = (task) => {
    if (!task) return ''
    const src = Array.isArray(task.images) ? task.images : []
    if (!src.length) {
      ElMessage.warning('这个任务没有可恢复的设计图')
      return ''
    }

    const mode = src.length > 1 ? 'multi' : 'single'
    const bucket = state.modes[mode]

    // 1) 图片：cosUrl 即恢复后的预览地址；tree 原样挂到 annotationTree
    const nextImages = src.map((it, i) => {
      const pageId = it.page_id || `p${i + 1}`
      const W = it.width || 0
      const H = it.height || 0
      return {
        uploadId: `restore_${pageId}_${Date.now()}_${i}`,
        file: null,
        preview: it.url,
        naturalWidth: W,
        naturalHeight: H,
        // 提交时优先 treeWidth（见 task.js），恢复时用原图尺寸兜底
        treeWidth: W,
        treeHeight: H,
        cosUrl: it.url,
        uploading: false,
        uploadError: '',
        description: it.page_desc || '',
        annotations: [],
        pageId,
        pageName: it.page_name || it.title || '',
        route: it.route || '',
        parentPageId: '',
        linkTrigger: 'node_click',
        board: mode === 'multi' ? boardPos(i) : null,
        annotationTree: it.tree || [],
      }
    })

    // 2) 素材：按 attached_assets 归属到各页（树节点 assetId 引用 asset_id）
    const nextAssets = []
    src.forEach((it, i) => {
      const ownerId = nextImages[i].uploadId
      ;(it.attached_assets || []).forEach((a, j) => {
        if (!a || !a.url) return
        nextAssets.push({
          assetId: a.asset_id || `asset_${i}_${j}`,
          ownerId,
          file: null,
          preview: a.url,
          naturalWidth: a.width || 0,
          naturalHeight: a.height || 0,
          label: a.label || `素材${j + 1}`,
          size: 0,
          mime: '',
          cosUrl: a.url,
          uploading: false,
          uploadError: '',
        })
      })
    })

    // 3) 页面关系：后端边表 → 关系板数据（悬空/自环边丢弃）
    const links = Array.isArray(task.page_links) ? task.page_links : []
    const pageIdSet = new Set(nextImages.map((i) => i.pageId))
    const nextLinks = links
      .filter((l) => pageIdSet.has(l.from_page_id) && pageIdSet.has(l.to_page_id) && l.from_page_id !== l.to_page_id)
      .map((l, i) => ({
        id: `restore_link_${Date.now()}_${i}`,
        fromPageId: l.from_page_id,
        toPageId: l.to_page_id,
        kind: TRIGGER_KIND[l.trigger] || 'route',
        trigger: l.trigger || 'node_click',
        nodeKey: l.node_key || '',
        label: l.label || '',
        desc: '',
      }))

    // 4) 落库：整桶替换目标模式
    bucket.images.splice(0, bucket.images.length, ...nextImages)
    bucket.assets.splice(0, bucket.assets.length, ...nextAssets)
    bucket.pageLinks.splice(0, bucket.pageLinks.length, ...nextLinks)

    // 5) 模式与配置回填（字典值必须与提交 targetMap 一致：Vue3 → Vue）
    state.mode = mode
    state.maxImages = modeMax[mode]
    state.activePageId = nextImages[0]?.pageId || ''
    state.isSubmitting = false
    state.config.framework = TARGET_LABEL[task.target] || 'Flutter'
    state.config.platform = task.platform || 'mobile'
    state.config.options = Array.isArray(task.options) && task.options.length
      ? [...task.options]
      : ['responsive']
    state.config.advanced = Array.isArray(task.advanced) ? [...task.advanced] : []
    state.config.componentLib = task.component_lib || ''

    return mode
  }

  return { restoreFromTask }
}
