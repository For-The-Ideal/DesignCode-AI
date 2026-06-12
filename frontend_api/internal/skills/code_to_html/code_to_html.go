package code_to_html

import (
	"context"
	"fmt"
	"frontend_api/pkg/ai"
	"frontend_api/pkg/logger"
	"os"
	"path/filepath"
	"strings"
)

// ═══════════════════════════════════════════════
//  CodeToHtmlSkill
//  职责：将 AI 生成的框架代码转换为 HTML 预览
//  输入：源代码 + 框架类型
//  输出：可供手机模拟器渲染的完整 HTML
// ═══════════════════════════════════════════════

// Skill 代码转 HTML 技能
type Skill struct {
	aiClient ai.Client
	logger   *logger.Logger
}

// Input 转换输入
type Input struct {
	SourceCode string `json:"source_code"` // 原始生成代码
	Framework  string `json:"framework"`   // flutter | react | vue3
	DSL        string `json:"dsl,omitempty"`
	Images     string `json:"images,omitempty"` // 原始设计稿 URL（逗号分隔）
}

// Output 转换输出
type Output struct {
	HTML string `json:"html"` // 转换后的 HTML 预览代码
}

// NewSkill 创建代码转 HTML 技能
func NewSkill(client ai.Client) *Skill {
	return &Skill{
		aiClient: client,
		logger:   logger.NewLogger("code-to-html"),
	}
}

// Name 返回技能名称
func (s *Skill) Name() string {
	return "CodeToHtmlSkill"
}

// Execute 执行转换
func (s *Skill) Execute(ctx context.Context, input interface{}) (interface{}, error) {
	gi, ok := input.(Input)
	if !ok {
		return nil, fmt.Errorf("code-to-html: 输入类型错误")
	}

	s.logger.Infof("[CodeToHtmlSkill] 开始转换, framework=%s, code_length=%d", gi.Framework, len(gi.SourceCode))
	s.logger.Infof("[CodeToHtmlSkill] 原始代码:\n%s", gi.SourceCode)

	prompt := s.loadPrompt()
	if prompt == "" {
		return nil, fmt.Errorf("code-to-html: 提示词文件未找到")
	}

	userMsg := fmt.Sprintf("SOURCE_CODE:\n%s\n\nFRAMEWORK: %s", gi.SourceCode, gi.Framework)
	if gi.DSL != "" {
		userMsg += fmt.Sprintf("\n\nDSL: %s", gi.DSL)
	}
	if gi.Images != "" {
		userMsg += fmt.Sprintf("\n\nSCREENSHOTS: %s", gi.Images)
	}

	messages := []ai.Message{
		{Role: "system", Content: prompt},
		{Role: "user", Content: userMsg},
	}

	content, err := s.aiClient.Chat(ctx, messages, ai.WithMaxTokens(8192))
	if err != nil {
		return nil, fmt.Errorf("code-to-html: AI 调用失败: %w", err)
	}

	// 清理 AI 输出：去除 markdown 代码块包裹
	content = strings.TrimSpace(content)
	// 去掉开头的 ```html / ``` / ```HTML
	if strings.HasPrefix(content, "```") {
		if idx := strings.Index(content, "\n"); idx != -1 {
			content = content[idx+1:]
		}
	}
	// 去掉结尾的 ```
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)

	s.logger.Infof("[CodeToHtmlSkill] 转换完成, html=%d chars", len(content))
	return Output{HTML: content}, nil
}

func (s *Skill) loadPrompt() string {
	paths := []string{
		"prompts/code-to-html/generate.txt",
		"../prompts/code-to-html/generate.txt",
		filepath.Join("..", "prompts", "code-to-html", "generate.txt"),
	}
	for _, p := range paths {
		data, err := os.ReadFile(p)
		if err == nil {
			return string(data)
		}
	}
	s.logger.Warn("[CodeToHtmlSkill] 未找到 prompts/code-to-html/generate.txt")
	return ""
}
