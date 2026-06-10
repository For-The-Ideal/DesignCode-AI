package routes

import (
	"frontend_api/controllers"

	"github.com/gin-gonic/gin"
)

// InitSSERoutes 初始化 SSE 流式代码生成路由
func InitSSERoutes(rg *gin.RouterGroup) {
	sseController := &controllers.SSEController{}

	sse := rg.Group("/")
	{
		sse.POST("/sse", sseController.Handle)
	}
}
