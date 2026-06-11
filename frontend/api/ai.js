import { localUrl } from '~/config/index.js'

export const aiApi = {
  /**
   * 建立 SSE 长连接（订阅 Broker，不传业务参数）
   */
  async connect (signal = null) {
    return fetch(`${localUrl}/api/v1/ai/sse`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'text/event-stream',
      },
      body: '{}',
      signal,
    })
  },
}
