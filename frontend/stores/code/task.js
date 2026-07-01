import { commonApi } from '~/api/common'
import { ElMessage } from 'element-plus'

export const createTaskActions = (ctx) => {
  const { state, runningCount, maxConcurrent } = ctx

  /** 框架 → 后端 target 映射 */
  const targetMap = { Flutter: 'flutter', React: 'react', Vue: 'vue3' }

  /** 活跃轮询定时器 { [taskId]: intervalId } */
  const pollingTimers = {}

  // ═══ 服务端任务统计 ═══
  const fetchUserTasks = async () => {
    try {
      const res = await commonApi.getUserTasks()
      if (res.code !== 200 || !res.data) return

      state.taskCounts = {
        pending: res.data.pending || 0,
        running: res.data.running || 0,
        success: res.data.success || 0,
        failed: res.data.failed || 0,
      }
    } catch {
      // 静默
    }
  }

  // ═══ 创建任务 ═══
  const createTask = async () => {
    const target = targetMap[state.config.framework] || 'flutter'

    // 并发检查
    if (runningCount.value >= maxConcurrent) {
      ElMessage.warning(`最多同时运行 ${maxConcurrent} 个任务，请等待其中某个完成后再试`)
      return null
    }

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
      const result = await commonApi.generateUi(payload)
      if (!result || result.code !== 200 || !result.data?.task_id) {
        ElMessage.error(result?.message || 'AI 启动失败，请稍后重试')
        state.isSubmitting = false
        return null
      }

      const taskId = result.data.task_id
      state.tasks[taskId] = {
        id: taskId,
        target,
        status: 'pending',
        progress: 0,
        currentStep: '',
        result: null,
        images: [...state.images],
      }

      await fetchUserTasks()
      startPolling(taskId)

      return taskId
    } catch (error) {
      console.error('[code store] 创建任务失败:', error)
      ElMessage.error(error.message || '生成失败，请稍后重试')
      state.isSubmitting = false
      return null
    }
  }

  // ═══ 轮询单任务 ═══
  const pollTask = async (taskId) => {
    try {
      const res = await commonApi.getTaskById(taskId)
      if (!res || res.code !== 200 || !res.data) return

      const data = res.data
      const task = state.tasks[taskId]
      if (!task) return

      task.status = data.status
      task.progress = data.progress || 0
      task.currentStep = data.current_step || ''
      task.target = data.target || task.target
      task.images = data.images || task.images
      if (data.result) task.result = data.result

      if (data.status === 'success' || data.status === 'failed') {
        fetchUserTasks()
        onTaskDone(taskId, data.status)
      }

      return task
    } catch {
      // 网络错误静默
    }
  }

  // ═══ 内部：启动轮询 ═══
  const startPolling = (taskId) => {
    pollingTimers[taskId] = setInterval(async () => {
      const updated = await pollTask(taskId)
      // pollTask 内部已调用 onTaskDone，这里只处理异常退出
      if (!updated) {
        clearInterval(pollingTimers[taskId])
        delete pollingTimers[taskId]
        checkAllDone()
      }
    }, 2000)
  }

  // ═══ 内部：任务结束 ═══
  const onTaskDone = (taskId, status) => {
    clearInterval(pollingTimers[taskId])
    delete pollingTimers[taskId]

    if (status === 'success') ElMessage.success('任务完成')
    if (status === 'failed') ElMessage.error('任务失败')

    checkAllDone()
  }

  // ═══ 内部：检查是否全部结束 ═══
  const checkAllDone = () => {
    const hasActive = Object.values(state.tasks).some(
      t => t.status === 'pending' || t.status === 'running'
    )
    if (!hasActive) state.isSubmitting = false
  }

  // ═══ 停止所有轮询（清空时调用） ═══
  const stopAllPolling = () => {
    Object.keys(pollingTimers).forEach((id) => {
      clearInterval(pollingTimers[id])
    })
    for (const k of Object.keys(pollingTimers)) delete pollingTimers[k]
  }

  return { createTask, pollTask, fetchUserTasks, stopAllPolling }
}
