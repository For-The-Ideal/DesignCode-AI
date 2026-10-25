package model

// ═══════════════════════════════════════════════
//  UINode — 结构编辑器提交的树节点
//  仅声明生成 + 回显/重开用到的字段；前端提交的多余键
//  （parentId/childrenIds/aiFilled/sortIndex 等）
//  在 json.Unmarshal 时被自动忽略（可由 children 嵌套重建），便于前端演进。
//  坐标 bbox 为原图像素。
// ═══════════════════════════════════════════════

// BBox 节点在原图上的像素框
type BBox struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	W float64 `json:"w"`
	H float64 `json:"h"`
}

// UINode 树节点（嵌套结构，Children 为子节点）
type UINode struct {
	NodeKey      string   `json:"nodeKey"`
	Name         string   `json:"name"`         // 顶层区域枚举 header/sidebar/main/... 或子节点标签
	NodeType     string   `json:"nodeType"`     // 编辑器节点类型（容器/叶子等），重开回显所需
	Layout       string   `json:"layout"`       // row | column | grid | ""
	LayoutSource string   `json:"layoutSource"` // layout 来源：ai | manual（手动不被 AI 覆盖，重开需保留）
	VariantIndex int      `json:"variantIndex"` // 同一区域枚举的多块序号（回显时区分同名块）
	Summary      string   `json:"summary"`      // 语义描述
	StyleHints   string   `json:"styleHints"`   // 样式提示
	Behavior     string   `json:"behavior"`     // 交互行为
	Status       string   `json:"status"`       // pending | ai_draft | confirmed | empty
	BBox         BBox     `json:"bbox"`
	Children     []UINode `json:"children"`
}
