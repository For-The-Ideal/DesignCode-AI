package model

import "time"

// ═══════════════════════════════════════════════
//  TaskTemplate — 预置模板
//  对应数据库表：task_template
// ═══════════════════════════════════════════════

// TaskTemplate 预置模板
type TaskTemplate struct {
	ID           int       `json:"id" gorm:"primaryKey"`
	Name         string    `json:"name" gorm:"size:255"`                      // 模板名称（如：电商首页模板）
	Framework    string    `json:"framework" gorm:"size:20;index"`            // 框架类型：flutter | vue3 | react
	TemplateCode string    `json:"template_code" gorm:"column:template_code"` // AI 生成的原始代码
	PreviewCode  string    `json:"preview_code" gorm:"column:preview_code"`   // 转换后的 HTML 预览代码
	Thumbnail    string    `json:"thumbnail" gorm:"size:500"`                 // 预览截图 URL（可选）
	Description  string    `json:"description" gorm:"type:text"`              // 模板描述
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TableName 指定表名
func (TaskTemplate) TableName() string {
	return "task_template"
}
