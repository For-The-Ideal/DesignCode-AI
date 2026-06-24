/**
 * SSR 请求日志工具
 * 记录每次 ajaxData 代理的请求参数（解密后）和响应结果
 */
import { appendFileSync, existsSync, mkdirSync } from 'node:fs'
import { resolve } from 'node:path'

// 日志目录：项目根目录/logs/ssr-api.log
const logDir = resolve(process.cwd(), 'logs')
const logFile = resolve(logDir, 'ssr-api.log')

if (!existsSync(logDir)) mkdirSync(logDir, { recursive: true })

function ts() {
  return new Date().toISOString().replace('T', ' ').slice(0, 19)
}

function write(line) {
  try {
    appendFileSync(logFile, line + '\n', 'utf-8')
  } catch {}
}

/**
 * 写日志，同时 console.log
 * @param {'REQ'|'RES'|'ERR'} type
 * @param {object} detail
 */
export function serverLog(type, detail = {}) {
  const { method = 'POST', url = '', params, status, duration, message, isBrowser } = detail
  const parts = [`[${ts()}]`, type === 'ERR' ? 'ERROR' : 'INFO', type]

  if (type === 'REQ') {
    parts.push(isBrowser ? 'BROWSER' : 'SSR')
    parts.push(method.toUpperCase(), url)
    parts.push(JSON.stringify(params))
  } else if (type === 'RES') {
    parts.push(method.toUpperCase(), url, `${status}`, `${duration}ms`)
  } else if (type === 'ERR') {
    parts.push(url, message)
  }

  const line = parts.join(' | ')
  console.log(`[SSR-LOG] ${line}`)
  write(line)
}
