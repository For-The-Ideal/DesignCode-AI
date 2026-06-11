import { ref, reactive } from 'vue'
import { aiApi } from '~/api/ai'
import { ElMessage } from 'element-plus'

/**
 * useSSE — 代码生成 SSE 客户端 composable（模块级单例）
 *
 * status 值	含义	页面可做的事
  'idle'	空闲	显示"连接 SSE"按钮
  'connecting'	连接中	显示加载态
  'streaming'	SSE响应中 ✅	显示流式进度
  'done'	SSE响应完成 ✅	显示完成提示
  'error'	单次连接失败（会自动重试）	显示"重连中..."
  'maxRetries'	最大次数失败 ✅	显示"连接失败，请手动重试"

 */

// ═══ 模块级共享状态 ═══
const status = ref('idle')
const error = ref('')
const retry = ref(0)

const sseData = reactive({
  templateCode: '',
  previewCode: '',
  id: null,        // 对齐 auth.go 模板接口 { template_code, preview_code, id }
  score: 0,
  dimensions: [],
})

let abortController = null
let reader = null
let active = false
const maxRetries = 3
const retryDelay = 1500
const maxRetryDelay = 30000
const timeout = 120000

export function useSSE () {

  // ═══ 日志 ═══
  const log = (msg, ...args) => {
    const ts = new Date().toLocaleTimeString()
    console.log(`[SSE ${ts}] ${msg}`, ...args)
  }

  // ═══ SSE 帧解析 ═══
  const parseSSEChunk = (chunk) => {
    const lines = chunk.split('\n')
    let eventType = 'message'
    let dataLines = []

    for (const line of lines) {
      if (line.startsWith('event: ')) {
        eventType = line.slice(7).trim()
      } else if (line.startsWith('data: ')) {
        dataLines.push(line.slice(6))
      } else if (line === '') {
        if (dataLines.length > 0) {
          const payload = dataLines.join('\n')
          log(`<- event:${eventType}`, payload.length > 80 ? payload.slice(0, 80) + '...' : payload)
          handleEvent(eventType, payload)
        }
        eventType = 'message'
        dataLines = []
      }
    }
  }

  // ═══ 事件分发 ═══
  function handleEvent (type, data) {
    switch (type) {
      case 'connected':
        status.value = 'streaming'
        log('SSE 连接已确认，等待 Broker 推送...')
        break
      case 'message':
        sseData.templateCode += data
        break
      case 'preview':
        sseData.previewCode = data
        log('预览 HTML 已更新')
        break
      case 'score':
        try {
          const parsed = JSON.parse(data)
          sseData.score = parsed.score
          sseData.dimensions = parsed.dimensions
          log('评分数据已接收', parsed)
        } catch { log('评分解析失败', data) }
        break
      case 'done':
        log('SSE 流传输完成，共接收', sseData.templateCode.length, '个字符')
        try {
          const doneData = JSON.parse(data)
          sseData.id = doneData.id || null
        } catch { /* data 非 JSON 时忽略 */ }
        status.value = 'idle'
        break
      case 'error':
        log('服务端错误', data)
        break
    }
  }

  // ═══ 建立 SSE 长连接（按 taskID 订阅） ═══
  const connect = async (taskId) => {
    if (status.value === 'connecting' || status.value === 'streaming') {
      log('已有活跃连接，跳过:', status.value)
      return
    }

    log('========== 建立 SSE 连接（订阅 Broker）==========')

    active = true
    status.value = 'connecting'
    error.value = ''
    retry.value = 0
    Object.assign(sseData, { templateCode: '', previewCode: '', id: null, score: 0, dimensions: [] })
    abortController = new AbortController()

    const timeoutId = setTimeout(() => {
      log('连接超时，中止')
      abortController.abort()
    }, timeout)

    try {
      const response = await aiApi.connect(taskId, abortController.signal)
      clearTimeout(timeoutId)

      if (!response.ok) {
        log('HTTP 错误:', response.status, response.statusText)
        throw new Error(`SSE 连接失败: HTTP ${response.status}`)
      }

      log('HTTP 连接成功, content-type:', response.headers.get('content-type'))

      retry.value = 0
      reader = response.body.getReader()
      const decoder = new TextDecoder()
      let buffer = ''

      while (true) {
        const { done, value } = await reader.read()
        if (done) {
          log('流读取完成')
          break
        }

        const text = decoder.decode(value, { stream: true })
        buffer += text

        const parts = buffer.split('\n\n')
        buffer = parts.pop()

        for (const part of parts) {
          if (part.trim()) {
            parseSSEChunk(part + '\n\n')
          }
        }
      }

      if (buffer.trim()) {
        parseSSEChunk(buffer)
      }

      status.value = 'done'
      log('========== SSE 连接结束（流自然关闭）==========')
    } catch (err) {
      clearTimeout(timeoutId)

      if (err?.name === 'AbortError') {
        log('连接被中止')
        status.value = 'idle'
        return
      }

      log('========== SSE 连接异常 ==========', err)
      error.value = err?.message || 'SSE 连接异常'
      status.value = 'error'

      if (active && retry.value < maxRetries) {
        retry.value++
        const delay = Math.min(retryDelay * Math.pow(2, retry.value - 1), maxRetryDelay)
        ElMessage.error(`SSE 自动重连 ${retry.value}/${maxRetries}，等待 ${delay}ms...`)
        await new Promise(r => setTimeout(r, delay))
        if (active) {
          connect()
        }
      } else if (retry.value >= maxRetries) {
        ElMessage.error('已达最大重连次数，放弃连接')
        error.value += ' — 已达最大重连次数'
        status.value = 'maxRetries'
      }
    }
  }

  // SSE 连接是否可用（空闲 / 流式传输中 / 已完成）
  function isAvailable () {
    return status.value === 'idle' || status.value === 'streaming' || status.value === 'done'
  }

  // 正在连接中，或正在流式传输中
  function isAlive () {
    return status.value === 'connecting' || status.value === 'streaming'
  }

  const disconnect = () => {
    log('主动断开')
    active = false
    if (reader) {
      reader.cancel().catch(() => {})
      reader = null
    }
    if (abortController) {
      abortController.abort()
      abortController = null
    }
    status.value = 'idle'
    error.value = ''
  }

  return { status, error, retry, sseData, connect, disconnect, isAvailable, isAlive }
}
