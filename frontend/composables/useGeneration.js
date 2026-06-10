import { ref } from 'vue'
import { loginApi } from '~/api/login'
import { useSSE } from './useSSE'

/**
 * 代码生成流程管理 composable
 *
 * 职责：
 *   1. initTemplateData() — 调用接口加载模板 → 流式代码输出 + 预览 HTML
 *
 * SSE 实例、回调映射、连接守卫 全部收拢在 useSSE.js 中。
 */
export function useGeneration () {
  // ═══ 响应式状态 ═══
  const template = ref({
    templateCode: '',
    previewCode: '',
    id: 0,
  })

  const generating = ref(false)

  // ═══ 内部状态（模板初始化用） ═══
  let fullTemplateCode = ''
  let dataLoaded = false

  let streamInterval = null
  let currentIndex = 0
  let generationStarted = false
  let isPaused = false

  // ═══ SSE 客户端（立即创建，不连线） ═══
  const sse = useSSE({ template, generating })

  // ═══ 流式渲染（本地模拟分片输出） ═══

  const stopStreamInterval = () => {
    if (streamInterval) { clearInterval(streamInterval); streamInterval = null }
  }

  const startStreaming = (fromIndex) => {
    if (!fullTemplateCode) return
    isPaused = false
    let index = fromIndex
    streamInterval = setInterval(() => {
      if (index < fullTemplateCode.length) {
        template.value.templateCode += fullTemplateCode.slice(index, index + 12)
        index += 12
        currentIndex = index
      } else {
        stopStreamInterval()
        isPaused = false
      }
    }, 25)
  }

  const resumeStreaming = () => {
    if (isPaused && currentIndex < fullTemplateCode.length) {
      startStreaming(currentIndex)
    }
  }

  const pauseStreaming = () => {
    if (streamInterval) {
      stopStreamInterval()
      isPaused = true
    }
  }

  const startGeneration = () => {
    if (generationStarted && !isPaused) return
    if (!generationStarted) {
      generationStarted = true
      template.value.templateCode = ''
      currentIndex = 0
    }
    if (isPaused) {
      resumeStreaming()
    } else {
      startStreaming(0)
    }
  }

  // ═══ 模板初始化 ═══

  const initTemplateData = async () => {
    if (!dataLoaded) {
      const id = 1
      try {
        const { data } = await loginApi.template({ template: id })
        if (data.id) {
          template.value = {
            templateCode: data.template_code || '',
            previewCode: data.preview_code || '',
            id: data.id || id,
          }
          fullTemplateCode = template.value.templateCode
          dataLoaded = true
        }
      } catch {
        // 接口不可用时静默降级
      }
    }
    if (fullTemplateCode) {
      startGeneration()
    }
  }

  const setGenerating = (value) => {
    generating.value = value
  }

  const cleanup = () => {
    stopStreamInterval()
    sse.disconnect()
  }

  return {
    template,
    generating,
    sseStatus: sse.status,
    sseError: sse.error,
    sseRetry: sse.retry,
    isAvailable: sse.isAvailable,
    isAlive: sse.isAlive,
    connectSSE: sse.connect,
    initTemplateData,
    startGeneration,
    pauseStreaming,
    cleanup,
    setGenerating,
  }
}
