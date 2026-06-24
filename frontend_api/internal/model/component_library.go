package model

import "time"

// ═══════════════════════════════════════════════
//  ComponentLibrary — 组件库配置表
//  对应数据库表：component_library
//  存储各框架可用的 UI 组件库 / CSS 框架列表
// ═══════════════════════════════════════════════

// ComponentLibrary 组件库 / UI 框架
type ComponentLibrary struct {
	ID          int       `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"size:100;not null"`           // 组件库名称，如 Element Plus
	Framework   string    `json:"framework" gorm:"size:20;index;not null"` // 所属框架：vue2 | vue3 | react | flutter
	Platform    string    `json:"platform" gorm:"size:20;default:'all'"`   // 适用平台：mobile | desktop | all
	Category    string    `json:"category" gorm:"size:50"`                 // 分类：ui-lib | css-framework
	Description string    `json:"description" gorm:"type:text"`            // 描述说明
	DocsURL     string    `json:"docs_url" gorm:"size:500"`                // 文档地址
	IsActive    bool      `json:"is_active" gorm:"default:true"`           // 是否启用
	SortOrder   int       `json:"sort_order" gorm:"default:0"`             // 排序序号
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName 指定表名
func (ComponentLibrary) TableName() string {
	return "component_library"
}

// ── 预置种子数据 ──────────────────────────────

// SeedComponentLibraries 返回初始化种子数据
func SeedComponentLibraries() []ComponentLibrary {
	now := time.Now()
	return []ComponentLibrary{
		// ── Vue2 ──
		{Name: "Element UI", Framework: "vue2", Platform: "desktop", Category: "ui-lib", Description: "Element UI 是一套为开发者、设计师和产品经理准备的基于 Vue 2.0 的桌面端组件库", DocsURL: "https://element.eleme.io", SortOrder: 1, IsActive: true, CreatedAt: now, UpdatedAt: now},
		{Name: "Vant 2", Framework: "vue2", Platform: "mobile", Category: "ui-lib", Description: "轻量、可靠的移动端 Vue 组件库，Vant 2 适用于 Vue 2", DocsURL: "https://vant-contrib.gitee.io/vant/v2", SortOrder: 2, IsActive: true, CreatedAt: now, UpdatedAt: now},
		{Name: "View UI", Framework: "vue2", Platform: "desktop", Category: "ui-lib", Description: "基于 Vue.js 2 的企业级 UI 组件库（原名 iView）", DocsURL: "https://www.iviewui.com", SortOrder: 3, IsActive: true, CreatedAt: now, UpdatedAt: now},
		{Name: "Ant Design Vue 1.x", Framework: "vue2", Platform: "desktop", Category: "ui-lib", Description: "Ant Design 的 Vue 实现，遵循 Ant Design 设计规范", DocsURL: "https://www.antdv.com/docs/vue/introduce-cn", SortOrder: 4, IsActive: true, CreatedAt: now, UpdatedAt: now},
		{Name: "Mint UI", Framework: "vue2", Platform: "mobile", Category: "ui-lib", Description: "基于 Vue.js 的移动端组件库", DocsURL: "https://mint-ui.github.io", SortOrder: 5, IsActive: true, CreatedAt: now, UpdatedAt: now},

		// ── Vue3 ──
		{Name: "Element Plus", Framework: "vue3", Platform: "desktop", Category: "ui-lib", Description: "基于 Vue 3 的桌面端组件库，Element UI 的 Vue 3 升级版", DocsURL: "https://element-plus.org", SortOrder: 1, IsActive: true, CreatedAt: now, UpdatedAt: now},
		{Name: "Ant Design Vue 4.x", Framework: "vue3", Platform: "desktop", Category: "ui-lib", Description: "Ant Design 的 Vue 3 实现，支持 Composition API", DocsURL: "https://www.antdv.com", SortOrder: 2, IsActive: true, CreatedAt: now, UpdatedAt: now},
		{Name: "Naive UI", Framework: "vue3", Platform: "desktop", Category: "ui-lib", Description: "一个 Vue 3 组件库，使用 TypeScript 编写，按需引入", DocsURL: "https://www.naiveui.com", SortOrder: 3, IsActive: true, CreatedAt: now, UpdatedAt: now},
		{Name: "Vant 4", Framework: "vue3", Platform: "mobile", Category: "ui-lib", Description: "轻量、可靠的移动端 Vue 3 组件库", DocsURL: "https://vant-ui.github.io", SortOrder: 4, IsActive: true, CreatedAt: now, UpdatedAt: now},
		{Name: "PrimeVue", Framework: "vue3", Platform: "all", Category: "ui-lib", Description: "丰富的 Vue 3 UI 组件库，同时支持桌面和移动端", DocsURL: "https://primevue.org", SortOrder: 5, IsActive: true, CreatedAt: now, UpdatedAt: now},
		{Name: "Arco Design Vue", Framework: "vue3", Platform: "desktop", Category: "ui-lib", Description: "字节跳动出品的中后台 Vue 3 组件库", DocsURL: "https://arco.design/vue", SortOrder: 6, IsActive: true, CreatedAt: now, UpdatedAt: now},
		{Name: "Quasar", Framework: "vue3", Platform: "all", Category: "ui-lib", Description: "基于 Vue 3 的高性能 UI 框架，一套代码适配多端", DocsURL: "https://quasar.dev", SortOrder: 7, IsActive: true, CreatedAt: now, UpdatedAt: now},

		// ── React ──
		{Name: "Ant Design", Framework: "react", Platform: "desktop", Category: "ui-lib", Description: "阿里巴巴中后台组件库，React 生态最流行的 UI 库", DocsURL: "https://ant.design", SortOrder: 1, IsActive: true, CreatedAt: now, UpdatedAt: now},
		{Name: "Material UI (MUI)", Framework: "react", Platform: "all", Category: "ui-lib", Description: "Google Material Design 的 React 实现，设计规范成熟", DocsURL: "https://mui.com", SortOrder: 2, IsActive: true, CreatedAt: now, UpdatedAt: now},
		{Name: "Arco Design React", Framework: "react", Platform: "desktop", Category: "ui-lib", Description: "字节跳动出品的企业级 React 组件库", DocsURL: "https://arco.design/react", SortOrder: 3, IsActive: true, CreatedAt: now, UpdatedAt: now},
		{Name: "Semi Design", Framework: "react", Platform: "desktop", Category: "ui-lib", Description: "字节跳动抖音前端与 UED 团队出品的中后台组件库", DocsURL: "https://semi.design", SortOrder: 4, IsActive: true, CreatedAt: now, UpdatedAt: now},
		{Name: "React Vant", Framework: "react", Platform: "mobile", Category: "ui-lib", Description: "移动端 React 组件库，Vant 的 React 版本", DocsURL: "https://react-vant.3lang.dev", SortOrder: 5, IsActive: true, CreatedAt: now, UpdatedAt: now},
		{Name: "Chakra UI", Framework: "react", Platform: "all", Category: "ui-lib", Description: "简洁、可访问的 React 组件库，注重开发体验", DocsURL: "https://chakra-ui.com", SortOrder: 6, IsActive: true, CreatedAt: now, UpdatedAt: now},
		{Name: "Next UI", Framework: "react", Platform: "all", Category: "ui-lib", Description: "基于 TailwindCSS 的 React UI 库，美观且可定制", DocsURL: "https://nextui.org", SortOrder: 7, IsActive: true, CreatedAt: now, UpdatedAt: now},

		// ── Flutter ──
		{Name: "Material Design", Framework: "flutter", Platform: "all", Category: "ui-lib", Description: "Flutter 默认内置的 Material Design 组件库，支持多端", DocsURL: "https://docs.flutter.dev/ui/widgets/material", SortOrder: 1, IsActive: true, CreatedAt: now, UpdatedAt: now},
		{Name: "Cupertino", Framework: "flutter", Platform: "mobile", Category: "ui-lib", Description: "Flutter 内置的 iOS 风格组件库，模拟原生 iOS 设计", DocsURL: "https://docs.flutter.dev/ui/widgets/cupertino", SortOrder: 2, IsActive: true, CreatedAt: now, UpdatedAt: now},
		{Name: "Flutter Element", Framework: "flutter", Platform: "desktop", Category: "ui-lib", Description: "类 Element 风格的 Flutter 组件库", DocsURL: "https://github.com/flutter-studio/flutter-element", SortOrder: 3, IsActive: true, CreatedAt: now, UpdatedAt: now},
		{Name: "GetWidget", Framework: "flutter", Platform: "all", Category: "ui-lib", Description: "Flutter 开源 UI 库，提供 1000+ 预置组件", DocsURL: "https://getwidget.dev", SortOrder: 4, IsActive: true, CreatedAt: now, UpdatedAt: now},
	}
}
