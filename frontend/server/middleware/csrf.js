/**
 * CSRF 防护中间件（Double Submit Cookie）
 *
 * 1. 下发：任意 GET 请求（页面 HTML / SSR 导航）若无 csrf-token cookie，则生成并下发，
 *    客户端 JS 可读取（httpOnly:false）后放入请求头 X-CSRF-Token。
 * 2. 校验：仅 Sec-Fetch-Site: cross-site 的写请求（POST/PUT/DELETE/PATCH）
 *    要求 cookie 与请求头完全一致。同源请求不属于 CSRF 威胁面；
 *    SSR 内部 $fetch / 服务端脚本调用没有 Sec-Fetch-Site
 *    （注意 undici fetch 会自动带 Sec-Fetch-Mode: cors，不能以此判定浏览器），放行。
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

  // 仅真正跨站的浏览器请求需要校验；同源请求与服务端调用（无 Sec-Fetch-Site）放行
  if (getHeader(event, 'sec-fetch-site') !== 'cross-site') return

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
