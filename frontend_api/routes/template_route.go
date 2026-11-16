package routes

import (
	"frontend_api/internal/handler"

	"github.com/gin-gonic/gin"
)

// InitTemplateRoutes 初始化模板相关路由
func InitTemplateRoutes(v1 *gin.RouterGroup, taskHandler *handler.TaskHandler) {
	// GET /api/v1/template     → 不带 id：随机返回一条预置模板
	// GET /api/v1/template/:id → 带 id：返回指定模板，查不到回落默认模板 1
	v1.GET("/template", taskHandler.GetTemplate)
	v1.GET("/template/:id", taskHandler.GetTemplate)
}
