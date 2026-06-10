package routes

import (
	"frontend_api/controllers"
	"frontend_api/middleware"

	"github.com/gin-gonic/gin"
)

// InitGenerateRoutes 初始化 AI 代码生成路由
func InitGenerateRoutes(rg *gin.RouterGroup) {
	generateController := &controllers.GenerateController{}

	generate := rg.Group("/generate")
	generate.Use(middleware.AuthMiddleware())
	{
		generate.POST("/send", generateController.Generate)
	}
}
