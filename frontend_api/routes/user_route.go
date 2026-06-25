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
//	GET  /api/v1/user/info           → 获取当前用户信息（含积分）
//	POST /api/v1/user/update         → 更新用户信息
//	POST /api/v1/user/updatePassword → 修改密码
func InitUserRoutes(v1 *gin.RouterGroup, userHandler *handler.UserHandler) {
	userV1 := v1.Group("/user")
	userV1.Use(middleware.AuthMiddleware())
	{
		// GET /api/v1/user/info → 获取当前用户信息（含积分）
		userV1.GET("/info", userHandler.GetUserInfo)

		// POST /api/v1/user/update → 更新用户信息
		userV1.POST("/update", userHandler.UpdateUserInfo)

		// POST /api/v1/user/updatePassword → 修改密码
		userV1.POST("/updatePassword", userHandler.UpdateUserPassword)
	}
}
