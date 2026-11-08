package model

import "time"

// ═══════════════════════════════════════════════
//  生成结果领域模型
// ═══════════════════════════════════════════════

// PageResult 多图模式下的单页生成结果。
// 单图任务也会写入 Pages（长度 1），旧字段 Code/Preview 始终同步为入口页内容，
// 保证未升级的前端/旧数据读取路径不受影响。
type PageResult struct {
	PageID        string         `json:"page_id"`
	PageName      string         `json:"page_name"`
	Route         string         `json:"route,omitempty"`
	DSL           string         `json:"dsl,omitempty"`
	Code          string         `json:"code"`
	Preview       string         `json:"preview,omitempty"`
	Status        string         `json:"status"` // success | failed
	ErrorMsg      string         `json:"error_msg,omitempty"`
	QualityReport *QualityReport `json:"quality_report,omitempty"` // 本页的确定性质量检查报告
}

// NavEdge 合成预览后的实际跳转边（以最终产物为准，供详情页回显/审计）
type NavEdge struct {
	FromPageID string `json:"from_page_id"`
	ToPageID   string `json:"to_page_id"`
	Trigger    string `json:"trigger"`
	NodeKey    string `json:"node_key,omitempty"`
	Route      string `json:"route,omitempty"`
}

// Result 生成结果
type Result struct {
	ID               uint      `json:"id" gorm:"primaryKey"`
	TaskID           string    `json:"task_id" gorm:"size:255;index"`
	UserID           uint      `json:"user_id"`                                          // 用户 ID
	DSL              string    `json:"dsl"`                                              // 视觉分析中间 DSL (JSON)
	Code             string    `json:"code"`                                             // 最终生成的代码
	Preview          string    `json:"preview"`                                          // 预览 HTML（vue3/react 内联渲染；flutter 可作降级兜底）
	IframeURL        string    `json:"iframe_url" gorm:"size:500"`                       // 构建后的可嵌入预览地址（flutter 真机预览走 iframe src，与 Preview 二选一）
	CodeType         string    `json:"code_type"`                                        // flutter | vue3 | react
	GenerationStatus string    `json:"generation_status" gorm:"size:20;default:success"` // success | failed
	PointsDeducted   int       `json:"points_deducted"`                                  // 本次扣除的积分数
	ErrorMsg         string    `json:"error_msg,omitempty"`
	CreatedAt        time.Time `json:"created_at"`

	// ── 多页产物（单页任务同样填充，Pages 长度 1）──
	Pages      []PageResult `json:"pages,omitempty" gorm:"column:pages;serializer:json"`
	AppPreview string       `json:"app_preview,omitempty" gorm:"type:longtext"` // 合成后的单文件 SPA HTML（多页预览/跳转用）
	NavMap     []NavEdge    `json:"nav_map,omitempty" gorm:"column:nav_map;serializer:json"`

	// 入口页的确定性质量检查报告（多页时各页报告见 Pages[i].QualityReport）
	QualityReport *QualityReport `json:"quality_report,omitempty" gorm:"column:quality_report;serializer:json;type:longtext"`
}

// QualityIssue 单条规则检查结果
//
// 只记录「未通过」的规则；通过的规则只体现在 Passed/Total 计数里。
type QualityIssue struct {
	Rule     string `json:"rule"`     // R1..R5，与 prompts/*/generate.txt 的硬性规则编号一一对应
	Severity string `json:"severity"` // critical | warning
	Detail   string `json:"detail"`   // 可读的问题描述（含触发位置/样例）
}

// QualityReport 生成代码的确定性质量检查报告。
//
// 与模型自评分不同：完全由代码扫描得出，同一份代码结果恒定，
// 前端只展示 pass/warning/fail 三态与问题明细，不展示百分制分数。
type QualityReport struct {
	Passed int            `json:"passed"` // 通过的规则数
	Total  int            `json:"total"`  // 规则总数
	Status string         `json:"status"` // pass | warning | fail
	Issues []QualityIssue `json:"issues,omitempty"`
}
