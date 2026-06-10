import { ref, readonly } from 'vue'

// ═══ 模块级共享状态（所有组件共享同一份） ═══
const generating = ref(false)
const progress = ref(0)

let timer = null

/**
 * 代码生成流程管理 composable
 *
 * 职责：集中管理 "生成中" 状态 + 进度，支持模拟定时器与 SSE 注入
 *
 * 用法：
 *   const { generating, progress, start, setProgress, complete, reset } = useGeneration()
 *
 * 后期 SSE 接入：SSE 回调中直接调用 setProgress(val)，替代模拟定时器
 */
export function useGeneration() {
  /** 开始生成（含模拟进度，SSE 接入后可移除） */
  function start() {
    generating.value = true
    progress.value = 0
    if (timer) clearInterval(timer)

    let current = 0
    const stepDuration = 8000 / 100 // 8 秒走完 100%
    timer = setInterval(() => {
      current++
      if (current <= 100) {
        progress.value = current
      } else {
        stopTimer()
      }
    }, stepDuration)
  }

  /** SSE / 外部手动设置进度 0-100 */
  function setProgress(val) {
    progress.value = Math.min(Math.max(val, 0), 100)
  }

  /** 生成完成 */
  function complete() {
    stopTimer()
    progress.value = 100
    setTimeout(() => {
      generating.value = false
      progress.value = 0
    }, 600)
  }

  /** 重置（失败 / 取消） */
  function reset() {
    stopTimer()
    generating.value = false
    progress.value = 0
  }

  function stopTimer() {
    if (timer) { clearInterval(timer); timer = null }
  }

  return {
    generating: readonly(generating),
    progress: readonly(progress),
    start,
    setProgress,
    complete,
    reset,
  }
}
