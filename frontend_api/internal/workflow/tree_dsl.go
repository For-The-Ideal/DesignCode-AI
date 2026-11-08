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

// PageDSL 单页 DSL 构建结果
//
//	OK=false 表示该页没有结构树，调用方应对该页回退 Vision。
type PageDSL struct {
	Page model.ImageItem
	DSL  string
	OK   bool
}

// OrderedPages 按「入口页优先 + sort_order 升序」排列页面。
// 不修改传入切片，避免影响任务快照。
func OrderedPages(images []model.ImageItem) []model.ImageItem {
	ordered := make([]model.ImageItem, len(images))
	copy(ordered, images)
	sort.SliceStable(ordered, func(i, j int) bool {
		if ordered[i].IsEntry != ordered[j].IsEntry {
			return ordered[i].IsEntry
		}
		return ordered[i].SortOrder < ordered[j].SortOrder
	})
	return ordered
}

// BuildPerPageDSL 逐页构建 DSL（每页独立，平铺合并已废弃）。
// 多图模式下每张图是一页，各自的树只描述自己的页面。
func BuildPerPageDSL(images []model.ImageItem) []PageDSL {
	ordered := OrderedPages(images)
	out := make([]PageDSL, 0, len(ordered))
	for _, img := range ordered {
		dsl, ok := BuildPageDSL(img)
		out = append(out, PageDSL{Page: img, DSL: dsl, OK: ok})
	}
	return out
}

// BuildPageDSL 把单页的树转成 DSL JSON；无树返回 ok=false。
func BuildPageDSL(img model.ImageItem) (string, bool) {
	if len(img.Tree) == 0 {
		return "", false
	}
	cw, ch := img.Width, img.Height
	if cw <= 0 || ch <= 0 {
		cw, ch = maxExtent(img.Tree)
	}

	comps := make([]map[string]any, 0, len(img.Tree))
	for _, n := range img.Tree {
		comps = append(comps, nodeToComponent(n))
	}

	name := img.PageName
	if name == "" {
		name = "UI页面"
	}
	dsl := map[string]any{
		"page_id":     img.PageID,
		"page_name":   name,
		"layout_type": inferLayout(topBoxes(img.Tree)),
		"canvas":      map[string]any{"width": cw, "height": ch},
		"components":  comps,
	}
	if img.PageDesc != "" {
		dsl["page_desc"] = img.PageDesc
	}
	if img.Route != "" {
		dsl["route"] = img.Route
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
		"type": componentType(n),
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
	// 关联附图素材：生成方能从所属图片的 attached_assets 里按 assetId 取到可用图片
	if n.AssetID != "" {
		comp["assetId"] = n.AssetID
	}
	// 页面跳转：目标 page_id，生成方应译为 data-nav 属性（多页 App 模式）
	if n.LinkTo != "" {
		comp["linkTo"] = n.LinkTo
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

// componentType 决定 DSL type：≥2 级元素类型（component）优先，其次顶层区域枚举，最后 container。
func componentType(n model.UINode) string {
	if t, ok := mapElementType(n.Component); ok {
		// 选了「图片」却又有子级：说明这块是容器（底图 + 浮层），
		// 继续输出 image 会让生成侧把子级当成图片附带内容而忽略。
		if t == "image" && len(n.Children) > 0 {
			return "container"
		}
		return t
	}
	t := mapRegionType(n.Name)
	// banner 整块只有自己时是一张整图（image）；一旦拆出子级，它就是容纳底图/文案/按钮的容器。
	if t == "image" && len(n.Children) > 0 {
		return "container"
	}
	return t
}

// mapElementType 前端元素类型 → analyze.txt 的 type 词表；ok=false 表示未选/未知，交给区域名兜底。
func mapElementType(component string) (string, bool) {
	switch component {
	case "group", "listItem", "tabs", "tab":
		// DSL 词表没有这几个粒度：能继续下钻的仍按 container，子组件自行表达语义
		return "container", true
	case "card":
		return "card", true
	case "list":
		return "list", true
	case "grid":
		return "grid", true
	case "swiper":
		return "carousel", true
	case "form":
		return "form", true
	case "upload":
		// 上传控件：已选图片预览 + 添加槽 + 配额文案；与顶层区域 upload（整体区块）区分
		return "upload", true
	case "title", "text":
		return "text", true
	case "image":
		return "image", true
	case "icon":
		return "icon", true
	case "button", "link":
		// link 是文字链接，生成侧按可点按钮落
		return "button", true
	case "input", "select":
		return "input", true
	case "switch", "checkbox":
		return "checkbox", true
	case "avatar":
		return "avatar", true
	case "badge", "tag":
		return "badge", true
	case "divider":
		return "divider", true
	default:
		return "", false
	}
}

// mapRegionType 顶层区域枚举 → DSL type 词表；未知/子节点统一按 container 占位。
func mapRegionType(name string) string {
	switch name {
	case "header", "search":
		// search 多为「返回 + 输入框 + 取消」的横向条，与 navbar 同构；region 字段仍保留 search 原值
		return "navbar"
	case "footer", "tabbar":
		// footer 包含底部操作条/版权栏，沿用 tabbar 的「底部横向条」语义
		return "tabbar"
	case "banner":
		return "image"
	case "modal":
		return "modal"
	default: // sidebar/main/upload/floating/other/空 及子节点（upload 是容器：上传控件 + 已传预览）
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
