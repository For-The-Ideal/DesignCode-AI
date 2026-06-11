package vision

import (
	"context"
	"encoding/json"
	"fmt"
	"frontend_api/internal/model"
	"frontend_api/pkg/logger"
	"os"
	"path/filepath"
)

// ═══════════════════════════════════════════════
//  VisionAnalyzeSkill
//  职责：将设计稿图片分析为结构化 DSL (JSON)
//  输出格式：
//    {
//      "page_name": "LoginPage",
//      "layout": "column",
//      "components": [
//        {"type":"image","id":"logo"},
//        {"type":"input","id":"username"},
//        {"type":"password","id":"password"},
//        {"type":"button","id":"login"}
//      ]
//    }
// ═══════════════════════════════════════════════

// Input 视觉分析输入
type Input struct {
	Images []model.TaskImage `json:"images"`
}

// Output 视觉分析输出
type Output struct {
	DSL string `json:"dsl"` // 标准 DSL JSON 字符串
}

// Skill 视觉分析技能
type Skill struct {
	logger *logger.Logger
}

// NewSkill 创建视觉分析技能
func NewSkill() *Skill {
	return &Skill{
		logger: logger.NewLogger("vision-analyze"),
	}
}

// Name 返回技能名称
func (s *Skill) Name() string {
	return "VisionAnalyzeSkill"
}

// Execute 执行视觉分析
// 当前为 Mock 实现，返回标准 LoginPage DSL
// TODO: 接入 AI 模型进行真实图片分析
func (s *Skill) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	s.logger.Info("[VisionAnalyzeSkill] 开始分析设计稿...")

	vi, ok := input.(Input)
	if !ok {
		return nil, fmt.Errorf("vision: 输入类型错误")
	}

	if len(vi.Images) == 0 {
		return nil, fmt.Errorf("vision: 至少需要一张图片")
	}

	// 读取分析 prompt
	promptContent := s.loadPrompt()
	_ = promptContent // TODO: 传给 AI 模型

	// Mock: 返回标准 LoginPage DSL
	dsl := model.DSL{
		PageName: "LoginPage",
		Layout:   "column",
		Components: []model.DSLComponent{
			{Type: "image", ID: "logo"},
			{Type: "input", ID: "username"},
			{Type: "password", ID: "password"},
			{Type: "button", ID: "login"},
		},
	}

	dslBytes, _ := json.MarshalIndent(dsl, "", "  ")
	s.logger.Infof("[VisionAnalyzeSkill] 分析完成，DSL: %s", string(dslBytes))

	return Output{DSL: string(dslBytes)}, nil
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
