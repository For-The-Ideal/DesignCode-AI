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
//	POST  /api/v1/task/create  → 创建 UI 生成任务
//	GET   /api/v1/task/:id     → 查询单个任务状态
//	GET   /api/v1/task/status  → 查询当前用户任务状态统计
//	POST  /api/v1/task/list    → 查询当前用户任务列表
func InitTaskRoutes(v1 *gin.RouterGroup,
	taskHandler *handler.TaskHandler,
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

		// DELETE /api/v1/task/:id → 删除任务
		task.DELETE("/:id", taskHandler.DeleteTask)

		// POST /api/v1/task/list → 获取任务列表
		task.POST("/list", taskHandler.GetTaskList)
	}
}
