import { ElMessage } from 'element-plus'
import { taskApi } from '~/api/task'
import { useUserStore } from '~/stores/user'

export const createTaskActions = (ctx) => {
  const { state, runningCount, maxConcurrent, } = ctx

  const targetMap = { Flutter: 'flutter', React: 'react', Vue: 'vue3' }

  // 序列化描述：整体描述在前，标注框逐条输出 [标签 x%,y%,w%,h] 描述
  // 坐标为相对原图的 0~1 比例，后端 Vision 可直接换算像素与 DSLComponent.position 对齐
  const buildImageDesc = (img) => {
    const parts = []
    if (img.description) parts.push(img.description.trim())
    ;(img.annotations || []).forEach((a) => {
      const label = a.label ? a.label.trim() : '区域'
      const pos = [
        Math.round(a.x * 100),
        Math.round(a.y * 100),
        Math.round(a.w * 100),
        Math.round(a.h * 100),
      ].join(',') + '%'
      parts.push(`[${label} ${pos}] ${(a.desc || '').trim()}`.trim())
    })
    return parts.join('\n')
  }

  // ═══ 创建任务 ═══
  const createTask = async () => {
    const target = targetMap[state.config.framework] || 'flutter'

    if (runningCount.value >= maxConcurrent) {
      ElMessage.warning(`最多同时运行 ${maxConcurrent} 个任务，请等待其中某个完成后再试`)
      return null
    }

    if (!state.images.length) {
      ElMessage.warning('请先上传图片')
      return null
    }

    // 所有图片必须上传成功（cosUrl 就绪），否则会提交空 url
    if (!state.images.every(img => !!img.cosUrl)) {
      ElMessage.warning('有图片尚未上传完成或上传失败，请稍候或重试上传')
      return null
    }

    if (state.isSubmitting) return null

    const payload = {
      target,
      platform: state.config.platform,
      options: state.config.options,
      advanced: state.config.advanced,
      component_lib: state.config.options.includes('component') ? state.config.componentLib : '',
      images: state.images.map((img, i) => ({
        url: img.cosUrl,
        desc: buildImageDesc(img),
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

      // 清空图片 + 重置配置
      state.images = []
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
