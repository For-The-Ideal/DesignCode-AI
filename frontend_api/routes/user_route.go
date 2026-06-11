package routes

import (
	"frontend_api/internal/handler"
	"frontend_api/middleware"

	"github.com/gin-gonic/gin"
)

// InitUserRoutes 初始化用户相关路由
//
// 所有接口均需登录（AuthMiddleware 自动注入 user_id）
//
// 请求路径：
//
//	GET   /api/v1/user/info   → 获取当前用户信息
func InitUserRoutes(v1 *gin.RouterGroup, userHandler *handler.UserHandler) {
	userV1 := v1.Group("/user")
	userV1.Use(middleware.AuthMiddleware())
	{
		// GET /api/v1/user/info → 获取当前用户信息
		userV1.GET("/info", userHandler.GetUserInfo)
	}
}
