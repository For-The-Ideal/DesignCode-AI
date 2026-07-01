package routes

import (
	"frontend_api/internal/handler"
	"frontend_api/middleware"

	"github.com/gin-gonic/gin"
)

// InitTaskRoutes 初始化任务相关路由
//
// 请求路径：
//
//	POST  /api/v1/task/create         → 创建 UI 生成任务
//	GET   /api/v1/task/:id            → 查询单个任务状态
//	GET   /api/v1/task/status         → 查询当前用户任务状态统计
//	GET   /api/v1/task/:id/events     → SSE 事件流（按任务订阅）
//	POST  /api/v1/task/sse            → SSE Broker 连接
func InitTaskRoutes(v1 *gin.RouterGroup,
	taskHandler *handler.TaskHandler,
	sseHandler *handler.SSEHandler,
) {
	task := v1.Group("/task")
	task.Use(middleware.AuthMiddleware())
	{
		// POST /api/v1/task/create → 创建生成任务
		task.POST("/create", taskHandler.CreateTask)

		// GET /api/v1/task/:id → 查询单个任务状态
		task.GET("/:id", taskHandler.GetTask)

		// GET /api/v1/task/status → 查询当前用户任务状态统计
		task.GET("/status", taskHandler.GetUserTaskStatus)

		// GET /api/v1/task/:id/events → SSE 事件流
		task.GET("/:id/events", sseHandler.StreamTask)

		// POST /api/v1/task/sse → SSE Broker 连接
		task.POST("/sse", sseHandler.StreamBroker)
	}
}
