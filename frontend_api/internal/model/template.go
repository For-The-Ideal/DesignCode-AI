package model

// ═══════════════════════════════════════════════
//  TaskTemplate — 预置模板
//  对应数据库表：task_template
// ═══════════════════════════════════════════════

// TaskTemplate 预置模板
type TaskTemplate struct {
	ID           int    `json:"id" gorm:"primaryKey"`
	TemplateCode string `json:"template_code" gorm:"column:template_code"`
	PreviewCode  string `json:"preview_code" gorm:"column:preview_code"`
}

// TableName 指定表名
func (TaskTemplate) TableName() string {
	return "task_template"
}
