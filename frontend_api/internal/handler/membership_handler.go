package handler

import (
	"errors"
	"frontend_api/internal/model"
	"frontend_api/pkg/mysql"
	"frontend_api/utils"
	"math/rand/v2"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

	// 3. 加积分、同步等级、创建订单：同一事务，避免部分成功
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

	var creditsAfter int
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&user).UpdateColumn("credits", gorm.Expr("credits + ?", plan.CreditsPerMonth)).Error; err != nil {
			return err
		}
		if err := tx.Model(&user).UpdateColumn("level", plan.Level).Error; err != nil {
			return err
		}
		if err := tx.Create(&order).Error; err != nil {
			return err
		}
		// 事务内重查，返回增加后的真实积分
		var latest model.User
		if err := tx.Select("credits").First(&latest, user.ID).Error; err != nil {
			return err
		}
		creditsAfter = latest.Credits
		return nil
	})
	if err != nil {
		utils.InternalError(c, "升级失败，请稍后重试")
		return
	}

	utils.Success(c, gin.H{
		"level":    plan.Level,
		"credits":  creditsAfter,
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

	// 3. 加积分、创建订单：同一事务，避免部分成功
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

	var creditsAfter int
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&user).UpdateColumn("credits", gorm.Expr("credits + ?", pkg.Credits)).Error; err != nil {
			return err
		}
		if err := tx.Create(&order).Error; err != nil {
			return err
		}
		// 事务内重查，返回增加后的真实积分
		var latest model.User
		if err := tx.Select("credits").First(&latest, user.ID).Error; err != nil {
			return err
		}
		creditsAfter = latest.Credits
		return nil
	})
	if err != nil {
		utils.InternalError(c, "积分购买失败，请稍后重试")
		return
	}

	utils.Success(c, gin.H{
		"order_no": order.OrderNo,
		"credits":  creditsAfter,
		"added":    pkg.Credits,
	}, "积分购买成功")
}

// errOrderStatusInvalid 回调到达时订单已不是待支付状态
var errOrderStatusInvalid = errors.New("订单状态异常")

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
	now := time.Now()

	var order model.PaymentOrder
	err := db.Transaction(func(tx *gorm.DB) error {
		// 行锁：同一订单的并发回调只能有一个完成 pending → paid
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("order_no = ?", req.OrderNo).First(&order).Error; err != nil {
			return err
		}
		if order.Status != model.OrderPending {
			return errOrderStatusInvalid
		}
		if err := tx.Model(&order).Updates(map[string]interface{}{
			"status":   model.OrderPaid,
			"trade_no": req.TradeNo,
			"paid_at":  now,
		}).Error; err != nil {
			return err
		}
		// 按订单类型发放权益（membership 的权益在 UpgradeMembership 已发放，不重复加）
		if order.OrderType == "credits" {
			if err := tx.Model(&model.User{}).Where("id = ?", order.UserID).
				UpdateColumn("credits", gorm.Expr("credits + ?", order.CreditsAmount)).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			utils.Error(c, http.StatusNotFound, "订单不存在")
		case errors.Is(err, errOrderStatusInvalid):
			utils.Error(c, http.StatusBadRequest, "订单状态异常")
		default:
			utils.InternalError(c, "回调处理失败")
		}
		return
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
		b[i] = digits[rand.IntN(len(digits))]
	}
	return string(b)
}
