package model

// ═══════════════════════════════════════════════
//  DSL — 核心中间语言
//  Vision 输出 DSL，Generator 只认 DSL
// ═══════════════════════════════════════════════

// DSLComponent DSL 组件描述
type DSLComponent struct {
	Type string `json:"type"` // image | input | password | button | text | card | ...
	ID   string `json:"id"`   // 组件唯一标识
}

// DSL DSL 中间语言（Vision → Generator 的桥梁）
type DSL struct {
	PageName   string         `json:"page_name"`
	Layout     string         `json:"layout"` // column | row | grid
	Components []DSLComponent `json:"components"`
}
