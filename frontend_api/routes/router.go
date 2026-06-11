package routes

import (
	"frontend_api/internal/handler"

	"github.com/gin-gonic/gin"
)

// InitV1Routes 初始化全部路由入口
//
// 功能模块拆分：
//
//	auth_route.go     — 认证相关（captcha / login / register）
//	user_route.go     — 用户相关（info）
//	generate_route.go — 代码生成相关（generate-ui / task / upload / SSE）
//
// 所有路由统一以 /api/v1/ 开头
func InitV1Routes(r *gin.Engine,
	taskHandler *handler.TaskHandler,
	sseHandler *handler.SSEHandler,
	uploadHandler *handler.UploadHandler,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
) {
	// /api/v1 路由组
	v1 := r.Group("/api/v1")

	// 认证路由
	InitAuthRoutes(v1, authHandler)

	// 用户路由（需登录）
	InitUserRoutes(v1, userHandler)

	// 代码生成路由
	InitGenerateRoutes(v1, taskHandler, sseHandler, uploadHandler)
}
