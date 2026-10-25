package routes

import (
	"frontend_api/internal/handler"

	"github.com/gin-gonic/gin"
)

// InitV1Routes 初始化全部路由入口
//
// 功能模块拆分：
//
//	auth_route.go        — 认证相关（captcha / login / register）
//	user_route.go        — 用户相关（info）
//	config_route.go      — 生成配置字典（config/options）
//	membership_route.go  — 会员相关（plans / upgrade / buy-credits / callback）
//	sse_route.go         — SSE 事件流（/sse/:id/events / sse/user）（代码/进度流式推送 / 任务状态变更推送）
//	task_route.go        — 任务相关（create / :id / status / list）
//	upload_route.go      — 上传相关（upload）
//	template_route.go    — 模板相关（template/:id）
//	ui_project_route.go  — UI 图转代码项目（ui-project/analyze 顶层区域识别）
//
// 所有路由统一以 /api/v1/ 开头
func InitV1Routes(r *gin.Engine,
	taskHandler *handler.TaskHandler,
	sseHandler *handler.SSEHandler,
	uploadHandler *handler.UploadHandler,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	adminHandler *handler.AdminHandler,
	membershipHandler *handler.MembershipHandler,
	uiProjectHandler *handler.UiProjectHandler,
	configHandler *handler.ConfigHandler,
) {
	// /api/v1 路由组
	v1 := r.Group("/api/v1")

	// 认证路由
	InitAuthRoutes(v1, authHandler)

	// 用户路由（需登录）
	InitUserRoutes(v1, userHandler)

	// 会员路由
	InitMembershipRoutes(v1, membershipHandler)

	// SSE 事件流路由
	InitSSERoutes(v1, sseHandler)

	// 任务路由
	InitTaskRoutes(v1, taskHandler)

	// 上传路由
	InitUploadRoutes(v1, uploadHandler)

	// 模板路由
	InitTemplateRoutes(v1, taskHandler)

	// 管理路由
	InitAdminRoutes(v1, adminHandler)

	// UI 图转代码项目路由
	InitUiProjectRoutes(v1, uiProjectHandler)

	// 生成配置字典路由（公开）
	InitConfigRoutes(v1, configHandler)
}
