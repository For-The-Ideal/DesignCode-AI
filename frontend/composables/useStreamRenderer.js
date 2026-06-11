import { ref, onUnmounted } from 'vue'

/**
 * useStreamRenderer — 纯流式渲染引擎
 *
 * 职责：接收完整文本 → 逐步分片输出到 displayText（打字机效果）
 * 与数据来源（SSE / API / 本地）无关，只负责"渲染输出"这一层
 *
 * 用法：
 *   const { displayText, start, pause, resume, stop } = useStreamRenderer()
 *   start(longText)   // 开始逐字输出
 *   pause()           // 暂停（IntersectionObserver 离开可视区）
 *   resume()          // 恢复
 *   stop()            // 停止并清空
 */

export function useStreamRenderer () {
  // ═══ 响应式状态 ═══
  const displayText = ref('')
  const isPaused = ref(false)
  const isActive = ref(false)

  // ═══ 内部变量 ═══
  let streamInterval = null
  let fullText = ''
  let currentIndex = 0
  const charsPerTick = 12
  const tickMs = 25

  // ═══ 私有方法 ═══

  const stopInterval = () => {
    if (streamInterval) {
      clearInterval(streamInterval)
      streamInterval = null
    }
  }

  const tick = (fromIndex) => {
    isPaused.value = false
    isActive.value = true
    let index = fromIndex

    streamInterval = setInterval(() => {
      if (index < fullText.length) {
        displayText.value += fullText.slice(index, index + charsPerTick)
        index += charsPerTick
        currentIndex = index
      } else {
        stopInterval()
        isActive.value = false
      }
    }, tickMs)
  }

  // ═══ 公开方法 ═══

  /** 开始流式渲染全新文本（会先清空 displayText） */
  function start (text) {
    stopInterval()
    displayText.value = ''
    fullText = text
    currentIndex = 0
    isPaused.value = false
    if (text) {
      tick(0)
    }
  }

  /** 暂停流式渲染（用于离开可视区） */
  function pause () {
    if (streamInterval) {
      stopInterval()
      isPaused.value = true
      isActive.value = false
    }
  }

  /** 恢复流式渲染（用于重新进入可视区） */
  function resume () {
    if (isPaused.value && currentIndex < fullText.length) {
      tick(currentIndex)
    }
  }

  /** 停止并清空所有状态 */
  function stop () {
    stopInterval()
    displayText.value = ''
    fullText = ''
    currentIndex = 0
    isPaused.value = false
    isActive.value = false
  }

  /** 立即追加文本到 displayText（不经过流式定时器） */
  function append (chunk) {
    displayText.value += chunk
  }

  /** 一次性输出全部剩余文本，跳过流式效果 */
  function flush () {
    stopInterval()
    displayText.value = fullText
    isActive.value = false
  }

  // ═══ 生命周期 ═══

  onUnmounted(() => {
    stopInterval()
  })

  return {
    displayText,
    isPaused,
    isActive,
    start,
    pause,
    resume,
    stop,
    append,
    flush,
  }
}
