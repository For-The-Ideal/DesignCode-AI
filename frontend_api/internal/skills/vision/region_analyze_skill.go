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
	"header": true, "sidebar": true, "main": true, "footer": true,
	"banner": true, "modal": true, "floating": true, "other": true,
}

// layout 合法值（与前端 LAYOUT_VALUES 一致）
var layoutValues = map[string]bool{"row": true, "column": true, "grid": true}

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

// loadPrompt 加载 prompts/vision/region.txt
func (s *RegionSkill) loadPrompt() string {
	paths := []string{
		"prompts/vision/region.txt",
		"../prompts/vision/region.txt",
		filepath.Join("..", "prompts", "vision", "region.txt"),
	}
	for _, p := range paths {
		data, err := os.ReadFile(p)
		if err == nil {
			return string(data)
		}
	}
	return ""
}
