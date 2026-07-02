import { ElMessage } from 'element-plus'
import { taskApi } from '~/api/task'
import { useUserStore } from '~/stores/user'

export const createTaskActions = (ctx) => {
  const { state, runningCount, maxConcurrent, } = ctx

  const targetMap = { Flutter: 'flutter', React: 'react', Vue: 'vue3' }

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

    if (state.isSubmitting) return null

    const payload = {
      target,
      platform: state.config.platform,
      options: state.config.options,
      advanced: state.config.advanced,
      component_lib: state.config.options.includes('component') ? state.config.componentLib : '',
      images: state.images.map((img, i) => ({
        url: img.cosUrl,
        desc: img.description || '',
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
