/**
 * CSRF Token 生成与验证工具
 *
 * 防御原理（Double Submit Cookie 模式）：
 *   1. 服务端生成 token，通过 cookie 下发给浏览器（httpOnly: false，JS 可读取）
 *   2. 客户端发起 POST/PUT/DELETE/PATCH 请求时，从 cookie 读取 token 放入请求头 X-CSRF-Token
 *   3. 服务端中间件校验：请求头中的 token 必须与 cookie 中的 token 一致
 *
 * 安全性：
 *   - 跨站攻击者能让浏览器自动带上 cookie，但无法读取 cookie 值
 *   - 因此无法伪造 X-CSRF-Token 请求头，请求会被拦截
 *
 * 在 server/ 目录下自动导入，直接使用：
 *   generateCsrfToken()
 *   validateCsrfToken(event)
 */

import { randomBytes, createHmac } from 'node:crypto'

/** Cookie 名称 */
const CSRF_COOKIE_NAME = 'csrf-token'

/** 请求头名称 */
const CSRF_HEADER_NAME = 'x-csrf-token'

/** Token 有效期（24 小时，单位毫秒） */
const TOKEN_MAX_AGE_MS = 24 * 60 * 60 * 1000

/** Cookie maxAge（秒） */
const COOKIE_MAX_AGE = 24 * 60 * 60

/**
 * 获取 CSRF 签名密钥
 */
function getSecret() {
  const defaultSecret = 'default-yysecret-ab709a98-42789411';
  try {
    const config = useRuntimeConfig()
    return config.csrfSecret || defaultSecret
  } catch {
    return process.env.NUXT_CSRF_SECRET || defaultSecret
  }
}

/**
 * 生成 CSRF Token
 * 格式：timestamp.randomBytes.hmacSignature
 */
export function generateCsrfToken() {
  const secret = getSecret()
  const timestamp = Date.now().toString(36)
  const random = randomBytes(16).toString('hex')
  const payload = `${timestamp}.${random}`
  const signature = createHmac('sha256', secret).update(payload).digest('hex').slice(0, 16)
  return `${payload}.${signature}`
}

/**
 * 验证 CSRF Token 签名和有效期
 */
export function verifyCsrfToken(token) {
  if (!token || typeof token !== 'string') return false

  const parts = token.split('.')
  if (parts.length !== 3) return false

  const [timestamp, random, signature] = parts
  const secret = getSecret()

  // 验证签名
  const payload = `${timestamp}.${random}`
  const expected = createHmac('sha256', secret).update(payload).digest('hex').slice(0, 16)
  if (signature !== expected) return false

  // 验证有效期
  const tokenTime = parseInt(timestamp, 36)
  if (Date.now() - tokenTime > TOKEN_MAX_AGE_MS) return false

  return true
}

/**
 * 设置 CSRF Cookie
 */
export function setCsrfCookie(event, token) {
  setCookie(event, CSRF_COOKIE_NAME, token, {
    httpOnly: false,
    secure: process.env.NODE_ENV === 'production',
    sameSite: 'lax',
    path: '/',
    maxAge: COOKIE_MAX_AGE,
  })
}

/**
 * 从请求中读取 CSRF Cookie
 */
export function getCsrfCookieToken(event) {
  return getCookie(event, CSRF_COOKIE_NAME) || ''
}

/**
 * 从请求头中读取 CSRF Token
 */
export function getCsrfHeaderToken(event) {
  return getHeader(event, CSRF_HEADER_NAME) || ''
}

/**
 * 导出常量供其他模块使用
 */
export const CSRF_CONSTANTS = {
  COOKIE_NAME: CSRF_COOKIE_NAME,
  HEADER_NAME: CSRF_HEADER_NAME,
}
