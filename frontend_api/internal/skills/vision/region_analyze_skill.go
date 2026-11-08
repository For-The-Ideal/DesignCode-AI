package vision

import (
	"context"
	"encoding/json"
	"fmt"
	"frontend_api/pkg/ai"
	"frontend_api/pkg/logger"
	"os"
	"path/filepath"
	"strings"
)

// ═══════════════════════════════════════════════
//  RegionAnalyzeSkill
//  职责：把整张设计稿识别为「顶层区域目录」（供前端框选顶层时按 name 命中自动填充描述）
//  使用千问 Qwen-VL 多模态模型（capability=read）
//  与 VisionAnalyzeSkill 并存：后者出整页 DSL，本 Skill 只出顶层区域语义目录
// ═══════════════════════════════════════════════

// 顶层区域枚举（与前端 utils/uiTree.js TOP_REGION_ENUM 严格一致）
var topRegionEnum = map[string]bool{
	"header": true, "search": true, "banner": true, "main": true, "sidebar": true,
	"footer": true, "tabbar": true, "modal": true, "floating": true, "other": true,
}

// layout 合法值（与前端 LAYOUT_VALUES 一致；stack=叠放，子元素按坐标叠加）
var layoutValues = map[string]bool{"row": true, "column": true, "grid": true, "stack": true}

// RegionInput 顶层识别输入
type RegionInput struct {
	ImageURL string `json:"imageUrl"` // 单张设计稿 COS URL
}

// Region 顶层区域描述（规范化后的输出单元）
type Region struct {
	Name         string `json:"name"`
	Summary      string `json:"summary"`
	StyleHints   string `json:"styleHints"`
	Behavior     string `json:"behavior"`
	ChildrenHint string `json:"childrenHint"`
	LayoutHint   string `json:"layoutHint"`
}

// RegionOutput 顶层识别输出
type RegionOutput struct {
	Regions []Region `json:"regions"`
}

// RegionSkill 顶层区域识别技能
type RegionSkill struct {
	aiClient ai.Client
	logger   *logger.Logger
}

// NewRegionSkill 创建顶层区域识别技能
func NewRegionSkill(client ai.Client) *RegionSkill {
	return &RegionSkill{
		aiClient: client,
		logger:   logger.NewLogger("region-analyze"),
	}
}

// Name 返回技能名称
func (s *RegionSkill) Name() string {
	return "RegionAnalyzeSkill"
}

// Execute 执行顶层区域识别
func (s *RegionSkill) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	ri, ok := input.(RegionInput)
	if !ok {
		return nil, fmt.Errorf("region: 输入类型错误")
	}
	if strings.TrimSpace(ri.ImageURL) == "" {
		return nil, fmt.Errorf("region: 缺少图片地址")
	}
	if s.aiClient == nil {
		return nil, fmt.Errorf("region: 视觉AI客户端未配置")
	}

	promptContent := s.loadPrompt()
	if promptContent == "" {
		return nil, fmt.Errorf("region: 提示词文件未找到")
	}

	s.logger.Infof("[RegionAnalyzeSkill] 开始识别顶层区域: %s", ri.ImageURL)

	content := []ai.ContentPart{
		ai.ImagePart(ri.ImageURL),
		ai.TextPart(promptContent),
	}
	messages := []ai.Message{
		{Role: "user", Content: content},
	}

	resp, err := s.aiClient.Chat(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("region: AI 识别失败: %w", err)
	}
	s.logger.Infof("[RegionAnalyzeSkill] AI 返回: %s", resp)

	regions, err := parseRegions(resp)
	if err != nil {
		return nil, fmt.Errorf("region: 解析AI返回失败: %w", err)
	}

	return RegionOutput{Regions: regions}, nil
}

