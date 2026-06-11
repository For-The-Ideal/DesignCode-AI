import { reactive, watch, onUnmounted } from 'vue'
import { useSSE } from './useSSE'
import { useStreamRenderer } from './useStreamRenderer'
import { commonApi } from '~/api/common'

/**
 * useGeneration — 代码生成流程编排 composable
 *
 * 职责：
 *   1. 协调 SSE 数据源 → 渲染输出
 *   2. 模板数据初始化（API）→ 流式渲染
 *   3. 向外暴露统一的 template 接口
 *
 * 数据流：
 *   SSE 推送  ──→ useSSE.sseData ──→ watch ──→ template.templateCode
 *   API 模板  ──→ loginApi ──→ useStreamRenderer.start() ──→ watch ──→ template.templateCode
 *
 * 两个数据源通过 status 判断互斥：SSE streaming 时不处理渲染器输出
 */

export function useGeneration () {
  // ═══ 子系统 ═══
  const sse = useSSE()
  const renderer = useStreamRenderer()

  // ═══ 统一数据出口 ═══
  // 结构对齐 auth.go Template 接口: { template_code, preview_code, id }
  const template = reactive({
    templateCode: '',
    previewCode: '',
    id: null,
  })

  // 模板加载场景的辅助状态
  let templatePreviewCode = ''
  let dataLoaded = false

  // ═══ SSE 数据 → template =================================================

  // SSE 流式推送代码 → 直接写入 template（SSE 本身就是流式的）
  watch(() => sse.sseData.templateCode, (code) => {
    if (sse.status.value === 'streaming') {
      template.templateCode = code
    }
  })

  // SSE 推送预览 HTML → 写入 template
  watch(() => sse.sseData.previewCode, (code) => {
    if (code) {
      template.previewCode = code
    }
  })

  // SSE 推送 id → 写入 template（对齐 auth.go 模板接口）
  watch(() => sse.sseData.id, (val) => {
    if (val != null) {
      template.id = val
    }
  })

  // 渲染器输出 → template（仅模板加载场景，SSE 不活跃时）
  watch(() => renderer.displayText.value, (text) => {
    if (sse.status.value !== 'streaming') {
      template.templateCode = text
    }
  })

  // ═══ 模板数据初始化（首页 SSEGenerator 用）══════════════════════════

  /**
   * 加载模板 → 流式渲染
   *
   * 首次调用：请求 API 获取模板数据 → renderer.start(code) → 逐字输出
   * 再次进入可视区：renderer.resume() 继续
   * 已输出完毕：无操作
   */
  const initTemplateData = async (id = 1) => {
    // SSE 活跃时跳过（避免覆盖实时生成的数据）
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
      // 已加载但暂停中（离开可视区后再次进入）
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

    // ── 模板加载（SSEGenerator 用）──
    initTemplateData,
    pauseStreaming: renderer.pause,
    resumeStreaming: renderer.resume,

    // ── SSE（透传给 code.vue）──
    sseStatus: sse.status,
    sseError: sse.error,
    sseData: sse.sseData,
    connectSSE: sse.connect,
    disconnectSSE: sse.disconnect,
    isAvailable: sse.isAvailable,
    isAlive: sse.isAlive,

    // ── 生命周期 ──
    cleanup,
  }
}
