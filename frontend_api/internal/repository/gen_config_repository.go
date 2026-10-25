package repository

import (
	"frontend_api/internal/model"
	"frontend_api/pkg/mysql"
	"log"
)

// ═══════════════════════════════════════════════
//  GenConfig Repository - 生成配置字典
//  提供：分组查询（供前端动态渲染）+ 首次启动灌入默认数据
// ═══════════════════════════════════════════════

// GenConfigRepository 配置字典存储
type GenConfigRepository struct{}

// NewGenConfigRepository 创建配置字典存储
func NewGenConfigRepository() *GenConfigRepository {
	return &GenConfigRepository{}
}

// ── 响应结构 ──

// ConfigOptionItem 单个配置项
type ConfigOptionItem struct {
	Label    string `json:"label"`
	Value    string `json:"value"`
	Icon     string `json:"icon,omitempty"`
	Platform string `json:"platform,omitempty"` // 仅组件库：mobile|desktop|both
}

// GenConfigResponse 分组后的完整字典
type GenConfigResponse struct {
	Frameworks    []ConfigOptionItem            `json:"frameworks"`
	Platforms     []ConfigOptionItem            `json:"platforms"`
	Options       []ConfigOptionItem            `json:"options"`
	Advanced      []ConfigOptionItem            `json:"advanced"`
	Statuses      []ConfigOptionItem            `json:"statuses"`
	ComponentLibs map[string][]ConfigOptionItem `json:"component_libs"` // key = 框架
}

// GetOptions 读取并按分类分组
func (r *GenConfigRepository) GetOptions() (*GenConfigResponse, error) {
	resp := &GenConfigResponse{
		Frameworks:    []ConfigOptionItem{},
		Platforms:     []ConfigOptionItem{},
		Options:       []ConfigOptionItem{},
		Advanced:      []ConfigOptionItem{},
		Statuses:      []ConfigOptionItem{},
		ComponentLibs: map[string][]ConfigOptionItem{},
	}
	db := mysql.GetDB()
	if db == nil {
		return resp, nil
	}

	var rows []model.GenConfigOption
	if err := db.Where("is_active = ?", true).Order("sort_order asc, id asc").Find(&rows).Error; err != nil {
		return nil, err
	}

	for _, row := range rows {
		item := ConfigOptionItem{Label: row.Label, Value: row.Value, Icon: row.Icon, Platform: row.PlatformScope}
		switch row.Category {
		case model.ConfigCategoryFramework:
			resp.Frameworks = append(resp.Frameworks, item)
		case model.ConfigCategoryPlatform:
			resp.Platforms = append(resp.Platforms, item)
		case model.ConfigCategoryOption:
			resp.Options = append(resp.Options, item)
		case model.ConfigCategoryAdvanced:
			resp.Advanced = append(resp.Advanced, item)
		case model.ConfigCategoryStatus:
			resp.Statuses = append(resp.Statuses, item)
		case model.ConfigCategoryComponentLib:
			resp.ComponentLibs[row.Framework] = append(resp.ComponentLibs[row.Framework], item)
		}
	}
	return resp, nil
}

// SeedIfEmpty 表为空时灌入默认字典（幂等，只在完全没有记录时执行一次）
func (r *GenConfigRepository) SeedIfEmpty() error {
	db := mysql.GetDB()
	if db == nil {
		return nil
	}
	var cnt int64
	if err := db.Model(&model.GenConfigOption{}).Count(&cnt).Error; err != nil {
		return err
	}
	if cnt > 0 {
		return nil
	}
	defaults := defaultGenConfig()
	if err := db.Create(&defaults).Error; err != nil {
		return err
	}
	log.Printf("[GenConfig] 已初始化生成配置字典 %d 条", len(defaults))
	return nil
}

