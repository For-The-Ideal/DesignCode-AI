package model

// ═══════════════════════════════════════════════
//  DSL — 核心中间语言
//  Vision 输出 DSL，Generator 只认 DSL
//  components[].type / id 是生成器的关键依赖
//  position / style / children 为像素测量扩展字段（可选）
// ═══════════════════════════════════════════════

// Position 组件位置与尺寸（像素测量）
type Position struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

// Spacing 内外边距
type Spacing struct {
	Top    int `json:"top"`
	Right  int `json:"right"`
	Bottom int `json:"bottom"`
	Left   int `json:"left"`
}

// ComponentStyle 组件样式
type ComponentStyle struct {
	BackgroundColor string  `json:"backgroundColor,omitempty"`
	Color           string  `json:"color,omitempty"`
	FontSize        int     `json:"fontSize,omitempty"`
	FontWeight      string  `json:"fontWeight,omitempty"`
	BorderRadius    int     `json:"borderRadius,omitempty"`
	BoxShadow       string  `json:"boxShadow,omitempty"`
	Margin          Spacing `json:"margin"`
	Padding         Spacing `json:"padding"`
}

// DSLComponent DSL 组件描述
// 语义字段(type/id)是生成器必须的，测量字段(position/style/children)为可选扩展
type DSLComponent struct {
	Type     string          `json:"type"`               // image | input | button | text | card | navbar | ...
	ID       string          `json:"id"`                 // 组件唯一标识
	Text     string          `json:"text"`               // 文字内容（必填，无则空串）
	Position *Position       `json:"position,omitempty"` // 像素位置
	Style    *ComponentStyle `json:"style,omitempty"`    // 样式
	Children []DSLComponent  `json:"children,omitempty"` // 子组件（容器类）
}

// Canvas 画布尺寸
// 移动端固定宽度 375px，高度按视觉比例估算
// 桌面端固定宽度 1440px
type Canvas struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// TypographyScale 排版层级
type TypographyScale struct {
	Size   int `json:"size"`
	Weight int `json:"weight"`
}

// Typography 全局排版
type Typography struct {
	H1      TypographyScale `json:"h1"`
	H2      TypographyScale `json:"h2"`
	Body    TypographyScale `json:"body"`
	Caption TypographyScale `json:"caption"`
}

// DSL DSL 中间语言（Vision → Generator 的桥梁）
type DSL struct {
	PageName         string         `json:"page_name"`
	Layout           string         `json:"layout_type"`                // column | row | grid | absolute | mixed
	Canvas           *Canvas        `json:"canvas,omitempty"`           // 画布尺寸（移动端 375，桌面端 1440）
	PageBackground   string         `json:"pageBackground,omitempty"`   // 页面背景色
	PageTextColor    string         `json:"pageTextColor,omitempty"`    // 页面默认文字颜色
	PageBaseFontSize int            `json:"pageBaseFontSize,omitempty"` // 页面基准字号
	Components       []DSLComponent `json:"components"`
	GlobalColors     []string       `json:"globalColors,omitempty"` // 全局色板
	Typography       *Typography    `json:"typography,omitempty"`   // 全局排版
}
