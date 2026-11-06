/**
 * CSRF 防护中间件（Double Submit Cookie）
 *
 * 1. 下发：任意 GET 请求（页面 HTML / SSR 导航）若无 csrf-token cookie，则生成并下发，
 *    客户端 JS 可读取（httpOnly:false）后放入请求头 X-CSRF-Token。
 * 2. 校验：浏览器发起的写请求（POST/PUT/DELETE/PATCH）要求 cookie 与请求头完全一致。
 *    带 Sec-Fetch-Mode 头即为浏览器发起；SSR 内部 $fetch / 服务端脚本调用无该头，
 *    来源受信放行——CSRF 威胁模型只覆盖“浏览器自动携带 cookie”的场景。
 *
 * 注：生产 Nginx 存在 /api/ 直连 Go 后端的路径，Go 侧 middleware/csrf.go 做同规则校验，
 *    两层同时生效才能完整覆盖。
 */

const SAFE_METHODS = new Set(['GET', 'HEAD', 'OPTIONS'])

export default defineEventHandler((event) => {
  const method = getMethod(event)

  // 安全方法：仅确保 cookie 已下发，不校验
  if (SAFE_METHODS.has(method)) {
    if (!getCsrfCookieToken(event)) {
      setCsrfCookie(event, generateCsrfToken())
    }
    return
  }

  // 非浏览器发起（无 Sec-Fetch-Mode）：SSR 代理 / 脚本调用，放行
  if (!getHeader(event, 'sec-fetch-mode')) return

  const cookieToken = getCsrfCookieToken(event)
  const headerToken = getCsrfHeaderToken(event)
  if (
    !cookieToken ||
    !headerToken ||
    cookieToken !== headerToken ||
    !verifyCsrfToken(cookieToken)
  ) {
    throw createError({
      statusCode: 403,
      statusMessage: 'CSRF 校验失败，请刷新页面后重试',
    })
  }
})
