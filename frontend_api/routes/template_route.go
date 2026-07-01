package routes

import (
	"frontend_api/internal/handler"

	"github.com/gin-gonic/gin"
)

// InitTemplateRoutes 初始化模板相关路由
func InitTemplateRoutes(v1 *gin.RouterGroup, taskHandler *handler.TaskHandler) {
	// GET /api/v1/template/:id → 获取预置模板
	v1.GET("/template/:id", taskHandler.GetTemplate)
}
