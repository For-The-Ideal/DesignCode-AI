package workflow

import (
	"encoding/json"
	"frontend_api/internal/model"
	"math"
	"sort"
)

// ═══════════════════════════════════════════════
//  tree_dsl：把结构编辑器提交的树（做法 B）转成 DSL
//  产出与 prompts/vision/analyze.txt 同 schema 的纯 JSON，
//  因此 Generator / CodeToHtml 无需感知数据来源。
//  有 tree → 本地拼装、跳过 Vision 大模型识别；
//  无 tree → 由调用方回退到 VisionAnalyzeSkill。
// ═══════════════════════════════════════════════

// BuildDSLFromTrees 汇总所有带树的图片，构建 DSL JSON 字符串。
// 返回 ok=false 表示没有任何树，调用方应回退 Vision。
func BuildDSLFromTrees(images []model.ImageItem) (string, bool) {
	// 按 sort_order 升序消费图片（多图页序由此字段权威决定，不依赖数组顺序）
	ordered := make([]model.ImageItem, len(images))
	copy(ordered, images)
	sort.SliceStable(ordered, func(i, j int) bool {
		return ordered[i].SortOrder < ordered[j].SortOrder
	})

	var tops []model.UINode
	cw, ch := 0, 0
	for _, img := range ordered {
		if len(img.Tree) == 0 {
			continue
		}
		tops = append(tops, img.Tree...)
		if cw == 0 && img.Width > 0 && img.Height > 0 {
			cw, ch = img.Width, img.Height
		}
	}
	if len(tops) == 0 {
		return "", false
	}
	// 宽高缺失时按框的最大外接矩形兜底
	if cw <= 0 || ch <= 0 {
		cw, ch = maxExtent(tops)
	}

	comps := make([]map[string]any, 0, len(tops))
	for _, n := range tops {
		comps = append(comps, nodeToComponent(n))
	}

	dsl := map[string]any{
		"page_name":   "UI页面",
		"layout_type": inferLayout(topBoxes(tops)),
		"canvas":      map[string]any{"width": cw, "height": ch},
		"components":  comps,
	}
	b, err := json.MarshalIndent(dsl, "", "  ")
	if err != nil {
		return "", false
	}
	return string(b), true
}

// nodeToComponent 递归把树节点映射为 DSL 组件。
// 对齐 analyze.txt 的必备键（type/id/text/position/children），
// 并附带 region/layout/desc/styleHints/behavior 供生成模型参考。
func nodeToComponent(n model.UINode) map[string]any {
	comp := map[string]any{
		"type": mapRegionType(n.Name),
		"id":   componentID(n),
		"text": "",
		"position": map[string]any{
			"x":      int(math.Round(n.BBox.X)),
			"y":      int(math.Round(n.BBox.Y)),
			"width":  int(math.Round(n.BBox.W)),
			"height": int(math.Round(n.BBox.H)),
		},
	}
	if n.Name != "" {
		comp["region"] = n.Name
	}
	if n.Layout != "" {
		comp["layout"] = n.Layout
	}
	if n.Summary != "" {
		comp["desc"] = n.Summary
	}
	if n.StyleHints != "" {
		comp["styleHints"] = n.StyleHints
	}
	if n.Behavior != "" {
		comp["behavior"] = n.Behavior
	}
	if len(n.Children) > 0 {
		kids := make([]map[string]any, 0, len(n.Children))
		for _, c := range n.Children {
			kids = append(kids, nodeToComponent(c))
		}
		comp["children"] = kids
	}
	return comp
}

// mapRegionType 顶层区域枚举 → DSL type 词表；未知/子节点统一按 container 占位。
func mapRegionType(name string) string {
	switch name {
	case "header":
		return "navbar"
	case "footer":
		return "tabbar"
	case "banner":
		return "image"
	case "modal":
		return "modal"
	default: // sidebar/main/floating/other/空 及子节点
		return "container"
	}
}

// componentID 优先用可读的区域名，重复时补 nodeKey 保证唯一。
func componentID(n model.UINode) string {
	if n.Name != "" {
		return n.Name + "_" + n.NodeKey
	}
	return n.NodeKey
}

// topBoxes 取顶层节点框（用于推断整页排版方向）。
func topBoxes(tops []model.UINode) []model.BBox {
	boxes := make([]model.BBox, 0, len(tops))
	for _, n := range tops {
		boxes = append(boxes, n.BBox)
	}
	return boxes
}

// maxExtent 递归求所有节点的最大外接右下角，作为画布尺寸兜底。
func maxExtent(tops []model.UINode) (int, int) {
	mr, mb := 0.0, 0.0
	var walk func(nodes []model.UINode)
	walk = func(nodes []model.UINode) {
		for _, n := range nodes {
			if r := n.BBox.X + n.BBox.W; r > mr {
				mr = r
			}
			if b := n.BBox.Y + n.BBox.H; b > mb {
				mb = b
			}
			walk(n.Children)
		}
	}
	walk(tops)
	return int(math.Round(mr)), int(math.Round(mb))
}

// inferLayout 依据同级框的中心点散布推断排版方向（与前端 inferLayout 同规则）。
func inferLayout(boxes []model.BBox) string {
	valid := make([]model.BBox, 0, len(boxes))
	for _, b := range boxes {
		if b.W > 0 && b.H > 0 {
			valid = append(valid, b)
		}
	}
	if len(valid) <= 1 {
		return "column"
	}
	const threshold = 20.0
	minX, maxX := valid[0].X, valid[0].X
	minY, maxY := valid[0].Y, valid[0].Y
	for _, b := range valid {
		cx, cy := b.X+b.W/2, b.Y+b.H/2
		if cx < minX {
			minX = cx
		}
		if cx > maxX {
			maxX = cx
		}
		if cy < minY {
			minY = cy
		}
		if cy > maxY {
			maxY = cy
		}
	}
	xSpread := maxX-minX > threshold
	ySpread := maxY-minY > threshold
	switch {
	case xSpread && ySpread:
		return "grid"
	case xSpread:
		return "row"
	case ySpread:
		return "column"
	default:
		return "column"
	}
}
