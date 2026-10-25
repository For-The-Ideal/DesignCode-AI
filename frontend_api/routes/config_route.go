package routes

import (
	"frontend_api/internal/handler"

	"github.com/gin-gonic/gin"
)

// InitConfigRoutes 初始化生成配置字典路由（公开，无需登录）
//
// 请求路径：
//
//	GET /api/v1/config/options → 获取框架/平台/选项/高级/组件库字典
func InitConfigRoutes(v1 *gin.RouterGroup, configHandler *handler.ConfigHandler) {
	config := v1.Group("/config")
	{
		// GET /api/v1/config/options → 获取生成配置字典
		config.GET("/options", configHandler.GetOptions)
	}
}