// parseRegions 从 AI 文本中解析并规范化区域列表。
// 兼容：包裹 ```json、返回顶层是数组或 {regions:[...]}、name 非法归为 other、layout 非法清空。
func parseRegions(resp string) ([]Region, error) {
	jsonStr := extractJSON(resp)
	if jsonStr == "" {
		return nil, fmt.Errorf("未从返回中找到 JSON")
	}

	// 优先按对象 { "regions": [...] } 解析
	var envelope struct {
		Regions []Region `json:"regions"`
	}
	if err := json.Unmarshal([]byte(jsonStr), &envelope); err == nil && envelope.Regions != nil {
		return normalizeRegions(envelope.Regions), nil
	}

	// 退化：顶层直接是数组
	var arr []Region
	if err := json.Unmarshal([]byte(jsonStr), &arr); err == nil {
		return normalizeRegions(arr), nil
	}

	return nil, fmt.Errorf("返回结构不符合预期")
}

// normalizeRegions 逐条校验：name 必须命中枚举（否则归 other），去重同名，layoutHint 必须合法。
func normalizeRegions(in []Region) []Region {
	out := make([]Region, 0, len(in))
	seen := make(map[string]bool)
	for _, r := range in {
		name := strings.ToLower(strings.TrimSpace(r.Name))
		if !topRegionEnum[name] {
			name = "other"
		}
		if seen[name] {
			continue // 同名只保留第一条（前端按 name 精确匹配，多块由前端加序号）
		}
		seen[name] = true

		layout := strings.ToLower(strings.TrimSpace(r.LayoutHint))
		if !layoutValues[layout] {
			layout = ""
		}

		out = append(out, Region{
			Name:         name,
			Summary:      strings.TrimSpace(r.Summary),
			StyleHints:   strings.TrimSpace(r.StyleHints),
			Behavior:     strings.TrimSpace(r.Behavior),
			ChildrenHint: strings.TrimSpace(r.ChildrenHint),
			LayoutHint:   layout,
		})
	}
	return out
}

// extractJSON 去掉可能的 ```json ... ``` 包裹后，截取顶层 JSON 片段。
// 顶层数组形态（当前 prompt 约定）以 [ 开头，优先整体截 [ ... ]；否则退而截 { ... }。
func extractJSON(s string) string {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "```") {
		if i := strings.IndexByte(s, '\n'); i >= 0 {
			s = s[i+1:]
		}
		s = strings.TrimSuffix(strings.TrimSpace(s), "```")
	}
	s = strings.TrimSpace(s)
	// 顶层直接是数组：从首个 [ 截到末个 ]（对象数组内部含 {，不能再用 { 定位）
	if strings.HasPrefix(s, "[") {
		if end := strings.LastIndexByte(s, ']'); end > 0 {
			return s[:end+1]
		}
		return ""
	}
	start := strings.IndexByte(s, '{')
	end := strings.LastIndexByte(s, '}')
	if start >= 0 && end > start {
		return s[start : end+1]
	}
	return ""
}

// RegionDescribeInput 单区域 AI 描述输入
// 参数来自标注编辑器：主图 URL + 整页区域树（所有顶层区域及全部子级，含描述与附图引用）
type RegionDescribeInput struct {
	ImageURL string       `json:"imageUrl"`
	Tree     []RegionNode `json:"tree"` // 整页区域树，isFocus 标记本次要描述的节点
}

// RegionNode 区域树节点（children 递归任意层级，无子级为空数组）
// NodeKey：前端节点唯一标识，回填时按它精确匹配（不依赖 name，避免重名/改名）
type RegionNode struct {
	NodeKey    string           `json:"nodeKey"`
	Name       string           `json:"name"`
	Summary    string           `json:"summary"`
	StyleHints string           `json:"styleHints"`
	Layout     string           `json:"layout"`
	Behavior   string           `json:"behavior"`
	BBox       *RegionBBox      `json:"bbox"`
	Asset      *RegionAssetInfo `json:"asset"` // 该区域关联的附图，无引用为 null
	IsFocus    bool             `json:"isFocus"`
	Children   []RegionNode     `json:"children"`
}

// RegionBBox 区域在原图中的像素范围
type RegionBBox struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	W float64 `json:"w"`
	H float64 `json:"h"`
}

// RegionAssetInfo 附图素材（URL 为 COS 地址）
type RegionAssetInfo struct {
	Label string `json:"label"`
	URL   string `json:"url"`
}

