package model

import "time"

// ═══════════════════════════════════════════════
//  GenConfigOption — 生成配置字典
//  对应数据库表：gen_config_option
//  统一存放：目标框架 / 目标平台 / 生成选项 / 高级选项 / 组件库
//  其中 component_lib 额外用 framework + platform_scope 归属与区分 PC/手机端
// ═══════════════════════════════════════════════

// 配置分类
const (
	ConfigCategoryFramework    = "framework"     // 目标框架
	ConfigCategoryPlatform     = "platform"      // 目标平台
	ConfigCategoryOption       = "option"        // 生成选项
	ConfigCategoryAdvanced     = "advanced"      // 高级选项
	ConfigCategoryComponentLib = "component_lib" // 组件库
	ConfigCategoryStatus       = "status"        // 任务状态（供列表/筛选回显 label）
)

// 组件库适用平台范围
const (
	PlatformScopeMobile  = "mobile"  // 仅手机端
	PlatformScopeDesktop = "desktop" // 仅 PC 桌面端
	PlatformScopeBoth    = "both"    // 通用
)

// GenConfigOption 生成配置项
type GenConfigOption struct {
	ID            uint      `gorm:"primarykey" json:"id"`
	Category      string    `gorm:"size:20;not null;index:idx_gen_config_category_framework" json:"category"` // framework|platform|option|advanced|component_lib
	Label         string    `gorm:"size:50;not null" json:"label"`                                            // 展示文案
	Value         string    `gorm:"size:50;not null" json:"value"`                                            // 提交给生成逻辑的值
	Icon          string    `gorm:"size:50" json:"icon"`                                                      // 图标（平台用）
	Framework     string    `gorm:"size:20;index:idx_gen_config_category_framework" json:"framework"`         // 仅 component_lib：Flutter|React|Vue
	PlatformScope string    `gorm:"size:10;default:both" json:"platform_scope"`                               // 仅 component_lib：mobile|desktop|both
	SortOrder     int       `gorm:"not null;default:0" json:"sort_order"`
	IsActive      bool      `gorm:"not null;default:true" json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// TableName 指定表名
func (GenConfigOption) TableName() string {
	return "gen_config_option"
}
