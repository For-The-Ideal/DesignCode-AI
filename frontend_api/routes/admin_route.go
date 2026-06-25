package routes

import (
	"frontend_api/internal/handler"

	"github.com/gin-gonic/gin"
)

// InitAdminRoutes 初始化管理员相关路由
//
// 请求路径：
//
//	POST /api/v1/admin/create-account → 管理员创建用户账号
func InitAdminRoutes(v1 *gin.RouterGroup, adminHandler *handler.AdminHandler) {
	admin := v1.Group("/admin")
	{
		// POST /api/v1/admin/create-account → 创建用户账号
		admin.POST("/create-account", adminHandler.CreateAccount)
	}
}
