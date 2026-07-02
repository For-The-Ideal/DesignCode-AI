package model

import (
	"time"

	"gorm.io/gorm"
)

// ═══════════════════════════════════════════════
//  gorm-base — GORM 基础字段模板（显式 JSON tag）
//
//  因为 gorm.Model 不包含 json 标签，直接嵌入会导致
//  序列化输出 ID / CreatedAt / UpdatedAt / DeletedAt
//  故三张表均显式声明，统一 snake_case 输出。
// ═══════════════════════════════════════════════

// ═══════════════════════════════════════════════
//  会员 / 支付领域模型
// ═══════════════════════════════════════════════

// 用户/套餐等级（0=免费 1=专业 2=旗舰）
const (
	LevelFree    = 0
	LevelPro     = 1
	LevelPremium = 2
)

var levelLabels = map[int]string{0: "免费版", 1: "专业版", 2: "旗舰版"}

// LevelLabel 等级 → 中文名
func LevelLabel(lv int) string {
	if s, ok := levelLabels[lv]; ok {
		return s
	}
	return "免费版"
}

// OrderStatus 订单状态
type OrderStatus string

const (
	OrderPending  OrderStatus = "pending"  // 待支付
	OrderPaid     OrderStatus = "paid"     // 已支付
	OrderExpired  OrderStatus = "expired"  // 已过期（未支付）
	OrderRefunded OrderStatus = "refunded" // 已退款
)

// PaymentMethod 支付方式
type PaymentMethod string

const (
	PayWechat   PaymentMethod = "wechat"
	PayAlipay   PaymentMethod = "alipay"
	PayCredits  PaymentMethod = "credits"  // 积分兑换
	PayInternal PaymentMethod = "internal" // 内部操作（管理员充值）
)

// ── MembershipPlan 套餐定义 ──────────────────

// MembershipPlan 会员套餐
type MembershipPlan struct {
	ID              uint           `gorm:"primarykey" json:"id"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
	Level           int            `gorm:"uniqueIndex;not null;default:0" json:"level"` // 0=免费 1=专业 2=旗舰
	Name            string         `gorm:"size:50;not null" json:"name"`
	CreditsPerMonth int            `gorm:"not null;default:0" json:"credits_per_month"`
	Price           float64        `gorm:"not null;default:0" json:"price"`
	OriginalPrice   float64        `gorm:"not null;default:0" json:"original_price"`
	IsActive        bool           `gorm:"not null;default:true" json:"is_active"` // 是否启用
	SortOrder       int            `gorm:"not null;default:0" json:"sort_order"`   // 排序
}

// ── PaymentOrder 支付订单（预留支付集成）──

// PaymentOrder 支付订单
type PaymentOrder struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	OrderNo       string         `gorm:"size:64;uniqueIndex;not null" json:"order_no"` // 订单号
	UserID        uint           `gorm:"index;not null" json:"user_id"`
	OrderType     string         `gorm:"size:20;not null" json:"order_type"`       // membership / credits
	Level         int            `gorm:"not null;default:0" json:"level"`          // 套餐等级（membership 订单）
	CreditsAmount int            `gorm:"not null;default:0" json:"credits_amount"` // 积分数量（credits 订单）
	Amount        float64        `gorm:"not null;default:0" json:"amount"`         // 金额（元）
	PaymentMethod PaymentMethod  `gorm:"size:20;not null;default:credits" json:"payment_method"`
	Status        OrderStatus    `gorm:"size:20;not null;default:pending" json:"status"`
	TradeNo       string         `gorm:"size:128" json:"trade_no"` // 第三方交易号（预留）
	PaidAt        *time.Time     `json:"paid_at"`
	ExpireAt      *time.Time     `json:"expire_at"` // 订单过期时间
}

// ── 请求 / 响应结构体 ──────────────────────────

// CreditsPackage 积分充值套餐
type CreditsPackage struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	Name          string         `gorm:"size:100;not null" json:"name"`
	Credits       int            `gorm:"not null;default:0" json:"credits"`
	Price         float64        `gorm:"not null;default:0" json:"price"`
	OriginalPrice float64        `gorm:"not null;default:0" json:"original_price"`
}

func (CreditsPackage) TableName() string { return "credit_packages" }

// UpgradeRequest 升级会员请求
type UpgradeRequest struct {
	Level         int           `json:"level" binding:"required"`
	PaymentMethod PaymentMethod `json:"payment_method"` // wechat / alipay / credits，默认 credits
}

// BuyCreditsRequest 购买积分请求
type BuyCreditsRequest struct {
	PackageID     int           `json:"package_id" binding:"required"`
	PaymentMethod PaymentMethod `json:"payment_method"` // wechat / alipay，默认留空（待支付接入后使用）
}

// ── 支付回调请求（预留）────────────────────

// PaymentCallbackReq 支付回调
type PaymentCallbackReq struct {
	OrderNo string  `json:"order_no" binding:"required"`
	TradeNo string  `json:"trade_no" binding:"required"`
	Status  string  `json:"status" binding:"required"` // success / fail
	Amount  float64 `json:"amount"`
	Sign    string  `json:"sign"`
}
