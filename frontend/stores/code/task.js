import { nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import { storeToRefs } from 'pinia'
import { taskApi } from '~/api/task'
import { useUserStore } from '~/stores/user'
import { checkAnnotations } from '~/utils/annotationCheck'

export const createTaskActions = (ctx) => {
  const { state, images, assets, pageLinks, ensurePageLinks, runningCount, maxConcurrent, } = ctx

  const targetMap = { Flutter: 'flutter', React: 'react', Vue: 'vue3' }

  // ═══ 创建任务 ═══
  const createTask = async () => {
    // 先把排队的响应式变更刷完（工作台编辑会通过 watch 写回 annotationTree），再校验结构
    await nextTick()
    const target = targetMap[state.config.framework] || 'flutter'

    if (runningCount.value >= maxConcurrent) {
      ElMessage.warning(`最多同时运行 ${maxConcurrent} 个任务，请等待其中某个完成后再试`)
      return null
    }

    if (!images.value.length) {
      ElMessage.warning('请先上传图片')
      return null
    }

    // 结构完整性：与 AI 描述共用 checkAnnotations（子级缺名 / 顶层未选类型）
    const checks = images.value.map((img) => checkAnnotations(img.annotationTree || []))
    // 缺名点名到页：多图时只说总数，用户不知道去哪一页补
    const unnamedPages = checks
      .map((c, i) => ({ c, i }))
      .filter((x) => x.c.unnamedChildren.length)
    if (unnamedPages.length) {
      // 页面只报序号不报页名，避免把文件名/品牌词带进提示
      const total = unnamedPages.reduce((sum, x) => sum + x.c.unnamedChildren.length, 0)
      const pages = unnamedPages.map((x) => `第 ${x.i + 1} 页`).join('、')
      ElMessage.warning(`还有 ${total} 个子区域没填区域名（${pages}）：进入对应页面标注，树里标红的「缺名」框，在右侧「区域名」补上后再生成`)
      return null
    }
    const missingType = checks.reduce((sum, c) => sum + c.untypedTop.length, 0)
    if (missingType) {
      ElMessage.warning(`有 ${missingType} 个区域还没选类型，请进入标注补全后再生成`)
      return null
    }

    // 所有图片必须上传成功（cosUrl 就绪），否则会提交空 url
    if (!images.value.every(img => !!img.cosUrl)) {
      ElMessage.warning('有图片尚未上传完成或上传失败，请稍候或重试上传')
      return null
    }

    // 附图必须就绪且可归属：上传中/失败会丢素材，无归属会丢关联
    if (assets.value.some(a => a.uploading)) {
      ElMessage.warning('有附图尚未上传完成，请稍候')
      return null
    }
    if (assets.value.some(a => a.uploadError || !a.cosUrl)) {
      ElMessage.warning('有附图上传失败，请重试或删除后再提交')
      return null
    }
    const dangling = assets.value.filter(a => !images.value.some(i => i.uploadId === a.ownerId))
    if (dangling.length) {
      ElMessage.warning(`有 ${dangling.length} 张附图未归属任何主图，请在素材上重新选择归属页`)
      return null
    }

    if (state.isSubmitting) return null

    // 页面关系 → 提交边表：以关系板维护的显式边为准（旧 parentPageId 先迁移）；
    // 每条显式边自动补一条反向 back 边（已有显式反向边时不重复），无效引用丢弃。
    ensurePageLinks()
    const pageIdSet = new Set(images.value.map(i => i.pageId))
    // 每页现存节点的 nodeKey 集合：源组件被删掉后关系降级为整页触发，避免提交悬空引用
    const nodeKeysOf = {}
    images.value.forEach((img) => {
      const keys = new Set()
      const walk = (list) => (list || []).forEach((n) => {
        if (n.nodeKey) keys.add(n.nodeKey)
        walk(n.children)
      })
      walk(img.annotationTree)
      nodeKeysOf[img.pageId] = keys
    })
    const validLinks = pageLinks.value.filter(l =>
      pageIdSet.has(l.fromPageId) && pageIdSet.has(l.toPageId) && l.fromPageId !== l.toPageId)
    const submittedLinks = []
    validLinks.forEach((l) => {
      submittedLinks.push({
        from_page_id: l.fromPageId,
        to_page_id: l.toPageId,
        trigger: l.trigger || 'node_click',
        // 触发组件：空 = 整页触发
        node_key: l.nodeKey && nodeKeysOf[l.fromPageId]?.has(l.nodeKey) ? l.nodeKey : undefined,
        // 说明优先用组件名（关系板写入的 label），其次用户填写的说明
        label: l.label || l.desc || undefined,
      })
      const hasBack = validLinks.some(x => x.fromPageId === l.toPageId && x.toPageId === l.fromPageId)
      if (!hasBack) submittedLinks.push({ from_page_id: l.toPageId, to_page_id: l.fromPageId, trigger: 'back' })
    })

    const payload = {
      target,
      platform: state.config.platform,
      options: state.config.options,
      advanced: state.config.advanced,
      component_lib: state.config.options.includes('component') ? state.config.componentLib : '',
      images: images.value.map((img, i) => ({
        url: img.cosUrl,
        // 每图简短标题：优先文件名去扩展名，回退顶层区域名（供列表/详情回显，不需携带整棵 tree）
        title: (img.file?.name || '').replace(/\.[^.]+$/, '') || img.annotationTree?.[0]?.name || undefined,
        // 结构编辑器确认后的树（含坐标/描述/层级），后端据此跳过 Vision 直接生成
        tree: img.annotationTree || undefined,
        // 该主图归属的附图素材（按上传先后编号；标注框通过 asset_id 引用）
        attached_assets: (() => {
          const list = assets.value
            .filter(a => a.ownerId === img.uploadId && a.cosUrl)
            .map((a, j) => ({
              asset_id: a.assetId,
              url: a.cosUrl,
              label: a.label,
              width: a.naturalWidth || undefined,
              height: a.naturalHeight || undefined,
              sort_order: j + 1,
            }))
          return list.length ? list : undefined
        })(),
        // 原图像素尺寸（naturalWidth 刷新后丢失，回退持久化的 treeWidth/Height）
        width: img.treeWidth || img.naturalWidth || undefined,
        height: img.treeHeight || img.naturalHeight || undefined,
        sort_order: i + 1,
        // ── 页面身份（单图同样提交，后端按一页处理）──
        page_id: img.pageId || undefined,
        page_name: img.pageName || (img.file?.name || '').replace(/\.[^.]+$/, '') || undefined,
        page_desc: img.description || undefined,
        route: img.route || undefined,
        // 入口页恒为列表第 1 张：排序后不需要再同步任何标记
        is_entry: i === 0 || undefined,
      })),
      // 页面跳转边（多图模式；单图无关系时省略）
      page_links: submittedLinks.length ? submittedLinks : undefined,
    }

    state.isSubmitting = true

    try {
      const result = await taskApi.taskCreate(payload)
      if (!result || result.code !== 200 || !result.data?.task_id) {
        ElMessage.error(result?.message || 'AI 启动失败，请稍后重试')
        state.isSubmitting = false
        return null
      }

      const userStore = useUserStore()
      const {userInfo} = storeToRefs(userStore)
      console.log('[code store] 创建任务成功:', userInfo.value)
      await userStore.setUserInfo({
        ...userInfo.value,
        credits: userInfo.value.credits - payload.images.length,
        credits_used: userInfo.value.credits_used + payload.images.length,
      })

      ElMessage.success('任务已创建，请在任务列表中查看进度')
      state.isSubmitting = false

      // 清空当前模式的主图 + 附图 + 页面关系（另一模式的数据保留）
      state.modes[state.mode].images.splice(0)
      state.modes[state.mode].assets.splice(0)
      state.modes[state.mode].pageLinks.splice(0)
      state.activePageId = ''
      state.config = {
        framework: 'Flutter',
        platform: 'mobile',
        options: ['responsive'],
        advanced: [],
        componentLib: '',
      }

      return result.data.task_id
    } catch (error) {
      console.error('[code store] 创建任务失败:', error)
      ElMessage.error(error.message || '生成失败，请稍后重试')
      state.isSubmitting = false
      return null
    }
  }

  return { createTask }
}
