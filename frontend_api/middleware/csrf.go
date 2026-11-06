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
//   - 无 Sec-Fetch-Mode 头的请求放行：Nuxt SSR 代理与脚本调用不带该头，
//     而 CSRF 威胁模型只覆盖“浏览器自动携带 cookie”的场景；
//   - 浏览器发起的写请求：csrf-token cookie 必须存在，且与 X-CSRF-Token 请求头完全一致，否则 403。
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

		// 仅浏览器发起的请求需要校验（跨站攻击请求必然带 Sec-Fetch-Mode）
		if strings.TrimSpace(c.GetHeader("Sec-Fetch-Mode")) == "" {
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
