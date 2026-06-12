import { reactive, ref, computed, watch, onUnmounted } from 'vue'
import { useSSE } from './useSSE'
import { useStreamRenderer } from './useStreamRenderer'
import { commonApi } from '~/api/common'

/**
 * useGeneration — 代码生成流程编排 composable
 *
 * 职责：
 *   1. 任务生命周期管理（创建、恢复、清理）
 *   2. 协调 SSE 数据源 → 渲染输出
 *   3. 模板数据初始化（API）→ 流式渲染
 *   4. 向外暴露统一的 template 接口
 *
 * 数据流：
 *   SSE 推送  ──→ useSSE.sseData ──→ watch ──→ template.templateCode
 *   API 模板  ──→ loginApi ──→ useStreamRenderer.start() ──→ watch ──→ template.templateCode
 *
 * 两个数据源通过 status 判断互斥：SSE streaming 时不处理渲染器输出
 */

// ═══ 模块级任务状态 ═══
const activeTaskId = ref(null)
const activeTaskFramework = ref('')   // flutter | vue3 | react
const taskStatus = ref('idle')        // idle | pending | running | success | failed
const taskProgress = ref(0)
const taskCurrentStep = ref('')
const taskErrorMsg = ref('')
const RESTORE_KEY = 'gen_active_task'

const isBusy = computed(() =>
  taskStatus.value === 'pending' || taskStatus.value === 'running'
)

export function useGeneration () {
  // ═══ 子系统 ═══
  const sse = useSSE()
  const renderer = useStreamRenderer()

  // ═══ 统一数据出口 ═══
  const template = reactive({
    templateCode: '',
    previewCode: '',
    id: null,
  })

  // 模板加载场景的辅助状态
  let templatePreviewCode = ''
  let dataLoaded = false

  // ═══ SSE 数据 → template =================================================

  watch(() => sse.sseData.templateCode, (code) => {
    if (sse.status.value === 'streaming') {
      template.templateCode = code
    }
  })

  watch(() => sse.sseData.previewCode, (code) => {
    if (code) {
      template.previewCode = code
    }
  })

  watch(() => sse.sseData.id, (val) => {
    if (val != null) {
      template.id = val
    }
  })

  // ═══ SSE progress → 任务进度 ═══
  watch(() => sse.sseData.progress, (val) => {
    taskProgress.value = val
  })

  watch(() => sse.sseData.currentStep, (val) => {
    taskCurrentStep.value = val
  })

  // ═══ SSE 状态 → 任务状态 ═══
  watch(sse.status, (val) => {
    if (val === 'idle' && taskStatus.value === 'running') {
      // SSE 流结束，任务完成
      taskStatus.value = 'success'
      taskProgress.value = 100
      localStorage.removeItem(RESTORE_KEY)
    }
  })

  // 渲染器输出 → template（仅模板加载场景，SSE 不活跃时）
  watch(() => renderer.displayText.value, (text) => {
    if (sse.status.value !== 'streaming') {
      template.templateCode = text
    }
  })

  // ═══ 任务生命周期 ═════════════════════════════════════

  /**
   * 保存当前任务（生成成功后调用）
   * 写入 localStorage 以便页面刷新后恢复
   */
  const saveActiveTask = (taskId, framework) => {
    activeTaskId.value = taskId
    activeTaskFramework.value = framework
    taskStatus.value = 'running'
    // 清空上一次残留数据
    template.templateCode = ''
    template.previewCode = ''
    template.id = null
    localStorage.setItem(RESTORE_KEY, JSON.stringify({ taskId, framework }))
  }

  /**
   * 清除当前任务（用户手动清除或任务结束）
   */
  const clearActiveTask = () => {
    activeTaskId.value = null
    activeTaskFramework.value = ''
    taskStatus.value = 'idle'
    taskProgress.value = 0
    taskCurrentStep.value = ''
    taskErrorMsg.value = ''
    sse.disconnect()
    localStorage.removeItem(RESTORE_KEY)
  }

  /**
   * 从 localStorage 恢复任务状态
   * 页面挂载时调用，用于用户刷新后恢复进度
   */
  const restoreTask = async () => {
    const saved = localStorage.getItem(RESTORE_KEY)
    if (!saved) return null

    try {
      const { taskId, framework } = JSON.parse(saved)
      activeTaskId.value = taskId
      activeTaskFramework.value = framework || ''

      const res = await commonApi.getTaskById(taskId)
      if (!res || !res.data) {
        localStorage.removeItem(RESTORE_KEY)
        return null
      }

      const data = res.data
      taskProgress.value = data.progress || 0
      taskCurrentStep.value = data.current_step || ''
      activeTaskFramework.value = data.target || framework || ''

      if (data.status === 'pending' || data.status === 'running') {
        // 任务还在执行 → 重连 SSE
        taskStatus.value = data.status
        sse.connect(taskId)
        return { taskId, status: data.status, framework: activeTaskFramework.value, images: data.images || [] }

      } else if (data.status === 'success' && data.result) {
        // 任务已完成 → 展示结果
        taskStatus.value = 'success'
        template.templateCode = data.result.code || ''
        template.previewCode = data.result.preview || ''
        template.id = data.result.id || null
        return { taskId, status: 'success', framework: activeTaskFramework.value, images: data.images || [] }

      } else if (data.status === 'failed') {
        taskStatus.value = 'failed'
        taskErrorMsg.value = '任务执行失败'
        return { taskId, status: 'failed', framework: activeTaskFramework.value, images: data.images || [] }
      }
    } catch {
      localStorage.removeItem(RESTORE_KEY)
    }
    return null
  }

  // ═══ 模板数据初始化 ═════════════════════════════════════

  const initTemplateData = async (id = 1) => {
    if (sse.status.value === 'streaming') return

    if (!dataLoaded) {
      try {
        const { data } = await commonApi.getTemplate({ id: id })
        if (data.id) {
          template.id = data.id
          templatePreviewCode = data.preview_code || ''
          template.previewCode = templatePreviewCode
          if (data.template_code) {
            renderer.start(data.template_code)
          }
          dataLoaded = true
        }
      } catch {
        // 接口不可用时静默降级
      }
    } else if (renderer.isPaused.value) {
      renderer.resume()
    } else if (!renderer.isActive.value) {
      // 流式已结束，无需操作
    }
  }

  // ═══ 生命周期 ═══

  const cleanup = () => {
    renderer.stop()
  }

  onUnmounted(() => {
    cleanup()
  })

  // ═══ 对外接口 ═══

  return {
    // ── 统一数据出口 ──
    template,

    // ── 任务生命周期 ──
    activeTaskId,
    taskStatus,
    taskProgress,
    taskCurrentStep,
    taskErrorMsg,
    isBusy,
    saveActiveTask,
    clearActiveTask,
    restoreTask,

    // ── 模板加载（SSEGenerator 用）──
    initTemplateData,
    pauseStreaming: renderer.pause,
    resumeStreaming: renderer.resume,

    // ── SSE（透传给 code.vue）──
    sseStatus: sse.status,
    sseError: sse.error,
    sseData: sse.sseData,
    connectSSE: (taskId) => sse.connect(taskId),
    disconnectSSE: sse.disconnect,
    isAvailable: sse.isAvailable,
    isAlive: sse.isAlive,

    // ── 生命周期 ──
    cleanup,
  }
}
