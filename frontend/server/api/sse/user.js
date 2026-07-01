/**
 * 用户级 SSE 代理端点
 * 前端请求 /api/sse/user → 代理到 Go 后端 /api/v1/task/sse
 * 流式透传，用于接收用户所有任务的状态变更
 */
import { createError } from 'h3'
import { serverLog } from '../../utils/logger.js'

export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig()
  const goUrl = `${config.public.apiBase}/api/v1/task/sse`
  const startTime = Date.now()

  serverLog('REQ', { method: 'GET', url: '/api/sse/user', isBrowser: false })

  try {
    const response = await fetch(goUrl, {
      method: 'POST',
      headers: {
        'Accept': 'text/event-stream',
        ...(event.node.req.headers.cookie ? { cookie: event.node.req.headers.cookie } : {}),
      },
    })

    if (!response.ok) {
      serverLog('ERR', { method: 'GET', url: '/api/sse/user', status: response.status, message: 'Go 后端连接失败', duration: Date.now() - startTime })
      throw createError({ statusCode: response.status, message: 'SSE 连接失败' })
    }

    setResponseHeaders(event, {
      'Content-Type': 'text/event-stream',
      'Cache-Control': 'no-cache',
      'Connection': 'keep-alive',
    })

    const stream = new ReadableStream({
      async start(controller) {
        const reader = response.body.getReader()
        try {
          while (true) {
            const { done, value } = await reader.read()
            if (done) break
            controller.enqueue(value)
          }
        } catch (err) {
          serverLog('ERR', { method: 'GET', url: '/api/sse/user', message: `流中断: ${err.message}` })
        } finally {
          controller.close()
        }
      },
    })

    return stream
  } catch (err) {
    serverLog('ERR', { method: 'GET', url: '/api/sse/user', status: 502, message: err.message, duration: Date.now() - startTime })
    throw createError({ statusCode: 502, message: 'SSE 代理异常' })
  }
})
