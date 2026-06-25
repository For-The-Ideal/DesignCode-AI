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
//  GeneratorSkill - Vue3
//  职责：将 DSL 描述转换为 Vue 3 代码
// ═══════════════════════════════════════════════

// Vue3Skill Vue3 代码生成技能
type Vue3Skill struct {
	aiClient ai.Client
	logger   *logger.Logger
}

// NewVue3Skill 创建 Vue3 生成技能
func NewVue3Skill(client ai.Client) *Vue3Skill {
	return &Vue3Skill{
		aiClient: client,
		logger:   logger.NewLogger("vue3-gen"),
	}
}

// Name 返回技能名称
func (s *Vue3Skill) Name() string {
	return "Vue3GeneratorSkill"
}

// Execute 执行代码生成
func (s *Vue3Skill) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	gi, ok := input.(Input)
	if !ok {
		return nil, fmt.Errorf("vue3: 输入类型错误")
	}

	s.logger.Infof("[Vue3Skill] 开始生成 Vue3 代码, DSL length=%d", len(gi.DSL))

	prompt := s.loadPrompt()
	if prompt == "" {
		return nil, fmt.Errorf("vue3: 提示词文件未找到")
	}

	messages := []ai.Message{
		{Role: "system", Content: prompt},
		{Role: "user", Content: fmt.Sprintf("Platform: %s\nOptions: %v\nComponentLib: %s\nDSL:\n%s", gi.Platform, gi.Options, gi.ComponentLib, gi.DSL)},
	}

	content, err := s.aiClient.Chat(ctx, messages, ai.WithDeepSeekThinking("high"), ai.WithMaxTokens(8192))
	if err != nil {
		return nil, fmt.Errorf("vue3: AI 调用失败: %w", err)
	}

	s.logger.Infof("[Vue3Skill] 生成完成, code=%d chars", len(content))
	return Output{Code: content, Preview: "", Score: 0}, nil
}

func (s *Vue3Skill) loadPrompt() string {
	paths := []string{
		"prompts/vue3/generate.txt",
		"../prompts/vue3/generate.txt",
		filepath.Join("..", "prompts", "vue3", "generate.txt"),
	}
	for _, p := range paths {
		data, err := os.ReadFile(p)
		if err == nil {
			return string(data)
		}
	}
	return ""
}
