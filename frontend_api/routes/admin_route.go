package routes

import (
	"frontend_api/internal/handler"
	"frontend_api/middleware"

	"github.com/gin-gonic/gin"
)

// InitAdminRoutes 初始化管理员相关路由
//
// 请求路径：
//
//	POST /api/v1/admin/create-account → 管理员创建用户账号
func InitAdminRoutes(v1 *gin.RouterGroup, adminHandler *handler.AdminHandler) {
	admin := v1.Group("/admin")
	// 登录 + 管理员邮箱白名单（server.admin_emails），两者缺一不可
	admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
	{
		// POST /api/v1/admin/create-account → 创建用户账号
		admin.POST("/create-account", adminHandler.CreateAccount)
	}
}
