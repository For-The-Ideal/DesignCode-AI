package middleware

import (
	"crypto/subtle"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CSRFMiddleware CSRF 防护（Double Submit Cookie）
//
// 规则：
//   - 安全方法（GET/HEAD/OPTIONS）放行；
//   - 仅 Sec-Fetch-Site: cross-site 的写请求才校验：
//     同源请求（same-origin/same-site）不属于 CSRF 威胁面；
//     Nuxt SSR 代理（undici fetch 会自动带 Sec-Fetch-Mode: cors，
//     但没有 Sec-Fetch-Site）属于服务端请求，同样放行；
//   - 跨站写请求：csrf-token cookie 必须存在，且与 X-CSRF-Token 请求头完全一致，否则 403。
//
// 说明：生产 Nginx 存在 /api/ 直连 Go 的路径（绕过 Nuxt 层校验），
// 因此该中间件与 Nuxt server/middleware/csrf.js 必须同时启用。
func CSRFMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		switch c.Request.Method {
		case http.MethodGet, http.MethodHead, http.MethodOptions:
			c.Next()
			return
		}

		// 仅真正跨站的浏览器请求需要校验（跨站攻击必然带 Sec-Fetch-Site: cross-site）
		if strings.TrimSpace(c.GetHeader("Sec-Fetch-Site")) != "cross-site" {
			c.Next()
			return
		}

		cookieToken, err := c.Cookie("csrf-token")
		headerToken := c.GetHeader("X-CSRF-Token")
		if err != nil || cookieToken == "" || headerToken == "" ||
			len(cookieToken) > 256 || len(headerToken) > 256 ||
			subtle.ConstantTimeCompare([]byte(cookieToken), []byte(headerToken)) != 1 {
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "CSRF 校验失败，请刷新页面后重试"})
			c.Abort()
			return
		}

		c.Next()
	}
}
