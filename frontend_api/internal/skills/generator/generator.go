package generator

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"frontend_api/pkg/ai"
)

// ═══════════════════════════════════════════════
//  GeneratorSkill 通用类型定义
// ═══════════════════════════════════════════════

// Input 生成器通用输入
type Input struct {
	DSL          string   `json:"dsl"`
	Prompt       string   `json:"prompt,omitempty"`   // 从 prompts/{target}/generate.txt 加载
	Images       []string `json:"images,omitempty"`   // base64 编码的图片（通用模型兼容）
	Platform     string   `json:"platform,omitempty"` // mobile | desktop | tablet
	Options      []string `json:"options,omitempty"`
	ComponentLib string   `json:"component_lib,omitempty"`
}

// Output 生成器通用输出
//
// 注意：生成器不再返回分数。质量由生成后的确定性规则扫描（internal/quality）
type Output struct {
	Code    string `json:"code"`
	Preview string `json:"preview"`
}

// Skill 生成器技能通用接口
type Skill interface {
	Name() string
	Execute(ctx context.Context, input interface{}) (interface{}, error)
}

// GeneratorFunc 将 DSL + target 映射到对应的生成器技能
func GeneratorFunc(target string, client ai.Client) Skill {
	switch target {
	case "flutter":
		return NewFlutterSkill(client)
	case "vue3":
		return NewVue3Skill(client)
	case "react":
		return NewReactSkill(client)
	default:
		return NewFlutterSkill(client)
	}
}

// ═══════════════════════════════════════════════
//  Prompt 拼装：公共段（base/generate.txt）+ 框架段（<target>/generate.txt）
//
//  合并只为维护方便，运行时仍是两段拼接：
//  公共段在前（输出契约 + R1-R5 + token 约束），框架段在后（技术栈差异），
//  两段冲突时以公共段为准。
// ═══════════════════════════════════════════════

// LoadPrompt 读取并拼装指定框架的生成 prompt；两段都缺失时返回空串。
func LoadPrompt(framework string) string {
	base := readPromptFile(filepath.Join("base", "generate.txt"))
	fw := readPromptFile(filepath.Join(framework, "generate.txt"))
	return joinPrompt(base, fw)
}

// joinPrompt 公共段在前、框架段在后，空段自动跳过
func joinPrompt(base, framework string) string {
	base = strings.TrimSpace(base)
	framework = strings.TrimSpace(framework)
	switch {
	case base == "":
		return framework
	case framework == "":
		return base
	default:
		return base + "\n\n" + framework
	}
}

// readPromptFile 依次尝试多个相对路径，兼容 go run（frontend_api 根目录）与测试（包目录）
func readPromptFile(rel string) string {
	paths := []string{
		filepath.Join("prompts", rel),
		filepath.Join("..", "prompts", rel),
		filepath.Join("..", "..", "prompts", rel),
		filepath.Join("..", "..", "..", "prompts", rel),
	}
	for _, p := range paths {
		if data, err := os.ReadFile(p); err == nil {
			return string(data)
		}
	}
	return ""
}
