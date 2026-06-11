package routes

import (
	"frontend_api/internal/handler"

	"github.com/gin-gonic/gin"
)

// InitGenerateRoutes 初始化代码生成相关路由（v1 新架构）
//
// 流程：创建任务 → 查询状态（SSE轮询/刷新恢复）→ SSE 事件流接收进度
//
// 请求路径：
//
//	POST  /api/v1/generate-ui        → 创建 UI 生成任务
//	      请求：{"target":"flutter","images":[{"url":"...","desc":"..."}]}
//	      响应：{"task_id":"xxx","status":"pending"}
//
//	GET   /api/v1/task/:id            → 查询任务状态
//	      用途：页面刷新后恢复进度，判断 can_sse 决定是否重连 SSE
//	      响应：{"task_id":"xxx","status":"running","progress":60,
//	             "current_step":"FlutterGenerateSkill","can_sse":true,"result":null}
//
//	GET   /api/v1/task/:id/events     → SSE 事件流
//	      前端：new EventSource("/api/v1/task/{taskId}/events")
//	      事件格式：
//	        event: progress  data: {"progress":20,"step":"VisionAnalyzeSkill"}
//	        event: done      data: {"progress":100,"step":"Done"}
//	        event: error     data: {"message":"AI接口超时"}
//
//	POST  /api/v1/upload              → 图片上传（占位，待接入 COS）
//	      请求：multipart/form-data  image
//	      响应：{"url":"","filename":"xxx","size":123}
func InitGenerateRoutes(v1 *gin.RouterGroup,
	taskHandler *handler.TaskHandler,
	sseHandler *handler.SSEHandler,
	uploadHandler *handler.UploadHandler,
) {
	// POST /api/v1/generate-ui → 创建 UI 生成任务
	v1.POST("/generate-ui", taskHandler.CreateTask)

	// GET /api/v1/task/:id → 查询任务状态
	v1.GET("/task/:id", taskHandler.GetTask)

	// GET /api/v1/task/:id/events → SSE 事件流
	v1.GET("/task/:id/events", sseHandler.StreamTask)

	// POST /api/v1/upload → 图片上传（占位）
	v1.POST("/upload", uploadHandler.Upload)

	// GET /api/v1/template → 获取预置模板（?id=1）
	v1.GET("/template", taskHandler.GetTemplate)
}