// defaultGenConfig 默认字典数据（组件库按 PC/手机/通用 标注 platform_scope）
func defaultGenConfig() []model.GenConfigOption {
	opt := func(cat, label, value, icon string, sort int) model.GenConfigOption {
		return model.GenConfigOption{
			Category: cat, Label: label, Value: value, Icon: icon,
			PlatformScope: model.PlatformScopeBoth, SortOrder: sort, IsActive: true,
		}
	}
	lib := func(fw, label, value, scope string, sort int) model.GenConfigOption {
		return model.GenConfigOption{
			Category: model.ConfigCategoryComponentLib, Framework: fw,
			Label: label, Value: value, PlatformScope: scope, SortOrder: sort, IsActive: true,
		}
	}
	return []model.GenConfigOption{
		// 目标框架
		opt(model.ConfigCategoryFramework, "Flutter", "Flutter", "", 1),
		opt(model.ConfigCategoryFramework, "React", "React", "", 2),
		opt(model.ConfigCategoryFramework, "Vue", "Vue", "", 3),
		// 目标平台
		opt(model.ConfigCategoryPlatform, "手机端", "mobile", "fa-solid fa-mobile-screen-button", 1),
		opt(model.ConfigCategoryPlatform, "桌面端", "desktop", "fa-solid fa-desktop", 2),
		// 生成选项
		opt(model.ConfigCategoryOption, "生成响应式布局", "responsive", "", 1),
		opt(model.ConfigCategoryOption, "添加注释", "comment", "", 2),
		opt(model.ConfigCategoryOption, "使用组件库", "component", "", 3),
		// 高级选项
		opt(model.ConfigCategoryAdvanced, "优化性能", "perf", "", 1),
		opt(model.ConfigCategoryAdvanced, "生成文档", "docs", "", 2),
		// 任务状态（供回显 label，颜色由前端按 value 的 CSS 类控制）
		opt(model.ConfigCategoryStatus, "已完成", "success", "", 1),
		opt(model.ConfigCategoryStatus, "生成中", "running", "", 2),
		opt(model.ConfigCategoryStatus, "排队中", "pending", "", 3),
		opt(model.ConfigCategoryStatus, "失败", "failed", "", 4),
		// 组件库 · Flutter
		lib("Flutter", "Material Design", "Material Design", model.PlatformScopeBoth, 1),
		lib("Flutter", "Cupertino", "Cupertino", model.PlatformScopeMobile, 2),
		lib("Flutter", "Flutter Element", "Flutter Element", model.PlatformScopeBoth, 3),
		lib("Flutter", "GetWidget", "GetWidget", model.PlatformScopeBoth, 4),
		// 组件库 · React
		lib("React", "Ant Design", "Ant Design", model.PlatformScopeDesktop, 1),
		lib("React", "Material UI (MUI)", "Material UI (MUI)", model.PlatformScopeDesktop, 2),
		lib("React", "Arco Design React", "Arco Design React", model.PlatformScopeDesktop, 3),
		lib("React", "Semi Design", "Semi Design", model.PlatformScopeDesktop, 4),
		lib("React", "React Vant", "React Vant", model.PlatformScopeMobile, 5),
		lib("React", "Chakra UI", "Chakra UI", model.PlatformScopeBoth, 6),
		lib("React", "Next UI", "Next UI", model.PlatformScopeDesktop, 7),
		// 组件库 · Vue
		lib("Vue", "Element Plus", "Element Plus", model.PlatformScopeDesktop, 1),
		lib("Vue", "Ant Design Vue 4.x", "Ant Design Vue 4.x", model.PlatformScopeDesktop, 2),
		lib("Vue", "Naive UI", "Naive UI", model.PlatformScopeDesktop, 3),
		lib("Vue", "Vant 4", "Vant 4", model.PlatformScopeMobile, 4),
		lib("Vue", "PrimeVue", "PrimeVue", model.PlatformScopeDesktop, 5),
		lib("Vue", "Arco Design Vue", "Arco Design Vue", model.PlatformScopeDesktop, 6),
		lib("Vue", "Quasar", "Quasar", model.PlatformScopeBoth, 7),
	}
}
