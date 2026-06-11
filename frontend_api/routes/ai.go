package routes

import (
	"frontend_api/controllers"

	"github.com/gin-gonic/gin"
)

// InitAiRoutes 初始化 SSE链接
func InitAiRoutes(rg *gin.RouterGroup) {
	aiController := &controllers.AiController{}

	ai := rg.Group("/ai")
	{
		ai.POST("/sse", aiController.AiSse)
	}
}
