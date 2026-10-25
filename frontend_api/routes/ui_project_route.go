package routes

import (
	"frontend_api/internal/handler"

	"github.com/gin-gonic/gin"
)

// InitUiProjectRoutes 初始化 UI 图转代码项目相关路由
func InitUiProjectRoutes(v1 *gin.RouterGroup, uiProjectHandler *handler.UiProjectHandler) {
	// POST /api/v1/ui-project/analyze → 顶层区域识别（进编辑器时按需调用）
	v1.POST("/ui-project/analyze", uiProjectHandler.AnalyzeRegions)
}
