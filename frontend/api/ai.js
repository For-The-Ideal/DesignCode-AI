/**
 * SSE 长连接——走 Nuxt 本地代理，统一日志和地址
 * /api/sse/:taskId → Go /api/v1/task/:taskId/events
 */
export const aiApi = {
  async connect(taskId, signal = null) {
    return fetch(`/api/sse/${taskId}`, {
      method: 'GET',
      headers: { 'Accept': 'text/event-stream' },
      signal,
    })
  },
}
