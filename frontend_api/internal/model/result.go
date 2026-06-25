package model

import "time"

// ═══════════════════════════════════════════════
//  生成结果领域模型
// ═══════════════════════════════════════════════

// Result 生成结果
type Result struct {
	ID               uint      `json:"id" gorm:"primaryKey"`
	TaskID           string    `json:"task_id" gorm:"size:255;index"`
	UserID           uint      `json:"user_id"`   // 用户 ID
	DSL              string    `json:"dsl"`       // 视觉分析中间 DSL (JSON)
	Code             string    `json:"code"`      // 最终生成的代码
	Preview          string    `json:"preview"`   // 预览 HTML
	CodeType         string    `json:"code_type"` // flutter | vue3 | react
	Score            int       `json:"score"`
	GenerationStatus string    `json:"generation_status" gorm:"size:20;default:success"` // success | failed
	PointsDeducted   int       `json:"points_deducted"`                                  // 本次扣除的积分数
	ErrorMsg         string    `json:"error_msg,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
}
