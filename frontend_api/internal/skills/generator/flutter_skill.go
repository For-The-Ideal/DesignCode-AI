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
//  GeneratorSkill - Flutter
//  职责：将 DSL 描述转换为 Flutter/Dart 代码
// ═══════════════════════════════════════════════

// FlutterSkill Flutter 代码生成技能
type FlutterSkill struct {
	aiClient ai.Client
	logger   *logger.Logger
}

// NewFlutterSkill 创建 Flutter 生成技能
func NewFlutterSkill(client ai.Client) *FlutterSkill {
	return &FlutterSkill{
		aiClient: client,
		logger:   logger.NewLogger("flutter-gen"),
	}
}

// Name 返回技能名称
func (s *FlutterSkill) Name() string {
	return "FlutterGeneratorSkill"
}

// Execute 执行代码生成
func (s *FlutterSkill) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	gi, ok := input.(Input)
	if !ok {
		return nil, fmt.Errorf("flutter: 输入类型错误")
	}

	s.logger.Infof("[FlutterSkill] 开始生成 Flutter 代码, DSL length=%d", len(gi.DSL))

	prompt := s.loadPrompt()
	if prompt == "" {
		return nil, fmt.Errorf("flutter: 提示词文件未找到")
	}

	messages := []ai.Message{
		{Role: "system", Content: prompt},
		{Role: "user", Content: fmt.Sprintf("Platform: %s\nDSL:\n%s", gi.Platform, gi.DSL)},
	}

	content, err := s.aiClient.Chat(ctx, messages, ai.WithDeepSeekThinking("high"), ai.WithMaxTokens(8192))
	if err != nil {
		return nil, fmt.Errorf("flutter: AI 调用失败: %w", err)
	}

	s.logger.Infof("[FlutterSkill] 生成完成, code=%d chars", len(content))
	return Output{Code: content, Preview: "", Score: 0}, nil
}

func (s *FlutterSkill) loadPrompt() string {
	paths := []string{
		"prompts/flutter/generate.txt",
		"../prompts/flutter/generate.txt",
		filepath.Join("..", "prompts", "flutter", "generate.txt"),
	}
	for _, p := range paths {
		data, err := os.ReadFile(p)
		if err == nil {
			return string(data)
		}
	}
	return ""
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
