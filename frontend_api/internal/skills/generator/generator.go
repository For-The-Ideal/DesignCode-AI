package generator

import "context"

// ═══════════════════════════════════════════════
//  GeneratorSkill 通用类型定义
// ═══════════════════════════════════════════════

// Input 生成器通用输入
type Input struct {
	DSL    string `json:"dsl"`
	Prompt string `json:"prompt,omitempty"` // 从 prompts/{target}/generate.txt 加载
}

// Output 生成器通用输出
type Output struct {
	Code    string `json:"code"`
	Preview string `json:"preview"`
	Score   int    `json:"score"`
}

// Skill 生成器技能通用接口
type Skill interface {
	Name() string
	Execute(ctx context.Context, input interface{}) (interface{}, error)
}

// GeneratorFunc 将 DSL + target 映射到对应的生成器技能
func GeneratorFunc(target string) Skill {
	switch target {
	case "flutter":
		return NewFlutterSkill()
	case "vue3":
		return NewVue3Skill()
	case "react":
		return NewReactSkill()
	default:
		return NewFlutterSkill()
	}
}