// RegionNodeDesc AI 返回的单个节点描述（内部解析用：nodeKey + 描述文本）
type RegionNodeDesc struct {
	NodeKey     string `json:"nodeKey"`
	Description string `json:"description"`
}

// RegionDescribeOutput AI 描述输出：与请求 tree 完全同构的区域树，每个节点的 summary 已填入 AI 描述
type RegionDescribeOutput struct {
	Tree []RegionNode `json:"tree"`
}

// Describe 为单个区域生成一句话描述（多模态：整图 + 区域信息）
func (s *RegionSkill) Describe(ctx context.Context, in RegionDescribeInput) (RegionDescribeOutput, error) {
	if strings.TrimSpace(in.ImageURL) == "" {
		return RegionDescribeOutput{}, fmt.Errorf("describe: 缺少图片地址")
	}
	if s.aiClient == nil {
		return RegionDescribeOutput{}, fmt.Errorf("describe: 视觉AI客户端未配置")
	}
	tpl := s.loadPromptFile("describe.txt")
	if tpl == "" {
		return RegionDescribeOutput{}, fmt.Errorf("describe: 提示词文件未找到")
	}

	// 整页区域树/附图信息拼到提示词尾部（AI 需要为每个节点生成描述）
	var info strings.Builder
	info.WriteString("\n\n【页面区域结构】\n以下为整页全部区域（每个顶层节点含其全部子级），每行开头 [nodeKey=xxx] 是该节点的唯一标识（输出时原样照抄），缩进表示父子层级；行尾「← 当前区域」为编辑器当前选中项（仅参考）。\n")
	writeRegionTree(&info, in.Tree, 0)

	s.logger.Infof("[RegionAnalyzeSkill] 生成区域描述: %s focus=%s", in.ImageURL, regionFocusName(in.Tree))

	content := []ai.ContentPart{
		ai.ImagePart(in.ImageURL),
		ai.TextPart(tpl + info.String()),
	}
	messages := []ai.Message{
		{Role: "user", Content: content},
	}

	resp, err := s.aiClient.Chat(ctx, messages)
	if err != nil {
		return RegionDescribeOutput{}, fmt.Errorf("describe: AI 生成失败: %w", err)
	}
	s.logger.Infof("[RegionAnalyzeSkill] AI 返回: %s", resp)

	list := parseDescribeNodes(resp)
	if len(list) == 0 {
		return RegionDescribeOutput{}, fmt.Errorf("describe: AI 返回格式不符合要求（需 nodes 数组，每项仅含 nodeKey 与 description）")
	}
	// nodeKey → AI 描述；按 tree 同构回填（编造的 key 不会命中，自然被忽略）
	descMap := make(map[string]string, len(list))
	for _, d := range list {
		descMap[d.NodeKey] = d.Description
	}
	filledTree, filled := fillRegionTree(in.Tree, descMap)
	if filled == 0 {
		return RegionDescribeOutput{}, fmt.Errorf("describe: AI 返回的 nodeKey 与区域树不匹配")
	}
	s.logger.Infof("[RegionAnalyzeSkill] 描述完成: %s 回填 %d/%d 个节点", in.ImageURL, filled, countRegionNodes(in.Tree))
	return RegionDescribeOutput{Tree: filledTree}, nil
}

// fillRegionTree 按 nodeKey 把 AI 描述写回树（与请求 tree 同构），返回新树与命中数
func fillRegionTree(nodes []RegionNode, descMap map[string]string) ([]RegionNode, int) {
	out := make([]RegionNode, len(nodes))
	hit := 0
	for i, n := range nodes {
		if d := strings.TrimSpace(descMap[n.NodeKey]); d != "" {
			n.Summary = d
			hit++
		}
		children, childHit := fillRegionTree(n.Children, descMap)
		n.Children = children
		hit += childHit
		out[i] = n
	}
	return out, hit
}

// countRegionNodes 统计整棵树的节点数
func countRegionNodes(nodes []RegionNode) int {
	total := 0
	for _, n := range nodes {
		total += 1 + countRegionNodes(n.Children)
	}
	return total
}

