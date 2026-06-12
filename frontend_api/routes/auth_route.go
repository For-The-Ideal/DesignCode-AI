package routes

import (
	"frontend_api/internal/handler"

	"github.com/gin-gonic/gin"
)

// InitAuthRoutes 初始化认证相关路由
//
// 请求路径：
//
//	GET   /api/v1/auth/captcha   → 获取图形验证码
//	POST  /api/v1/auth/login     → 用户登录
//	POST  /api/v1/auth/register  → 用户注册
//	POST  /api/v1/auth/logout    → 退出登录
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
	}
}
