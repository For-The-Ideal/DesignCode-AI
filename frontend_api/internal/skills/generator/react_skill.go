package generator

import (
	"context"
	"fmt"
	"frontend_api/pkg/ai"
	"frontend_api/pkg/logger"
	"os"
	"path/filepath"
)

// ═══════════════════════════════════════════════
//  GeneratorSkill - React
//  职责：将 DSL 描述转换为 React + TypeScript 代码
// ═══════════════════════════════════════════════

// ReactSkill React 代码生成技能
type ReactSkill struct {
	aiClient ai.Client
	logger   *logger.Logger
}

// NewReactSkill 创建 React 生成技能
func NewReactSkill(client ai.Client) *ReactSkill {
	return &ReactSkill{
		aiClient: client,
		logger:   logger.NewLogger("react-gen"),
	}
}

// Name 返回技能名称
func (s *ReactSkill) Name() string {
	return "ReactGeneratorSkill"
}

// Execute 执行代码生成
func (s *ReactSkill) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	gi, ok := input.(Input)
	if !ok {
		return nil, fmt.Errorf("react: 输入类型错误")
	}

	s.logger.Infof("[ReactSkill] 开始生成 React 代码, DSL length=%d", len(gi.DSL))

	prompt := s.loadPrompt()
	if prompt == "" {
		return nil, fmt.Errorf("react: 提示词文件未找到")
	}

	messages := []ai.Message{
		{Role: "system", Content: prompt},
		{Role: "user", Content: fmt.Sprintf("Platform: %s\nDSL:\n%s", gi.Platform, gi.DSL)},
	}

	content, err := s.aiClient.Chat(ctx, messages, ai.WithDeepSeekThinking("high"), ai.WithMaxTokens(8192))
	if err != nil {
		return nil, fmt.Errorf("react: AI 调用失败: %w", err)
	}

	s.logger.Infof("[ReactSkill] 生成完成, code=%d chars", len(content))
	return Output{Code: content, Preview: "", Score: 0}, nil
}

func (s *ReactSkill) loadPrompt() string {
	paths := []string{
		"prompts/react/generate.txt",
		"../prompts/react/generate.txt",
		filepath.Join("..", "prompts", "react", "generate.txt"),
	}
	for _, p := range paths {
		data, err := os.ReadFile(p)
		if err == nil {
			return string(data)
		}
	}
	return ""
}
