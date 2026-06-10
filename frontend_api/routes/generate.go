package routes

import (
	"frontend_api/controllers"

	"github.com/gin-gonic/gin"
)

// InitGenerateRoutes 初始化 AI 代码生成路由（阻塞式）
func InitGenerateRoutes(rg *gin.RouterGroup) {
	generateController := &controllers.GenerateController{}

	generate := rg.Group("/generate")
	{
		generate.POST("/send", generateController.Generate)
	}
}
