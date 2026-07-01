/**
 * SSR 请求日志工具
 * 按类型拆分为三个文件，便于分别监控
 *   ssr-req.log — 请求日志（REQ）
 *   ssr-res.log — 成功响应（RES）
 *   ssr-err.log — 错误（ERR）
 */
import { appendFileSync, existsSync, mkdirSync } from 'node:fs'
import { resolve } from 'node:path'

const logDir = resolve(process.cwd(), 'logs')

// 按类型分文件
const LOG_FILES = {
  REQ: resolve(logDir, 'ssr-req.log'),
  RES: resolve(logDir, 'ssr-res.log'),
  ERR: resolve(logDir, 'ssr-err.log'),
}

if (!existsSync(logDir)) mkdirSync(logDir, { recursive: true })

function ts() {
  return new Date().toISOString().replace('T', ' ').slice(0, 19)
}

function write(type, line) {
  try {
    appendFileSync(LOG_FILES[type], line + '\n', 'utf-8')
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
  write(type, line)
}
