package skills

import "context"

// ═══════════════════════════════════════════════
//  Skill 通用接口
//  每个 Skill 负责一个独立的能力单元
// ═══════════════════════════════════════════════

// Skill 技能接口
type Skill interface {
	// Name 返回技能名称
	Name() string
	// Execute 执行技能，返回输出
	Execute(ctx context.Context, input interface{}) (interface{}, error)
}
