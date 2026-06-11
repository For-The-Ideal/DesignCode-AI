import { localUrl } from '~/config/index.js'

export const aiApi = {
  /**
   * 建立 SSE 长连接（按 taskID 订阅事件）
   */
  async connect (taskId, signal = null) {
    return fetch(`${localUrl}/api/v1/task/${taskId}/events`, {
      method: 'GET',
      headers: {
        'Accept': 'text/event-stream',
      },
      signal,
    })
  },
}
