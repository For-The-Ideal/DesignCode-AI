/**
 * SSE 代理端点
 * 前端请求 /api/sse/:taskId → 代理到 Go 后端 /api/v1/sse/:taskId/events
 * 流式透传，不经过加密解密
 */
import { createError } from 'h3'
import { serverLog } from '../../utils/logger.js'

export default defineEventHandler(async (event) => {
  const taskId = getRouterParam(event, 'taskId')
  if (!taskId) {
    throw createError({ statusCode: 400, message: '缺少 taskId' })
  }
  const requestUrl = useRuntimeConfig()
  console.log(requestUrl,'requestUrl--requestUrl')
  const goUrl = `${requestUrl.public.apiBase}/api/v1/sse/${taskId}/events`
  const startTime = Date.now()
  console.log(goUrl,'g')
  serverLog('REQ', { method: 'GET', url: `/api/sse/${taskId}`, params: { taskId }, isBrowser: false })

  try {
    const response = await fetch(goUrl, {
      headers: {
        'Accept': 'text/event-stream',
        ...(event.node.req.headers.cookie ? { cookie: event.node.req.headers.cookie } : {}),
      },
    })

    if (!response.ok) {
      serverLog('ERR', { method: 'GET', url: `/api/sse/${taskId}`, status: response.status, message: 'Go 后端连接失败', duration: Date.now() - startTime })
      throw createError({ statusCode: response.status, message: 'SSE 连接失败' })
    }

    // 设置 SSE 响应头，流式透传
    setResponseHeaders(event, {
      'Content-Type': 'text/event-stream',
      'Cache-Control': 'no-cache',
      'Connection': 'keep-alive',
    })

    // 流式转发
    const reader = response.body.getReader()
    const stream = new ReadableStream({
      async start(controller) {
        try {
          while (true) {
            const { done, value } = await reader.read()
            if (done) break
            controller.enqueue(value)
          }
        } catch (err) {
          serverLog('ERR', { method: 'GET', url: `/api/sse/${taskId}`, message: `流中断: ${err.message}` })
        } finally {
          controller.close()
        }
      },
    })

    return stream
  } catch (err) {
    serverLog('ERR', { method: 'GET', url: `/api/sse/${taskId}`, status: 502, message: err.message, duration: Date.now() - startTime })
    throw createError({ statusCode: 502, message: 'SSE 代理异常' })
  }
})
