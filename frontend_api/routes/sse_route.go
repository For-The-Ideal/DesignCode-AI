package routes

import (
	"frontend_api/internal/handler"
	"frontend_api/middleware"

	"github.com/gin-gonic/gin"
)

// InitSSERoutes 初始化 SSE 事件流路由
//
//	GET  /api/v1/sse/:id/events  → 按任务订阅 SSE（代码/进度流式推送）
//	POST /api/v1/sse/user         → 用户级 SSE Broker（任务状态变更推送）
func InitSSERoutes(v1 *gin.RouterGroup, sseHandler *handler.SSEHandler) {
	task := v1.Group("/sse")
	task.Use(middleware.AuthMiddleware())
	{
		// GET /api/v1/sse/:id/events → SSE 事件流（按任务订阅）（代码/进度流式推送）
		task.GET("/:id/events", sseHandler.StreamTask)

		// POST /api/v1/sse/user → 用户级 SSE Broker（任务状态变更推送）
		task.POST("/user", sseHandler.StreamBroker)
	}
}
