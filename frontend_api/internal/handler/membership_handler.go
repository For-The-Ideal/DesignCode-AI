package handler

import (
	"frontend_api/internal/model"
	"frontend_api/pkg/mysql"
	"frontend_api/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ═══════════════════════════════════════════════
//  MembershipHandler — 会员 / 积分充值
//  当前：积分兑换；预留：微信 / 支付宝回调
// ═══════════════════════════════════════════════

// MembershipHandler 会员处理器
type MembershipHandler struct{}

// NewMembershipHandler 创建会员处理器
func NewMembershipHandler() *MembershipHandler {
	return &MembershipHandler{}
}

// 套餐等级 → 功能列表（纯展示，不存 DB）
var planFeatures = map[int][]string{
	0: {"100次/月生成", "基础组件库", "标准响应速度"},
	1: {"500次/月生成", "全组件库", "高清导出", "优先支持"},
	2: {"2000次/月生成", "全组件库", "高清导出", "API 接入", "skills支持", "专属客服"},
}

// ═══════════════════════════════════════════════
//  GET /api/v1/membership/plans
// ═══════════════════════════════════════════════

// GetPlans 获取套餐列表（公开接口）
func (h *MembershipHandler) GetPlans(c *gin.Context) {
	db := mysql.GetDB()
	var rows []model.MembershipPlan
	db.Where("is_active = ?", true).Order("sort_order asc").Find(&rows)

	// 注入功能列表
	type planOut struct {
		model.MembershipPlan
		Features []string `json:"features"`
	}
	plans := make([]planOut, len(rows))
	for i, p := range rows {
		plans[i] = planOut{MembershipPlan: p, Features: planFeatures[p.Level]}
	}

	var packages []model.CreditsPackage
	db.Order("id asc").Find(&packages)
	utils.Success(c, gin.H{
		"plans":    plans,
		"packages": packages,
	}, "获取套餐列表成功")
}

// ═══════════════════════════════════════════════
//  POST /api/v1/membership/upgrade
// ═══════════════════════════════════════════════

// UpgradeMembership 升级会员
func (h *MembershipHandler) UpgradeMembership(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var req model.UpgradeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数校验失败: "+err.Error())
		return
	}
	if req.Level == model.LevelFree {
		utils.BadRequest(c, "不能升级到免费版")
		return
	}

	db := mysql.GetDB()

	// 1. 查套餐
	var plan model.MembershipPlan
	if err := db.Where("level = ? AND is_active = ?", req.Level, true).First(&plan).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "套餐不存在或已下架")
		return
	}

	// 2. 查用户
	var user model.User
	if err := db.First(&user, userID).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "用户不存在")
		return
	}

	// 3. 增加积分（套餐每月赠送积分）
	if err := db.Model(&user).UpdateColumn("credits", gorm.Expr("credits + ?", plan.CreditsPerMonth)).Error; err != nil {
		utils.InternalError(c, "积分增加失败")
		return
	}

	// 3.5 同步用户等级
	db.Model(&user).UpdateColumn("level", plan.Level)

	// 4. 创建订单记录
	order := model.PaymentOrder{
		OrderNo:       genOrderNo("MEM"),
		UserID:        user.ID,
		OrderType:     "membership",
		Level:         plan.Level,
		CreditsAmount: plan.CreditsPerMonth,
		Amount:        plan.Price,
		Status:        model.OrderPaid,
	}
	now := time.Now()
	order.PaidAt = &now
	db.Create(&order)

	utils.Success(c, gin.H{
		"level":    plan.Level,
		"credits":  user.Credits + plan.CreditsPerMonth,
		"order_no": order.OrderNo,
	}, "升级成功")
}

// ═══════════════════════════════════════════════
//  POST /api/v1/membership/buy-credits
// ═══════════════════════════════════════════════

// BuyCredits 购买积分
func (h *MembershipHandler) BuyCredits(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var req model.BuyCreditsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数校验失败: "+err.Error())
		return
	}

	db := mysql.GetDB()

	// 1. 查积分包
	var pkg model.CreditsPackage
	if err := db.Where("id = ?", req.PackageID).First(&pkg).Error; err != nil {
		utils.BadRequest(c, "积分套餐不存在")
		return
	}

	// 2. 查用户
	var user model.User
	if err := db.First(&user, userID).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "用户不存在")
		return
	}

	// 3. 增加积分
	if err := db.Model(&user).UpdateColumn("credits", gorm.Expr("credits + ?", pkg.Credits)).Error; err != nil {
		utils.InternalError(c, "积分增加失败")
		return
	}

	// 4. 创建订单记录
	order := model.PaymentOrder{
		OrderNo:       genOrderNo("CRD"),
		UserID:        user.ID,
		OrderType:     "credits",
		CreditsAmount: pkg.Credits,
		Amount:        pkg.Price,
		Status:        model.OrderPaid,
	}
	now := time.Now()
	order.PaidAt = &now
	db.Create(&order)

	utils.Success(c, gin.H{
		"order_no": order.OrderNo,
		"credits":  user.Credits + pkg.Credits,
		"added":    pkg.Credits,
	}, "积分购买成功")
}

// ═══════════════════════════════════════════════
//  POST /api/v1/membership/callback
// ═══════════════════════════════════════════════

// PaymentCallback 支付回调（预留）
func (h *MembershipHandler) PaymentCallback(c *gin.Context) {
	var req model.PaymentCallbackReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数校验失败: "+err.Error())
		return
	}

	if req.Status != "success" {
		utils.Success(c, gin.H{}, "已记录")
		return
	}

	db := mysql.GetDB()
	var order model.PaymentOrder
	if err := db.Where("order_no = ?", req.OrderNo).First(&order).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "订单不存在")
		return
	}

	if order.Status != model.OrderPending {
		utils.Error(c, http.StatusBadRequest, "订单状态异常")
		return
	}

	// 标记已支付
	now := time.Now()
	db.Model(&order).Updates(map[string]interface{}{
		"status":   model.OrderPaid,
		"trade_no": req.TradeNo,
		"paid_at":  now,
	})

	// 按订单类型发放权益
	var user model.User
	if err := db.First(&user, order.UserID).Error; err == nil {
		if order.OrderType == "credits" {
			db.Model(&user).UpdateColumn("credits", gorm.Expr("credits + ?", order.CreditsAmount))
		} else if order.OrderType == "membership" {
			// 激活会员（忽略，UpgradeMembership 已处理）
		}
	}

	utils.Success(c, gin.H{}, "支付成功")
}

// ── 工具函数 ──────────────────────────────────

func genOrderNo(prefix string) string {
	return prefix + time.Now().Format("20060102150405") + randomDigit(6)
}

func randomDigit(n int) string {
	const digits = "0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = digits[time.Now().UnixNano()%10]
		time.Sleep(time.Nanosecond) // 防同纳秒重复
	}
	return string(b)
}
