package routes

import (
	"frontend_api/internal/handler"
	"frontend_api/middleware"

	"github.com/gin-gonic/gin"
)

// InitMembershipRoutes 初始化会员相关路由
//
// 请求路径：
//
//	GET  /api/v1/membership/plans       → 套餐列表（公开）
//	POST /api/v1/membership/upgrade     → 升级会员（需登录）
//	POST /api/v1/membership/buy-credits → 购买积分（需登录）
//	POST /api/v1/membership/callback    → 支付回调（预留，公开）
func InitMembershipRoutes(v1 *gin.RouterGroup, h *handler.MembershipHandler) {
	membership := v1.Group("/membership")
	{
		// 公开接口
		membership.GET("/plans", h.GetPlans)

		// 支付回调（第三方调用，不走 cookie 认证 —— 预留安全签名校验）
		membership.POST("/callback", h.PaymentCallback)

		// 需登录
		authed := membership.Group("")
		authed.Use(middleware.AuthMiddleware())
		{
			//升级套餐
			authed.POST("/upgrade", h.UpgradeMembership)
			//追加积分
			authed.POST("/buy-credits", h.BuyCredits)
		}
	}
}
