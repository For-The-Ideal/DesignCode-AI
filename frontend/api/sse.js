/**
 * SSE 长连接——走 Nuxt 本地代理，统一日志和地址
 * /api/sse/:taskId → Go /api/v1/sse/:taskId/events
 * /api/sse/user → Go /api/v1/sse/user/events
 */
export const sseApi = {
  async connect(taskId, signal = null) {
    return fetch(`/api/sse/${taskId}`, {
      method: 'GET',
      headers: { 'Accept': 'text/event-stream' },
      signal,
    })
  },

  async connectUser(signal = null) {
    return fetch(`/api/sse/user`, {
      method: 'GET',
      headers: { 'Accept': 'text/event-stream' },
      signal,
    })
  },
}
