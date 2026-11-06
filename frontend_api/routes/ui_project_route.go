package routes

import (
	"frontend_api/internal/handler"
	"frontend_api/middleware"

	"github.com/gin-gonic/gin"
)

// InitUiProjectRoutes 初始化 UI 图转代码项目相关路由
func InitUiProjectRoutes(v1 *gin.RouterGroup, uiProjectHandler *handler.UiProjectHandler) {
	// 两个接口都会调用视觉大模型，必须登录（防止匿名烧 AI 额度）
	p := v1.Group("/ui-project")
	p.Use(middleware.AuthMiddleware())
	{
		// POST /api/v1/ui-project/analyze → 顶层区域识别（进编辑器时按需调用）
		p.POST("/analyze", uiProjectHandler.AnalyzeRegions)
		// POST /api/v1/ui-project/describe → 单区域 AI 描述（右侧面板「AI 描述」按钮）
		p.POST("/describe", uiProjectHandler.DescribeRegion)
	}
}