// writeRegionTree 把区域树写成缩进文本（含布局/位置/描述/视觉/行为/附图，深度限制防异常数据）
func writeRegionTree(b *strings.Builder, nodes []RegionNode, depth int) {
	if depth > 8 {
		return
	}
	for _, n := range nodes {
		line := strings.Repeat("  ", depth) + "- [nodeKey=" + n.NodeKey + "] " + regionTextOrDefault(n.Name, "未命名")
		if l := strings.TrimSpace(n.Layout); l != "" {
			line += "｜布局 " + layoutText(l)
		}
		if bb := n.BBox; bb != nil {
			line += fmt.Sprintf("｜位置 x=%.0f y=%.0f w=%.0f h=%.0f", bb.X, bb.Y, bb.W, bb.H)
		}
		if s := strings.TrimSpace(n.Summary); s != "" {
			line += "｜描述: " + s
		}
		if s := strings.TrimSpace(n.StyleHints); s != "" {
			line += "｜视觉: " + s
		}
		if s := strings.TrimSpace(n.Behavior); s != "" {
			line += "｜行为: " + s
		}
		if a := n.Asset; a != nil && strings.TrimSpace(a.URL) != "" {
			line += "｜附图: " + regionTextOrDefault(a.Label, "未命名素材") + "(" + strings.TrimSpace(a.URL) + ")"
		}
		if n.IsFocus {
			line += "  ← 当前区域"
		}
		b.WriteString(line + "\n")
		writeRegionTree(b, n.Children, depth+1)
	}
}

// layoutText 把 layout 代码翻成生成模型能直接理解的中文（stack 必须点明可重叠）
func layoutText(v string) string {
	switch v {
	case "row":
		return "横排 row"
	case "column":
		return "竖排 column"
	case "grid":
		return "网格 grid"
	case "stack":
		return "叠放 stack（子元素不排队，按原图坐标叠加，允许重叠）"
	}
	return v
}

// regionTextOrDefault 去空白后为空时给默认值
func regionTextOrDefault(v, def string) string {
	if s := strings.TrimSpace(v); s != "" {
		return s
	}
	return def
}

// regionFocusName 取 isFocus 节点的名称（日志用）
func regionFocusName(nodes []RegionNode) string {
	for _, n := range nodes {
		if n.IsFocus {
			return regionTextOrDefault(n.Name, "未命名")
		}
		if name := regionFocusName(n.Children); name != "" {
			return name
		}
	}
	return ""
}

// parseDescribeNodes 从 AI 返回中提取整树描述列表（每项 nodeKey + description）。
// 严格模式：响应必须是 {"nodes":[{"nodeKey","description"}]}；
// 字段乱写、返回整棵树、纯文本等一律判失败，不做文本兜底。
func parseDescribeNodes(resp string) []RegionNodeDesc {
	jsonStr := extractJSON(resp)
	if jsonStr == "" {
		return nil
	}
	var out struct {
		Nodes []RegionNodeDesc `json:"nodes"`
	}
	if err := json.Unmarshal([]byte(jsonStr), &out); err != nil {
		return nil
	}
	list := make([]RegionNodeDesc, 0, len(out.Nodes))
	for _, n := range out.Nodes {
		key := strings.TrimSpace(n.NodeKey)
		desc := strings.TrimSpace(n.Description)
		if key == "" || desc == "" {
			continue
		}
		list = append(list, RegionNodeDesc{NodeKey: key, Description: desc})
	}
	return list
}

// loadPrompt 加载 prompts/vision/region.txt
func (s *RegionSkill) loadPrompt() string {
	return s.loadPromptFile("region.txt")
}

// loadPromptFile 从 prompts/vision 目录加载指定提示词（兼容不同工作目录）
func (s *RegionSkill) loadPromptFile(name string) string {
	paths := []string{
		"prompts/vision/" + name,
		"../prompts/vision/" + name,
		filepath.Join("..", "prompts", "vision", name),
	}
	for _, p := range paths {
		data, err := os.ReadFile(p)
		if err == nil {
			return string(data)
		}
	}
	return ""
}
