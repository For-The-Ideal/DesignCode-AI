import { ElMessage } from 'element-plus'
import { storeToRefs } from 'pinia'
import { taskApi } from '~/api/task'
import { useUserStore } from '~/stores/user'

export const createTaskActions = (ctx) => {
  const { state, images, assets, runningCount, maxConcurrent, } = ctx

  const targetMap = { Flutter: 'flutter', React: 'react', Vue: 'vue3' }

  // ═══ 创建任务 ═══
  const createTask = async () => {
    const target = targetMap[state.config.framework] || 'flutter'

    if (runningCount.value >= maxConcurrent) {
      ElMessage.warning(`最多同时运行 ${maxConcurrent} 个任务，请等待其中某个完成后再试`)
      return null
    }

    if (!images.value.length) {
      ElMessage.warning('请先上传图片')
      return null
    }

    // 结构完整性：顶层区域必须已选类型（标注允许带待补保存，生成前必须补齐）
    const missingType = images.value.reduce(
      (sum, img) => sum + (img.annotationTree || []).filter((n) => !n.name).length, 0)
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
      })),
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

      // 清空当前模式的主图 + 附图（另一模式的数据保留）
      state.modes[state.mode].images.splice(0)
      state.modes[state.mode].assets.splice(0)
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
