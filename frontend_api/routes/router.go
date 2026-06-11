package routes

import (
	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		// 认证模块路由 (登录、注册)
		InitAuthRoutes(api)

		// 用户模块路由 (个人信息等)
		InitUserRoutes(api)

		// AI 代码生成路由（阻塞式）
		InitGenerateRoutes(api)

		// SSE 流式代码生成路由
		InitAiRoutes(api)
	}
}
