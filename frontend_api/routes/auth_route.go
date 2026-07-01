package routes

import (
	"frontend_api/internal/handler"

	"github.com/gin-gonic/gin"
)

// InitAuthRoutes 初始化认证相关路由
//
// 请求路径：
//
//	GET   /api/v1/auth/captcha          → 获取图形验证码
//	POST  /api/v1/auth/login            → 用户登录
//	POST  /api/v1/auth/register         → 用户注册
//	POST  /api/v1/auth/logout           → 退出登录
//	POST  /api/v1/auth/forgot-password  → 忘记密码（生成重置令牌并发送邮件）
//	POST  /api/v1/auth/reset-password   → 重置密码
func InitAuthRoutes(v1 *gin.RouterGroup, authHandler *handler.AuthHandler) {
	auth := v1.Group("/auth")
	{
		// GET /api/v1/auth/captcha → 获取图形验证码
		auth.GET("/captcha", authHandler.Captcha)

		// POST /api/v1/auth/login → 用户登录
		auth.POST("/login", authHandler.Login)

		// POST /api/v1/auth/register → 用户注册
		auth.POST("/register", authHandler.Register)

		// POST /api/v1/auth/logout → 退出登录
		auth.POST("/logout", authHandler.Logout)

		// POST /api/v1/auth/forgot-password → 忘记密码
		auth.POST("/forgot-password", authHandler.ForgotPassword)

		// POST /api/v1/auth/reset-password → 重置密码
		auth.POST("/reset-password", authHandler.ResetPassword)
	}
}
