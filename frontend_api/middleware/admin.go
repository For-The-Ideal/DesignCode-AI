package middleware

import (
	"frontend_api/config"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AdminMiddleware 管理员权限校验（必须挂在 AuthMiddleware 之后）
//
// 规则：JWT 中的 email 必须命中配置 server.admin_emails 白名单；
// 白名单为空时拒绝所有请求（安全默认，避免忘记配置导致后台裸奔）。
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		email, _ := c.Get("email")
		addr, _ := email.(string)

		if !isAdminEmail(addr) {
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "无管理员权限"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// isAdminEmail 判断邮箱是否在管理员白名单中（忽略大小写与首尾空白）
func isAdminEmail(email string) bool {
	email = strings.TrimSpace(email)
	if email == "" || config.AppConfig == nil {
		return false
	}
	for _, allowed := range config.AppConfig.Server.AdminEmails {
		if strings.EqualFold(strings.TrimSpace(allowed), email) {
			return true
		}
	}
	return false
}
