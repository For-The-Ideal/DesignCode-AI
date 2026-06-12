package vision

import (
	"context"
	"encoding/json"
	"fmt"
	"frontend_api/internal/model"
	"frontend_api/pkg/ai"
	"frontend_api/pkg/logger"
	"os"
	"path/filepath"
)

// ═══════════════════════════════════════════════
//  VisionAnalyzeSkill
//  职责：将设计稿图片分析为结构化 DSL (JSON)
//  使用千问 Qwen3-VL-Plus 多模态模型
// ═══════════════════════════════════════════════

// Input 视觉分析输入
type Input struct {
	Images []string `json:"images"` // COS URL 数组
}

// Output 视觉分析输出
type Output struct {
	DSL string `json:"dsl"` // 标准 DSL JSON 字符串
}

// Skill 视觉分析技能
type Skill struct {
	aiClient ai.Client
	logger   *logger.Logger
}

// NewSkill 创建视觉分析技能
func NewSkill(client ai.Client) *Skill {
	return &Skill{
		aiClient: client,
		logger:   logger.NewLogger("vision-analyze"),
	}
}

// Name 返回技能名称
func (s *Skill) Name() string {
	return "VisionAnalyzeSkill"
}

// Execute 执行视觉分析
func (s *Skill) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	vi, ok := input.(Input)
	if !ok {
		return nil, fmt.Errorf("vision: 输入类型错误")
	}

	if len(vi.Images) == 0 {
		return nil, fmt.Errorf("vision: 至少需要一张图片")
	}

	s.logger.Infof("[VisionAnalyzeSkill] 开始分析设计稿 (%d 张图片)...", len(vi.Images))

	// 读取分析 prompt
	promptContent := s.loadPrompt()
	if promptContent == "" {
		return nil, fmt.Errorf("vision: 提示词文件未找到")
	}

	// 构建多模态消息：图片 + 分析指令
	content := make([]ai.ContentPart, 0, len(vi.Images)+1)
	for _, imgURL := range vi.Images {
		content = append(content, ai.ImagePart(imgURL))
	}
	content = append(content, ai.TextPart(promptContent))

	messages := []ai.Message{
		{Role: "user", Content: content},
	}

	if s.aiClient == nil {
		return nil, fmt.Errorf("vision: AI 客户端未配置")
	}

	resp, err := s.aiClient.Chat(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("vision: AI 分析失败: %w", err)
	}

	s.logger.Infof("[VisionAnalyzeSkill] AI 返回: %s", resp)

	// 校验返回的是合法 DSL JSON
	var dsl model.DSL
	if err := json.Unmarshal([]byte(resp), &dsl); err != nil {
		s.logger.Infof("[VisionAnalyzeSkill] AI 返回非标准 DSL JSON: %v", err)
		// 仍然原样返回，由上层处理
	}

	return Output{DSL: resp}, nil
}

// loadPrompt 加载 prompts/vision/analyze.txt
func (s *Skill) loadPrompt() string {
	paths := []string{
		"prompts/vision/analyze.txt",
		"../prompts/vision/analyze.txt",
		filepath.Join("..", "prompts", "vision", "analyze.txt"),
	}
	for _, p := range paths {
		data, err := os.ReadFile(p)
		if err == nil {
			return string(data)
		}
	}
	return ""
}
