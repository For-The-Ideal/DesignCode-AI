const BASE_URL = 'http://localhost:8888'

export const sseApi = {
  /**
   * 建立 SSE 长连接（订阅 Broker，不传业务参数）
   */
  async connect (signal = null) {
    return fetch(`${BASE_URL}/api/sse`, {
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
